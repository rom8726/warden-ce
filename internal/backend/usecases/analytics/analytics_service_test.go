package analytics

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/rom8726/warden/internal/domain"
	mockcontract "github.com/rom8726/warden/test_mocks/internal_/backend/contract"
)

func TestNew(t *testing.T) {
	t.Parallel()

	// Create mocks
	mockReleaseRepo := mockcontract.NewMockReleaseRepository(t)
	mockReleaseStatsRepo := mockcontract.NewMockReleaseStatsRepository(t)
	mockEventRepo := mockcontract.NewMockEventRepository(t)
	mockProjectsRepo := mockcontract.NewMockProjectsRepository(t)
	mockIssuesRepo := mockcontract.NewMockIssuesRepository(t)

	// Create service
	service := New(
		mockReleaseRepo,
		mockReleaseStatsRepo,
		mockEventRepo,
		mockProjectsRepo,
		mockIssuesRepo,
	)

	// Verify service was created correctly
	require.NotNil(t, service)
	require.Equal(t, mockReleaseRepo, service.releaseRepo)
	require.Equal(t, mockReleaseStatsRepo, service.releaseStatsRepo)
	require.Equal(t, mockEventRepo, service.eventRepo)
	require.Equal(t, mockProjectsRepo, service.projectsRepo)
	require.Equal(t, mockIssuesRepo, service.issuesRepo)
}

func TestListReleases(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		setupMocks func(
			mockReleaseRepo *mockcontract.MockReleaseRepository,
			mockReleaseStatsRepo *mockcontract.MockReleaseStatsRepository,
		)
		projectID        domain.ProjectID
		expectedReleases []domain.Release
		expectedStats    map[string]domain.ReleaseStats
		expectedError    bool
		errorContains    string
	}{
		{
			name: "Success with stats",
			setupMocks: func(
				mockReleaseRepo *mockcontract.MockReleaseRepository,
				mockReleaseStatsRepo *mockcontract.MockReleaseStatsRepository,
			) {
				releases := []domain.Release{
					{
						ID:        1,
						ProjectID: 1,
						Version:   "1.0.0",
						CreatedAt: time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
					},
					{
						ID:        2,
						ProjectID: 1,
						Version:   "1.1.0",
						CreatedAt: time.Date(2023, 1, 2, 0, 0, 0, 0, time.UTC),
					},
				}

				mockReleaseRepo.EXPECT().ListByProject(
					mock.Anything,
					domain.ProjectID(1),
				).Return(releases, nil)

				stats1 := domain.ReleaseStats{
					ReleaseID:              1,
					Release:                "1.0.0",
					KnownIssuesTotal:       5,
					NewIssuesTotal:         3,
					ResolvedInVersionTotal: 2,
					UsersAffected:          10,
				}

				stats2 := domain.ReleaseStats{
					ReleaseID:              2,
					Release:                "1.1.0",
					KnownIssuesTotal:       8,
					NewIssuesTotal:         2,
					ResolvedInVersionTotal: 5,
					UsersAffected:          15,
				}

				mockReleaseStatsRepo.EXPECT().GetByProjectAndRelease(
					mock.Anything,
					domain.ProjectID(1),
					"1.0.0",
				).Return(stats1, nil)

				mockReleaseStatsRepo.EXPECT().GetByProjectAndRelease(
					mock.Anything,
					domain.ProjectID(1),
					"1.1.0",
				).Return(stats2, nil)
			},
			projectID: domain.ProjectID(1),
			expectedReleases: []domain.Release{
				{
					ID:        1,
					ProjectID: 1,
					Version:   "1.0.0",
					CreatedAt: time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
				},
				{
					ID:        2,
					ProjectID: 1,
					Version:   "1.1.0",
					CreatedAt: time.Date(2023, 1, 2, 0, 0, 0, 0, time.UTC),
				},
			},
			expectedStats: map[string]domain.ReleaseStats{
				"1.0.0": {
					ReleaseID:              1,
					Release:                "1.0.0",
					KnownIssuesTotal:       5,
					NewIssuesTotal:         3,
					ResolvedInVersionTotal: 2,
					UsersAffected:          10,
				},
				"1.1.0": {
					ReleaseID:              2,
					Release:                "1.1.0",
					KnownIssuesTotal:       8,
					NewIssuesTotal:         2,
					ResolvedInVersionTotal: 5,
					UsersAffected:          15,
				},
			},
			expectedError: false,
		},
		{
			name: "Success with missing stats",
			setupMocks: func(
				mockReleaseRepo *mockcontract.MockReleaseRepository,
				mockReleaseStatsRepo *mockcontract.MockReleaseStatsRepository,
			) {
				releases := []domain.Release{
					{
						ID:        1,
						ProjectID: 1,
						Version:   "1.0.0",
						CreatedAt: time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
					},
				}

				mockReleaseRepo.EXPECT().ListByProject(
					mock.Anything,
					domain.ProjectID(1),
				).Return(releases, nil)

				mockReleaseStatsRepo.EXPECT().GetByProjectAndRelease(
					mock.Anything,
					domain.ProjectID(1),
					"1.0.0",
				).Return(domain.ReleaseStats{}, errors.New("stats not found"))
			},
			projectID: domain.ProjectID(1),
			expectedReleases: []domain.Release{
				{
					ID:        1,
					ProjectID: 1,
					Version:   "1.0.0",
					CreatedAt: time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
				},
			},
			expectedStats: map[string]domain.ReleaseStats{},
			expectedError: false,
		},
		{
			name: "Database error",
			setupMocks: func(
				mockReleaseRepo *mockcontract.MockReleaseRepository,
				mockReleaseStatsRepo *mockcontract.MockReleaseStatsRepository,
			) {
				mockReleaseRepo.EXPECT().ListByProject(
					mock.Anything,
					domain.ProjectID(999),
				).Return(nil, errors.New("database error"))
			},
			projectID:        domain.ProjectID(999),
			expectedReleases: nil,
			expectedStats:    nil,
			expectedError:    true,
			errorContains:    "database error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// Create mocks
			mockReleaseRepo := mockcontract.NewMockReleaseRepository(t)
			mockReleaseStatsRepo := mockcontract.NewMockReleaseStatsRepository(t)
			mockEventRepo := mockcontract.NewMockEventRepository(t)
			mockProjectsRepo := mockcontract.NewMockProjectsRepository(t)
			mockIssuesRepo := mockcontract.NewMockIssuesRepository(t)

			// Setup mocks
			tt.setupMocks(mockReleaseRepo, mockReleaseStatsRepo)

			// Create service
			service := New(
				mockReleaseRepo,
				mockReleaseStatsRepo,
				mockEventRepo,
				mockProjectsRepo,
				mockIssuesRepo,
			)

			// Call method
			releases, stats, err := service.ListReleases(context.Background(), tt.projectID)

			// Check result
			if tt.expectedError {
				require.Error(t, err)
				if tt.errorContains != "" {
					require.Contains(t, err.Error(), tt.errorContains)
				}
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.expectedReleases, releases)
				require.Equal(t, tt.expectedStats, stats)
			}
		})
	}
}

func TestGetReleaseDetails(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		setupMocks func(
			mockReleaseRepo *mockcontract.MockReleaseRepository,
			mockReleaseStatsRepo *mockcontract.MockReleaseStatsRepository,
			mockEventRepo *mockcontract.MockEventRepository,
			mockIssuesRepo *mockcontract.MockIssuesRepository,
		)
		projectID                       domain.ProjectID
		version                         string
		topLimit                        uint
		expectedReleaseAnalyticsDetails domain.ReleaseAnalyticsDetails
		expectedError                   bool
		errorContains                   string
	}{
		{
			name: "Success",
			setupMocks: func(
				mockReleaseRepo *mockcontract.MockReleaseRepository,
				mockReleaseStatsRepo *mockcontract.MockReleaseStatsRepository,
				mockEventRepo *mockcontract.MockEventRepository,
				mockIssuesRepo *mockcontract.MockIssuesRepository,
			) {
				release := domain.Release{
					ID:        1,
					ProjectID: 1,
					Version:   "1.0.0",
					CreatedAt: time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
				}

				stats := domain.ReleaseStats{
					ReleaseID:              1,
					Release:                "1.0.0",
					KnownIssuesTotal:       5,
					NewIssuesTotal:         3,
					ResolvedInVersionTotal: 2,
					UsersAffected:          10,
				}

				topIssues := []string{"issue1", "issue2"}

				byPlatform := map[string]uint{"web": 8, "mobile": 2}
				byBrowser := map[string]uint{"chrome": 5, "firefox": 3}
				byOS := map[string]uint{"windows": 4, "macos": 3, "linux": 3}
				byDeviceArch := map[string]uint{"x86_64": 6, "arm64": 3, "x86": 2}
				byRuntimeName := map[string]uint{"node": 5, "python": 3, "java": 2, "go": 1}

				mockReleaseRepo.EXPECT().GetByProjectAndVersion(
					mock.Anything,
					domain.ProjectID(1),
					"1.0.0",
				).Return(release, nil)

				mockReleaseStatsRepo.EXPECT().GetByProjectAndRelease(
					mock.Anything,
					domain.ProjectID(1),
					"1.0.0",
				).Return(stats, nil)

				mockEventRepo.EXPECT().TopIssuesByRelease(
					mock.Anything,
					domain.ProjectID(1),
					"1.0.0",
					uint(10),
				).Return(topIssues, nil)

				mockEventRepo.EXPECT().AggregateBySegment(
					mock.Anything,
					domain.ProjectID(1),
					"1.0.0",
					domain.SegmentNamePlatform,
				).Return(byPlatform, nil)

				mockEventRepo.EXPECT().AggregateBySegment(
					mock.Anything,
					domain.ProjectID(1),
					"1.0.0",
					domain.SegmentNameBrowserName,
				).Return(byBrowser, nil)

				mockEventRepo.EXPECT().AggregateBySegment(
					mock.Anything,
					domain.ProjectID(1),
					"1.0.0",
					domain.SegmentNameOSName,
				).Return(byOS, nil)

				mockEventRepo.EXPECT().AggregateBySegment(
					mock.Anything,
					domain.ProjectID(1),
					"1.0.0",
					domain.SegmentNameDeviceArch,
				).Return(byDeviceArch, nil)

				mockEventRepo.EXPECT().AggregateBySegment(
					mock.Anything,
					domain.ProjectID(1),
					"1.0.0",
					domain.SegmentNameRuntimeName,
				).Return(byRuntimeName, nil)

				mockIssuesRepo.EXPECT().ListByFingerprints(
					mock.Anything,
					[]string{"issue1", "issue2"},
				).Return([]domain.Issue{{ID: 1, Title: "issue1"}, {ID: 2, Title: "issue2"}}, nil)
			},
			projectID: domain.ProjectID(1),
			version:   "1.0.0",
			topLimit:  10,
			expectedReleaseAnalyticsDetails: domain.ReleaseAnalyticsDetails{
				Release: domain.Release{
					ID:        1,
					ProjectID: 1,
					Version:   "1.0.0",
					CreatedAt: time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
				},
				Stats: domain.ReleaseStats{
					ReleaseID:              1,
					Release:                "1.0.0",
					KnownIssuesTotal:       5,
					NewIssuesTotal:         3,
					ResolvedInVersionTotal: 2,
					UsersAffected:          10,
				},
				TopIssues:     []domain.Issue{{ID: 1, Title: "issue1"}, {ID: 2, Title: "issue2"}},
				ByPlatform:    map[string]uint{"web": 8, "mobile": 2},
				ByBrowser:     map[string]uint{"chrome": 5, "firefox": 3},
				ByOS:          map[string]uint{"windows": 4, "macos": 3, "linux": 3},
				ByDeviceArch:  map[string]uint{"x86_64": 6, "arm64": 3, "x86": 2},
				ByRuntimeName: map[string]uint{"node": 5, "python": 3, "java": 2, "go": 1},
			},
			expectedError: false,
		},
		{
			name: "Release not found",
			setupMocks: func(
				mockReleaseRepo *mockcontract.MockReleaseRepository,
				mockReleaseStatsRepo *mockcontract.MockReleaseStatsRepository,
				mockEventRepo *mockcontract.MockEventRepository,
				mockIssuesRepo *mockcontract.MockIssuesRepository,
			) {
				mockReleaseRepo.EXPECT().GetByProjectAndVersion(
					mock.Anything,
					domain.ProjectID(1),
					"2.0.0",
				).Return(domain.Release{}, domain.ErrEntityNotFound)
			},
			projectID:                       domain.ProjectID(1),
			version:                         "2.0.0",
			topLimit:                        10,
			expectedReleaseAnalyticsDetails: domain.ReleaseAnalyticsDetails{},
			expectedError:                   true,
			errorContains:                   "not found",
		},
		{
			name: "Stats not found - should continue with empty stats",
			setupMocks: func(
				mockReleaseRepo *mockcontract.MockReleaseRepository,
				mockReleaseStatsRepo *mockcontract.MockReleaseStatsRepository,
				mockEventRepo *mockcontract.MockEventRepository,
				mockIssuesRepo *mockcontract.MockIssuesRepository,
			) {
				release := domain.Release{
					ID:        1,
					ProjectID: 1,
					Version:   "1.0.0",
					CreatedAt: time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
				}

				topIssues := []string{"issue1"}

				byPlatform := map[string]uint{"web": 5}
				byBrowser := map[string]uint{"chrome": 3}
				byOS := map[string]uint{"windows": 2}
				byDeviceArch := map[string]uint{"x86_64": 4, "arm64": 1}
				byRuntimeName := map[string]uint{"node": 3, "python": 2}

				mockReleaseRepo.EXPECT().GetByProjectAndVersion(
					mock.Anything,
					domain.ProjectID(1),
					"1.0.0",
				).Return(release, nil)

				mockReleaseStatsRepo.EXPECT().GetByProjectAndRelease(
					mock.Anything,
					domain.ProjectID(1),
					"1.0.0",
				).Return(domain.ReleaseStats{}, domain.ErrEntityNotFound)

				mockEventRepo.EXPECT().TopIssuesByRelease(
					mock.Anything,
					domain.ProjectID(1),
					"1.0.0",
					uint(10),
				).Return(topIssues, nil)

				mockEventRepo.EXPECT().AggregateBySegment(
					mock.Anything,
					domain.ProjectID(1),
					"1.0.0",
					domain.SegmentNamePlatform,
				).Return(byPlatform, nil)

				mockEventRepo.EXPECT().AggregateBySegment(
					mock.Anything,
					domain.ProjectID(1),
					"1.0.0",
					domain.SegmentNameBrowserName,
				).Return(byBrowser, nil)

				mockEventRepo.EXPECT().AggregateBySegment(
					mock.Anything,
					domain.ProjectID(1),
					"1.0.0",
					domain.SegmentNameOSName,
				).Return(byOS, nil)

				mockEventRepo.EXPECT().AggregateBySegment(
					mock.Anything,
					domain.ProjectID(1),
					"1.0.0",
					domain.SegmentNameDeviceArch,
				).Return(byDeviceArch, nil)

				mockEventRepo.EXPECT().AggregateBySegment(
					mock.Anything,
					domain.ProjectID(1),
					"1.0.0",
					domain.SegmentNameRuntimeName,
				).Return(byRuntimeName, nil)

				mockIssuesRepo.EXPECT().ListByFingerprints(
					mock.Anything,
					[]string{"issue1"},
				).Return([]domain.Issue{{ID: 1, Title: "issue1"}}, nil)
			},
			projectID: domain.ProjectID(1),
			version:   "1.0.0",
			topLimit:  10,
			expectedReleaseAnalyticsDetails: domain.ReleaseAnalyticsDetails{
				Release: domain.Release{
					ID:        1,
					ProjectID: 1,
					Version:   "1.0.0",
					CreatedAt: time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
				},
				Stats:         domain.ReleaseStats{},
				TopIssues:     []domain.Issue{{ID: 1, Title: "issue1"}},
				ByPlatform:    map[string]uint{"web": 5},
				ByBrowser:     map[string]uint{"chrome": 3},
				ByOS:          map[string]uint{"windows": 2},
				ByDeviceArch:  map[string]uint{"x86_64": 4, "arm64": 1},
				ByRuntimeName: map[string]uint{"node": 3, "python": 2},
			},
			expectedError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// Create mocks
			mockReleaseRepo := mockcontract.NewMockReleaseRepository(t)
			mockReleaseStatsRepo := mockcontract.NewMockReleaseStatsRepository(t)
			mockEventRepo := mockcontract.NewMockEventRepository(t)
			mockProjectsRepo := mockcontract.NewMockProjectsRepository(t)
			mockIssuesRepo := mockcontract.NewMockIssuesRepository(t)

			// Setup mocks
			tt.setupMocks(mockReleaseRepo, mockReleaseStatsRepo, mockEventRepo, mockIssuesRepo)

			// Create service
			service := New(
				mockReleaseRepo,
				mockReleaseStatsRepo,
				mockEventRepo,
				mockProjectsRepo,
				mockIssuesRepo,
			)

			// Call method
			result, err := service.GetReleaseDetails(
				context.Background(),
				tt.projectID,
				tt.version,
				tt.topLimit,
			)

			// Check result
			if tt.expectedError {
				require.Error(t, err)
				if tt.errorContains != "" {
					require.Contains(t, err.Error(), tt.errorContains)
				}
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.expectedReleaseAnalyticsDetails, result)
			}
		})
	}
}

func TestCompareReleases(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		setupMocks func(
			mockReleaseStatsRepo *mockcontract.MockReleaseStatsRepository,
		)
		projectID          domain.ProjectID
		baseVersion        string
		targetVersion      string
		expectedComparison domain.ReleaseComparison
		expectedError      bool
		errorContains      string
	}{
		{
			name: "Success",
			setupMocks: func(
				mockReleaseStatsRepo *mockcontract.MockReleaseStatsRepository,
			) {
				baseStats := domain.ReleaseStats{
					ReleaseID:              1,
					Release:                "1.0.0",
					KnownIssuesTotal:       5,
					NewIssuesTotal:         3,
					RegressionsTotal:       1,
					ResolvedInVersionTotal: 2,
					FixedNewInVersionTotal: 1,
					FixedOldInVersionTotal: 1,
					UsersAffected:          10,
				}

				targetStats := domain.ReleaseStats{
					ReleaseID:              2,
					Release:                "1.1.0",
					KnownIssuesTotal:       8,
					NewIssuesTotal:         2,
					RegressionsTotal:       0,
					ResolvedInVersionTotal: 5,
					FixedNewInVersionTotal: 2,
					FixedOldInVersionTotal: 3,
					UsersAffected:          15,
				}

				mockReleaseStatsRepo.EXPECT().GetByProjectAndRelease(
					mock.Anything,
					domain.ProjectID(1),
					"1.0.0",
				).Return(baseStats, nil)

				mockReleaseStatsRepo.EXPECT().GetByProjectAndRelease(
					mock.Anything,
					domain.ProjectID(1),
					"1.1.0",
				).Return(targetStats, nil)
			},
			projectID:     domain.ProjectID(1),
			baseVersion:   "1.0.0",
			targetVersion: "1.1.0",
			expectedComparison: domain.ReleaseComparison{
				BaseRelease: domain.ReleaseStats{
					ReleaseID:              1,
					Release:                "1.0.0",
					KnownIssuesTotal:       5,
					NewIssuesTotal:         3,
					RegressionsTotal:       1,
					ResolvedInVersionTotal: 2,
					FixedNewInVersionTotal: 1,
					FixedOldInVersionTotal: 1,
					UsersAffected:          10,
				},
				TargetRelease: domain.ReleaseStats{
					ReleaseID:              2,
					Release:                "1.1.0",
					KnownIssuesTotal:       8,
					NewIssuesTotal:         2,
					RegressionsTotal:       0,
					ResolvedInVersionTotal: 5,
					FixedNewInVersionTotal: 2,
					FixedOldInVersionTotal: 3,
					UsersAffected:          15,
				},
				Delta: map[string]uint{
					"known_issues_total":         3,
					"new_issues_total":           1,
					"regressions_total":          1,
					"resolved_in_version_total":  3,
					"fixed_new_in_version_total": 1,
					"fixed_old_in_version_total": 2,
					"users_affected":             5,
				},
			},
			expectedError: false,
		},
		{
			name: "Base release not found",
			setupMocks: func(
				mockReleaseStatsRepo *mockcontract.MockReleaseStatsRepository,
			) {
				mockReleaseStatsRepo.EXPECT().GetByProjectAndRelease(
					mock.Anything,
					domain.ProjectID(1),
					"1.0.0",
				).Return(domain.ReleaseStats{}, domain.ErrEntityNotFound)
			},
			projectID:          domain.ProjectID(1),
			baseVersion:        "1.0.0",
			targetVersion:      "1.1.0",
			expectedComparison: domain.ReleaseComparison{},
			expectedError:      true,
			errorContains:      "not found",
		},
		{
			name: "Target release not found",
			setupMocks: func(
				mockReleaseStatsRepo *mockcontract.MockReleaseStatsRepository,
			) {
				baseStats := domain.ReleaseStats{
					ReleaseID: 1,
					Release:   "1.0.0",
				}

				mockReleaseStatsRepo.EXPECT().GetByProjectAndRelease(
					mock.Anything,
					domain.ProjectID(1),
					"1.0.0",
				).Return(baseStats, nil)

				mockReleaseStatsRepo.EXPECT().GetByProjectAndRelease(
					mock.Anything,
					domain.ProjectID(1),
					"1.1.0",
				).Return(domain.ReleaseStats{}, domain.ErrEntityNotFound)
			},
			projectID:          domain.ProjectID(1),
			baseVersion:        "1.0.0",
			targetVersion:      "1.1.0",
			expectedComparison: domain.ReleaseComparison{},
			expectedError:      true,
			errorContains:      "not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// Create mocks
			mockReleaseRepo := mockcontract.NewMockReleaseRepository(t)
			mockReleaseStatsRepo := mockcontract.NewMockReleaseStatsRepository(t)
			mockEventRepo := mockcontract.NewMockEventRepository(t)
			mockProjectsRepo := mockcontract.NewMockProjectsRepository(t)
			mockIssuesRepo := mockcontract.NewMockIssuesRepository(t)

			// Setup mocks
			tt.setupMocks(mockReleaseStatsRepo)

			// Create service
			service := New(
				mockReleaseRepo,
				mockReleaseStatsRepo,
				mockEventRepo,
				mockProjectsRepo,
				mockIssuesRepo,
			)

			// Call method
			result, err := service.CompareReleases(
				context.Background(),
				tt.projectID,
				tt.baseVersion,
				tt.targetVersion,
			)

			// Check result
			if tt.expectedError {
				require.Error(t, err)
				if tt.errorContains != "" {
					require.Contains(t, err.Error(), tt.errorContains)
				}
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.expectedComparison, result)
			}
		})
	}
}

func TestGetUserSegments(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		setupMocks func(
			mockEventRepo *mockcontract.MockEventRepository,
		)
		projectID                     domain.ProjectID
		release                       string
		expectedUserSegmentsAnalytics domain.UserSegmentsAnalytics
		expectedError                 bool
		errorContains                 string
	}{
		{
			name: "Success",
			setupMocks: func(
				mockEventRepo *mockcontract.MockEventRepository,
			) {
				platform := map[string]uint{"web": 8, "android": 2}
				browser := map[string]uint{"Chrome": 5, "Firefox": 3}
				os := map[string]uint{"Windows": 4, "macOS": 3, "Linux": 3}
				deviceArch := map[string]uint{"x86_64": 6, "arm64": 3, "x86": 2}
				runtimeName := map[string]uint{"node": 5, "python": 3, "java": 2, "go": 1}

				mockEventRepo.EXPECT().AggregateBySegment(
					mock.Anything,
					domain.ProjectID(1),
					"1.0.0",
					domain.SegmentNamePlatform,
				).Return(platform, nil)

				mockEventRepo.EXPECT().AggregateBySegment(
					mock.Anything,
					domain.ProjectID(1),
					"1.0.0",
					domain.SegmentNameBrowserName,
				).Return(browser, nil)

				mockEventRepo.EXPECT().AggregateBySegment(
					mock.Anything,
					domain.ProjectID(1),
					"1.0.0",
					domain.SegmentNameOSName,
				).Return(os, nil)

				mockEventRepo.EXPECT().AggregateBySegment(
					mock.Anything,
					domain.ProjectID(1),
					"1.0.0",
					domain.SegmentNameDeviceArch,
				).Return(deviceArch, nil)

				mockEventRepo.EXPECT().AggregateBySegment(
					mock.Anything,
					domain.ProjectID(1),
					"1.0.0",
					domain.SegmentNameRuntimeName,
				).Return(runtimeName, nil)
			},
			projectID: domain.ProjectID(1),
			release:   "1.0.0",
			expectedUserSegmentsAnalytics: domain.UserSegmentsAnalytics{
				Platform: domain.UserSegmentsAggregation{
					domain.SegmentPlatformWeb:     8,
					domain.SegmentPlatformAndroid: 2,
				},
				Browser: domain.UserSegmentsAggregation{
					domain.SegmentBrowserChrome:  5,
					domain.SegmentBrowserFirefox: 3,
				},
				OS: domain.UserSegmentsAggregation{
					domain.SegmentOSWindows: 4,
					domain.SegmentOSMacOS:   3,
					domain.SegmentOSLinux:   3,
				},
				DeviceArch: domain.UserSegmentsAggregation{
					"x86_64": 6,
					"arm64":  3,
					"x86":    2,
				},
				RuntimeName: domain.UserSegmentsAggregation{
					"node":   5,
					"python": 3,
					"java":   2,
					"go":     1,
				},
			},
			expectedError: false,
		},
		{
			name: "Platform aggregation error",
			setupMocks: func(
				mockEventRepo *mockcontract.MockEventRepository,
			) {
				mockEventRepo.EXPECT().AggregateBySegment(
					mock.Anything,
					domain.ProjectID(1),
					"1.0.0",
					domain.SegmentNamePlatform,
				).Return(nil, errors.New("aggregation error"))
			},
			projectID:                     domain.ProjectID(1),
			release:                       "1.0.0",
			expectedUserSegmentsAnalytics: domain.UserSegmentsAnalytics{},
			expectedError:                 true,
			errorContains:                 "aggregation error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// Create mocks
			mockReleaseRepo := mockcontract.NewMockReleaseRepository(t)
			mockReleaseStatsRepo := mockcontract.NewMockReleaseStatsRepository(t)
			mockEventRepo := mockcontract.NewMockEventRepository(t)
			mockProjectsRepo := mockcontract.NewMockProjectsRepository(t)
			mockIssuesRepo := mockcontract.NewMockIssuesRepository(t)

			// Setup mocks
			tt.setupMocks(mockEventRepo)

			// Create service
			service := New(
				mockReleaseRepo,
				mockReleaseStatsRepo,
				mockEventRepo,
				mockProjectsRepo,
				mockIssuesRepo,
			)

			// Call method
			result, err := service.GetUserSegments(
				context.Background(),
				tt.projectID,
				tt.release,
			)

			// Check result
			if tt.expectedError {
				require.Error(t, err)
				if tt.errorContains != "" {
					require.Contains(t, err.Error(), tt.errorContains)
				}
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.expectedUserSegmentsAnalytics, result)
			}
		})
	}
}

func TestGetErrorsByTime(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		setupMocks func(
			mockEventRepo *mockcontract.MockEventRepository,
		)
		projectID          domain.ProjectID
		release            string
		period             time.Duration
		granularity        time.Duration
		levels             []domain.IssueLevel
		groupBy            domain.EventTimeseriesGroup
		expectedTimeseries []domain.Timeseries
		expectedError      bool
		errorContains      string
	}{
		{
			name: "Success with release filter",
			setupMocks: func(
				mockEventRepo *mockcontract.MockEventRepository,
			) {
				timeseries := []domain.Timeseries{
					{
						Name:        "error",
						Period:      domain.Period{Interval: 24 * time.Hour, Granularity: time.Hour},
						Occurrences: []uint{10, 15},
					},
					{
						Name:        "warning",
						Period:      domain.Period{Interval: 24 * time.Hour, Granularity: time.Hour},
						Occurrences: []uint{5, 8},
					},
				}

				mockEventRepo.EXPECT().Timeseries(
					mock.Anything,
					mock.MatchedBy(func(filter *domain.EventTimeseriesFilter) bool {
						return filter.ProjectID != nil && *filter.ProjectID == domain.ProjectID(1) &&
							len(filter.Levels) == 2 &&
							filter.Period.Interval == 24*time.Hour &&
							filter.Period.Granularity == time.Hour &&
							filter.GroupBy == domain.EventTimeseriesGroupLevel &&
							filter.Release != nil && *filter.Release == "1.0.0"
					}),
				).Return(timeseries, nil)
			},
			projectID:   domain.ProjectID(1),
			release:     "1.0.0",
			period:      24 * time.Hour,
			granularity: time.Hour,
			levels:      []domain.IssueLevel{domain.IssueLevelError, domain.IssueLevelWarning},
			groupBy:     domain.EventTimeseriesGroupLevel,
			expectedTimeseries: []domain.Timeseries{
				{
					Name:        "error",
					Period:      domain.Period{Interval: 24 * time.Hour, Granularity: time.Hour},
					Occurrences: []uint{10, 15},
				},
				{
					Name:        "warning",
					Period:      domain.Period{Interval: 24 * time.Hour, Granularity: time.Hour},
					Occurrences: []uint{5, 8},
				},
			},
			expectedError: false,
		},
		{
			name: "Success without release filter",
			setupMocks: func(
				mockEventRepo *mockcontract.MockEventRepository,
			) {
				timeseries := []domain.Timeseries{
					{
						Name:        "total",
						Period:      domain.Period{Interval: 7 * 24 * time.Hour, Granularity: 24 * time.Hour},
						Occurrences: []uint{20},
					},
				}

				mockEventRepo.EXPECT().Timeseries(
					mock.Anything,
					mock.MatchedBy(func(filter *domain.EventTimeseriesFilter) bool {
						return filter.ProjectID != nil && *filter.ProjectID == domain.ProjectID(1) &&
							filter.Release == nil
					}),
				).Return(timeseries, nil)
			},
			projectID:   domain.ProjectID(1),
			release:     "",
			period:      7 * 24 * time.Hour,
			granularity: 24 * time.Hour,
			levels:      nil,
			groupBy:     domain.EventTimeseriesGroupNone,
			expectedTimeseries: []domain.Timeseries{
				{
					Name:        "total",
					Period:      domain.Period{Interval: 7 * 24 * time.Hour, Granularity: 24 * time.Hour},
					Occurrences: []uint{20},
				},
			},
			expectedError: false,
		},
		{
			name: "Database error",
			setupMocks: func(
				mockEventRepo *mockcontract.MockEventRepository,
			) {
				mockEventRepo.EXPECT().Timeseries(
					mock.Anything,
					mock.Anything,
				).Return(nil, errors.New("database error"))
			},
			projectID:          domain.ProjectID(1),
			release:            "1.0.0",
			period:             24 * time.Hour,
			granularity:        time.Hour,
			levels:             []domain.IssueLevel{domain.IssueLevelError},
			groupBy:            domain.EventTimeseriesGroupLevel,
			expectedTimeseries: nil,
			expectedError:      true,
			errorContains:      "database error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// Create mocks
			mockReleaseRepo := mockcontract.NewMockReleaseRepository(t)
			mockReleaseStatsRepo := mockcontract.NewMockReleaseStatsRepository(t)
			mockEventRepo := mockcontract.NewMockEventRepository(t)
			mockProjectsRepo := mockcontract.NewMockProjectsRepository(t)
			mockIssuesRepo := mockcontract.NewMockIssuesRepository(t)

			// Setup mocks
			tt.setupMocks(mockEventRepo)

			// Create service
			service := New(
				mockReleaseRepo,
				mockReleaseStatsRepo,
				mockEventRepo,
				mockProjectsRepo,
				mockIssuesRepo,
			)

			// Call method
			result, err := service.GetErrorsByTime(
				context.Background(),
				tt.projectID,
				tt.release,
				tt.period,
				tt.granularity,
				tt.levels,
				tt.groupBy,
			)

			// Check result
			if tt.expectedError {
				require.Error(t, err)
				if tt.errorContains != "" {
					require.Contains(t, err.Error(), tt.errorContains)
				}
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.expectedTimeseries, result)
			}
		})
	}
}

func TestDiffUint(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		a        uint
		b        uint
		expected uint
	}{
		{
			name:     "Positive difference",
			a:        10,
			b:        5,
			expected: 5,
		},
		{
			name:     "Negative difference (absolute)",
			a:        5,
			b:        10,
			expected: 5,
		},
		{
			name:     "Equal values",
			a:        10,
			b:        10,
			expected: 0,
		},
		{
			name:     "Zero values",
			a:        0,
			b:        0,
			expected: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			result := diffUint(tt.a, tt.b)
			require.Equal(t, tt.expected, result)
		})
	}
}
