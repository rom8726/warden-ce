package jobs

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/rom8726/warden/internal/scheduler/contract"
	"github.com/rom8726/warden/internal/scheduler/scheduler"
)

const maxNotificationAge = time.Hour * 24 * 3

var _ scheduler.Job = (*NotificationsQueueCleanerJob)(nil)

type NotificationsQueueCleanerJob struct {
	repo contract.NotificationsQueueRepository
}

func NewNotificationsQueueCleaner(repo contract.NotificationsQueueRepository) *NotificationsQueueCleanerJob {
	return &NotificationsQueueCleanerJob{repo: repo}
}

func (n *NotificationsQueueCleanerJob) Name() string {
	return "notifications_queue_cleaner"
}

func (n *NotificationsQueueCleanerJob) Run(ctx context.Context) error {
	start := time.Now()
	slog.Info("run notifications queue cleaner job", "job", n.Name())

	cleaned, err := n.repo.DeleteOld(ctx, maxNotificationAge, 100)
	if err != nil {
		slog.Error("delete old notifications failed", "error", err, "job", n.Name())

		return fmt.Errorf("delete old notifications: %w", err)
	}

	slog.Info("DONE run notifications queue cleaner job", "duration",
		time.Since(start), "job", n.Name(), "cleaned", cleaned)

	return nil
}
