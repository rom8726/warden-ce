//nolint:gosec // it's ok here
package analytics

import (
	"context"
	"log/slog"
	"slices"
	"time"

	"github.com/rom8726/warden/internal/domain"
	"github.com/rom8726/warden/internal/scheduler/contract"
)

type AnalyticsService struct {
	releaseRepo      contract.ReleaseRepository
	releaseStatsRepo contract.ReleaseStatsRepository
	eventRepo        contract.EventRepository
	projectsRepo     contract.ProjectsRepository
	issuesRepo       contract.IssuesRepository
}

func New(
	releaseRepo contract.ReleaseRepository,
	releaseStatsRepo contract.ReleaseStatsRepository,
	eventRepo contract.EventRepository,
	projectsRepo contract.ProjectsRepository,
	issuesRepo contract.IssuesRepository,
) *AnalyticsService {
	return &AnalyticsService{
		releaseRepo:      releaseRepo,
		releaseStatsRepo: releaseStatsRepo,
		eventRepo:        eventRepo,
		projectsRepo:     projectsRepo,
		issuesRepo:       issuesRepo,
	}
}

//nolint:gocyclo // need refactoring
func (s *AnalyticsService) RecalculateReleaseStatsForAllProjects(ctx context.Context) error {
	projects, err := s.projectsRepo.List(ctx)
	if err != nil {
		return err
	}

	for _, project := range projects {
		releases, err := s.releaseRepo.ListByProject(ctx, project.ID)
		if err != nil {
			slog.Warn("analytics_stats: failed to list releases", "project", project.ID, "err", err)

			continue
		}
		for _, rel := range releases {
			// --- Known issues ---
			totalIssues, err := s.eventRepo.AggregateBySegment(ctx, project.ID, rel.Version, "group_hash")
			if err != nil {
				slog.Warn("analytics_stats: failed to count total issues", "release", rel.ID, "err", err)

				continue
			}
			knownIssuesTotal := uint(len(totalIssues))

			// --- New issues ---
			newIssues, err := s.issuesRepo.NewIssuesForRelease(ctx, project.ID, rel.Version)
			if err != nil {
				slog.Warn("analytics_stats: failed to get new issues", "release", rel.ID, "err", err)

				continue
			}
			newIssuesTotal := uint(len(newIssues))

			// --- Resolved in version ---
			resolvedIssues, err := s.issuesRepo.ResolvedInRelease(ctx, project.ID, rel.Version)
			if err != nil {
				slog.Warn("analytics_stats: failed to get resolved issues", "release", rel.ID, "err", err)

				continue
			}
			resolvedInVersionTotal := uint(len(resolvedIssues))

			// --- Fixed new/old in version ---
			fixedNew := 0
			fixedOld := 0
			newSet := make(map[string]struct{}, len(newIssues))
			for _, n := range newIssues {
				newSet[n] = struct{}{}
			}
			for _, rid := range resolvedIssues {
				if _, isNew := newSet[rid]; isNew {
					fixedNew++
				} else {
					fixedOld++
				}
			}

			// --- Fix times ---
			fixTimes, err := s.issuesRepo.FixTimesForRelease(ctx, project.ID, rel.Version)
			if err != nil {
				slog.Warn("analytics_stats: failed to get fix times", "release", rel.ID, "err", err)

				continue
			}

			var times []float64
			var avgFix, medianFix, p95Fix *time.Duration

			for _, d := range fixTimes {
				times = append(times, d.Seconds())
			}
			if len(times) > 0 {
				slices.Sort(times)
				total := 0.0
				for _, t := range times {
					total += t
				}

				avg := time.Duration(total/float64(len(times))) * time.Second
				avgFix = &avg

				median := time.Duration(times[len(times)/2]) * time.Second
				medianFix = &median

				p95 := time.Duration(times[int(float64(len(times))*0.95)]) * time.Second
				p95Fix = &p95
			}

			// --- Users affected ---
			usersAffectedAgg, err := s.eventRepo.AggregateBySegment(ctx, project.ID, rel.Version, "user_id")
			if err != nil {
				slog.Warn("analytics_stats: failed to count users affected", "release", rel.ID, "err", err)

				continue
			}
			usersAffected := uint(len(usersAffectedAgg))

			// --- Severity distribution ---
			severityAgg, err := s.eventRepo.AggregateBySegment(ctx, project.ID, rel.Version, "level")
			if err != nil {
				slog.Warn("analytics_stats: failed to aggregate severity", "release", rel.ID, "err", err)

				continue
			}

			// --- Regressions ---
			regressions, err := s.issuesRepo.RegressionsForRelease(ctx, project.ID, rel.Version)
			if err != nil {
				slog.Warn("analytics_stats: failed to get regressions", "release", rel.ID, "err", err)

				continue
			}
			regressionsTotal := uint(len(regressions))

			stats := domain.ReleaseStats{
				ProjectID:              rel.ProjectID,
				ReleaseID:              rel.ID,
				Release:                rel.Version,
				GeneratedAt:            time.Now(),
				KnownIssuesTotal:       knownIssuesTotal,
				NewIssuesTotal:         newIssuesTotal,
				RegressionsTotal:       regressionsTotal,
				ResolvedInVersionTotal: resolvedInVersionTotal,
				FixedNewInVersionTotal: uint(fixedNew),
				FixedOldInVersionTotal: uint(fixedOld),
				AvgFixTime:             avgFix,
				MedianFixTime:          medianFix,
				P95FixTime:             p95Fix,
				SeverityDistribution:   severityAgg,
				UsersAffected:          usersAffected,
			}
			err = s.releaseStatsRepo.Create(ctx, stats)
			if err != nil {
				slog.Warn("analytics_stats: failed to save stats", "release", rel.ID, "err", err)
			}
		}
	}

	return nil
}
