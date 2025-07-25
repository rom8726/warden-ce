package analytics

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/rom8726/warden/internal/domain"
	mockcontract "github.com/rom8726/warden/test_mocks/internal_/scheduler/contract"
)

func TestRecalculateReleaseStatsForAllProjects(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		setupMocks func(
			mockProjectsRepo *mockcontract.MockProjectsRepository,
			mockReleaseRepo *mockcontract.MockReleaseRepository,
			mockEventRepo *mockcontract.MockEventRepository,
			mockIssuesRepo *mockcontract.MockIssuesRepository,
			mockReleaseStatsRepo *mockcontract.MockReleaseStatsRepository,
		)
		expectedError bool
		errorContains string
	}{
		{
			name: "Success",
			setupMocks: func(
				mockProjectsRepo *mockcontract.MockProjectsRepository,
				mockReleaseRepo *mockcontract.MockReleaseRepository,
				mockEventRepo *mockcontract.MockEventRepository,
				mockIssuesRepo *mockcontract.MockIssuesRepository,
				mockReleaseStatsRepo *mockcontract.MockReleaseStatsRepository,
			) {
				projects := []domain.ProjectExtended{
					{Project: domain.Project{ID: 1, Name: "Project 1"}},
					{Project: domain.Project{ID: 2, Name: "Project 2"}},
				}

				releases := []domain.Release{
					{ID: 1, ProjectID: 1, Version: "1.0.0"},
					{ID: 2, ProjectID: 1, Version: "1.1.0"},
				}

				// Mock projects list
				mockProjectsRepo.EXPECT().List(mock.Anything).Return(projects, nil)

				// Mock releases for project 1
				mockReleaseRepo.EXPECT().ListByProject(mock.Anything, domain.ProjectID(1)).Return(releases, nil)

				// Mock total issues aggregation
				totalIssues := map[string]uint{"issue1": 5, "issue2": 3}
				mockEventRepo.EXPECT().AggregateBySegment(
					mock.Anything,
					domain.ProjectID(1),
					"1.0.0",
					domain.SegmentName("group_hash"),
				).Return(totalIssues, nil)

				// Mock new issues
				newIssues := []string{"issue1", "issue2"}
				mockIssuesRepo.EXPECT().NewIssuesForRelease(
					mock.Anything,
					domain.ProjectID(1),
					"1.0.0",
				).Return(newIssues, nil)

				// Mock resolved issues
				resolvedIssues := []string{"issue1"}
				mockIssuesRepo.EXPECT().ResolvedInRelease(
					mock.Anything,
					domain.ProjectID(1),
					"1.0.0",
				).Return(resolvedIssues, nil)

				// Mock fix times
				fixTimes := map[string]time.Duration{"issue1": time.Hour, "issue2": 2 * time.Hour}
				mockIssuesRepo.EXPECT().FixTimesForRelease(
					mock.Anything,
					domain.ProjectID(1),
					"1.0.0",
				).Return(fixTimes, nil)

				// Mock users affected
				usersAffected := map[string]uint{"user1": 3, "user2": 2}
				mockEventRepo.EXPECT().AggregateBySegment(
					mock.Anything,
					domain.ProjectID(1),
					"1.0.0",
					domain.SegmentName("user_id"),
				).Return(usersAffected, nil)

				// Mock severity distribution
				severityDist := map[string]uint{"error": 5, "warning": 3}
				mockEventRepo.EXPECT().AggregateBySegment(
					mock.Anything,
					domain.ProjectID(1),
					"1.0.0",
					domain.SegmentName("level"),
				).Return(severityDist, nil)

				// Mock regressions
				regressions := []string{"regression1"}
				mockIssuesRepo.EXPECT().RegressionsForRelease(
					mock.Anything,
					domain.ProjectID(1),
					"1.0.0",
				).Return(regressions, nil)

				// Mock stats creation
				mockReleaseStatsRepo.EXPECT().Create(
					mock.Anything,
					mock.MatchedBy(func(stats domain.ReleaseStats) bool {
						return stats.ProjectID == 1 &&
							stats.ReleaseID == 1 &&
							stats.Release == "1.0.0" &&
							stats.KnownIssuesTotal == 2 &&
							stats.NewIssuesTotal == 2 &&
							stats.RegressionsTotal == 1 &&
							stats.ResolvedInVersionTotal == 1 &&
							stats.FixedNewInVersionTotal == 1 &&
							stats.FixedOldInVersionTotal == 0 &&
							stats.UsersAffected == 2
					}),
				).Return(nil)

				// Mock releases for project 2 (empty)
				mockReleaseRepo.EXPECT().ListByProject(mock.Anything, domain.ProjectID(2)).Return([]domain.Release{}, nil)

				// Моки для релиза 1.1.0 (аналогично 1.0.0, но с другими значениями)
				mockEventRepo.EXPECT().AggregateBySegment(
					mock.Anything,
					domain.ProjectID(1),
					"1.1.0",
					domain.SegmentName("group_hash"),
				).Return(map[string]uint{"issue3": 2}, nil)

				mockIssuesRepo.EXPECT().NewIssuesForRelease(
					mock.Anything,
					domain.ProjectID(1),
					"1.1.0",
				).Return([]string{"issue3"}, nil)

				mockIssuesRepo.EXPECT().ResolvedInRelease(
					mock.Anything,
					domain.ProjectID(1),
					"1.1.0",
				).Return([]string{"issue3"}, nil)

				mockIssuesRepo.EXPECT().FixTimesForRelease(
					mock.Anything,
					domain.ProjectID(1),
					"1.1.0",
				).Return(map[string]time.Duration{"issue3": 3 * time.Hour}, nil)

				mockEventRepo.EXPECT().AggregateBySegment(
					mock.Anything,
					domain.ProjectID(1),
					"1.1.0",
					domain.SegmentName("user_id"),
				).Return(map[string]uint{"user3": 1}, nil)

				mockEventRepo.EXPECT().AggregateBySegment(
					mock.Anything,
					domain.ProjectID(1),
					"1.1.0",
					domain.SegmentName("level"),
				).Return(map[string]uint{"error": 2}, nil)

				mockIssuesRepo.EXPECT().RegressionsForRelease(
					mock.Anything,
					domain.ProjectID(1),
					"1.1.0",
				).Return([]string{}, nil)

				mockReleaseStatsRepo.EXPECT().Create(
					mock.Anything,
					mock.MatchedBy(func(stats domain.ReleaseStats) bool {
						return stats.ProjectID == 1 &&
							stats.ReleaseID == 2 &&
							stats.Release == "1.1.0" &&
							stats.KnownIssuesTotal == 1 &&
							stats.NewIssuesTotal == 1 &&
							stats.RegressionsTotal == 0 &&
							stats.ResolvedInVersionTotal == 1 &&
							stats.FixedNewInVersionTotal == 1 &&
							stats.FixedOldInVersionTotal == 0 &&
							stats.UsersAffected == 1
					}),
				).Return(nil)
			},
			expectedError: false,
		},
		{
			name: "Projects list error",
			setupMocks: func(
				mockProjectsRepo *mockcontract.MockProjectsRepository,
				mockReleaseRepo *mockcontract.MockReleaseRepository,
				mockEventRepo *mockcontract.MockEventRepository,
				mockIssuesRepo *mockcontract.MockIssuesRepository,
				mockReleaseStatsRepo *mockcontract.MockReleaseStatsRepository,
			) {
				mockProjectsRepo.EXPECT().List(mock.Anything).Return(nil, errors.New("database error"))
			},
			expectedError: true,
			errorContains: "database error",
		},
		{
			name: "Releases list error - should continue",
			setupMocks: func(
				mockProjectsRepo *mockcontract.MockProjectsRepository,
				mockReleaseRepo *mockcontract.MockReleaseRepository,
				mockEventRepo *mockcontract.MockEventRepository,
				mockIssuesRepo *mockcontract.MockIssuesRepository,
				mockReleaseStatsRepo *mockcontract.MockReleaseStatsRepository,
			) {
				projects := []domain.ProjectExtended{
					{Project: domain.Project{ID: 1, Name: "Project 1"}},
				}

				mockProjectsRepo.EXPECT().List(mock.Anything).Return(projects, nil)
				mockReleaseRepo.EXPECT().ListByProject(mock.Anything, domain.ProjectID(1)).Return(nil, errors.New("releases error"))
			},
			expectedError: false, // Should continue despite error
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
			tt.setupMocks(mockProjectsRepo, mockReleaseRepo, mockEventRepo, mockIssuesRepo, mockReleaseStatsRepo)

			// Create service
			service := New(
				mockReleaseRepo,
				mockReleaseStatsRepo,
				mockEventRepo,
				mockProjectsRepo,
				mockIssuesRepo,
			)

			// Call method
			err := service.RecalculateReleaseStatsForAllProjects(context.Background())

			// Check result
			if tt.expectedError {
				require.Error(t, err)
				if tt.errorContains != "" {
					require.Contains(t, err.Error(), tt.errorContains)
				}
			} else {
				require.NoError(t, err)
			}
		})
	}
}
