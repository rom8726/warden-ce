package notificator

import (
	"context"
	"encoding/json"
	"log/slog"
	"sync"
	"time"

	"github.com/rom8726/di"

	"github.com/rom8726/warden/internal/domain"
	"github.com/rom8726/warden/internal/user-notificator/contract"
	"github.com/rom8726/warden/pkg/resilience"
)

const (
	defaultBatchSize   = 100
	defaultInterval    = time.Minute
	defaultWorkerCount = 4
)

var _ di.Servicer = (*Service)(nil)

type Channel interface {
	Type() domain.NotificationType
	Send(
		ctx context.Context,
		issue *domain.Issue,
		project *domain.Project,
		config json.RawMessage,
		isRegress bool,
	) error
}

type notificationResult struct {
	notificationID domain.UserNotificationID
	skipped        bool
	skipReason     string
	err            error
}

type Service struct {
	userNotificationsUseCase contract.UserNotificationsUseCase
	emailService             contract.Emailer
	userRepo                 contract.UsersRepository

	ctx       context.Context
	cancelCtx func()

	batchSize   uint
	interval    time.Duration
	workerCount int

	circuitBreaker resilience.CircuitBreaker
}

func New(
	userNotificationsUseCase contract.UserNotificationsUseCase,
	emailService contract.Emailer,
	userRepo contract.UsersRepository,
	workerCount int,
) *Service {
	if workerCount == 0 {
		workerCount = defaultWorkerCount
	}

	ctx, cancel := context.WithCancel(context.Background())

	return &Service{
		userNotificationsUseCase: userNotificationsUseCase,
		emailService:             emailService,
		userRepo:                 userRepo,
		ctx:                      ctx,
		cancelCtx:                cancel,
		batchSize:                defaultBatchSize,
		interval:                 defaultInterval,
		workerCount:              max2Ints(workerCount, 1),
		circuitBreaker:           resilience.NewNotificationCircuitBreaker(),
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

// ProcessOutbox processes pending email notifications in the outbox.
func (s *Service) ProcessOutbox(ctx context.Context) {
	for {
		if s.ctx.Err() != nil {
			slog.Error("context error", "error", s.ctx.Err())

			break
		}

		sent := 0
		notifications, err := s.userNotificationsUseCase.TakePendingEmailNotifications(ctx, s.batchSize)
		if err != nil {
			slog.Error("take pending email notifications failed", "error", err)

			break
		}

		if len(notifications) == 0 {
			slog.Debug("no pending email notifications")

			break
		}

		slog.Debug("got pending email notifications", "count", len(notifications))

		// Create channels for parallel processing
		notificationChan := make(chan *domain.UserNotification, len(notifications))
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
				err = s.userNotificationsUseCase.MarkEmailAsFailed(ctx, result.notificationID, result.skipReason)
				if err != nil {
					slog.Error("mark notification as failed",
						"error", err, "notification_id", result.notificationID)
				}
			} else {
				sent++
			}
		}

		if sent > 0 {
			slog.Info("sent email notifications", "sent", sent)
		}
	}
}

func (s *Service) checkAndNotify(
	ctx context.Context,
	notification *domain.UserNotification,
) (skipped bool, skipReason string) {
	// Get user to get email
	user, err := s.userRepo.GetByID(ctx, notification.UserID)
	if err != nil {
		slog.Error("user not found for notification", "user_id", notification.UserID, "error", err)

		return true, "user not found"
	}
	if user.Email == "" {
		slog.Warn("user has no email, skip notification", "user_id", notification.UserID)

		return true, "user has no email"
	}

	var content domain.UserNotificationContent
	if err := json.Unmarshal(notification.Content, &content); err != nil {
		slog.Error("failed to unmarshal notification content", "error", err, "notification_id", notification.ID)

		return true, "invalid content format"
	}

	err = resilience.WithCircuitBreakerAndRetry(
		ctx,
		s.circuitBreaker,
		func(ctx context.Context) error {
			return s.emailService.SendUserNotificationEmail(ctx, user.Email, notification.Type, content)
		},
		resilience.DefaultRetryOptions()...,
	)

	if err != nil {
		slog.Error("send email notification failed",
			"error", err, "notification_id", notification.ID)

		err = s.userNotificationsUseCase.MarkEmailAsFailed(ctx, notification.ID, err.Error())
		if err != nil {
			slog.Error("mark notification as failed",
				"error", err, "notification_id", notification.ID)
		}
	} else {
		slog.Debug("sent email notification",
			"notification_id", notification.ID)

		err = s.userNotificationsUseCase.MarkEmailAsSent(ctx, notification.ID)
		if err != nil {
			slog.Error("mark notification as sent failed",
				"error", err, "notification_id", notification.ID)
		}
	}

	return false, ""
}

func (s *Service) worker(
	ctx context.Context,
	wg *sync.WaitGroup,
	notificationChan <-chan *domain.UserNotification,
	resultChan chan<- notificationResult,
) {
	defer wg.Done()

	for notification := range notificationChan {
		skipped, skipReason := s.checkAndNotify(ctx, notification)
		resultChan <- notificationResult{
			notificationID: notification.ID,
			skipped:        skipped,
			skipReason:     skipReason,
		}
	}
}

func max2Ints(a, b int) int {
	if a > b {
		return a
	}

	return b
}
