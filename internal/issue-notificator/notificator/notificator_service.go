package notificator

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"math"
	"sync"
	"time"

	"github.com/rom8726/di"

	"github.com/rom8726/warden/internal/domain"
	"github.com/rom8726/warden/internal/issue-notificator/contract"
	"github.com/rom8726/warden/pkg/db"
	"github.com/rom8726/warden/pkg/resilience"
)

const (
	defaultBatchSize   = 100
	defaultInterval    = time.Minute
	defaultWorkerCount = 4
)

var _ di.Servicer = (*Service)(nil)

type notificationResult struct {
	notificationID domain.NotificationID
	skipped        bool
	skipReason     string
	err            error
}

type Service struct {
	txManager db.TxManager

	channelsMap          map[domain.NotificationType]Channel
	notificationsUseCase contract.NotificationsUseCase
	issuesRepo           contract.IssuesRepository
	projectsRepo         contract.ProjectsRepository

	ctx       context.Context
	cancelCtx func()

	batchSize   uint
	interval    time.Duration
	workerCount int

	circuitBreaker resilience.CircuitBreaker
}

func New(
	channels []Channel,
	txManager db.TxManager,
	notificationsUseCase contract.NotificationsUseCase,
	issuesRepo contract.IssuesRepository,
	projectsRepo contract.ProjectsRepository,
	workerCount int,
) *Service {
	if workerCount <= 0 {
		workerCount = defaultWorkerCount
	}

	ctx, cancel := context.WithCancel(context.Background())

	channelsMap := make(map[domain.NotificationType]Channel, len(channels))
	for i := range channels {
		channel := channels[i]
		channelsMap[channel.Type()] = channel
	}

	return &Service{
		channelsMap:          channelsMap,
		txManager:            txManager,
		notificationsUseCase: notificationsUseCase,
		issuesRepo:           issuesRepo,
		projectsRepo:         projectsRepo,
		ctx:                  ctx,
		cancelCtx:            cancel,
		batchSize:            defaultBatchSize,
		interval:             defaultInterval,
		workerCount:          max2Ints(workerCount, 1),
		circuitBreaker:       resilience.NewNotificationCircuitBreaker(),
	}
}

// Start starts the worker.
func (s *Service) Start(context.Context) error {
	go s.run()

	slog.Info("Worker started")

	return nil
}

// Stop stops the worker.
func (s *Service) Stop(context.Context) error {
	s.cancelCtx()

	return nil
}

// run is the main loop of the worker.
func (s *Service) run() {
	for {
		select {
		case <-s.ctx.Done():
			return
		case <-time.After(s.interval):
			s.ProcessOutbox(s.ctx)
		}
	}
}

// ProcessOutbox processes pending emails in the outbox.
//
//nolint:nestif // refactor!
func (s *Service) ProcessOutbox(ctx context.Context) {
	for {
		if s.ctx.Err() != nil {
			slog.Error("context error", "error", s.ctx.Err())

			break
		}

		sent := 0

		notifications, err := s.notificationsUseCase.TakePendingNotificationsWithSettings(ctx, s.batchSize)
		if err != nil {
			slog.Error("take pending notifications failed", "error", err)

			break
		}

		if len(notifications) == 0 {
			slog.Debug("no pending notifications")

			break
		}

		slog.Debug("got pending notifications", "count", len(notifications))

		// Create channels for parallel processing
		notificationChan := make(chan *domain.NotificationWithSettings, len(notifications))
		resultChan := make(chan notificationResult, len(notifications))

		// Start worker goroutines
		var wg sync.WaitGroup
		for i := 0; i < s.workerCount; i++ {
			wg.Add(1)
			go s.worker(ctx, &wg, notificationChan, resultChan)
		}

		go func() {
			defer close(resultChan)

			// Send notifications to workers
			for i := range notifications {
				notification := notifications[i]
				notificationChan <- &notification
			}
			close(notificationChan)

			// Wait for all workers to complete
			wg.Wait()
		}()

		// Process results
		for result := range resultChan {
			if result.err != nil {
				slog.Error("check and notify failed",
					"error", result.err, "notification_id", result.notificationID)

				continue
			}

			if result.skipped {
				slog.Debug("notification skipped",
					"notification_id", result.notificationID, "reason", result.skipReason)
				err = s.notificationsUseCase.MarkNotificationAsSkipped(ctx, result.notificationID, result.skipReason)
				if err != nil {
					slog.Error("mark notification as skipped failed",
						"error", err, "notification_id", result.notificationID)
				}
			} else {
				sent++
			}
		}

		if sent > 0 {
			slog.Info("sent notifications", "sent", sent)
		}
	}
}

func (s *Service) SendTestNotification(
	ctx context.Context,
	notificationSettingID domain.NotificationSettingID,
) error {
	issue := domain.Issue{
		ID:          math.MaxUint,
		ProjectID:   math.MaxUint,
		Fingerprint: "1234567890",
		Source:      domain.SourceEvent,
		Status:      domain.IssueStatusUnresolved,
		Title:       "Test notification",
		Level:       domain.IssueLevelError,
		Platform:    "elixir",
		FirstSeen:   time.Now(),
		LastSeen:    time.Now(),
		TotalEvents: 1,
	}

	project := domain.Project{
		ID:   math.MaxUint,
		Name: "Test project",
	}

	setting, err := s.notificationsUseCase.GetNotificationSetting(ctx, notificationSettingID)
	if err != nil {
		return fmt.Errorf("get notification setting: %w", err)
	}

	channel := s.channelsMap[setting.Type]
	if channel == nil {
		return errors.New("channel not found")
	}

	err = channel.Send(ctx, &issue, &project, setting.Config, false)
	if err != nil {
		return fmt.Errorf("send notification failed: %w", err)
	}

	return nil
}

func (s *Service) checkAndNotify(
	ctx context.Context,
	notification *domain.NotificationWithSettings,
) (skipped bool, skipReason string, err error) {
	issue, err := s.issuesRepo.GetByID(ctx, notification.IssueID)
	if err != nil {
		slog.Error("get issue failed", "error", err, "issue_id", notification.IssueID)
		if errors.Is(err, domain.ErrEntityNotFound) {
			return true, "issue not found", nil
		}

		return false, "", err
	}

	project, err := s.projectsRepo.GetByID(ctx, issue.ProjectID)
	if err != nil {
		slog.Error("get project failed", "error", err, "project_id", issue.ProjectID)
		if errors.Is(err, domain.ErrEntityNotFound) {
			return true, "project not found", nil
		}

		return false, "", err
	}

	settings := filterSettings(notification)
	if len(settings) == 0 {
		return true, "no settings", nil
	}

	for _, setting := range settings {
		channel := s.channelsMap[setting.Type]
		if channel == nil {
			continue
		}

		err := resilience.WithCircuitBreakerAndRetry(
			ctx,
			s.circuitBreaker,
			func(ctx context.Context) error {
				return channel.Send(ctx, &issue, &project, setting.Config, notification.WasReactivated)
			},
			resilience.DefaultRetryOptions()...,
		)

		if err != nil {
			slog.Error("send notification failed",
				"error", err, "channel", channel.Type())

			err = s.notificationsUseCase.MarkNotificationAsFailed(ctx, notification.ID, err.Error())
			if err != nil {
				slog.Error("mark notification as failed",
					"error", err, "notification_id", notification.ID)
			}
		} else {
			slog.Debug("sent notification",
				"notification_id", notification.ID, "channel", channel.Type())

			err = s.notificationsUseCase.MarkNotificationAsSent(ctx, notification.ID)
			if err != nil {
				slog.Error("mark notification as sent failed",
					"error", err, "notification_id", notification.ID)
			}
		}
	}

	return false, "", nil
}

func (s *Service) worker(
	ctx context.Context,
	wg *sync.WaitGroup,
	notificationChan <-chan *domain.NotificationWithSettings,
	resultChan chan<- notificationResult,
) {
	defer wg.Done()

	for notification := range notificationChan {
		skipped, skipReason, err := s.checkAndNotify(ctx, notification)
		resultChan <- notificationResult{
			notificationID: notification.ID,
			skipped:        skipped,
			skipReason:     skipReason,
			err:            err,
		}
	}
}

func filterSettings(notification *domain.NotificationWithSettings) []domain.NotificationSetting {
	var availableSettings []domain.NotificationSetting

	for _, setting := range notification.Settings {
		if !setting.Enabled {
			continue
		}

		var ok bool
		for _, rule := range setting.Rules {
			if rule.EventLevel != "" && (rule.EventLevel != notification.Level) {
				continue
			}

			sendForNew := rule.IsNewError != nil && *rule.IsNewError && notification.IsNew
			sendForRegress := rule.IsRegression != nil && *rule.IsRegression && notification.WasReactivated

			ok = sendForNew || sendForRegress
			if ok {
				break
			}
		}

		if ok {
			availableSettings = append(availableSettings, setting)
		}
	}

	return availableSettings
}

func max2Ints(a, b int) int {
	if a > b {
		return a
	}

	return b
}
