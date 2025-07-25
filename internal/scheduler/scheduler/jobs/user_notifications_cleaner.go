package jobs

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/rom8726/warden/internal/scheduler/contract"
	"github.com/rom8726/warden/internal/scheduler/scheduler"
)

const maxUserNotificationAge = time.Hour * 24 * 30 // 30 days

var _ scheduler.Job = (*UserNotificationsCleanerJob)(nil)

type UserNotificationsCleanerJob struct {
	userNotificationsUseCase contract.UserNotificationsUseCase
}

func NewUserNotificationsCleaner(
	userNotificationsUseCase contract.UserNotificationsUseCase,
) *UserNotificationsCleanerJob {
	return &UserNotificationsCleanerJob{
		userNotificationsUseCase: userNotificationsUseCase,
	}
}

func (n *UserNotificationsCleanerJob) Name() string {
	return "user_notifications_cleaner"
}

func (n *UserNotificationsCleanerJob) Run(ctx context.Context) error {
	start := time.Now()
	slog.Info("run user notifications cleaner job", "job", n.Name())

	deleted, err := n.userNotificationsUseCase.DeleteOldNotifications(ctx, maxUserNotificationAge, 1000)
	if err != nil {
		slog.Error("delete old user notifications failed", "error", err, "job", n.Name())

		return fmt.Errorf("delete old user notifications: %w", err)
	}

	slog.Info("DONE run user notifications cleaner job", "duration",
		time.Since(start), "job", n.Name(), "deleted", deleted)

	return nil
}
