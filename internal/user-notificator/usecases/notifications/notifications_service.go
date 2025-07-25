package notifications

import (
	"context"
	"fmt"

	"github.com/rom8726/warden/internal/domain"
	"github.com/rom8726/warden/internal/user-notificator/contract"
)

type Service struct {
	userNotificationsRepo contract.UserNotificationsRepository
}

func New(
	userNotificationsRepo contract.UserNotificationsRepository,
) *Service {
	return &Service{
		userNotificationsRepo: userNotificationsRepo,
	}
}

func (s *Service) TakePendingEmailNotifications(ctx context.Context, limit uint) ([]domain.UserNotification, error) {
	notifications, err := s.userNotificationsRepo.GetPendingEmailNotifications(ctx, limit)
	if err != nil {
		return nil, fmt.Errorf("get pending email notifications: %w", err)
	}

	return notifications, nil
}

func (s *Service) MarkEmailAsSent(ctx context.Context, id domain.UserNotificationID) error {
	err := s.userNotificationsRepo.MarkEmailAsSent(ctx, id)
	if err != nil {
		return fmt.Errorf("mark email as sent: %w", err)
	}

	return nil
}

func (s *Service) MarkEmailAsFailed(ctx context.Context, id domain.UserNotificationID, reason string) error {
	err := s.userNotificationsRepo.MarkEmailAsFailed(ctx, id, reason)
	if err != nil {
		return fmt.Errorf("mark email as failed: %w", err)
	}

	return nil
}
