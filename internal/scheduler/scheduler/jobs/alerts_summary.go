package jobs

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/rom8726/warden/internal/scheduler/contract"
	"github.com/rom8726/warden/internal/scheduler/scheduler"
)

var _ scheduler.Job = (*AlertsSummary)(nil)

type AlertsSummary struct {
	issuesRepo contract.IssuesRepository
	emailer    contract.Emailer
}

func NewAlertsSummary(
	issuesRepo contract.IssuesRepository,
	emailer contract.Emailer,
) *AlertsSummary {
	return &AlertsSummary{
		issuesRepo: issuesRepo,
		emailer:    emailer,
	}
}

func (a AlertsSummary) Name() string {
	return "alerts_summary"
}

func (a AlertsSummary) Run(ctx context.Context) error {
	start := time.Now()
	slog.Info("run alerts summary job", "job", a.Name())

	list, err := a.issuesRepo.ListUnresolved(ctx)
	if err != nil {
		return fmt.Errorf("list unresolved issues: %w", err)
	}

	slog.Info("Fetched unresolved issues", "count", len(list), "job", a.Name())

	if len(list) > 0 {
		err = a.emailer.SendUnresolvedIssuesSummaryEmail(ctx, list)
		if err != nil {
			slog.Error("send unresolved issues summary emails failed", "error", err)

			return fmt.Errorf("send emails: %w", err)
		}
	}

	slog.Info("DONE run alerts summary job", "duration", time.Since(start), "job", a.Name())

	return nil
}
