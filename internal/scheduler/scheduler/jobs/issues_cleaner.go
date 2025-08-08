package jobs

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/rom8726/warden/internal/scheduler/contract"
	"github.com/rom8726/warden/internal/scheduler/scheduler"
)

const maxIssueAge = time.Hour * 24 * 90

var _ scheduler.Job = (*IssuesCleanerJob)(nil)

type IssuesCleanerJob struct {
	repo contract.IssuesRepository
}

func NewIssuesCleanerJob(repo contract.IssuesRepository) *IssuesCleanerJob {
	return &IssuesCleanerJob{repo: repo}
}

func (n *IssuesCleanerJob) Name() string {
	return "issues_cleaner"
}

func (n *IssuesCleanerJob) Run(ctx context.Context) error {
	start := time.Now()
	slog.Info("run issues cleaner job", "job", n.Name())

	cleaned, err := n.repo.DeleteOld(ctx, maxIssueAge, 100)
	if err != nil {
		slog.Error("delete old issues failed", "error", err, "job", n.Name())

		return fmt.Errorf("delete old issues: %w", err)
	}

	slog.Info("DONE run issues cleaner job", "duration",
		time.Since(start), "job", n.Name(), "cleaned", cleaned)

	return nil
}
