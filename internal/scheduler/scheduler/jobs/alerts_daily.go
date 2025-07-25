package jobs

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/rom8726/warden/internal/scheduler/contract"
	"github.com/rom8726/warden/internal/scheduler/scheduler"
)

var _ scheduler.Job = (*AlertsDailyJob)(nil)

type AlertsDailyJob struct {
	issuesRepo     contract.IssuesRepository
	notifQueueRepo contract.NotificationsQueueRepository

	strategy Strategy
}

func NewAlertsDaily(
	issuesRepo contract.IssuesRepository,
	notifQueueRepo contract.NotificationsQueueRepository,
) *AlertsDailyJob {
	return &AlertsDailyJob{
		issuesRepo:     issuesRepo,
		notifQueueRepo: notifQueueRepo,
		strategy:       newThroughOneStrategy(6),
	}
}

func (a *AlertsDailyJob) Name() string {
	return "alerts_daily"
}

func (a *AlertsDailyJob) Run(ctx context.Context) error {
	start := time.Now()
	slog.Info("run alerts daily job", "job", a.Name())

	list, err := a.issuesRepo.ListUnresolved(ctx)
	if err != nil {
		return fmt.Errorf("list unresolved issues: %w", err)
	}

	slog.Info("Fetched unresolved issues", "count", len(list), "job", a.Name())

	cnt := 0
	for i := range list {
		issue := list[i]
		isNew := issue.ResolvedAt == nil

		var send bool
		if issue.LastNotificationAt == nil {
			send = true
		} else {
			today := time.Now().Truncate(24 * time.Hour)
			lastSentAtDay := issue.LastNotificationAt.Truncate(24 * time.Hour)

			var firstDateDay time.Time
			if isNew { // new issue
				firstDateDay = issue.CreatedAt.Truncate(24 * time.Hour)
			} else { // regress issue
				firstDateDay = issue.ResolvedAt.Truncate(24 * time.Hour)
			}

			send = a.canSend(firstDateDay, lastSentAtDay, today)
		}

		if send {
			err = a.notifQueueRepo.AddNotification(
				ctx,
				issue.ProjectID,
				issue.ID,
				issue.Level,
				isNew,
				!isNew,
			)
			if err != nil {
				slog.Error("add notification to queue failed",
					"error", err, "issue_id", issue.ID, "job", a.Name())
			} else {
				cnt++
			}
		}
	}

	slog.Info("DONE run alerts daily job",
		"count", cnt, "elapsed", time.Since(start).String(), "job", a.Name())

	return nil
}

func (a *AlertsDailyJob) canSend(firstDay, lastSentAtDay, today time.Time) bool {
	if firstDay == lastSentAtDay || lastSentAtDay == today {
		return false
	}

	diff := int(today.Sub(firstDay).Hours() / 24)

	return a.strategy.Present(diff)
}
