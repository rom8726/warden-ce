package notifications

import (
	"context"
	"fmt"

	"github.com/rom8726/warden/internal/domain"
	"github.com/rom8726/warden/internal/issue-notificator/contract"
	"github.com/rom8726/warden/pkg/db"
)

type Service struct {
	txManager                db.TxManager
	notificationSettingsRepo contract.NotificationSettingsRepository
	notificationsQueueRepo   contract.NotificationsQueueRepository
	issuesRepo               contract.IssuesRepository
}

func New(
	txManager db.TxManager,
	notificationSettingsRepo contract.NotificationSettingsRepository,
	notificationsQueueRepo contract.NotificationsQueueRepository,
	issuesRepo contract.IssuesRepository,
) *Service {
	return &Service{
		txManager:                txManager,
		notificationSettingsRepo: notificationSettingsRepo,
		notificationsQueueRepo:   notificationsQueueRepo,
		issuesRepo:               issuesRepo,
	}
}

// GetNotificationSetting gets a notification setting by ID.
func (s *Service) GetNotificationSetting(
	ctx context.Context,
	id domain.NotificationSettingID,
) (domain.NotificationSetting, error) {
	setting, err := s.notificationSettingsRepo.GetSettingByID(ctx, id)
	if err != nil {
		return domain.NotificationSetting{}, fmt.Errorf("get notification setting: %w", err)
	}

	return setting, nil
}

func (s *Service) TakePendingNotificationsWithSettings(
	ctx context.Context,
	limit uint,
) ([]domain.NotificationWithSettings, error) {
	notifications, err := s.notificationsQueueRepo.TakePending(ctx, limit)
	if err != nil {
		return nil, fmt.Errorf("take pending notifications: %w", err)
	}

	projectsMap := make(map[domain.ProjectID][]domain.NotificationSetting)
	for i := range notifications {
		notification := notifications[i]
		projectsMap[notification.ProjectID] = nil
	}

	for projectID := range projectsMap {
		settingsWithRules, err := s.notificationSettingsRepo.ListSettings(ctx, projectID)
		if err != nil {
			return nil, fmt.Errorf("list notification settings: %w", err)
		}

		projectsMap[projectID] = settingsWithRules
	}

	result := make([]domain.NotificationWithSettings, 0, len(notifications))
	for i := range notifications {
		notification := notifications[i]
		settings := projectsMap[notification.ProjectID]
		result = append(result, domain.NotificationWithSettings{
			Notification: notification,
			Settings:     settings,
		})
	}

	return result, nil
}

func (s *Service) MarkNotificationAsSent(ctx context.Context, id domain.NotificationID) error {
	ntf, err := s.notificationsQueueRepo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("get notification by ID: %w", err)
	}

	return s.txManager.ReadCommitted(ctx, func(ctx context.Context) error {
		if err := s.issuesRepo.MarkAsNotified(ctx, ntf.IssueID); err != nil {
			return fmt.Errorf("mark issue as notified: %w", err)
		}

		return s.notificationsQueueRepo.MarkAsSent(ctx, id)
	})
}

func (s *Service) MarkNotificationAsFailed(ctx context.Context, id domain.NotificationID, reason string) error {
	if _, err := s.notificationsQueueRepo.GetByID(ctx, id); err != nil {
		return fmt.Errorf("get notification by ID: %w", err)
	}

	return s.notificationsQueueRepo.MarkAsFailed(ctx, id, reason)
}

func (s *Service) MarkNotificationAsSkipped(ctx context.Context, id domain.NotificationID, reason string) error {
	if _, err := s.notificationsQueueRepo.GetByID(ctx, id); err != nil {
		return fmt.Errorf("get notification by ID: %w", err)
	}

	return s.notificationsQueueRepo.MarkAsSkipped(ctx, id, reason)
}
