package usernotifications

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/rom8726/warden/internal/backend/contract"
	"github.com/rom8726/warden/internal/domain"
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

func (s *Service) CreateNotification(
	ctx context.Context,
	userID domain.UserID,
	notificationType domain.UserNotificationType,
	content domain.UserNotificationContent,
) error {
	contentJSON, err := json.Marshal(content)
	if err != nil {
		return fmt.Errorf("marshal notification content: %w", err)
	}

	_, err = s.userNotificationsRepo.Create(ctx, userID, notificationType, contentJSON)
	if err != nil {
		return fmt.Errorf("create notification: %w", err)
	}

	return nil
}

func (s *Service) GetUserNotifications(
	ctx context.Context,
	userID domain.UserID,
	limit, offset uint,
) ([]domain.UserNotification, error) {
	if limit == 0 {
		limit = 50 // Default limit
	}
	if limit > 100 {
		limit = 100 // Max limit
	}

	notifications, err := s.userNotificationsRepo.GetByUserID(ctx, userID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("get user notifications: %w", err)
	}

	return notifications, nil
}

func (s *Service) GetUnreadCount(ctx context.Context, userID domain.UserID) (uint, error) {
	count, err := s.userNotificationsRepo.GetUnreadCount(ctx, userID)
	if err != nil {
		return 0, fmt.Errorf("get unread count: %w", err)
	}

	return count, nil
}

func (s *Service) MarkAsRead(ctx context.Context, notificationID domain.UserNotificationID) error {
	_, err := s.userNotificationsRepo.GetByID(ctx, notificationID)
	if err != nil {
		return fmt.Errorf("get notification by ID: %w", err)
	}

	err = s.userNotificationsRepo.MarkAsRead(ctx, notificationID)
	if err != nil {
		return fmt.Errorf("mark as read: %w", err)
	}

	return nil
}

func (s *Service) MarkAllAsRead(ctx context.Context, userID domain.UserID) error {
	err := s.userNotificationsRepo.MarkAllAsRead(ctx, userID)
	if err != nil {
		return fmt.Errorf("mark all as read: %w", err)
	}

	return nil
}

func (s *Service) DeleteOldNotifications(
	ctx context.Context,
	maxAge time.Duration,
	limit uint,
) (uint, error) {
	deleted, err := s.userNotificationsRepo.DeleteOld(ctx, maxAge, limit)
	if err != nil {
		return 0, fmt.Errorf("delete old notifications: %w", err)
	}

	return deleted, nil
}
