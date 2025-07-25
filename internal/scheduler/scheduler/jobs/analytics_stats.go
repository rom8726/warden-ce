package jobs

import (
	"context"

	"github.com/rom8726/warden/internal/scheduler/contract"
	"github.com/rom8726/warden/internal/scheduler/scheduler"
)

var _ scheduler.Job = (*AnalyticsStatsJob)(nil)

type AnalyticsStatsJob struct {
	analyticsService contract.AnalyticsUseCase
}

func (job *AnalyticsStatsJob) Name() string {
	return "analytics_stats"
}

func NewAnalyticsStats(analyticsService contract.AnalyticsUseCase) *AnalyticsStatsJob {
	return &AnalyticsStatsJob{
		analyticsService: analyticsService,
	}
}

func (job *AnalyticsStatsJob) Run(ctx context.Context) error {
	return job.analyticsService.RecalculateReleaseStatsForAllProjects(ctx)
}
