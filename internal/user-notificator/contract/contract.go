package contract

import (
	"context"

	"github.com/rom8726/warden/internal/domain"
)

type UserNotificationsUseCase interface {
	TakePendingEmailNotifications(ctx context.Context, limit uint) ([]domain.UserNotification, error)
	MarkEmailAsSent(ctx context.Context, id domain.UserNotificationID) error
	MarkEmailAsFailed(ctx context.Context, id domain.UserNotificationID, reason string) error
}

type Emailer interface {
	SendUserNotificationEmail(
		ctx context.Context,
		toEmail string,
		notifType domain.UserNotificationType,
		content domain.UserNotificationContent,
	) error
}

type UsersRepository interface {
	GetByID(ctx context.Context, id domain.UserID) (domain.User, error)
}

type UserNotificationsRepository interface {
	GetPendingEmailNotifications(ctx context.Context, limit uint) ([]domain.UserNotification, error)
	MarkEmailAsSent(ctx context.Context, id domain.UserNotificationID) error
	MarkEmailAsFailed(ctx context.Context, id domain.UserNotificationID, reason string) error
}
