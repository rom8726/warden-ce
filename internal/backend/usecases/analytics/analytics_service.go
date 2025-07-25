//nolint:gosec // it's ok here
package analytics

import (
	"context"
	"errors"
	"time"

	"github.com/rom8726/warden/internal/backend/contract"
	"github.com/rom8726/warden/internal/domain"
)

type AnalyticsRelease struct {
	Release domain.Release
	Stats   *domain.ReleaseStats
}

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

// ListReleases returns a list of releases and a map of metrics by version.
func (s *AnalyticsService) ListReleases(
	ctx context.Context,
	projectID domain.ProjectID,
) ([]domain.Release, map[string]domain.ReleaseStats, error) {
	releases, err := s.releaseRepo.ListByProject(ctx, projectID)
	if err != nil {
		return nil, nil, err
	}
	statsMap := make(map[string]domain.ReleaseStats, len(releases))
	for _, rel := range releases {
		stats, err := s.releaseStatsRepo.GetByProjectAndRelease(ctx, projectID, rel.Version)
		if err == nil {
			statsMap[rel.Version] = stats
		}
	}

	return releases, statsMap, nil
}

// GetReleaseDetails returns release, metrics, top issues, and aggregations by platform/browser/OS.
func (s *AnalyticsService) GetReleaseDetails(
	ctx context.Context,
	projectID domain.ProjectID,
	version string,
	topIssuesLimit uint,
) (domain.ReleaseAnalyticsDetails, error) {
	release, err := s.releaseRepo.GetByProjectAndVersion(ctx, projectID, version)
	if err != nil {
		return domain.ReleaseAnalyticsDetails{}, err
	}
	stats, err := s.releaseStatsRepo.GetByProjectAndRelease(ctx, projectID, version)
	if errors.Is(err, domain.ErrEntityNotFound) {
		stats = domain.ReleaseStats{}
	} else if err != nil {
		return domain.ReleaseAnalyticsDetails{}, err
	}
	topIssuesFingerprints, err := s.eventRepo.TopIssuesByRelease(ctx, projectID, version, topIssuesLimit)
	if err != nil {
		return domain.ReleaseAnalyticsDetails{}, err
	}
	topIssues, err := s.issuesRepo.ListByFingerprints(ctx, topIssuesFingerprints)
	if err != nil {
		return domain.ReleaseAnalyticsDetails{}, err
	}
	byPlatform, err := s.eventRepo.AggregateBySegment(ctx, projectID, version, domain.SegmentNamePlatform)
	if err != nil {
		return domain.ReleaseAnalyticsDetails{}, err
	}
	byBrowser, err := s.eventRepo.AggregateBySegment(ctx, projectID, version, domain.SegmentNameBrowserName)
	if err != nil {
		return domain.ReleaseAnalyticsDetails{}, err
	}
	byOS, err := s.eventRepo.AggregateBySegment(ctx, projectID, version, domain.SegmentNameOSName)
	if err != nil {
		return domain.ReleaseAnalyticsDetails{}, err
	}
	byDeviceArch, err := s.eventRepo.AggregateBySegment(ctx, projectID, version, domain.SegmentNameDeviceArch)
	if err != nil {
		return domain.ReleaseAnalyticsDetails{}, err
	}
	byRuntimeName, err := s.eventRepo.AggregateBySegment(ctx, projectID, version, domain.SegmentNameRuntimeName)
	if err != nil {
		return domain.ReleaseAnalyticsDetails{}, err
	}

	return domain.ReleaseAnalyticsDetails{
		Release:       release,
		Stats:         stats,
		TopIssues:     topIssues,
		ByPlatform:    byPlatform,
		ByBrowser:     byBrowser,
		ByOS:          byOS,
		ByDeviceArch:  byDeviceArch,
		ByRuntimeName: byRuntimeName,
	}, nil
}

// CompareReleases returns comparison of two releases and delta by key metrics.
func (s *AnalyticsService) CompareReleases(
	ctx context.Context,
	projectID domain.ProjectID,
	baseVersion, targetVersion string,
) (domain.ReleaseComparison, error) {
	baseStats, err := s.releaseStatsRepo.GetByProjectAndRelease(ctx, projectID, baseVersion)
	if err != nil {
		return domain.ReleaseComparison{}, err
	}
	targetStats, err := s.releaseStatsRepo.GetByProjectAndRelease(ctx, projectID, targetVersion)
	if err != nil {
		return domain.ReleaseComparison{}, err
	}
	delta := make(map[string]uint)
	delta["known_issues_total"] = diffUint(targetStats.KnownIssuesTotal, baseStats.KnownIssuesTotal)
	delta["new_issues_total"] = diffUint(targetStats.NewIssuesTotal, baseStats.NewIssuesTotal)
	delta["regressions_total"] = diffUint(targetStats.RegressionsTotal, baseStats.RegressionsTotal)
	delta["resolved_in_version_total"] = diffUint(targetStats.ResolvedInVersionTotal, baseStats.ResolvedInVersionTotal)
	delta["fixed_new_in_version_total"] = diffUint(targetStats.FixedNewInVersionTotal, baseStats.FixedNewInVersionTotal)
	delta["fixed_old_in_version_total"] = diffUint(targetStats.FixedOldInVersionTotal, baseStats.FixedOldInVersionTotal)
	delta["users_affected"] = diffUint(targetStats.UsersAffected, baseStats.UsersAffected)

	return domain.ReleaseComparison{
		BaseRelease:   baseStats,
		TargetRelease: targetStats,
		Delta:         delta,
	}, nil
}

// GetErrorsByTime returns time series of errors for a release, with optional level and grouping.
func (s *AnalyticsService) GetErrorsByTime(
	ctx context.Context,
	projectID domain.ProjectID,
	release string,
	period, granularity time.Duration,
	levels []domain.IssueLevel,
	groupBy domain.EventTimeseriesGroup,
) ([]domain.Timeseries, error) {
	filter := &domain.EventTimeseriesFilter{
		ProjectID: &projectID,
		Levels:    levels,
		Period: domain.Period{
			Interval:    period,
			Granularity: granularity,
		},
		GroupBy: groupBy,
	}
	if release != "" {
		filter.Release = &release
	}

	return s.eventRepo.Timeseries(ctx, filter)
}

// GetUserSegments returns aggregations by platform, browser, OS for a release.
func (s *AnalyticsService) GetUserSegments(
	ctx context.Context,
	projectID domain.ProjectID,
	release string,
) (domain.UserSegmentsAnalytics, error) {
	platformRaw, err := s.eventRepo.AggregateBySegment(
		ctx,
		projectID,
		release,
		domain.SegmentNamePlatform,
	)
	if err != nil {
		return domain.UserSegmentsAnalytics{}, err
	}
	browserRaw, err := s.eventRepo.AggregateBySegment(
		ctx,
		projectID,
		release,
		domain.SegmentNameBrowserName,
	)
	if err != nil {
		return domain.UserSegmentsAnalytics{}, err
	}
	osRaw, err := s.eventRepo.AggregateBySegment(
		ctx,
		projectID,
		release,
		domain.SegmentNameOSName,
	)
	if err != nil {
		return domain.UserSegmentsAnalytics{}, err
	}
	deviceArchRaw, err := s.eventRepo.AggregateBySegment(
		ctx,
		projectID,
		release,
		domain.SegmentNameDeviceArch,
	)
	if err != nil {
		return domain.UserSegmentsAnalytics{}, err
	}
	runtimeNameRaw, err := s.eventRepo.AggregateBySegment(
		ctx,
		projectID,
		release,
		domain.SegmentNameRuntimeName,
	)
	if err != nil {
		return domain.UserSegmentsAnalytics{}, err
	}
	platform := make(domain.UserSegmentsAggregation, len(platformRaw))
	for k, v := range platformRaw {
		platform[domain.UserSegmentKey(k)] = v
	}
	browser := make(domain.UserSegmentsAggregation, len(browserRaw))
	for k, v := range browserRaw {
		browser[domain.UserSegmentKey(k)] = v
	}
	osName := make(domain.UserSegmentsAggregation, len(osRaw))
	for k, v := range osRaw {
		osName[domain.UserSegmentKey(k)] = v
	}
	deviceArch := make(domain.UserSegmentsAggregation, len(deviceArchRaw))
	for k, v := range deviceArchRaw {
		deviceArch[domain.UserSegmentKey(k)] = v
	}
	runtimeName := make(domain.UserSegmentsAggregation, len(runtimeNameRaw))
	for k, v := range runtimeNameRaw {
		runtimeName[domain.UserSegmentKey(k)] = v
	}

	return domain.UserSegmentsAnalytics{
		Platform:    platform,
		Browser:     browser,
		OS:          osName,
		DeviceArch:  deviceArch,
		RuntimeName: runtimeName,
	}, nil
}

func diffUint(a, b uint) uint {
	if a > b {
		return a - b
	}

	return b - a
}
