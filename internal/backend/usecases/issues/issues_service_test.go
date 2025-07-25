package issues

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	wardencontext "github.com/rom8726/warden/internal/context"
	"github.com/rom8726/warden/internal/domain"
	mockcontract "github.com/rom8726/warden/test_mocks/internal_/backend/contract"
	mockdb "github.com/rom8726/warden/test_mocks/pkg/db"
)

func TestNew(t *testing.T) {
	t.Parallel()

	// Create mocks
	mockTxManager := mockdb.NewMockTxManager(t)
	mockIssuesRepo := mockcontract.NewMockIssuesRepository(t)
	mockProjectsRepo := mockcontract.NewMockProjectsRepository(t)
	mockEventsRepo := mockcontract.NewMockEventRepository(t)
	mockProjectsService := mockcontract.NewMockProjectsUseCase(t)
	mockResolutionsRepo := mockcontract.NewMockResolutionsRepository(t)
	mockUsersRepo := mockcontract.NewMockUsersRepository(t)
	mockTeamsRepo := mockcontract.NewMockTeamsRepository(t)
	mockUserNotificationsUseCase := mockcontract.NewMockUserNotificationsUseCase(t)

	// Create service
	service := New(
		mockTxManager,
		mockIssuesRepo,
		mockProjectsRepo,
		mockEventsRepo,
		mockProjectsService,
		mockResolutionsRepo,
		mockUsersRepo,
		mockTeamsRepo,
		mockUserNotificationsUseCase,
	)

	// Verify service was created correctly
	require.NotNil(t, service)
	require.Equal(t, mockTxManager, service.txManager)
	require.Equal(t, mockIssuesRepo, service.issuesRepo)
	require.Equal(t, mockProjectsRepo, service.projectsRepo)
	require.Equal(t, mockEventsRepo, service.eventsRepo)
	require.Equal(t, mockProjectsService, service.projectsService)
	require.Equal(t, mockResolutionsRepo, service.resolutionsRepo)
	require.Equal(t, mockUsersRepo, service.usersRepo)
}

func TestList(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		filter         *domain.ListIssuesFilter
		setupMocks     func(mockIssuesRepo *mockcontract.MockIssuesRepository)
		expectedIssues []domain.IssueExtended
		expectedCount  uint64
		expectedError  bool
		errorContains  string
	}{
		{
			name: "Success",
			filter: &domain.ListIssuesFilter{
				ProjectID: func() *domain.ProjectID {
					id := domain.ProjectID(1)

					return &id
				}(),
				Level: func() *domain.IssueLevel {
					level := domain.IssueLevelError

					return &level
				}(),
				Status: func() *domain.IssueStatus {
					status := domain.IssueStatusUnresolved

					return &status
				}(),
				TimeFrom: time.Now().Add(-24 * time.Hour),
				TimeTo:   time.Now(),
				PageNum:  1,
				PerPage:  10,
			},
			setupMocks: func(mockIssuesRepo *mockcontract.MockIssuesRepository) {
				mockIssuesRepo.EXPECT().ListExtended(mock.Anything, mock.AnythingOfType("*domain.ListIssuesFilter")).
					Return([]domain.IssueExtended{
						{
							Issue: domain.Issue{
								ID:        1,
								ProjectID: 1,
								Title:     "Test Issue",
							},
							ProjectName: "Test Project",
						},
					}, uint64(1), nil)
			},
			expectedIssues: []domain.IssueExtended{
				{
					Issue: domain.Issue{
						ID:        1,
						ProjectID: 1,
						Title:     "Test Issue",
					},
					ProjectName: "Test Project",
				},
			},
			expectedCount: 1,
			expectedError: false,
		},
		{
			name:   "Error",
			filter: &domain.ListIssuesFilter{},
			setupMocks: func(mockIssuesRepo *mockcontract.MockIssuesRepository) {
				mockIssuesRepo.EXPECT().ListExtended(mock.Anything, mock.AnythingOfType("*domain.ListIssuesFilter")).
					Return(nil, uint64(0), errors.New("database error"))
			},
			expectedIssues: nil,
			expectedCount:  0,
			expectedError:  true,
			errorContains:  "database error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// Create mocks
			mockTxManager := mockdb.NewMockTxManager(t)
			mockIssuesRepo := mockcontract.NewMockIssuesRepository(t)
			mockProjectsRepo := mockcontract.NewMockProjectsRepository(t)
			mockEventsRepo := mockcontract.NewMockEventRepository(t)
			mockProjectsService := mockcontract.NewMockProjectsUseCase(t)
			mockResolutionsRepo := mockcontract.NewMockResolutionsRepository(t)
			mockUsersRepo := mockcontract.NewMockUsersRepository(t)
			mockTeamsRepo := mockcontract.NewMockTeamsRepository(t)
			mockUserNotificationsUseCase := mockcontract.NewMockUserNotificationsUseCase(t)

			// Setup mocks
			tt.setupMocks(mockIssuesRepo)

			// Create service
			service := New(
				mockTxManager,
				mockIssuesRepo,
				mockProjectsRepo,
				mockEventsRepo,
				mockProjectsService,
				mockResolutionsRepo,
				mockUsersRepo,
				mockTeamsRepo,
				mockUserNotificationsUseCase,
			)

			// Call the method
			issues, count, err := service.List(context.Background(), tt.filter)

			// Check results
			if tt.expectedError {
				require.Error(t, err)
				if tt.errorContains != "" {
					require.Contains(t, err.Error(), tt.errorContains)
				}
				require.Nil(t, issues)
				require.Equal(t, uint64(0), count)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.expectedIssues, issues)
				require.Equal(t, tt.expectedCount, count)
			}
		})
	}
}

func TestRecentIssues(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		limit          uint
		setupMocks     func(mockIssuesRepo *mockcontract.MockIssuesRepository)
		expectedIssues []domain.IssueExtended
		expectedError  bool
		errorContains  string
	}{
		{
			name:  "Success",
			limit: 10,
			setupMocks: func(mockIssuesRepo *mockcontract.MockIssuesRepository) {
				mockIssuesRepo.EXPECT().RecentIssues(mock.Anything, uint(10)).
					Return([]domain.IssueExtended{
						{
							Issue: domain.Issue{
								ID:        1,
								ProjectID: 1,
								Title:     "Test Issue",
							},
							ProjectName: "Test Project",
						},
					}, nil)
			},
			expectedIssues: []domain.IssueExtended{
				{
					Issue: domain.Issue{
						ID:        1,
						ProjectID: 1,
						Title:     "Test Issue",
					},
					ProjectName: "Test Project",
				},
			},
			expectedError: false,
		},
		{
			name:  "Error",
			limit: 10,
			setupMocks: func(mockIssuesRepo *mockcontract.MockIssuesRepository) {
				mockIssuesRepo.EXPECT().RecentIssues(mock.Anything, uint(10)).
					Return(nil, errors.New("database error"))
			},
			expectedIssues: nil,
			expectedError:  true,
			errorContains:  "database error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// Create mocks
			mockTxManager := mockdb.NewMockTxManager(t)
			mockIssuesRepo := mockcontract.NewMockIssuesRepository(t)
			mockProjectsRepo := mockcontract.NewMockProjectsRepository(t)
			mockEventsRepo := mockcontract.NewMockEventRepository(t)
			mockProjectsService := mockcontract.NewMockProjectsUseCase(t)
			mockResolutionsRepo := mockcontract.NewMockResolutionsRepository(t)
			mockUsersRepo := mockcontract.NewMockUsersRepository(t)
			mockTeamsRepo := mockcontract.NewMockTeamsRepository(t)
			mockUserNotificationsUseCase := mockcontract.NewMockUserNotificationsUseCase(t)

			// Setup mocks
			tt.setupMocks(mockIssuesRepo)

			// Create service
			service := New(
				mockTxManager,
				mockIssuesRepo,
				mockProjectsRepo,
				mockEventsRepo,
				mockProjectsService,
				mockResolutionsRepo,
				mockUsersRepo,
				mockTeamsRepo,
				mockUserNotificationsUseCase,
			)

			// Call the method
			issues, err := service.RecentIssues(context.Background(), tt.limit)

			// Check results
			if tt.expectedError {
				require.Error(t, err)
				if tt.errorContains != "" {
					require.Contains(t, err.Error(), tt.errorContains)
				}
				require.Nil(t, issues)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.expectedIssues, issues)
			}
		})
	}
}

func TestGetByIDWithChildren(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		issueID    domain.IssueID
		setupMocks func(
			mockIssuesRepo *mockcontract.MockIssuesRepository,
			mockUsersRepo *mockcontract.MockUsersRepository,
			mockProjectsService *mockcontract.MockProjectsUseCase,
			mockProjectsRepo *mockcontract.MockProjectsRepository,
			mockEventsRepo *mockcontract.MockEventRepository,
		)
		setupContext  func(ctx context.Context) context.Context
		expectedIssue domain.IssueExtendedWithChildren
		expectedError bool
		errorContains string
	}{
		{
			name:    "Success",
			issueID: 1,
			setupMocks: func(
				mockIssuesRepo *mockcontract.MockIssuesRepository,
				mockUsersRepo *mockcontract.MockUsersRepository,
				mockProjectsService *mockcontract.MockProjectsUseCase,
				mockProjectsRepo *mockcontract.MockProjectsRepository,
				mockEventsRepo *mockcontract.MockEventRepository,
			) {
				// Setup issue repo
				mockIssuesRepo.EXPECT().GetByID(mock.Anything, domain.IssueID(1)).
					Return(domain.Issue{
						ID:          1,
						ProjectID:   1,
						Fingerprint: "test-fingerprint",
						Title:       "Test Issue",
					}, nil)

				// Setup users repo
				mockUsersRepo.EXPECT().GetByID(mock.Anything, domain.UserID(123)).
					Return(domain.User{
						ID:          123,
						Username:    "testuser",
						IsSuperuser: false,
					}, nil)

				// Setup projects service
				mockProjectsService.EXPECT().GetProjectsByUserID(mock.Anything, domain.UserID(123), false).
					Return([]domain.ProjectExtended{
						{
							Project: domain.Project{
								ID:   1,
								Name: "Test Project",
							},
						},
					}, nil)

				// Setup projects repo
				mockProjectsRepo.EXPECT().GetByID(mock.Anything, domain.ProjectID(1)).
					Return(domain.Project{
						ID:   1,
						Name: "Test Project",
					}, nil)

				// Setup events repo
				mockEventsRepo.EXPECT().FetchForIssue(mock.Anything, domain.ProjectID(1), "test-fingerprint", uint(10)).
					Return([]domain.Event{
						{
							ID:        "event-1",
							ProjectID: 1,
							Message:   "Test Event",
						},
					}, nil)
			},
			setupContext: func(ctx context.Context) context.Context {
				return wardencontext.WithUserID(ctx, 123)
			},
			expectedIssue: domain.IssueExtendedWithChildren{
				Issue: domain.Issue{
					ID:          1,
					ProjectID:   1,
					Fingerprint: "test-fingerprint",
					Title:       "Test Issue",
				},
				ProjectName: "Test Project",
				Events: []domain.Event{
					{
						ID:        "event-1",
						ProjectID: 1,
						Message:   "Test Event",
					},
				},
			},
			expectedError: false,
		},
		{
			name:    "Error getting issue",
			issueID: 1,
			setupMocks: func(
				mockIssuesRepo *mockcontract.MockIssuesRepository,
				mockUsersRepo *mockcontract.MockUsersRepository,
				mockProjectsService *mockcontract.MockProjectsUseCase,
				mockProjectsRepo *mockcontract.MockProjectsRepository,
				mockEventsRepo *mockcontract.MockEventRepository,
			) {
				mockIssuesRepo.EXPECT().GetByID(mock.Anything, domain.IssueID(1)).
					Return(domain.Issue{}, errors.New("issue not found"))
			},
			setupContext: func(ctx context.Context) context.Context {
				return wardencontext.WithUserID(ctx, 123)
			},
			expectedIssue: domain.IssueExtendedWithChildren{},
			expectedError: true,
			errorContains: "get issue by ID: issue not found",
		},
		{
			name:    "Error getting user",
			issueID: 1,
			setupMocks: func(
				mockIssuesRepo *mockcontract.MockIssuesRepository,
				mockUsersRepo *mockcontract.MockUsersRepository,
				mockProjectsService *mockcontract.MockProjectsUseCase,
				mockProjectsRepo *mockcontract.MockProjectsRepository,
				mockEventsRepo *mockcontract.MockEventRepository,
			) {
				mockIssuesRepo.EXPECT().GetByID(mock.Anything, domain.IssueID(1)).
					Return(domain.Issue{
						ID:        1,
						ProjectID: 1,
						Title:     "Test Issue",
					}, nil)

				mockUsersRepo.EXPECT().GetByID(mock.Anything, domain.UserID(123)).
					Return(domain.User{}, errors.New("user not found"))
			},
			setupContext: func(ctx context.Context) context.Context {
				return wardencontext.WithUserID(ctx, 123)
			},
			expectedIssue: domain.IssueExtendedWithChildren{},
			expectedError: true,
			errorContains: "get current user: user not found",
		},
		{
			name:    "Error getting user projects",
			issueID: 1,
			setupMocks: func(
				mockIssuesRepo *mockcontract.MockIssuesRepository,
				mockUsersRepo *mockcontract.MockUsersRepository,
				mockProjectsService *mockcontract.MockProjectsUseCase,
				mockProjectsRepo *mockcontract.MockProjectsRepository,
				mockEventsRepo *mockcontract.MockEventRepository,
			) {
				mockIssuesRepo.EXPECT().GetByID(mock.Anything, domain.IssueID(1)).
					Return(domain.Issue{
						ID:        1,
						ProjectID: 1,
						Title:     "Test Issue",
					}, nil)

				mockUsersRepo.EXPECT().GetByID(mock.Anything, domain.UserID(123)).
					Return(domain.User{
						ID:          123,
						Username:    "testuser",
						IsSuperuser: false,
					}, nil)

				mockProjectsService.EXPECT().GetProjectsByUserID(mock.Anything, domain.UserID(123), false).
					Return(nil, errors.New("error getting projects"))
			},
			setupContext: func(ctx context.Context) context.Context {
				return wardencontext.WithUserID(ctx, 123)
			},
			expectedIssue: domain.IssueExtendedWithChildren{},
			expectedError: true,
			errorContains: "get user projects: error getting projects",
		},
		{
			name:    "User does not have access to project",
			issueID: 1,
			setupMocks: func(
				mockIssuesRepo *mockcontract.MockIssuesRepository,
				mockUsersRepo *mockcontract.MockUsersRepository,
				mockProjectsService *mockcontract.MockProjectsUseCase,
				mockProjectsRepo *mockcontract.MockProjectsRepository,
				mockEventsRepo *mockcontract.MockEventRepository,
			) {
				mockIssuesRepo.EXPECT().GetByID(mock.Anything, domain.IssueID(1)).
					Return(domain.Issue{
						ID:        1,
						ProjectID: 1,
						Title:     "Test Issue",
					}, nil)

				mockUsersRepo.EXPECT().GetByID(mock.Anything, domain.UserID(123)).
					Return(domain.User{
						ID:          123,
						Username:    "testuser",
						IsSuperuser: false,
					}, nil)

				mockProjectsService.EXPECT().GetProjectsByUserID(mock.Anything, domain.UserID(123), false).
					Return([]domain.ProjectExtended{
						{
							Project: domain.Project{
								ID:   2, // Different project ID
								Name: "Another Project",
							},
						},
					}, nil)
			},
			setupContext: func(ctx context.Context) context.Context {
				return wardencontext.WithUserID(ctx, 123)
			},
			expectedIssue: domain.IssueExtendedWithChildren{},
			expectedError: true,
			errorContains: "permission denied",
		},
		{
			name:    "Error getting project",
			issueID: 1,
			setupMocks: func(
				mockIssuesRepo *mockcontract.MockIssuesRepository,
				mockUsersRepo *mockcontract.MockUsersRepository,
				mockProjectsService *mockcontract.MockProjectsUseCase,
				mockProjectsRepo *mockcontract.MockProjectsRepository,
				mockEventsRepo *mockcontract.MockEventRepository,
			) {
				mockIssuesRepo.EXPECT().GetByID(mock.Anything, domain.IssueID(1)).
					Return(domain.Issue{
						ID:        1,
						ProjectID: 1,
						Title:     "Test Issue",
					}, nil)

				mockUsersRepo.EXPECT().GetByID(mock.Anything, domain.UserID(123)).
					Return(domain.User{
						ID:          123,
						Username:    "testuser",
						IsSuperuser: false,
					}, nil)

				mockProjectsService.EXPECT().GetProjectsByUserID(mock.Anything, domain.UserID(123), false).
					Return([]domain.ProjectExtended{
						{
							Project: domain.Project{
								ID:   1,
								Name: "Test Project",
							},
						},
					}, nil)

				mockProjectsRepo.EXPECT().GetByID(mock.Anything, domain.ProjectID(1)).
					Return(domain.Project{}, errors.New("project not found"))
			},
			setupContext: func(ctx context.Context) context.Context {
				return wardencontext.WithUserID(ctx, 123)
			},
			expectedIssue: domain.IssueExtendedWithChildren{},
			expectedError: true,
			errorContains: "get project by ID: project not found",
		},
		{
			name:    "Error fetching events",
			issueID: 1,
			setupMocks: func(
				mockIssuesRepo *mockcontract.MockIssuesRepository,
				mockUsersRepo *mockcontract.MockUsersRepository,
				mockProjectsService *mockcontract.MockProjectsUseCase,
				mockProjectsRepo *mockcontract.MockProjectsRepository,
				mockEventsRepo *mockcontract.MockEventRepository,
			) {
				mockIssuesRepo.EXPECT().GetByID(mock.Anything, domain.IssueID(1)).
					Return(domain.Issue{
						ID:          1,
						ProjectID:   1,
						Fingerprint: "test-fingerprint",
						Title:       "Test Issue",
					}, nil)

				mockUsersRepo.EXPECT().GetByID(mock.Anything, domain.UserID(123)).
					Return(domain.User{
						ID:          123,
						Username:    "testuser",
						IsSuperuser: false,
					}, nil)

				mockProjectsService.EXPECT().GetProjectsByUserID(mock.Anything, domain.UserID(123), false).
					Return([]domain.ProjectExtended{
						{
							Project: domain.Project{
								ID:   1,
								Name: "Test Project",
							},
						},
					}, nil)

				mockProjectsRepo.EXPECT().GetByID(mock.Anything, domain.ProjectID(1)).
					Return(domain.Project{
						ID:   1,
						Name: "Test Project",
					}, nil)

				mockEventsRepo.EXPECT().FetchForIssue(mock.Anything, domain.ProjectID(1), "test-fingerprint", uint(10)).
					Return(nil, errors.New("error fetching events"))
			},
			setupContext: func(ctx context.Context) context.Context {
				return wardencontext.WithUserID(ctx, 123)
			},
			expectedIssue: domain.IssueExtendedWithChildren{},
			expectedError: true,
			errorContains: "fetch events for issue: error fetching events",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// Create mocks
			mockTxManager := mockdb.NewMockTxManager(t)
			mockIssuesRepo := mockcontract.NewMockIssuesRepository(t)
			mockProjectsRepo := mockcontract.NewMockProjectsRepository(t)
			mockEventsRepo := mockcontract.NewMockEventRepository(t)
			mockProjectsService := mockcontract.NewMockProjectsUseCase(t)
			mockResolutionsRepo := mockcontract.NewMockResolutionsRepository(t)
			mockUsersRepo := mockcontract.NewMockUsersRepository(t)
			mockTeamsRepo := mockcontract.NewMockTeamsRepository(t)
			mockUserNotificationsUseCase := mockcontract.NewMockUserNotificationsUseCase(t)

			// Setup mocks
			tt.setupMocks(
				mockIssuesRepo,
				mockUsersRepo,
				mockProjectsService,
				mockProjectsRepo,
				mockEventsRepo,
			)

			// Create service
			service := New(
				mockTxManager,
				mockIssuesRepo,
				mockProjectsRepo,
				mockEventsRepo,
				mockProjectsService,
				mockResolutionsRepo,
				mockUsersRepo,
				mockTeamsRepo,
				mockUserNotificationsUseCase,
			)

			// Setup context
			ctx := tt.setupContext(context.Background())

			// Call the method
			issue, err := service.GetByIDWithChildren(ctx, tt.issueID)

			// Check results
			if tt.expectedError {
				require.Error(t, err)
				if tt.errorContains != "" {
					require.Contains(t, err.Error(), tt.errorContains)
				}
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.expectedIssue.ID, issue.ID)
				require.Equal(t, tt.expectedIssue.ProjectID, issue.ProjectID)
				require.Equal(t, tt.expectedIssue.Fingerprint, issue.Fingerprint)
				require.Equal(t, tt.expectedIssue.Title, issue.Title)
				require.Equal(t, tt.expectedIssue.ProjectName, issue.ProjectName)
				require.Equal(t, len(tt.expectedIssue.Events), len(issue.Events))
				if len(tt.expectedIssue.Events) > 0 {
					require.Equal(t, tt.expectedIssue.Events[0].ID, issue.Events[0].ID)
				}
			}
		})
	}
}

func TestTimeseries(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		filter         *domain.IssueTimeseriesFilter
		setupMocks     func(mockIssuesRepo *mockcontract.MockIssuesRepository)
		expectedSeries []domain.Timeseries
		expectedError  bool
		errorContains  string
	}{
		{
			name: "Success",
			filter: &domain.IssueTimeseriesFilter{
				Period: domain.Period{
					Interval:    24 * time.Hour,
					Granularity: time.Hour,
				},
				ProjectID: func() *domain.ProjectID {
					id := domain.ProjectID(1)

					return &id
				}(),
				Levels:   []domain.IssueLevel{domain.IssueLevelError},
				Statuses: []domain.IssueStatus{domain.IssueStatusUnresolved},
				GroupBy:  domain.IssueTimeseriesGroupNone,
			},
			setupMocks: func(mockIssuesRepo *mockcontract.MockIssuesRepository) {
				mockIssuesRepo.EXPECT().Timeseries(mock.Anything, mock.AnythingOfType("*domain.IssueTimeseriesFilter")).
					Return([]domain.Timeseries{
						{
							Name: "test",
							Period: domain.Period{
								Interval:    24 * time.Hour,
								Granularity: time.Hour,
							},
							Occurrences: []uint{10},
						},
					}, nil)
			},
			expectedSeries: []domain.Timeseries{
				{
					Name: "test",
					Period: domain.Period{
						Interval:    24 * time.Hour,
						Granularity: time.Hour,
					},
					Occurrences: []uint{10},
				},
			},
			expectedError: false,
		},
		{
			name:   "Error",
			filter: &domain.IssueTimeseriesFilter{},
			setupMocks: func(mockIssuesRepo *mockcontract.MockIssuesRepository) {
				mockIssuesRepo.EXPECT().Timeseries(mock.Anything, mock.AnythingOfType("*domain.IssueTimeseriesFilter")).
					Return(nil, errors.New("database error"))
			},
			expectedSeries: nil,
			expectedError:  true,
			errorContains:  "database error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// Create mocks
			mockTxManager := mockdb.NewMockTxManager(t)
			mockIssuesRepo := mockcontract.NewMockIssuesRepository(t)
			mockProjectsRepo := mockcontract.NewMockProjectsRepository(t)
			mockEventsRepo := mockcontract.NewMockEventRepository(t)
			mockProjectsService := mockcontract.NewMockProjectsUseCase(t)
			mockResolutionsRepo := mockcontract.NewMockResolutionsRepository(t)
			mockUsersRepo := mockcontract.NewMockUsersRepository(t)
			mockTeamsRepo := mockcontract.NewMockTeamsRepository(t)
			mockUserNotificationsUseCase := mockcontract.NewMockUserNotificationsUseCase(t)

			// Setup mocks
			tt.setupMocks(mockIssuesRepo)

			// Create service
			service := New(
				mockTxManager,
				mockIssuesRepo,
				mockProjectsRepo,
				mockEventsRepo,
				mockProjectsService,
				mockResolutionsRepo,
				mockUsersRepo,
				mockTeamsRepo,
				mockUserNotificationsUseCase,
			)

			// Call the method
			series, err := service.Timeseries(context.Background(), tt.filter)

			// Check results
			if tt.expectedError {
				require.Error(t, err)
				if tt.errorContains != "" {
					require.Contains(t, err.Error(), tt.errorContains)
				}
				require.Nil(t, series)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.expectedSeries, series)
			}
		})
	}
}

func TestChangeStatus(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		issueID    domain.IssueID
		status     domain.IssueStatus
		setupMocks func(
			mockTxManager *mockdb.MockTxManager,
			mockIssuesRepo *mockcontract.MockIssuesRepository,
			mockUsersRepo *mockcontract.MockUsersRepository,
			mockProjectsService *mockcontract.MockProjectsUseCase,
			mockResolutionsRepo *mockcontract.MockResolutionsRepository,
			mockProjectsRepo *mockcontract.MockProjectsRepository,
		)
		setupContext  func(ctx context.Context) context.Context
		expectedError bool
		errorContains string
	}{
		{
			name:    "Success",
			issueID: 1,
			status:  domain.IssueStatusResolved,
			setupMocks: func(
				mockTxManager *mockdb.MockTxManager,
				mockIssuesRepo *mockcontract.MockIssuesRepository,
				mockUsersRepo *mockcontract.MockUsersRepository,
				mockProjectsService *mockcontract.MockProjectsUseCase,
				mockResolutionsRepo *mockcontract.MockResolutionsRepository,
				mockProjectsRepo *mockcontract.MockProjectsRepository,
			) {
				mockProjectsRepo.EXPECT().GetByID(
					mock.Anything,
					domain.ProjectID(1),
				).Return(domain.Project{}, nil)

				// Setup users repo
				mockUsersRepo.EXPECT().GetByID(mock.Anything, domain.UserID(123)).
					Return(domain.User{
						ID:          123,
						Username:    "testuser",
						IsSuperuser: false,
					}, nil)

				// Setup issues repo
				mockIssuesRepo.EXPECT().GetByID(mock.Anything, domain.IssueID(1)).
					Return(domain.Issue{
						ID:        1,
						ProjectID: 1,
						Title:     "Test Issue",
					}, nil)

				// Setup projects service
				mockProjectsService.EXPECT().GetProjectsByUserID(mock.Anything, domain.UserID(123), false).
					Return([]domain.ProjectExtended{
						{
							Project: domain.Project{
								ID:   1,
								Name: "Test Project",
							},
						},
					}, nil)

				// Setup transaction
				mockTxManager.On("RepeatableRead", mock.Anything, mock.AnythingOfType("func(context.Context) error")).
					Run(func(args mock.Arguments) {
						// Execute the transaction function
						txFunc := args.Get(1).(func(context.Context) error)
						err := txFunc(context.Background())
						require.NoError(t, err)
					}).Return(nil)

				// Setup resolutions repo
				userID := domain.UserID(123)
				mockResolutionsRepo.EXPECT().Create(mock.Anything, domain.ResolutionDTO{
					ProjectID:  1,
					IssueID:    1,
					Status:     domain.IssueStatusResolved,
					ResolvedBy: &userID,
					Comment:    "",
				}).Return(domain.Resolution{}, nil)

				// Setup issues repo for update
				mockIssuesRepo.EXPECT().UpdateStatus(mock.Anything, domain.IssueID(1), domain.IssueStatusResolved).
					Return(nil)
			},
			setupContext: func(ctx context.Context) context.Context {
				return wardencontext.WithUserID(ctx, 123)
			},
			expectedError: false,
		},
		{
			name:    "Error getting user",
			issueID: 1,
			status:  domain.IssueStatusResolved,
			setupMocks: func(
				mockTxManager *mockdb.MockTxManager,
				mockIssuesRepo *mockcontract.MockIssuesRepository,
				mockUsersRepo *mockcontract.MockUsersRepository,
				mockProjectsService *mockcontract.MockProjectsUseCase,
				mockResolutionsRepo *mockcontract.MockResolutionsRepository,
				mockProjectsRepo *mockcontract.MockProjectsRepository,
			) {
				mockUsersRepo.EXPECT().GetByID(mock.Anything, domain.UserID(123)).
					Return(domain.User{}, errors.New("user not found"))
			},
			setupContext: func(ctx context.Context) context.Context {
				return wardencontext.WithUserID(ctx, 123)
			},
			expectedError: true,
			errorContains: "get current user by ID: user not found",
		},
		{
			name:    "Error getting issue",
			issueID: 1,
			status:  domain.IssueStatusResolved,
			setupMocks: func(
				mockTxManager *mockdb.MockTxManager,
				mockIssuesRepo *mockcontract.MockIssuesRepository,
				mockUsersRepo *mockcontract.MockUsersRepository,
				mockProjectsService *mockcontract.MockProjectsUseCase,
				mockResolutionsRepo *mockcontract.MockResolutionsRepository,
				mockProjectsRepo *mockcontract.MockProjectsRepository,
			) {
				mockUsersRepo.EXPECT().GetByID(mock.Anything, domain.UserID(123)).
					Return(domain.User{
						ID:          123,
						Username:    "testuser",
						IsSuperuser: false,
					}, nil)

				mockIssuesRepo.EXPECT().GetByID(mock.Anything, domain.IssueID(1)).
					Return(domain.Issue{}, errors.New("issue not found"))
			},
			setupContext: func(ctx context.Context) context.Context {
				return wardencontext.WithUserID(ctx, 123)
			},
			expectedError: true,
			errorContains: "get issue by ID: issue not found",
		},
		{
			name:    "Error getting user projects",
			issueID: 1,
			status:  domain.IssueStatusResolved,
			setupMocks: func(
				mockTxManager *mockdb.MockTxManager,
				mockIssuesRepo *mockcontract.MockIssuesRepository,
				mockUsersRepo *mockcontract.MockUsersRepository,
				mockProjectsService *mockcontract.MockProjectsUseCase,
				mockResolutionsRepo *mockcontract.MockResolutionsRepository,
				mockProjectsRepo *mockcontract.MockProjectsRepository,
			) {
				mockProjectsRepo.EXPECT().GetByID(
					mock.Anything,
					domain.ProjectID(1),
				).Return(domain.Project{}, nil)

				mockUsersRepo.EXPECT().GetByID(mock.Anything, domain.UserID(123)).
					Return(domain.User{
						ID:          123,
						Username:    "testuser",
						IsSuperuser: false,
					}, nil)

				mockIssuesRepo.EXPECT().GetByID(mock.Anything, domain.IssueID(1)).
					Return(domain.Issue{
						ID:        1,
						ProjectID: 1,
						Title:     "Test Issue",
					}, nil)

				mockProjectsService.EXPECT().GetProjectsByUserID(mock.Anything, domain.UserID(123), false).
					Return(nil, errors.New("error getting projects"))
			},
			setupContext: func(ctx context.Context) context.Context {
				return wardencontext.WithUserID(ctx, 123)
			},
			expectedError: true,
			errorContains: "get user projects: error getting projects",
		},
		{
			name:    "User does not have access to issue",
			issueID: 1,
			status:  domain.IssueStatusResolved,
			setupMocks: func(
				mockTxManager *mockdb.MockTxManager,
				mockIssuesRepo *mockcontract.MockIssuesRepository,
				mockUsersRepo *mockcontract.MockUsersRepository,
				mockProjectsService *mockcontract.MockProjectsUseCase,
				mockResolutionsRepo *mockcontract.MockResolutionsRepository,
				mockProjectsRepo *mockcontract.MockProjectsRepository,
			) {
				mockProjectsRepo.EXPECT().GetByID(
					mock.Anything,
					domain.ProjectID(1),
				).Return(domain.Project{}, nil)

				mockUsersRepo.EXPECT().GetByID(mock.Anything, domain.UserID(123)).
					Return(domain.User{
						ID:          123,
						Username:    "testuser",
						IsSuperuser: false,
					}, nil)

				mockIssuesRepo.EXPECT().GetByID(mock.Anything, domain.IssueID(1)).
					Return(domain.Issue{
						ID:        1,
						ProjectID: 1,
						Title:     "Test Issue",
					}, nil)

				mockProjectsService.EXPECT().GetProjectsByUserID(mock.Anything, domain.UserID(123), false).
					Return([]domain.ProjectExtended{
						{
							Project: domain.Project{
								ID:   2, // Different project ID
								Name: "Another Project",
							},
						},
					}, nil)
			},
			setupContext: func(ctx context.Context) context.Context {
				return wardencontext.WithUserID(ctx, 123)
			},
			expectedError: true,
			errorContains: "user does not have access to this issue",
		},
		{
			name:    "Transaction error",
			issueID: 1,
			status:  domain.IssueStatusResolved,
			setupMocks: func(
				mockTxManager *mockdb.MockTxManager,
				mockIssuesRepo *mockcontract.MockIssuesRepository,
				mockUsersRepo *mockcontract.MockUsersRepository,
				mockProjectsService *mockcontract.MockProjectsUseCase,
				mockResolutionsRepo *mockcontract.MockResolutionsRepository,
				mockProjectsRepo *mockcontract.MockProjectsRepository,
			) {
				mockProjectsRepo.EXPECT().GetByID(
					mock.Anything,
					domain.ProjectID(1),
				).Return(domain.Project{}, nil)

				mockUsersRepo.EXPECT().GetByID(mock.Anything, domain.UserID(123)).
					Return(domain.User{
						ID:          123,
						Username:    "testuser",
						IsSuperuser: false,
					}, nil)

				mockIssuesRepo.EXPECT().GetByID(mock.Anything, domain.IssueID(1)).
					Return(domain.Issue{
						ID:        1,
						ProjectID: 1,
						Title:     "Test Issue",
					}, nil)

				mockProjectsService.EXPECT().GetProjectsByUserID(mock.Anything, domain.UserID(123), false).
					Return([]domain.ProjectExtended{
						{
							Project: domain.Project{
								ID:   1,
								Name: "Test Project",
							},
						},
					}, nil)

				mockTxManager.EXPECT().RepeatableRead(mock.Anything, mock.AnythingOfType("func(context.Context) error")).
					Return(errors.New("transaction error"))
			},
			setupContext: func(ctx context.Context) context.Context {
				return wardencontext.WithUserID(ctx, 123)
			},
			expectedError: true,
			errorContains: "change issue status: transaction error",
		},
		{
			name:    "Error creating resolution",
			issueID: 1,
			status:  domain.IssueStatusResolved,
			setupMocks: func(
				mockTxManager *mockdb.MockTxManager,
				mockIssuesRepo *mockcontract.MockIssuesRepository,
				mockUsersRepo *mockcontract.MockUsersRepository,
				mockProjectsService *mockcontract.MockProjectsUseCase,
				mockResolutionsRepo *mockcontract.MockResolutionsRepository,
				mockProjectsRepo *mockcontract.MockProjectsRepository,
			) {
				mockProjectsRepo.EXPECT().GetByID(
					mock.Anything,
					domain.ProjectID(1),
				).Return(domain.Project{}, nil)

				mockUsersRepo.EXPECT().GetByID(mock.Anything, domain.UserID(123)).
					Return(domain.User{
						ID:          123,
						Username:    "testuser",
						IsSuperuser: false,
					}, nil)

				mockIssuesRepo.EXPECT().GetByID(mock.Anything, domain.IssueID(1)).
					Return(domain.Issue{
						ID:        1,
						ProjectID: 1,
						Title:     "Test Issue",
					}, nil)

				mockProjectsService.EXPECT().GetProjectsByUserID(mock.Anything, domain.UserID(123), false).
					Return([]domain.ProjectExtended{
						{
							Project: domain.Project{
								ID:   1,
								Name: "Test Project",
							},
						},
					}, nil)

				mockTxManager.On("RepeatableRead", mock.Anything, mock.AnythingOfType("func(context.Context) error")).
					Run(func(args mock.Arguments) {
						// Execute the transaction function
						txFunc := args.Get(1).(func(context.Context) error)
						_ = txFunc(context.Background())
					}).Return(errors.New("change issue status: create resolution: error creating resolution"))

				userID := domain.UserID(123)
				mockResolutionsRepo.EXPECT().Create(mock.Anything, domain.ResolutionDTO{
					ProjectID:  1,
					IssueID:    1,
					Status:     domain.IssueStatusResolved,
					ResolvedBy: &userID,
					Comment:    "",
				}).Return(domain.Resolution{}, errors.New("error creating resolution"))
			},
			setupContext: func(ctx context.Context) context.Context {
				return wardencontext.WithUserID(ctx, 123)
			},
			expectedError: true,
			errorContains: "change issue status: create resolution: error creating resolution",
		},
		{
			name:    "Error updating issue status",
			issueID: 1,
			status:  domain.IssueStatusResolved,
			setupMocks: func(
				mockTxManager *mockdb.MockTxManager,
				mockIssuesRepo *mockcontract.MockIssuesRepository,
				mockUsersRepo *mockcontract.MockUsersRepository,
				mockProjectsService *mockcontract.MockProjectsUseCase,
				mockResolutionsRepo *mockcontract.MockResolutionsRepository,
				mockProjectsRepo *mockcontract.MockProjectsRepository,
			) {
				mockProjectsRepo.EXPECT().GetByID(
					mock.Anything,
					domain.ProjectID(1),
				).Return(domain.Project{}, nil)

				mockUsersRepo.EXPECT().GetByID(mock.Anything, domain.UserID(123)).
					Return(domain.User{
						ID:          123,
						Username:    "testuser",
						IsSuperuser: false,
					}, nil)

				mockIssuesRepo.EXPECT().GetByID(mock.Anything, domain.IssueID(1)).
					Return(domain.Issue{
						ID:        1,
						ProjectID: 1,
						Title:     "Test Issue",
					}, nil)

				mockProjectsService.EXPECT().GetProjectsByUserID(mock.Anything, domain.UserID(123), false).
					Return([]domain.ProjectExtended{
						{
							Project: domain.Project{
								ID:   1,
								Name: "Test Project",
							},
						},
					}, nil)

				mockTxManager.On("RepeatableRead", mock.Anything, mock.AnythingOfType("func(context.Context) error")).
					Run(func(args mock.Arguments) {
						// Execute the transaction function
						txFunc := args.Get(1).(func(context.Context) error)
						_ = txFunc(context.Background())
					}).Return(errors.New("change issue status: error updating status"))

				userID := domain.UserID(123)
				mockResolutionsRepo.EXPECT().Create(mock.Anything, domain.ResolutionDTO{
					ProjectID:  1,
					IssueID:    1,
					Status:     domain.IssueStatusResolved,
					ResolvedBy: &userID,
					Comment:    "",
				}).Return(domain.Resolution{}, nil)

				mockIssuesRepo.EXPECT().UpdateStatus(mock.Anything, domain.IssueID(1), domain.IssueStatusResolved).
					Return(errors.New("error updating status"))
			},
			setupContext: func(ctx context.Context) context.Context {
				return wardencontext.WithUserID(ctx, 123)
			},
			expectedError: true,
			errorContains: "change issue status: error updating status",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// Create mocks
			mockTxManager := mockdb.NewMockTxManager(t)
			mockIssuesRepo := mockcontract.NewMockIssuesRepository(t)
			mockProjectsRepo := mockcontract.NewMockProjectsRepository(t)
			mockEventsRepo := mockcontract.NewMockEventRepository(t)
			mockProjectsService := mockcontract.NewMockProjectsUseCase(t)
			mockResolutionsRepo := mockcontract.NewMockResolutionsRepository(t)
			mockUsersRepo := mockcontract.NewMockUsersRepository(t)
			mockTeamsRepo := mockcontract.NewMockTeamsRepository(t)
			mockUserNotificationsUseCase := mockcontract.NewMockUserNotificationsUseCase(t)

			// Setup mocks
			tt.setupMocks(
				mockTxManager,
				mockIssuesRepo,
				mockUsersRepo,
				mockProjectsService,
				mockResolutionsRepo,
				mockProjectsRepo,
			)

			// Create service
			service := New(
				mockTxManager,
				mockIssuesRepo,
				mockProjectsRepo,
				mockEventsRepo,
				mockProjectsService,
				mockResolutionsRepo,
				mockUsersRepo,
				mockTeamsRepo,
				mockUserNotificationsUseCase,
			)

			// Setup context
			ctx := tt.setupContext(context.Background())

			// Call the method
			err := service.ChangeStatus(ctx, tt.issueID, tt.status)

			// Check results
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
