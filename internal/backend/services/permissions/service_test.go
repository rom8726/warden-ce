//nolint:lll // ooook
package permissions

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	wardencontext "github.com/rom8726/warden/internal/context"
	"github.com/rom8726/warden/internal/domain"
	mockcontract "github.com/rom8726/warden/test_mocks/internal_/backend/contract"
)

func TestCanAccessProject(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		setupMocks    func(teamsUseCase *mockcontract.MockTeamsUseCase, projectRepo *mockcontract.MockProjectsRepository)
		setupContext  func(ctx context.Context) context.Context
		projectID     domain.ProjectID
		expectedError error
	}{
		{
			name: "Super user can access any project",
			setupMocks: func(teamsUseCase *mockcontract.MockTeamsUseCase, projectRepo *mockcontract.MockProjectsRepository) {
				projectRepo.EXPECT().GetByID(mock.Anything, domain.ProjectID(1)).Return(domain.Project{}, nil)
			},
			setupContext: func(ctx context.Context) context.Context {
				return wardencontext.WithIsSuper(ctx, true)
			},
			projectID:     1,
			expectedError: nil,
		},
		{
			name: "No user ID in context",
			setupMocks: func(teamsUseCase *mockcontract.MockTeamsUseCase, projectRepo *mockcontract.MockProjectsRepository) {
				projectRepo.EXPECT().GetByID(mock.Anything, domain.ProjectID(1)).Return(domain.Project{}, nil)
			},
			setupContext: func(ctx context.Context) context.Context {
				return ctx // No user ID in context
			},
			projectID:     1,
			expectedError: domain.ErrUserNotFound,
		},
		{
			name: "Error getting project",
			setupMocks: func(teamsUseCase *mockcontract.MockTeamsUseCase, projectRepo *mockcontract.MockProjectsRepository) {
				projectRepo.EXPECT().GetByID(mock.Anything, domain.ProjectID(1)).Return(domain.Project{}, errors.New("project not found"))
			},
			setupContext: func(ctx context.Context) context.Context {
				return wardencontext.WithUserID(ctx, 1)
			},
			projectID:     1,
			expectedError: errors.New("project not found"),
		},
		{
			name: "Project with no team is accessible to all users",
			setupMocks: func(teamsUseCase *mockcontract.MockTeamsUseCase, projectRepo *mockcontract.MockProjectsRepository) {
				projectRepo.EXPECT().GetByID(mock.Anything, domain.ProjectID(1)).Return(domain.Project{
					ID:   1,
					Name: "Test Project",
				}, nil)
			},
			setupContext: func(ctx context.Context) context.Context {
				return wardencontext.WithUserID(ctx, 1)
			},
			projectID:     1,
			expectedError: nil,
		},
		{
			name: "Error getting user teams",
			setupMocks: func(teamsUseCase *mockcontract.MockTeamsUseCase, projectRepo *mockcontract.MockProjectsRepository) {
				teamID := domain.TeamID(1)
				projectRepo.EXPECT().GetByID(mock.Anything, domain.ProjectID(1)).Return(domain.Project{
					ID:     1,
					Name:   "Test Project",
					TeamID: &teamID,
				}, nil)
				teamsUseCase.EXPECT().GetTeamsByUserID(mock.Anything, domain.UserID(1)).Return(nil, errors.New("error getting teams"))
			},
			setupContext: func(ctx context.Context) context.Context {
				return wardencontext.WithUserID(ctx, 1)
			},
			projectID:     1,
			expectedError: errors.New("error getting teams"),
		},
		{
			name: "User is member of project's team",
			setupMocks: func(teamsUseCase *mockcontract.MockTeamsUseCase, projectRepo *mockcontract.MockProjectsRepository) {
				teamID := domain.TeamID(1)
				projectRepo.EXPECT().GetByID(mock.Anything, domain.ProjectID(1)).Return(domain.Project{
					ID:     1,
					Name:   "Test Project",
					TeamID: &teamID,
				}, nil)
				teamsUseCase.EXPECT().GetTeamsByUserID(mock.Anything, domain.UserID(1)).Return([]domain.Team{
					{
						ID:   1,
						Name: "Test Team",
					},
				}, nil)
			},
			setupContext: func(ctx context.Context) context.Context {
				return wardencontext.WithUserID(ctx, 1)
			},
			projectID:     1,
			expectedError: nil,
		},
		{
			name: "User is not member of project's team",
			setupMocks: func(teamsUseCase *mockcontract.MockTeamsUseCase, projectRepo *mockcontract.MockProjectsRepository) {
				teamID := domain.TeamID(1)
				projectRepo.EXPECT().GetByID(mock.Anything, domain.ProjectID(1)).Return(domain.Project{
					ID:     1,
					Name:   "Test Project",
					TeamID: &teamID,
				}, nil)
				teamsUseCase.EXPECT().GetTeamsByUserID(mock.Anything, domain.UserID(1)).Return([]domain.Team{
					{
						ID:   2,
						Name: "Another Team",
					},
				}, nil)
			},
			setupContext: func(ctx context.Context) context.Context {
				return wardencontext.WithUserID(ctx, 1)
			},
			projectID:     1,
			expectedError: domain.ErrPermissionDenied,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// Create mocks
			teamsUseCase := mockcontract.NewMockTeamsUseCase(t)
			projectRepo := mockcontract.NewMockProjectsRepository(t)
			issueRepo := mockcontract.NewMockIssuesRepository(t)

			// Setup mocks
			tt.setupMocks(teamsUseCase, projectRepo)

			// Create service
			service := New(teamsUseCase, projectRepo, issueRepo)

			// Setup context
			ctx := tt.setupContext(context.Background())

			// Call the method
			err := service.CanAccessProject(ctx, tt.projectID)

			// Check the result
			if tt.expectedError != nil {
				require.EqualError(t, err, tt.expectedError.Error())
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestCanAccessIssue(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		setupMocks    func(teamsUseCase *mockcontract.MockTeamsUseCase, projectRepo *mockcontract.MockProjectsRepository, issueRepo *mockcontract.MockIssuesRepository)
		setupContext  func(ctx context.Context) context.Context
		issueID       domain.IssueID
		expectedError error
	}{
		{
			name: "Super user can access any issue",
			setupMocks: func(teamsUseCase *mockcontract.MockTeamsUseCase, projectRepo *mockcontract.MockProjectsRepository, issueRepo *mockcontract.MockIssuesRepository) {
				// No mocks needed for super user
			},
			setupContext: func(ctx context.Context) context.Context {
				return wardencontext.WithIsSuper(ctx, true)
			},
			issueID:       1,
			expectedError: nil,
		},
		{
			name: "Error getting issue",
			setupMocks: func(teamsUseCase *mockcontract.MockTeamsUseCase, projectRepo *mockcontract.MockProjectsRepository, issueRepo *mockcontract.MockIssuesRepository) {
				issueRepo.EXPECT().GetByID(mock.Anything, domain.IssueID(1)).Return(domain.Issue{}, errors.New("issue not found"))
			},
			setupContext: func(ctx context.Context) context.Context {
				return wardencontext.WithUserID(ctx, 1)
			},
			issueID:       1,
			expectedError: errors.New("issue not found"),
		},
		{
			name: "User can access issue if they can access the project",
			setupMocks: func(teamsUseCase *mockcontract.MockTeamsUseCase, projectRepo *mockcontract.MockProjectsRepository, issueRepo *mockcontract.MockIssuesRepository) {
				issueRepo.EXPECT().GetByID(mock.Anything, domain.IssueID(1)).Return(domain.Issue{
					ID:        1,
					ProjectID: 1,
				}, nil)
				projectRepo.EXPECT().GetByID(mock.Anything, domain.ProjectID(1)).Return(domain.Project{
					ID:   1,
					Name: "Test Project",
				}, nil)
			},
			setupContext: func(ctx context.Context) context.Context {
				return wardencontext.WithUserID(ctx, 1)
			},
			issueID:       1,
			expectedError: nil,
		},
		{
			name: "User cannot access issue if they cannot access the project",
			setupMocks: func(teamsUseCase *mockcontract.MockTeamsUseCase, projectRepo *mockcontract.MockProjectsRepository, issueRepo *mockcontract.MockIssuesRepository) {
				issueRepo.EXPECT().GetByID(mock.Anything, domain.IssueID(1)).Return(domain.Issue{
					ID:        1,
					ProjectID: 1,
				}, nil)
				teamID := domain.TeamID(1)
				projectRepo.EXPECT().GetByID(mock.Anything, domain.ProjectID(1)).Return(domain.Project{
					ID:     1,
					Name:   "Test Project",
					TeamID: &teamID,
				}, nil)
				teamsUseCase.EXPECT().GetTeamsByUserID(mock.Anything, domain.UserID(1)).Return([]domain.Team{
					{
						ID:   2,
						Name: "Another Team",
					},
				}, nil)
			},
			setupContext: func(ctx context.Context) context.Context {
				return wardencontext.WithUserID(ctx, 1)
			},
			issueID:       1,
			expectedError: domain.ErrPermissionDenied,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// Create mocks
			teamsUseCase := mockcontract.NewMockTeamsUseCase(t)
			projectRepo := mockcontract.NewMockProjectsRepository(t)
			issueRepo := mockcontract.NewMockIssuesRepository(t)

			// Setup mocks
			tt.setupMocks(teamsUseCase, projectRepo, issueRepo)

			// Create service
			service := New(teamsUseCase, projectRepo, issueRepo)

			// Setup context
			ctx := tt.setupContext(context.Background())

			// Call the method
			err := service.CanAccessIssue(ctx, tt.issueID)

			// Check the result
			if tt.expectedError != nil {
				assert.EqualError(t, err, tt.expectedError.Error())
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestCanManageProject(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name              string
		setupMocks        func(teamsUseCase *mockcontract.MockTeamsUseCase, projectRepo *mockcontract.MockProjectsRepository)
		setupContext      func(ctx context.Context) context.Context
		projectID         domain.ProjectID
		isIssueManagement bool
		expectedError     error
	}{
		{
			name: "Super user can manage any project",
			setupMocks: func(teamsUseCase *mockcontract.MockTeamsUseCase, projectRepo *mockcontract.MockProjectsRepository) {
				projectRepo.EXPECT().GetByID(mock.Anything, domain.ProjectID(1)).Return(domain.Project{}, nil)
			},
			setupContext: func(ctx context.Context) context.Context {
				return wardencontext.WithIsSuper(ctx, true)
			},
			projectID:     1,
			expectedError: nil,
		},
		{
			name: "No user ID in context",
			setupMocks: func(teamsUseCase *mockcontract.MockTeamsUseCase, projectRepo *mockcontract.MockProjectsRepository) {
				projectRepo.EXPECT().GetByID(mock.Anything, domain.ProjectID(1)).Return(domain.Project{}, nil)
			},
			setupContext: func(ctx context.Context) context.Context {
				return ctx // No user ID in context
			},
			projectID:     1,
			expectedError: domain.ErrUserNotFound,
		},
		{
			name: "Error getting project",
			setupMocks: func(teamsUseCase *mockcontract.MockTeamsUseCase, projectRepo *mockcontract.MockProjectsRepository) {
				projectRepo.EXPECT().GetByID(mock.Anything, domain.ProjectID(1)).Return(domain.Project{}, errors.New("project not found"))
			},
			setupContext: func(ctx context.Context) context.Context {
				return wardencontext.WithUserID(ctx, 1)
			},
			projectID:     1,
			expectedError: errors.New("project not found"),
		},
		{
			name: "Project with no team cannot be managed by regular users",
			setupMocks: func(teamsUseCase *mockcontract.MockTeamsUseCase, projectRepo *mockcontract.MockProjectsRepository) {
				projectRepo.EXPECT().GetByID(mock.Anything, domain.ProjectID(1)).Return(domain.Project{
					ID:   1,
					Name: "Test Project",
				}, nil)
			},
			setupContext: func(ctx context.Context) context.Context {
				return wardencontext.WithUserID(ctx, 1)
			},
			projectID:     1,
			expectedError: domain.ErrPermissionDenied,
		},
		{
			name: "Project with no team can manage issues",
			setupMocks: func(teamsUseCase *mockcontract.MockTeamsUseCase, projectRepo *mockcontract.MockProjectsRepository) {
				projectRepo.EXPECT().GetByID(mock.Anything, domain.ProjectID(1)).Return(domain.Project{
					ID:   1,
					Name: "Test Project",
				}, nil)
			},
			setupContext: func(ctx context.Context) context.Context {
				return wardencontext.WithUserID(ctx, 1)
			},
			projectID:         1,
			isIssueManagement: true,
			expectedError:     nil,
		},
		{
			name: "Error getting team members",
			setupMocks: func(teamsUseCase *mockcontract.MockTeamsUseCase, projectRepo *mockcontract.MockProjectsRepository) {
				teamID := domain.TeamID(1)
				projectRepo.EXPECT().GetByID(mock.Anything, domain.ProjectID(1)).Return(domain.Project{
					ID:     1,
					Name:   "Test Project",
					TeamID: &teamID,
				}, nil)
				teamsUseCase.EXPECT().GetMembers(mock.Anything, domain.TeamID(1)).Return(nil, errors.New("error getting members"))
			},
			setupContext: func(ctx context.Context) context.Context {
				return wardencontext.WithUserID(ctx, 1)
			},
			projectID:     1,
			expectedError: errors.New("error getting members"),
		},
		{
			name: "User is owner of the team",
			setupMocks: func(teamsUseCase *mockcontract.MockTeamsUseCase, projectRepo *mockcontract.MockProjectsRepository) {
				teamID := domain.TeamID(1)
				projectRepo.EXPECT().GetByID(mock.Anything, domain.ProjectID(1)).Return(domain.Project{
					ID:     1,
					Name:   "Test Project",
					TeamID: &teamID,
				}, nil)
				teamsUseCase.EXPECT().GetMembers(mock.Anything, domain.TeamID(1)).Return([]domain.TeamMember{
					{
						TeamID: 1,
						UserID: 1,
						Role:   domain.RoleOwner,
					},
				}, nil)
			},
			setupContext: func(ctx context.Context) context.Context {
				return wardencontext.WithUserID(ctx, 1)
			},
			projectID:     1,
			expectedError: nil,
		},
		{
			name: "User is admin of the team",
			setupMocks: func(teamsUseCase *mockcontract.MockTeamsUseCase, projectRepo *mockcontract.MockProjectsRepository) {
				teamID := domain.TeamID(1)
				projectRepo.EXPECT().GetByID(mock.Anything, domain.ProjectID(1)).Return(domain.Project{
					ID:     1,
					Name:   "Test Project",
					TeamID: &teamID,
				}, nil)
				teamsUseCase.EXPECT().GetMembers(mock.Anything, domain.TeamID(1)).Return([]domain.TeamMember{
					{
						TeamID: 1,
						UserID: 1,
						Role:   domain.RoleAdmin,
					},
				}, nil)
			},
			setupContext: func(ctx context.Context) context.Context {
				return wardencontext.WithUserID(ctx, 1)
			},
			projectID:     1,
			expectedError: nil,
		},
		{
			name: "User is regular member of the team",
			setupMocks: func(teamsUseCase *mockcontract.MockTeamsUseCase, projectRepo *mockcontract.MockProjectsRepository) {
				teamID := domain.TeamID(1)
				projectRepo.EXPECT().GetByID(mock.Anything, domain.ProjectID(1)).Return(domain.Project{
					ID:     1,
					Name:   "Test Project",
					TeamID: &teamID,
				}, nil)
				teamsUseCase.EXPECT().GetMembers(mock.Anything, domain.TeamID(1)).Return([]domain.TeamMember{
					{
						TeamID: 1,
						UserID: 1,
						Role:   domain.RoleMember,
					},
				}, nil)
			},
			setupContext: func(ctx context.Context) context.Context {
				return wardencontext.WithUserID(ctx, 1)
			},
			projectID:     1,
			expectedError: domain.ErrPermissionDenied,
		},
		{
			name: "User is not a member of the team",
			setupMocks: func(teamsUseCase *mockcontract.MockTeamsUseCase, projectRepo *mockcontract.MockProjectsRepository) {
				teamID := domain.TeamID(1)
				projectRepo.EXPECT().GetByID(mock.Anything, domain.ProjectID(1)).Return(domain.Project{
					ID:     1,
					Name:   "Test Project",
					TeamID: &teamID,
				}, nil)
				teamsUseCase.EXPECT().GetMembers(mock.Anything, domain.TeamID(1)).Return([]domain.TeamMember{
					{
						TeamID: 1,
						UserID: 2,
						Role:   domain.RoleOwner,
					},
				}, nil)
			},
			setupContext: func(ctx context.Context) context.Context {
				return wardencontext.WithUserID(ctx, 1)
			},
			projectID:     1,
			expectedError: domain.ErrPermissionDenied,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// Create mocks
			teamsUseCase := mockcontract.NewMockTeamsUseCase(t)
			projectRepo := mockcontract.NewMockProjectsRepository(t)
			issueRepo := mockcontract.NewMockIssuesRepository(t)

			// Setup mocks
			tt.setupMocks(teamsUseCase, projectRepo)

			// Create service
			service := New(teamsUseCase, projectRepo, issueRepo)

			// Setup context
			ctx := tt.setupContext(context.Background())

			// Call the method
			err := service.CanManageProject(ctx, tt.projectID, tt.isIssueManagement)

			// Check the result
			if tt.expectedError != nil {
				assert.EqualError(t, err, tt.expectedError.Error())
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestCanManageIssue(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		setupMocks    func(teamsUseCase *mockcontract.MockTeamsUseCase, projectRepo *mockcontract.MockProjectsRepository, issueRepo *mockcontract.MockIssuesRepository)
		setupContext  func(ctx context.Context) context.Context
		issueID       domain.IssueID
		expectedError error
	}{
		{
			name: "Super user can manage any issue",
			setupMocks: func(teamsUseCase *mockcontract.MockTeamsUseCase, projectRepo *mockcontract.MockProjectsRepository, issueRepo *mockcontract.MockIssuesRepository) {
				issueRepo.EXPECT().GetByID(mock.Anything, domain.IssueID(1)).Return(domain.Issue{ProjectID: 1}, nil)
			},
			setupContext: func(ctx context.Context) context.Context {
				return wardencontext.WithIsSuper(ctx, true)
			},
			issueID:       1,
			expectedError: nil,
		},
		{
			name: "Error getting issue",
			setupMocks: func(teamsUseCase *mockcontract.MockTeamsUseCase, projectRepo *mockcontract.MockProjectsRepository, issueRepo *mockcontract.MockIssuesRepository) {
				issueRepo.EXPECT().GetByID(mock.Anything, domain.IssueID(1)).Return(domain.Issue{}, errors.New("issue not found"))
			},
			setupContext: func(ctx context.Context) context.Context {
				return wardencontext.WithUserID(ctx, 1)
			},
			issueID:       1,
			expectedError: errors.New("issue not found"),
		},
		{
			name: "User can manage issue if they can manage the project",
			setupMocks: func(teamsUseCase *mockcontract.MockTeamsUseCase, projectRepo *mockcontract.MockProjectsRepository, issueRepo *mockcontract.MockIssuesRepository) {
				issueRepo.EXPECT().GetByID(mock.Anything, domain.IssueID(1)).Return(domain.Issue{
					ID:        1,
					ProjectID: 1,
				}, nil)
				teamID := domain.TeamID(1)
				projectRepo.EXPECT().GetByID(mock.Anything, domain.ProjectID(1)).Return(domain.Project{
					ID:     1,
					Name:   "Test Project",
					TeamID: &teamID,
				}, nil)
				teamsUseCase.EXPECT().GetMembers(mock.Anything, domain.TeamID(1)).Return([]domain.TeamMember{
					{
						TeamID: 1,
						UserID: 1,
						Role:   domain.RoleOwner,
					},
				}, nil)
			},
			setupContext: func(ctx context.Context) context.Context {
				return wardencontext.WithUserID(ctx, 1)
			},
			issueID:       1,
			expectedError: nil,
		},
		{
			name: "Regular member can manage issues",
			setupMocks: func(teamsUseCase *mockcontract.MockTeamsUseCase, projectRepo *mockcontract.MockProjectsRepository, issueRepo *mockcontract.MockIssuesRepository) {
				issueRepo.EXPECT().GetByID(mock.Anything, domain.IssueID(1)).Return(domain.Issue{
					ID:        1,
					ProjectID: 1,
				}, nil)
				teamID := domain.TeamID(1)
				projectRepo.EXPECT().GetByID(mock.Anything, domain.ProjectID(1)).Return(domain.Project{
					ID:     1,
					Name:   "Test Project",
					TeamID: &teamID,
				}, nil)
				teamsUseCase.EXPECT().GetMembers(mock.Anything, domain.TeamID(1)).Return([]domain.TeamMember{
					{
						TeamID: 1,
						UserID: 1,
						Role:   domain.RoleMember,
					},
				}, nil)
			},
			setupContext: func(ctx context.Context) context.Context {
				return wardencontext.WithUserID(ctx, 1)
			},
			issueID:       1,
			expectedError: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// Create mocks
			teamsUseCase := mockcontract.NewMockTeamsUseCase(t)
			projectRepo := mockcontract.NewMockProjectsRepository(t)
			issueRepo := mockcontract.NewMockIssuesRepository(t)

			// Setup mocks
			tt.setupMocks(teamsUseCase, projectRepo, issueRepo)

			// Create service
			service := New(teamsUseCase, projectRepo, issueRepo)

			// Setup context
			ctx := tt.setupContext(context.Background())

			// Call the method
			err := service.CanManageIssue(ctx, tt.issueID)

			// Check the result
			if tt.expectedError != nil {
				assert.EqualError(t, err, tt.expectedError.Error())
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestGetAccessibleProjects(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name             string
		setupMocks       func(teamsUseCase *mockcontract.MockTeamsUseCase)
		setupContext     func(ctx context.Context) context.Context
		inputProjects    []domain.ProjectExtended
		expectedProjects []domain.ProjectExtended
		expectedError    error
	}{
		{
			name: "Super user can access all projects",
			setupMocks: func(teamsUseCase *mockcontract.MockTeamsUseCase) {
				// No mocks needed for super user
			},
			setupContext: func(ctx context.Context) context.Context {
				return wardencontext.WithIsSuper(ctx, true)
			},
			inputProjects: []domain.ProjectExtended{
				{
					Project: domain.Project{
						ID:   1,
						Name: "Project 1",
					},
				},
				{
					Project: domain.Project{
						ID:   2,
						Name: "Project 2",
					},
				},
			},
			expectedProjects: []domain.ProjectExtended{
				{
					Project: domain.Project{
						ID:   1,
						Name: "Project 1",
					},
				},
				{
					Project: domain.Project{
						ID:   2,
						Name: "Project 2",
					},
				},
			},
			expectedError: nil,
		},
		{
			name: "No user ID in context",
			setupMocks: func(teamsUseCase *mockcontract.MockTeamsUseCase) {
				// No mocks needed for this case
			},
			setupContext: func(ctx context.Context) context.Context {
				return ctx // No user ID in context
			},
			inputProjects: []domain.ProjectExtended{
				{
					Project: domain.Project{
						ID:   1,
						Name: "Project 1",
					},
				},
			},
			expectedProjects: nil,
			expectedError:    domain.ErrUserNotFound,
		},
		{
			name: "Error getting user teams",
			setupMocks: func(teamsUseCase *mockcontract.MockTeamsUseCase) {
				teamsUseCase.EXPECT().GetTeamsByUserID(mock.Anything, domain.UserID(1)).Return(nil, errors.New("error getting teams"))
			},
			setupContext: func(ctx context.Context) context.Context {
				return wardencontext.WithUserID(ctx, 1)
			},
			inputProjects: []domain.ProjectExtended{
				{
					Project: domain.Project{
						ID:   1,
						Name: "Project 1",
					},
				},
			},
			expectedProjects: nil,
			expectedError:    errors.New("error getting teams"),
		},
		{
			name: "Filter projects by team membership",
			setupMocks: func(teamsUseCase *mockcontract.MockTeamsUseCase) {
				teamsUseCase.EXPECT().GetTeamsByUserID(mock.Anything, domain.UserID(1)).Return([]domain.Team{
					{
						ID:   1,
						Name: "Team 1",
					},
				}, nil)
			},
			setupContext: func(ctx context.Context) context.Context {
				return wardencontext.WithUserID(ctx, 1)
			},
			inputProjects: []domain.ProjectExtended{
				{
					Project: domain.Project{
						ID:   1,
						Name: "Project 1",
						TeamID: func() *domain.TeamID {
							id := domain.TeamID(1)

							return &id
						}(),
					},
				},
				{
					Project: domain.Project{
						ID:   2,
						Name: "Project 2",
						TeamID: func() *domain.TeamID {
							id := domain.TeamID(2)

							return &id
						}(),
					},
				},
				{
					Project: domain.Project{
						ID:   3,
						Name: "Project 3",
					},
				},
			},
			expectedProjects: []domain.ProjectExtended{
				{
					Project: domain.Project{
						ID:   1,
						Name: "Project 1",
						TeamID: func() *domain.TeamID {
							id := domain.TeamID(1)

							return &id
						}(),
					},
				},
				{
					Project: domain.Project{
						ID:   3,
						Name: "Project 3",
					},
				},
			},
			expectedError: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// Create mocks
			teamsUseCase := mockcontract.NewMockTeamsUseCase(t)
			projectRepo := mockcontract.NewMockProjectsRepository(t)
			issueRepo := mockcontract.NewMockIssuesRepository(t)

			// Setup mocks
			tt.setupMocks(teamsUseCase)

			// Create service
			service := New(teamsUseCase, projectRepo, issueRepo)

			// Setup context
			ctx := tt.setupContext(context.Background())

			// Call the method
			projects, err := service.GetAccessibleProjects(ctx, tt.inputProjects)

			// Check the result
			if tt.expectedError != nil {
				assert.EqualError(t, err, tt.expectedError.Error())
				assert.Nil(t, projects)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.expectedProjects, projects)
			}
		})
	}
}
