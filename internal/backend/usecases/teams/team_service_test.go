package teams

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	wardencontext "github.com/rom8726/warden/internal/context"
	"github.com/rom8726/warden/internal/domain"
	mockcontract "github.com/rom8726/warden/test_mocks/internal_/backend/contract"
	mockdb "github.com/rom8726/warden/test_mocks/pkg/db"
)

func TestNew(t *testing.T) {
	t.Parallel()

	mockTxManager := mockdb.NewMockTxManager(t)
	mockTeamsRepo := mockcontract.NewMockTeamsRepository(t)
	mockUsersRepo := mockcontract.NewMockUsersRepository(t)
	mockUserNotificationsUseCase := mockcontract.NewMockUserNotificationsUseCase(t)
	mockProjectsRepo := mockcontract.NewMockProjectsRepository(t)

	// Create service
	service := New(mockTxManager, mockTeamsRepo, mockUsersRepo, mockUserNotificationsUseCase, mockProjectsRepo)

	// Verify service was created correctly
	require.NotNil(t, service)
	require.Equal(t, mockTeamsRepo, service.teamsRepo)
	require.Equal(t, mockUsersRepo, service.usersRepo)
}

func TestCreate(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		setupMocks func(
			mockTeamsRepo *mockcontract.MockTeamsRepository,
			mockUsersRepo *mockcontract.MockUsersRepository,
		)
		teamName      string
		expectedTeam  domain.Team
		expectedError bool
		errorContains string
	}{
		{
			name: "Success",
			setupMocks: func(
				mockTeamsRepo *mockcontract.MockTeamsRepository,
				mockUsersRepo *mockcontract.MockUsersRepository,
			) {
				// Check if a user exists
				mockUsersRepo.EXPECT().ExistsByID(
					mock.Anything,
					domain.UserID(1),
				).Return(true, nil)

				// Create team
				mockTeamsRepo.EXPECT().Create(
					mock.Anything,
					domain.TeamDTO{Name: "Test Team"},
				).Return(domain.Team{
					ID:   1,
					Name: "Test Team",
				}, nil)

				// Add the current user as an owner
				mockTeamsRepo.EXPECT().AddMember(
					mock.Anything,
					domain.TeamID(1),
					domain.UserID(1),
					domain.RoleOwner,
				).Return(nil)

				mockTeamsRepo.EXPECT().GetByName(
					mock.Anything,
					"Test Team",
				).Return(domain.Team{}, domain.ErrEntityNotFound)

				// Get team by ID to refresh data
				mockTeamsRepo.EXPECT().GetByID(
					mock.Anything,
					domain.TeamID(1),
				).Return(domain.Team{
					ID:   1,
					Name: "Test Team",
					Members: []domain.TeamMember{
						{
							UserID: 1,
							Role:   domain.RoleOwner,
						},
					},
				}, nil)
			},
			teamName: "Test Team",
			expectedTeam: domain.Team{
				ID:   1,
				Name: "Test Team",
				Members: []domain.TeamMember{
					{
						UserID: 1,
						Role:   domain.RoleOwner,
					},
				},
			},
			expectedError: false,
		},
		{
			name: "User does not exist",
			setupMocks: func(
				mockTeamsRepo *mockcontract.MockTeamsRepository,
				mockUsersRepo *mockcontract.MockUsersRepository,
			) {
				// Check if user exists
				mockUsersRepo.EXPECT().ExistsByID(
					mock.Anything,
					domain.UserID(1),
				).Return(false, nil)
			},
			teamName:      "Test Team",
			expectedTeam:  domain.Team{},
			expectedError: true,
			errorContains: "not found",
		},
		{
			name: "Error checking if user exists",
			setupMocks: func(
				mockTeamsRepo *mockcontract.MockTeamsRepository,
				mockUsersRepo *mockcontract.MockUsersRepository,
			) {
				// Check if user exists
				mockUsersRepo.EXPECT().ExistsByID(
					mock.Anything,
					domain.UserID(1),
				).Return(false, errors.New("database error"))
			},
			teamName:      "Test Team",
			expectedTeam:  domain.Team{},
			expectedError: true,
			errorContains: "get current user by id",
		},
		{
			name: "Error creating team",
			setupMocks: func(
				mockTeamsRepo *mockcontract.MockTeamsRepository,
				mockUsersRepo *mockcontract.MockUsersRepository,
			) {
				// Check if a user exists
				mockUsersRepo.EXPECT().ExistsByID(
					mock.Anything,
					domain.UserID(1),
				).Return(true, nil)

				mockTeamsRepo.EXPECT().GetByName(
					mock.Anything,
					"Test Team",
				).Return(domain.Team{}, domain.ErrEntityNotFound)

				// Create team
				mockTeamsRepo.EXPECT().Create(
					mock.Anything,
					domain.TeamDTO{Name: "Test Team"},
				).Return(domain.Team{}, errors.New("database error"))
			},
			teamName:      "Test Team",
			expectedTeam:  domain.Team{},
			expectedError: true,
			errorContains: "database error",
		},
		{
			name: "Error adding member",
			setupMocks: func(
				mockTeamsRepo *mockcontract.MockTeamsRepository,
				mockUsersRepo *mockcontract.MockUsersRepository,
			) {
				// Check if user exists
				mockUsersRepo.EXPECT().ExistsByID(
					mock.Anything,
					domain.UserID(1),
				).Return(true, nil)

				mockTeamsRepo.EXPECT().GetByName(
					mock.Anything,
					"Test Team",
				).Return(domain.Team{}, domain.ErrEntityNotFound)

				// Create team
				mockTeamsRepo.EXPECT().Create(
					mock.Anything,
					domain.TeamDTO{Name: "Test Team"},
				).Return(domain.Team{
					ID:   1,
					Name: "Test Team",
				}, nil)

				// Add the current user as an owner
				mockTeamsRepo.EXPECT().AddMember(
					mock.Anything,
					domain.TeamID(1),
					domain.UserID(1),
					domain.RoleOwner,
				).Return(errors.New("database error"))
			},
			teamName:      "Test Team",
			expectedTeam:  domain.Team{},
			expectedError: true,
			errorContains: "add current user as owner",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// Create mocks
			mockTxManager := mockdb.NewMockTxManager(t)
			mockTeamsRepo := mockcontract.NewMockTeamsRepository(t)
			mockUsersRepo := mockcontract.NewMockUsersRepository(t)
			mockProjectsRepo := mockcontract.NewMockProjectsRepository(t)

			// Setup mocks
			tt.setupMocks(mockTeamsRepo, mockUsersRepo)

			// Create service
			mockUserNotificationsUseCase := mockcontract.NewMockUserNotificationsUseCase(t)
			service := New(mockTxManager, mockTeamsRepo, mockUsersRepo, mockUserNotificationsUseCase, mockProjectsRepo)

			// Create context with user ID
			ctx := wardencontext.WithUserID(context.Background(), 1)

			// Call method
			team, err := service.Create(ctx, domain.TeamDTO{Name: tt.teamName})

			// Check result
			if tt.expectedError {
				require.Error(t, err)
				if tt.errorContains != "" {
					require.Contains(t, err.Error(), tt.errorContains)
				}
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.expectedTeam, team)
			}
		})
	}
}

func TestGetByID(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		setupMocks func(
			mockTeamsRepo *mockcontract.MockTeamsRepository,
			mockUsersRepo *mockcontract.MockUsersRepository,
		)
		teamID        domain.TeamID
		expectedTeam  domain.Team
		expectedError bool
		errorContains string
	}{
		{
			name: "Success - Superuser",
			setupMocks: func(
				mockTeamsRepo *mockcontract.MockTeamsRepository,
				mockUsersRepo *mockcontract.MockUsersRepository,
			) {
				// Get current user
				mockUsersRepo.EXPECT().GetByID(
					mock.Anything,
					domain.UserID(1),
				).Return(domain.User{
					ID:          1,
					Username:    "admin",
					IsSuperuser: true,
				}, nil)

				// Get team
				mockTeamsRepo.EXPECT().GetByID(
					mock.Anything,
					domain.TeamID(1),
				).Return(domain.Team{
					ID:   1,
					Name: "Test Team",
					Members: []domain.TeamMember{
						{
							UserID: 2,
							Role:   domain.RoleMember,
						},
					},
				}, nil)
			},
			teamID: domain.TeamID(1),
			expectedTeam: domain.Team{
				ID:   1,
				Name: "Test Team",
				Members: []domain.TeamMember{
					{
						UserID: 2,
						Role:   domain.RoleMember,
					},
				},
			},
			expectedError: false,
		},
		{
			name: "Success - Team Member",
			setupMocks: func(
				mockTeamsRepo *mockcontract.MockTeamsRepository,
				mockUsersRepo *mockcontract.MockUsersRepository,
			) {
				// Get current user
				mockUsersRepo.EXPECT().GetByID(
					mock.Anything,
					domain.UserID(1),
				).Return(domain.User{
					ID:          1,
					Username:    "user1",
					IsSuperuser: false,
				}, nil)

				// Get team
				mockTeamsRepo.EXPECT().GetByID(
					mock.Anything,
					domain.TeamID(1),
				).Return(domain.Team{
					ID:   1,
					Name: "Test Team",
					Members: []domain.TeamMember{
						{
							UserID: 1,
							Role:   domain.RoleMember,
						},
					},
				}, nil)
			},
			teamID: domain.TeamID(1),
			expectedTeam: domain.Team{
				ID:   1,
				Name: "Test Team",
				Members: []domain.TeamMember{
					{
						UserID: 1,
						Role:   domain.RoleMember,
					},
				},
			},
			expectedError: false,
		},
		{
			name: "Error - Current User Not Found",
			setupMocks: func(
				mockTeamsRepo *mockcontract.MockTeamsRepository,
				mockUsersRepo *mockcontract.MockUsersRepository,
			) {
				// Get current user
				mockUsersRepo.EXPECT().GetByID(
					mock.Anything,
					domain.UserID(1),
				).Return(domain.User{}, domain.ErrEntityNotFound)
			},
			teamID:        domain.TeamID(1),
			expectedTeam:  domain.Team{},
			expectedError: true,
			errorContains: "get current user by id",
		},
		{
			name: "Error - Team Not Found",
			setupMocks: func(
				mockTeamsRepo *mockcontract.MockTeamsRepository,
				mockUsersRepo *mockcontract.MockUsersRepository,
			) {
				// Get current user
				mockUsersRepo.EXPECT().GetByID(
					mock.Anything,
					domain.UserID(1),
				).Return(domain.User{
					ID:          1,
					Username:    "user1",
					IsSuperuser: false,
				}, nil)

				// Get team
				mockTeamsRepo.EXPECT().GetByID(
					mock.Anything,
					domain.TeamID(999),
				).Return(domain.Team{}, domain.ErrEntityNotFound)
			},
			teamID:        domain.TeamID(999),
			expectedTeam:  domain.Team{},
			expectedError: true,
			errorContains: "not found",
		},
		{
			name: "Error - Not a Team Member and Not Superuser",
			setupMocks: func(
				mockTeamsRepo *mockcontract.MockTeamsRepository,
				mockUsersRepo *mockcontract.MockUsersRepository,
			) {
				// Get current user
				mockUsersRepo.EXPECT().GetByID(
					mock.Anything,
					domain.UserID(1),
				).Return(domain.User{
					ID:          1,
					Username:    "user1",
					IsSuperuser: false,
				}, nil)

				// Get team
				mockTeamsRepo.EXPECT().GetByID(
					mock.Anything,
					domain.TeamID(1),
				).Return(domain.Team{
					ID:   1,
					Name: "Test Team",
					Members: []domain.TeamMember{
						{
							UserID: 2,
							Role:   domain.RoleMember,
						},
					},
				}, nil)
			},
			teamID:        domain.TeamID(1),
			expectedTeam:  domain.Team{},
			expectedError: true,
			errorContains: "forbidden",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// Create mocks
			mockTxManager := mockdb.NewMockTxManager(t)
			mockTeamsRepo := mockcontract.NewMockTeamsRepository(t)
			mockUsersRepo := mockcontract.NewMockUsersRepository(t)
			mockProjectsRepo := mockcontract.NewMockProjectsRepository(t)

			// Setup mocks
			tt.setupMocks(mockTeamsRepo, mockUsersRepo)

			// Create service
			mockUserNotificationsUseCase := mockcontract.NewMockUserNotificationsUseCase(t)
			service := New(mockTxManager, mockTeamsRepo, mockUsersRepo, mockUserNotificationsUseCase, mockProjectsRepo)

			// Create context with user ID
			ctx := wardencontext.WithUserID(context.Background(), 1)

			// Call method
			team, err := service.GetByID(ctx, tt.teamID)

			// Check result
			if tt.expectedError {
				require.Error(t, err)
				if tt.errorContains != "" {
					require.Contains(t, err.Error(), tt.errorContains)
				}
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.expectedTeam, team)
			}
		})
	}
}

func TestList(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		setupMocks func(
			mockTeamsRepo *mockcontract.MockTeamsRepository,
			mockUsersRepo *mockcontract.MockUsersRepository,
		)
		expectedTeams []domain.Team
		expectedError bool
		errorContains string
	}{
		{
			name: "Success - Superuser",
			setupMocks: func(
				mockTeamsRepo *mockcontract.MockTeamsRepository,
				mockUsersRepo *mockcontract.MockUsersRepository,
			) {
				// Get current user
				mockUsersRepo.EXPECT().GetByID(
					mock.Anything,
					domain.UserID(1),
				).Return(domain.User{
					ID:          1,
					Username:    "admin",
					IsSuperuser: true,
				}, nil)

				// List all teams
				mockTeamsRepo.EXPECT().List(
					mock.Anything,
				).Return([]domain.Team{
					{
						ID:   1,
						Name: "Team 1",
						Members: []domain.TeamMember{
							{
								UserID: 2,
								Role:   domain.RoleMember,
							},
						},
					},
					{
						ID:   2,
						Name: "Team 2",
						Members: []domain.TeamMember{
							{
								UserID: 3,
								Role:   domain.RoleMember,
							},
						},
					},
				}, nil)
			},
			expectedTeams: []domain.Team{
				{
					ID:   1,
					Name: "Team 1",
					Members: []domain.TeamMember{
						{
							UserID: 2,
							Role:   domain.RoleMember,
						},
					},
				},
				{
					ID:   2,
					Name: "Team 2",
					Members: []domain.TeamMember{
						{
							UserID: 3,
							Role:   domain.RoleMember,
						},
					},
				},
			},
			expectedError: false,
		},
		{
			name: "Success - Regular User",
			setupMocks: func(
				mockTeamsRepo *mockcontract.MockTeamsRepository,
				mockUsersRepo *mockcontract.MockUsersRepository,
			) {
				// Get current user
				mockUsersRepo.EXPECT().GetByID(
					mock.Anything,
					domain.UserID(1),
				).Return(domain.User{
					ID:          1,
					Username:    "user1",
					IsSuperuser: false,
				}, nil)

				// Get teams by user ID
				mockTeamsRepo.EXPECT().GetTeamsByUserID(
					mock.Anything,
					domain.UserID(1),
				).Return([]domain.Team{
					{
						ID:   1,
						Name: "Team 1",
						Members: []domain.TeamMember{
							{
								UserID: 1,
								Role:   domain.RoleMember,
							},
						},
					},
				}, nil)
			},
			expectedTeams: []domain.Team{
				{
					ID:   1,
					Name: "Team 1",
					Members: []domain.TeamMember{
						{
							UserID: 1,
							Role:   domain.RoleMember,
						},
					},
				},
			},
			expectedError: false,
		},
		{
			name: "Error - Current User Not Found",
			setupMocks: func(
				mockTeamsRepo *mockcontract.MockTeamsRepository,
				mockUsersRepo *mockcontract.MockUsersRepository,
			) {
				// Get current user
				mockUsersRepo.EXPECT().GetByID(
					mock.Anything,
					domain.UserID(1),
				).Return(domain.User{}, domain.ErrEntityNotFound)
			},
			expectedTeams: nil,
			expectedError: true,
			errorContains: "get current user by id",
		},
		{
			name: "Error - Database Error When Listing Teams",
			setupMocks: func(
				mockTeamsRepo *mockcontract.MockTeamsRepository,
				mockUsersRepo *mockcontract.MockUsersRepository,
			) {
				// Get current user
				mockUsersRepo.EXPECT().GetByID(
					mock.Anything,
					domain.UserID(1),
				).Return(domain.User{
					ID:          1,
					Username:    "admin",
					IsSuperuser: true,
				}, nil)

				// List all teams
				mockTeamsRepo.EXPECT().List(
					mock.Anything,
				).Return(nil, errors.New("database error"))
			},
			expectedTeams: nil,
			expectedError: true,
			errorContains: "database error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// Create mocks
			mockTxManager := mockdb.NewMockTxManager(t)
			mockTeamsRepo := mockcontract.NewMockTeamsRepository(t)
			mockUsersRepo := mockcontract.NewMockUsersRepository(t)
			mockProjectsRepo := mockcontract.NewMockProjectsRepository(t)

			// Setup mocks
			tt.setupMocks(mockTeamsRepo, mockUsersRepo)

			// Create service
			mockUserNotificationsUseCase := mockcontract.NewMockUserNotificationsUseCase(t)
			service := New(mockTxManager, mockTeamsRepo, mockUsersRepo, mockUserNotificationsUseCase, mockProjectsRepo)

			// Create context with user ID
			ctx := wardencontext.WithUserID(context.Background(), 1)

			// Call method
			teams, err := service.List(ctx)

			// Check result
			if tt.expectedError {
				require.Error(t, err)
				if tt.errorContains != "" {
					require.Contains(t, err.Error(), tt.errorContains)
				}
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.expectedTeams, teams)
			}
		})
	}
}

func TestDelete(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		setupMocks func(
			mockTeamsRepo *mockcontract.MockTeamsRepository,
			mockUsersRepo *mockcontract.MockUsersRepository,
			mockProjectsRepo *mockcontract.MockProjectsRepository,
		)
		teamID        domain.TeamID
		expectedError bool
		errorContains string
	}{
		{
			name: "Success - Superuser",
			setupMocks: func(
				mockTeamsRepo *mockcontract.MockTeamsRepository,
				mockUsersRepo *mockcontract.MockUsersRepository,
				mockProjectsRepo *mockcontract.MockProjectsRepository,
			) {
				// Get current user
				mockUsersRepo.EXPECT().GetByID(
					mock.Anything,
					domain.UserID(1),
				).Return(domain.User{
					ID:          1,
					Username:    "admin",
					IsSuperuser: true,
				}, nil)

				// Get team
				mockTeamsRepo.EXPECT().GetByID(
					mock.Anything,
					domain.TeamID(1),
				).Return(domain.Team{
					ID:   1,
					Name: "Test Team",
					Members: []domain.TeamMember{
						{
							UserID: 2,
							Role:   domain.RoleMember,
						},
					},
				}, nil)

				// Check if team has projects
				mockProjectsRepo.EXPECT().List(
					mock.Anything,
				).Return([]domain.ProjectExtended{
					{
						Project: domain.Project{
							ID:   1,
							Name: "Project 1",
						},
					},
				}, nil)

				// Delete team
				mockTeamsRepo.EXPECT().Delete(
					mock.Anything,
					domain.TeamID(1),
				).Return(nil)
			},
			teamID:        domain.TeamID(1),
			expectedError: false,
		},
		{
			name: "Success - Team Owner",
			setupMocks: func(
				mockTeamsRepo *mockcontract.MockTeamsRepository,
				mockUsersRepo *mockcontract.MockUsersRepository,
				mockProjectsRepo *mockcontract.MockProjectsRepository,
			) {
				// Get current user
				mockUsersRepo.EXPECT().GetByID(
					mock.Anything,
					domain.UserID(1),
				).Return(domain.User{
					ID:          1,
					Username:    "user1",
					IsSuperuser: false,
				}, nil)

				// Get team
				mockTeamsRepo.EXPECT().GetByID(
					mock.Anything,
					domain.TeamID(1),
				).Return(domain.Team{
					ID:   1,
					Name: "Test Team",
					Members: []domain.TeamMember{
						{
							UserID: 1,
							Role:   domain.RoleOwner,
						},
					},
				}, nil)

				// Check if team has projects
				mockProjectsRepo.EXPECT().List(
					mock.Anything,
				).Return([]domain.ProjectExtended{
					{
						Project: domain.Project{
							ID:   1,
							Name: "Project 1",
						},
					},
				}, nil)

				// Delete team
				mockTeamsRepo.EXPECT().Delete(
					mock.Anything,
					domain.TeamID(1),
				).Return(nil)
			},
			teamID:        domain.TeamID(1),
			expectedError: false,
		},
		{
			name: "Error - Current User Not Found",
			setupMocks: func(
				mockTeamsRepo *mockcontract.MockTeamsRepository,
				mockUsersRepo *mockcontract.MockUsersRepository,
				mockProjectsRepo *mockcontract.MockProjectsRepository,
			) {
				// Get current user
				mockUsersRepo.EXPECT().GetByID(
					mock.Anything,
					domain.UserID(1),
				).Return(domain.User{}, domain.ErrEntityNotFound)
			},
			teamID:        domain.TeamID(1),
			expectedError: true,
			errorContains: "get current user by id",
		},
		{
			name: "Error - Team Not Found",
			setupMocks: func(
				mockTeamsRepo *mockcontract.MockTeamsRepository,
				mockUsersRepo *mockcontract.MockUsersRepository,
				mockProjectsRepo *mockcontract.MockProjectsRepository,
			) {
				// Get current user
				mockUsersRepo.EXPECT().GetByID(
					mock.Anything,
					domain.UserID(1),
				).Return(domain.User{
					ID:          1,
					Username:    "user1",
					IsSuperuser: false,
				}, nil)

				// Get team
				mockTeamsRepo.EXPECT().GetByID(
					mock.Anything,
					domain.TeamID(999),
				).Return(domain.Team{}, domain.ErrEntityNotFound)
			},
			teamID:        domain.TeamID(999),
			expectedError: true,
			errorContains: "not found",
		},
		{
			name: "Error - Not Owner and Not Superuser",
			setupMocks: func(
				mockTeamsRepo *mockcontract.MockTeamsRepository,
				mockUsersRepo *mockcontract.MockUsersRepository,
				mockProjectsRepo *mockcontract.MockProjectsRepository,
			) {
				// Get current user
				mockUsersRepo.EXPECT().GetByID(
					mock.Anything,
					domain.UserID(1),
				).Return(domain.User{
					ID:          1,
					Username:    "user1",
					IsSuperuser: false,
				}, nil)

				// Get team
				mockTeamsRepo.EXPECT().GetByID(
					mock.Anything,
					domain.TeamID(1),
				).Return(domain.Team{
					ID:   1,
					Name: "Test Team",
					Members: []domain.TeamMember{
						{
							UserID: 1,
							Role:   domain.RoleMember,
						},
					},
				}, nil)
			},
			teamID:        domain.TeamID(1),
			expectedError: true,
			errorContains: "forbidden",
		},
		{
			name: "Error - Database Error When Deleting Team",
			setupMocks: func(
				mockTeamsRepo *mockcontract.MockTeamsRepository,
				mockUsersRepo *mockcontract.MockUsersRepository,
				mockProjectsRepo *mockcontract.MockProjectsRepository,
			) {
				// Get current user
				mockUsersRepo.EXPECT().GetByID(
					mock.Anything,
					domain.UserID(1),
				).Return(domain.User{
					ID:          1,
					Username:    "admin",
					IsSuperuser: true,
				}, nil)

				// Get team
				mockTeamsRepo.EXPECT().GetByID(
					mock.Anything,
					domain.TeamID(1),
				).Return(domain.Team{
					ID:   1,
					Name: "Test Team",
					Members: []domain.TeamMember{
						{
							UserID: 2,
							Role:   domain.RoleMember,
						},
					},
				}, nil)

				// Check if team has projects
				mockProjectsRepo.EXPECT().List(
					mock.Anything,
				).Return([]domain.ProjectExtended{
					{
						Project: domain.Project{
							ID:   1,
							Name: "Project 1",
						},
					},
				}, nil)

				// Delete team
				mockTeamsRepo.EXPECT().Delete(
					mock.Anything,
					domain.TeamID(1),
				).Return(errors.New("database error"))
			},
			teamID:        domain.TeamID(1),
			expectedError: true,
			errorContains: "database error",
		},
		{
			name: "Error - Team Has Projects",
			setupMocks: func(
				mockTeamsRepo *mockcontract.MockTeamsRepository,
				mockUsersRepo *mockcontract.MockUsersRepository,
				mockProjectsRepo *mockcontract.MockProjectsRepository,
			) {
				// Get current user
				mockUsersRepo.EXPECT().GetByID(
					mock.Anything,
					domain.UserID(1),
				).Return(domain.User{
					ID:          1,
					Username:    "admin",
					IsSuperuser: true,
				}, nil)

				// Get team
				mockTeamsRepo.EXPECT().GetByID(
					mock.Anything,
					domain.TeamID(1),
				).Return(domain.Team{
					ID:   1,
					Name: "Test Team",
					Members: []domain.TeamMember{
						{
							UserID: 2,
							Role:   domain.RoleMember,
						},
					},
				}, nil)

				// Check if team has projects
				teamID := domain.TeamID(1)
				mockProjectsRepo.EXPECT().List(
					mock.Anything,
				).Return([]domain.ProjectExtended{
					{
						Project: domain.Project{
							ID:     1,
							Name:   "Project 1",
							TeamID: &teamID,
						},
					},
				}, nil)
			},
			teamID:        domain.TeamID(1),
			expectedError: true,
			errorContains: domain.ErrTeamHasProjects.Error(),
		},
		{
			name: "Error - Failed to Check Team Projects",
			setupMocks: func(
				mockTeamsRepo *mockcontract.MockTeamsRepository,
				mockUsersRepo *mockcontract.MockUsersRepository,
				mockProjectsRepo *mockcontract.MockProjectsRepository,
			) {
				// Get current user
				mockUsersRepo.EXPECT().GetByID(
					mock.Anything,
					domain.UserID(1),
				).Return(domain.User{
					ID:          1,
					Username:    "admin",
					IsSuperuser: true,
				}, nil)

				// Get team
				mockTeamsRepo.EXPECT().GetByID(
					mock.Anything,
					domain.TeamID(1),
				).Return(domain.Team{
					ID:   1,
					Name: "Test Team",
					Members: []domain.TeamMember{
						{
							UserID: 2,
							Role:   domain.RoleMember,
						},
					},
				}, nil)

				// Check if team has projects
				mockProjectsRepo.EXPECT().List(
					mock.Anything,
				).Return(nil, errors.New("failed to list projects"))
			},
			teamID:        domain.TeamID(1),
			expectedError: true,
			errorContains: "check team projects",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// Create mocks
			mockTxManager := mockdb.NewMockTxManager(t)
			mockTeamsRepo := mockcontract.NewMockTeamsRepository(t)
			mockUsersRepo := mockcontract.NewMockUsersRepository(t)
			mockProjectsRepo := mockcontract.NewMockProjectsRepository(t)

			// Setup mocks
			tt.setupMocks(mockTeamsRepo, mockUsersRepo, mockProjectsRepo)

			// Create service
			mockUserNotificationsUseCase := mockcontract.NewMockUserNotificationsUseCase(t)
			service := New(mockTxManager, mockTeamsRepo, mockUsersRepo, mockUserNotificationsUseCase, mockProjectsRepo)

			// Create context with user ID
			ctx := wardencontext.WithUserID(context.Background(), 1)

			// Call method
			err := service.Delete(ctx, tt.teamID)

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

func TestAddMember(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		setupMocks func(
			mockTxManager *mockdb.MockTxManager,
			mockTeamsRepo *mockcontract.MockTeamsRepository,
			mockUsersRepo *mockcontract.MockUsersRepository,
			mockUserNotificationsUseCase *mockcontract.MockUserNotificationsUseCase,
		)
		teamID        domain.TeamID
		userID        domain.UserID
		role          domain.Role
		expectedError bool
		errorContains string
	}{
		{
			name: "Success - Superuser",
			setupMocks: func(
				mockTxManager *mockdb.MockTxManager,
				mockTeamsRepo *mockcontract.MockTeamsRepository,
				mockUsersRepo *mockcontract.MockUsersRepository,
				mockUserNotificationsUseCase *mockcontract.MockUserNotificationsUseCase,
			) {
				// Get current user
				mockUsersRepo.EXPECT().GetByID(
					mock.Anything,
					domain.UserID(1),
				).Return(domain.User{
					ID:          1,
					Username:    "admin",
					IsSuperuser: true,
				}, nil)

				// Check if user to add exists
				mockUsersRepo.EXPECT().GetByID(
					mock.Anything,
					domain.UserID(2),
				).Return(domain.User{
					ID:       2,
					Username: "user2",
				}, nil)

				// Get team
				mockTeamsRepo.EXPECT().GetByID(
					mock.Anything,
					domain.TeamID(1),
				).Return(domain.Team{
					ID:   1,
					Name: "Test Team",
					Members: []domain.TeamMember{
						{
							UserID: 3,
							Role:   domain.RoleMember,
						},
					},
				}, nil)

				mockTxManager.EXPECT().ReadCommitted(mock.Anything, mock.Anything).Return(nil)
				//// Add member
				//mockTeamsRepo.EXPECT().AddMember(
				//	mock.Anything,
				//	domain.TeamID(1),
				//	domain.UserID(2),
				//	domain.RoleMember,
				//).Return(nil)
				//
				//// In AddMember success cases, mockUserNotificationsUseCase.On("CreateNotification", ...).Return(nil)
				//mockUserNotificationsUseCase.On("CreateNotification", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)
			},
			teamID:        domain.TeamID(1),
			userID:        domain.UserID(2),
			role:          domain.RoleMember,
			expectedError: false,
		},
		{
			name: "Success - Team Owner",
			setupMocks: func(
				mockTxManager *mockdb.MockTxManager,
				mockTeamsRepo *mockcontract.MockTeamsRepository,
				mockUsersRepo *mockcontract.MockUsersRepository,
				mockUserNotificationsUseCase *mockcontract.MockUserNotificationsUseCase,
			) {
				// Get current user
				mockUsersRepo.EXPECT().GetByID(
					mock.Anything,
					domain.UserID(1),
				).Return(domain.User{
					ID:          1,
					Username:    "user1",
					IsSuperuser: false,
				}, nil)

				// Check if user to add exists
				mockUsersRepo.EXPECT().GetByID(
					mock.Anything,
					domain.UserID(2),
				).Return(domain.User{
					ID:       2,
					Username: "user2",
				}, nil)

				// Get team
				mockTeamsRepo.EXPECT().GetByID(
					mock.Anything,
					domain.TeamID(1),
				).Return(domain.Team{
					ID:   1,
					Name: "Test Team",
					Members: []domain.TeamMember{
						{
							UserID: 1,
							Role:   domain.RoleOwner,
						},
					},
				}, nil)

				mockTxManager.EXPECT().ReadCommitted(mock.Anything, mock.Anything).Return(nil)
				//// Add member
				//mockTeamsRepo.EXPECT().AddMember(
				//	mock.Anything,
				//	domain.TeamID(1),
				//	domain.UserID(2),
				//	domain.RoleMember,
				//).Return(nil)
				//
				//// In AddMember success cases, mockUserNotificationsUseCase.On("CreateNotification", ...).Return(nil)
				//mockUserNotificationsUseCase.On("CreateNotification", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)
			},
			teamID:        domain.TeamID(1),
			userID:        domain.UserID(2),
			role:          domain.RoleMember,
			expectedError: false,
		},
		{
			name: "Success - Team Admin",
			setupMocks: func(
				mockTxManager *mockdb.MockTxManager,
				mockTeamsRepo *mockcontract.MockTeamsRepository,
				mockUsersRepo *mockcontract.MockUsersRepository,
				mockUserNotificationsUseCase *mockcontract.MockUserNotificationsUseCase,
			) {
				// Get current user
				mockUsersRepo.EXPECT().GetByID(
					mock.Anything,
					domain.UserID(1),
				).Return(domain.User{
					ID:          1,
					Username:    "user1",
					IsSuperuser: false,
				}, nil)

				// Check if user to add exists
				mockUsersRepo.EXPECT().GetByID(
					mock.Anything,
					domain.UserID(2),
				).Return(domain.User{
					ID:       2,
					Username: "user2",
				}, nil)

				// Get team
				mockTeamsRepo.EXPECT().GetByID(
					mock.Anything,
					domain.TeamID(1),
				).Return(domain.Team{
					ID:   1,
					Name: "Test Team",
					Members: []domain.TeamMember{
						{
							UserID: 1,
							Role:   domain.RoleAdmin,
						},
						{
							UserID: 2,
							Role:   domain.RoleMember,
						},
					},
				}, nil)

				mockTxManager.EXPECT().ReadCommitted(mock.Anything, mock.Anything).Return(nil)
				//// Add member
				//mockTeamsRepo.EXPECT().AddMember(
				//	mock.Anything,
				//	domain.TeamID(1),
				//	domain.UserID(2),
				//	domain.RoleMember,
				//).Return(nil)
				//
				//// In AddMember success cases, mockUserNotificationsUseCase.On("CreateNotification", ...).Return(nil)
				//mockUserNotificationsUseCase.On("CreateNotification", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)
			},
			teamID:        domain.TeamID(1),
			userID:        domain.UserID(2),
			role:          domain.RoleMember,
			expectedError: false,
		},
		{
			name: "Error - Current User Not Found",
			setupMocks: func(
				mockTxManager *mockdb.MockTxManager,
				mockTeamsRepo *mockcontract.MockTeamsRepository,
				mockUsersRepo *mockcontract.MockUsersRepository,
				mockUserNotificationsUseCase *mockcontract.MockUserNotificationsUseCase,
			) {
				// Get current user
				mockUsersRepo.EXPECT().GetByID(
					mock.Anything,
					domain.UserID(1),
				).Return(domain.User{}, domain.ErrEntityNotFound)
			},
			teamID:        domain.TeamID(1),
			userID:        domain.UserID(2),
			role:          domain.RoleMember,
			expectedError: true,
			errorContains: "get current user by id",
		},
		{
			name: "Error - User to Add Not Found",
			setupMocks: func(
				mockTxManager *mockdb.MockTxManager,
				mockTeamsRepo *mockcontract.MockTeamsRepository,
				mockUsersRepo *mockcontract.MockUsersRepository,
				mockUserNotificationsUseCase *mockcontract.MockUserNotificationsUseCase,
			) {
				// Get current user
				mockUsersRepo.EXPECT().GetByID(
					mock.Anything,
					domain.UserID(1),
				).Return(domain.User{
					ID:          1,
					Username:    "admin",
					IsSuperuser: true,
				}, nil)

				// Check if user to add exists
				mockUsersRepo.EXPECT().GetByID(
					mock.Anything,
					domain.UserID(999),
				).Return(domain.User{}, domain.ErrEntityNotFound)
			},
			teamID:        domain.TeamID(1),
			userID:        domain.UserID(999),
			role:          domain.RoleMember,
			expectedError: true,
			errorContains: "not found",
		},
		{
			name: "Error - Team Not Found",
			setupMocks: func(
				mockTxManager *mockdb.MockTxManager,
				mockTeamsRepo *mockcontract.MockTeamsRepository,
				mockUsersRepo *mockcontract.MockUsersRepository,
				mockUserNotificationsUseCase *mockcontract.MockUserNotificationsUseCase,
			) {
				// Get current user
				mockUsersRepo.EXPECT().GetByID(
					mock.Anything,
					domain.UserID(1),
				).Return(domain.User{
					ID:          1,
					Username:    "admin",
					IsSuperuser: true,
				}, nil)

				// Check if user to add exists
				mockUsersRepo.EXPECT().GetByID(
					mock.Anything,
					domain.UserID(2),
				).Return(domain.User{
					ID:       2,
					Username: "user2",
				}, nil)

				// Get team
				mockTeamsRepo.EXPECT().GetByID(
					mock.Anything,
					domain.TeamID(999),
				).Return(domain.Team{}, domain.ErrEntityNotFound)
			},
			teamID:        domain.TeamID(999),
			userID:        domain.UserID(2),
			role:          domain.RoleMember,
			expectedError: true,
			errorContains: "not found",
		},
		{
			name: "Error - Not Admin/Owner and Not Superuser",
			setupMocks: func(
				mockTxManager *mockdb.MockTxManager,
				mockTeamsRepo *mockcontract.MockTeamsRepository,
				mockUsersRepo *mockcontract.MockUsersRepository,
				mockUserNotificationsUseCase *mockcontract.MockUserNotificationsUseCase,
			) {
				// Get current user
				mockUsersRepo.EXPECT().GetByID(
					mock.Anything,
					domain.UserID(1),
				).Return(domain.User{
					ID:          1,
					Username:    "user1",
					IsSuperuser: false,
				}, nil)

				// Check if user to add exists
				mockUsersRepo.EXPECT().GetByID(
					mock.Anything,
					domain.UserID(2),
				).Return(domain.User{
					ID:       2,
					Username: "user2",
				}, nil)

				// Get team
				mockTeamsRepo.EXPECT().GetByID(
					mock.Anything,
					domain.TeamID(1),
				).Return(domain.Team{
					ID:   1,
					Name: "Test Team",
					Members: []domain.TeamMember{
						{
							UserID: 1,
							Role:   domain.RoleMember,
						},
					},
				}, nil)
			},
			teamID:        domain.TeamID(1),
			userID:        domain.UserID(2),
			role:          domain.RoleMember,
			expectedError: true,
			errorContains: "forbidden",
		},
		{
			name: "Error - Database Error When Adding Member",
			setupMocks: func(
				mockTxManager *mockdb.MockTxManager,
				mockTeamsRepo *mockcontract.MockTeamsRepository,
				mockUsersRepo *mockcontract.MockUsersRepository,
				mockUserNotificationsUseCase *mockcontract.MockUserNotificationsUseCase,
			) {
				// Get current user
				mockUsersRepo.EXPECT().GetByID(
					mock.Anything,
					domain.UserID(1),
				).Return(domain.User{
					ID:          1,
					Username:    "admin",
					IsSuperuser: true,
				}, nil)

				// Check if user to add exists
				mockUsersRepo.EXPECT().GetByID(
					mock.Anything,
					domain.UserID(2),
				).Return(domain.User{
					ID:       2,
					Username: "user2",
				}, nil)

				// Get team
				mockTeamsRepo.EXPECT().GetByID(
					mock.Anything,
					domain.TeamID(1),
				).Return(domain.Team{
					ID:   1,
					Name: "Test Team",
					Members: []domain.TeamMember{
						{
							UserID: 3,
							Role:   domain.RoleMember,
						},
					},
				}, nil)

				mockTxManager.EXPECT().ReadCommitted(mock.Anything, mock.Anything).Return(errors.New("database error"))
				//// Add member
				//mockTeamsRepo.EXPECT().AddMember(
				//	mock.Anything,
				//	domain.TeamID(1),
				//	domain.UserID(2),
				//	domain.RoleMember,
				//).Return(errors.New("database error"))
			},
			teamID:        domain.TeamID(1),
			userID:        domain.UserID(2),
			role:          domain.RoleMember,
			expectedError: true,
			errorContains: "database error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mockTxManager := mockdb.NewMockTxManager(t)
			mockTeamsRepo := mockcontract.NewMockTeamsRepository(t)
			mockUsersRepo := mockcontract.NewMockUsersRepository(t)
			mockUserNotificationsUseCase := mockcontract.NewMockUserNotificationsUseCase(t)
			mockProjectsRepo := mockcontract.NewMockProjectsRepository(t)

			// Setup mocks
			tt.setupMocks(mockTxManager, mockTeamsRepo, mockUsersRepo, mockUserNotificationsUseCase)

			service := New(mockTxManager, mockTeamsRepo, mockUsersRepo, mockUserNotificationsUseCase, mockProjectsRepo)

			ctx := wardencontext.WithUserID(context.Background(), 1)
			err := service.AddMember(ctx, tt.teamID, tt.userID, tt.role)

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

func TestRemoveMemberWithChecks(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		setupMocks func(
			mockTxManager *mockdb.MockTxManager,
			mockTeamsRepo *mockcontract.MockTeamsRepository,
			mockUsersRepo *mockcontract.MockUsersRepository,
			mockUserNotificationsUseCase *mockcontract.MockUserNotificationsUseCase,
		)
		teamID        domain.TeamID
		userID        domain.UserID
		expectedError bool
		errorContains string
	}{
		{
			name: "Success - Superuser",
			setupMocks: func(
				mockTxManager *mockdb.MockTxManager,
				mockTeamsRepo *mockcontract.MockTeamsRepository,
				mockUsersRepo *mockcontract.MockUsersRepository,
				mockUserNotificationsUseCase *mockcontract.MockUserNotificationsUseCase,
			) {
				// Get current user
				mockUsersRepo.EXPECT().GetByID(
					mock.Anything,
					domain.UserID(1),
				).Return(domain.User{
					ID:          1,
					Username:    "admin",
					IsSuperuser: true,
				}, nil)

				mockUsersRepo.EXPECT().GetByID(
					mock.Anything,
					domain.UserID(2),
				).Return(domain.User{
					ID:          2,
					Username:    "user2",
					IsSuperuser: false,
				}, nil)

				// Get team
				mockTeamsRepo.EXPECT().GetByID(
					mock.Anything,
					domain.TeamID(1),
				).Return(domain.Team{
					ID:   1,
					Name: "Test Team",
					Members: []domain.TeamMember{
						{
							UserID: 2,
							Role:   domain.RoleMember,
						},
						{
							UserID: 3,
							Role:   domain.RoleOwner,
						},
					},
				}, nil)

				mockTxManager.EXPECT().ReadCommitted(mock.Anything, mock.Anything).Return(nil)

				//// Remove member
				//mockTeamsRepo.EXPECT().RemoveMember(
				//	mock.Anything,
				//	domain.TeamID(1),
				//	domain.UserID(2),
				//).Return(nil)
				//
				//// In RemoveMemberWithChecks success cases, mockUserNotificationsUseCase.On("CreateNotification", ...).Return(nil)
				//mockUserNotificationsUseCase.On("CreateNotification", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)
			},
			teamID:        domain.TeamID(1),
			userID:        domain.UserID(2),
			expectedError: false,
		},
		{
			name: "Success - Team Owner removing member",
			setupMocks: func(
				mockTxManager *mockdb.MockTxManager,
				mockTeamsRepo *mockcontract.MockTeamsRepository,
				mockUsersRepo *mockcontract.MockUsersRepository,
				mockUserNotificationsUseCase *mockcontract.MockUserNotificationsUseCase,
			) {
				// Get current user
				mockUsersRepo.EXPECT().GetByID(
					mock.Anything,
					domain.UserID(1),
				).Return(domain.User{
					ID:          1,
					Username:    "user1",
					IsSuperuser: false,
				}, nil)

				mockUsersRepo.EXPECT().GetByID(
					mock.Anything,
					domain.UserID(2),
				).Return(domain.User{
					ID:          2,
					Username:    "user2",
					IsSuperuser: false,
				}, nil)

				// Get team
				mockTeamsRepo.EXPECT().GetByID(
					mock.Anything,
					domain.TeamID(1),
				).Return(domain.Team{
					ID:   1,
					Name: "Test Team",
					Members: []domain.TeamMember{
						{
							UserID: 1,
							Role:   domain.RoleOwner,
						},
						{
							UserID: 2,
							Role:   domain.RoleMember,
						},
					},
				}, nil)

				mockTxManager.EXPECT().ReadCommitted(mock.Anything, mock.Anything).Return(nil)

				//// Remove member
				//mockTeamsRepo.EXPECT().RemoveMember(
				//	mock.Anything,
				//	domain.TeamID(1),
				//	domain.UserID(2),
				//).Return(nil)
				//
				//// In RemoveMemberWithChecks success cases, mockUserNotificationsUseCase.On("CreateNotification", ...).Return(nil)
				//mockUserNotificationsUseCase.On("CreateNotification", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)
			},
			teamID:        domain.TeamID(1),
			userID:        domain.UserID(2),
			expectedError: false,
		},
		{
			name: "Success - Team Admin removing member",
			setupMocks: func(
				mockTxManager *mockdb.MockTxManager,
				mockTeamsRepo *mockcontract.MockTeamsRepository,
				mockUsersRepo *mockcontract.MockUsersRepository,
				mockUserNotificationsUseCase *mockcontract.MockUserNotificationsUseCase,
			) {
				// Get current user
				mockUsersRepo.EXPECT().GetByID(
					mock.Anything,
					domain.UserID(1),
				).Return(domain.User{
					ID:          1,
					Username:    "user1",
					IsSuperuser: false,
				}, nil)

				mockUsersRepo.EXPECT().GetByID(
					mock.Anything,
					domain.UserID(2),
				).Return(domain.User{
					ID:          2,
					Username:    "user2",
					IsSuperuser: false,
				}, nil)

				// Get team
				mockTeamsRepo.EXPECT().GetByID(
					mock.Anything,
					domain.TeamID(1),
				).Return(domain.Team{
					ID:   1,
					Name: "Test Team",
					Members: []domain.TeamMember{
						{
							UserID: 1,
							Role:   domain.RoleAdmin,
						},
						{
							UserID: 2,
							Role:   domain.RoleMember,
						},
					},
				}, nil)

				mockTxManager.EXPECT().ReadCommitted(mock.Anything, mock.Anything).Return(nil)

				//// Remove member
				//mockTeamsRepo.EXPECT().RemoveMember(
				//	mock.Anything,
				//	domain.TeamID(1),
				//	domain.UserID(2),
				//).Return(nil)
				//
				//// In RemoveMemberWithChecks success cases, mockUserNotificationsUseCase.On("CreateNotification", ...).Return(nil)
				//mockUserNotificationsUseCase.On("CreateNotification", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)
			},
			teamID:        domain.TeamID(1),
			userID:        domain.UserID(2),
			expectedError: false,
		},
		{
			name: "Error - Current User Not Found",
			setupMocks: func(
				mockTxManager *mockdb.MockTxManager,
				mockTeamsRepo *mockcontract.MockTeamsRepository,
				mockUsersRepo *mockcontract.MockUsersRepository,
				mockUserNotificationsUseCase *mockcontract.MockUserNotificationsUseCase,
			) {
				// Get current user
				mockUsersRepo.EXPECT().GetByID(
					mock.Anything,
					domain.UserID(1),
				).Return(domain.User{}, domain.ErrEntityNotFound)
			},
			teamID:        domain.TeamID(1),
			userID:        domain.UserID(2),
			expectedError: true,
			errorContains: "get current user by id",
		},
		{
			name: "Error - Team Not Found",
			setupMocks: func(
				mockTxManager *mockdb.MockTxManager,
				mockTeamsRepo *mockcontract.MockTeamsRepository,
				mockUsersRepo *mockcontract.MockUsersRepository,
				mockUserNotificationsUseCase *mockcontract.MockUserNotificationsUseCase,
			) {
				// Get current user
				mockUsersRepo.EXPECT().GetByID(
					mock.Anything,
					domain.UserID(1),
				).Return(domain.User{
					ID:          1,
					Username:    "admin",
					IsSuperuser: true,
				}, nil)

				// Get team
				mockTeamsRepo.EXPECT().GetByID(
					mock.Anything,
					domain.TeamID(999),
				).Return(domain.Team{}, domain.ErrEntityNotFound)
			},
			teamID:        domain.TeamID(999),
			userID:        domain.UserID(2),
			expectedError: true,
			errorContains: "not found",
		},
		{
			name: "Error - Not Admin/Owner and Not Superuser",
			setupMocks: func(
				mockTxManager *mockdb.MockTxManager,
				mockTeamsRepo *mockcontract.MockTeamsRepository,
				mockUsersRepo *mockcontract.MockUsersRepository,
				mockUserNotificationsUseCase *mockcontract.MockUserNotificationsUseCase,
			) {
				// Get current user
				mockUsersRepo.EXPECT().GetByID(
					mock.Anything,
					domain.UserID(1),
				).Return(domain.User{
					ID:          1,
					Username:    "user1",
					IsSuperuser: false,
				}, nil)

				// Get team
				mockTeamsRepo.EXPECT().GetByID(
					mock.Anything,
					domain.TeamID(1),
				).Return(domain.Team{
					ID:   1,
					Name: "Test Team",
					Members: []domain.TeamMember{
						{
							UserID: 1,
							Role:   domain.RoleMember,
						},
					},
				}, nil)
			},
			teamID:        domain.TeamID(1),
			userID:        domain.UserID(2),
			expectedError: true,
			errorContains: "forbidden",
		},
		{
			name: "Error - Non-Owner trying to remove Owner",
			setupMocks: func(
				mockTxManager *mockdb.MockTxManager,
				mockTeamsRepo *mockcontract.MockTeamsRepository,
				mockUsersRepo *mockcontract.MockUsersRepository,
				mockUserNotificationsUseCase *mockcontract.MockUserNotificationsUseCase,
			) {
				// Get current user
				mockUsersRepo.EXPECT().GetByID(
					mock.Anything,
					domain.UserID(1),
				).Return(domain.User{
					ID:          1,
					Username:    "user1",
					IsSuperuser: false,
				}, nil)

				// Get team
				mockTeamsRepo.EXPECT().GetByID(
					mock.Anything,
					domain.TeamID(1),
				).Return(domain.Team{
					ID:   1,
					Name: "Test Team",
					Members: []domain.TeamMember{
						{
							UserID: 1,
							Role:   domain.RoleAdmin,
						},
						{
							UserID: 2,
							Role:   domain.RoleOwner,
						},
					},
				}, nil)
			},
			teamID:        domain.TeamID(1),
			userID:        domain.UserID(2),
			expectedError: true,
			errorContains: "forbidden",
		},
		{
			name: "Error - Database Error When Removing Member",
			setupMocks: func(
				mockTxManager *mockdb.MockTxManager,
				mockTeamsRepo *mockcontract.MockTeamsRepository,
				mockUsersRepo *mockcontract.MockUsersRepository,
				mockUserNotificationsUseCase *mockcontract.MockUserNotificationsUseCase,
			) {
				// Get current user
				mockUsersRepo.EXPECT().GetByID(
					mock.Anything,
					domain.UserID(1),
				).Return(domain.User{
					ID:          1,
					Username:    "admin",
					IsSuperuser: true,
				}, nil)

				mockUsersRepo.EXPECT().GetByID(
					mock.Anything,
					domain.UserID(2),
				).Return(domain.User{
					ID:          2,
					Username:    "user2",
					IsSuperuser: false,
				}, nil)

				// Get team
				mockTeamsRepo.EXPECT().GetByID(
					mock.Anything,
					domain.TeamID(1),
				).Return(domain.Team{
					ID:   1,
					Name: "Test Team",
					Members: []domain.TeamMember{
						{
							UserID: 2,
							Role:   domain.RoleMember,
						},
					},
				}, nil)

				mockTxManager.EXPECT().ReadCommitted(mock.Anything, mock.Anything).Return(errors.New("database error"))
				//// Remove member
				//mockTeamsRepo.EXPECT().RemoveMember(
				//	mock.Anything,
				//	domain.TeamID(1),
				//	domain.UserID(2),
				//).Return(errors.New("database error"))
			},
			teamID:        domain.TeamID(1),
			userID:        domain.UserID(2),
			expectedError: true,
			errorContains: "database error",
		},
		{
			name: "Success - self-leave as owner, not last owner",
			setupMocks: func(
				mockTxManager *mockdb.MockTxManager,
				mockTeamsRepo *mockcontract.MockTeamsRepository,
				mockUsersRepo *mockcontract.MockUsersRepository,
				mockUserNotificationsUseCase *mockcontract.MockUserNotificationsUseCase,
			) {
				mockUsersRepo.EXPECT().GetByID(mock.Anything, domain.UserID(1)).Return(domain.User{ID: 1, Username: "owner1", IsSuperuser: false}, nil)
				mockTeamsRepo.EXPECT().GetByID(mock.Anything, domain.TeamID(1)).Return(domain.Team{
					ID:   1,
					Name: "Test Team",
					Members: []domain.TeamMember{
						{UserID: 1, Role: domain.RoleOwner},
						{UserID: 2, Role: domain.RoleOwner},
					},
				}, nil)
				mockUsersRepo.EXPECT().GetByID(mock.Anything, domain.UserID(1)).Return(domain.User{ID: 1, Username: "owner1", IsSuperuser: false}, nil)

				mockTxManager.EXPECT().ReadCommitted(mock.Anything, mock.Anything).Return(nil)
				//mockTeamsRepo.EXPECT().RemoveMember(mock.Anything, domain.TeamID(1), domain.UserID(1)).Return(nil)
				//mockUserNotificationsUseCase.On("CreateNotification", mock.Anything, domain.UserID(1), domain.UserNotificationTypeTeamRemoved, mock.Anything).Return(nil)
			},
			teamID:        domain.TeamID(1),
			userID:        domain.UserID(1),
			expectedError: false,
		},
		{
			name: "Error - self-leave as last owner",
			setupMocks: func(
				mockTxManager *mockdb.MockTxManager,
				mockTeamsRepo *mockcontract.MockTeamsRepository,
				mockUsersRepo *mockcontract.MockUsersRepository,
				mockUserNotificationsUseCase *mockcontract.MockUserNotificationsUseCase,
			) {
				mockUsersRepo.EXPECT().GetByID(mock.Anything, domain.UserID(1)).Return(domain.User{ID: 1, Username: "owner1", IsSuperuser: false}, nil)
				mockTeamsRepo.EXPECT().GetByID(mock.Anything, domain.TeamID(1)).Return(domain.Team{
					ID:   1,
					Name: "Test Team",
					Members: []domain.TeamMember{
						{UserID: 1, Role: domain.RoleOwner},
					},
				}, nil)
			},
			teamID:        domain.TeamID(1),
			userID:        domain.UserID(1),
			expectedError: true,
			errorContains: domain.ErrLastOwner.Error(),
		},
		{
			name: "Error - self-leave as member (not owner)",
			setupMocks: func(
				mockTxManager *mockdb.MockTxManager,
				mockTeamsRepo *mockcontract.MockTeamsRepository,
				mockUsersRepo *mockcontract.MockUsersRepository,
				mockUserNotificationsUseCase *mockcontract.MockUserNotificationsUseCase,
			) {
				mockUsersRepo.EXPECT().GetByID(mock.Anything, domain.UserID(1)).Return(domain.User{ID: 1, Username: "member1", IsSuperuser: false}, nil)
				mockTeamsRepo.EXPECT().GetByID(mock.Anything, domain.TeamID(1)).Return(domain.Team{
					ID:   1,
					Name: "Test Team",
					Members: []domain.TeamMember{
						{UserID: 1, Role: domain.RoleMember},
						{UserID: 2, Role: domain.RoleOwner},
					},
				}, nil)
				mockUsersRepo.EXPECT().GetByID(mock.Anything, domain.UserID(1)).Return(domain.User{ID: 1, Username: "member1", IsSuperuser: false}, nil)

				mockTxManager.EXPECT().ReadCommitted(mock.Anything, mock.Anything).Return(nil)
				//mockTeamsRepo.EXPECT().RemoveMember(mock.Anything, domain.TeamID(1), domain.UserID(1)).Return(nil)
				//mockUserNotificationsUseCase.On("CreateNotification", mock.Anything, domain.UserID(1), domain.UserNotificationTypeTeamRemoved, mock.Anything).Return(nil)
			},
			teamID:        domain.TeamID(1),
			userID:        domain.UserID(1),
			expectedError: false,
		},
		{
			name: "Error - self-leave as superuser",
			setupMocks: func(
				mockTxManager *mockdb.MockTxManager,
				mockTeamsRepo *mockcontract.MockTeamsRepository,
				mockUsersRepo *mockcontract.MockUsersRepository,
				mockUserNotificationsUseCase *mockcontract.MockUserNotificationsUseCase,
			) {
				mockUsersRepo.EXPECT().GetByID(mock.Anything, domain.UserID(1)).Return(domain.User{ID: 1, Username: "superuser", IsSuperuser: true}, nil)
				mockTeamsRepo.EXPECT().GetByID(mock.Anything, domain.TeamID(1)).Return(domain.Team{
					ID:   1,
					Name: "Test Team",
					Members: []domain.TeamMember{
						{UserID: 1, Role: domain.RoleOwner},
						{UserID: 2, Role: domain.RoleOwner},
					},
				}, nil)
				mockUsersRepo.EXPECT().GetByID(mock.Anything, domain.UserID(1)).Return(domain.User{ID: 1, Username: "superuser", IsSuperuser: true}, nil)
			},
			teamID:        domain.TeamID(1),
			userID:        domain.UserID(1),
			expectedError: true,
			errorContains: domain.ErrForbidden.Error(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mockTxManager := mockdb.NewMockTxManager(t)
			mockTeamsRepo := mockcontract.NewMockTeamsRepository(t)
			mockUsersRepo := mockcontract.NewMockUsersRepository(t)
			mockUserNotificationsUseCase := mockcontract.NewMockUserNotificationsUseCase(t)
			mockProjectsRepo := mockcontract.NewMockProjectsRepository(t)

			// Setup mocks
			tt.setupMocks(mockTxManager, mockTeamsRepo, mockUsersRepo, mockUserNotificationsUseCase)

			service := New(mockTxManager, mockTeamsRepo, mockUsersRepo, mockUserNotificationsUseCase, mockProjectsRepo)

			ctx := wardencontext.WithUserID(context.Background(), 1)
			err := service.RemoveMemberWithChecks(ctx, tt.teamID, tt.userID)

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

func TestGetTeamsByUserID(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		setupMocks    func(mockTeamsRepo *mockcontract.MockTeamsRepository)
		userID        domain.UserID
		expectedTeams []domain.Team
		expectedError bool
		errorContains string
	}{
		{
			name: "Success",
			setupMocks: func(mockTeamsRepo *mockcontract.MockTeamsRepository) {
				mockTeamsRepo.EXPECT().GetTeamsByUserID(
					mock.Anything,
					domain.UserID(1),
				).Return([]domain.Team{
					{
						ID:   1,
						Name: "Team 1",
						Members: []domain.TeamMember{
							{
								UserID: 1,
								Role:   domain.RoleMember,
							},
						},
					},
					{
						ID:   2,
						Name: "Team 2",
						Members: []domain.TeamMember{
							{
								UserID: 1,
								Role:   domain.RoleAdmin,
							},
						},
					},
				}, nil)
			},
			userID: domain.UserID(1),
			expectedTeams: []domain.Team{
				{
					ID:   1,
					Name: "Team 1",
					Members: []domain.TeamMember{
						{
							UserID: 1,
							Role:   domain.RoleMember,
						},
					},
				},
				{
					ID:   2,
					Name: "Team 2",
					Members: []domain.TeamMember{
						{
							UserID: 1,
							Role:   domain.RoleAdmin,
						},
					},
				},
			},
			expectedError: false,
		},
		{
			name: "No Teams",
			setupMocks: func(mockTeamsRepo *mockcontract.MockTeamsRepository) {
				mockTeamsRepo.EXPECT().GetTeamsByUserID(
					mock.Anything,
					domain.UserID(2),
				).Return([]domain.Team{}, nil)
			},
			userID:        domain.UserID(2),
			expectedTeams: []domain.Team{},
			expectedError: false,
		},
		{
			name: "Database Error",
			setupMocks: func(mockTeamsRepo *mockcontract.MockTeamsRepository) {
				mockTeamsRepo.EXPECT().GetTeamsByUserID(
					mock.Anything,
					domain.UserID(3),
				).Return(nil, errors.New("database error"))
			},
			userID:        domain.UserID(3),
			expectedTeams: nil,
			expectedError: true,
			errorContains: "database error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// Create mocks
			mockTxManager := mockdb.NewMockTxManager(t)
			mockTeamsRepo := mockcontract.NewMockTeamsRepository(t)
			mockUsersRepo := mockcontract.NewMockUsersRepository(t)
			mockProjectsRepo := mockcontract.NewMockProjectsRepository(t)

			// Setup mocks
			tt.setupMocks(mockTeamsRepo)

			// Create service
			mockUserNotificationsUseCase := mockcontract.NewMockUserNotificationsUseCase(t)
			service := New(mockTxManager, mockTeamsRepo, mockUsersRepo, mockUserNotificationsUseCase, mockProjectsRepo)

			// Call method
			teams, err := service.GetTeamsByUserID(context.Background(), tt.userID)

			// Check result
			if tt.expectedError {
				require.Error(t, err)
				if tt.errorContains != "" {
					require.Contains(t, err.Error(), tt.errorContains)
				}
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.expectedTeams, teams)
			}
		})
	}
}

func TestGetTeamByID(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		setupMocks    func(mockTeamsRepo *mockcontract.MockTeamsRepository)
		teamID        domain.TeamID
		expectedTeam  domain.Team
		expectedError bool
		errorContains string
	}{
		{
			name: "Success",
			setupMocks: func(mockTeamsRepo *mockcontract.MockTeamsRepository) {
				mockTeamsRepo.EXPECT().GetByID(
					mock.Anything,
					domain.TeamID(1),
				).Return(domain.Team{
					ID:   1,
					Name: "Team 1",
					Members: []domain.TeamMember{
						{
							UserID: 1,
							Role:   domain.RoleMember,
						},
					},
				}, nil)
			},
			teamID: domain.TeamID(1),
			expectedTeam: domain.Team{
				ID:   1,
				Name: "Team 1",
				Members: []domain.TeamMember{
					{
						UserID: 1,
						Role:   domain.RoleMember,
					},
				},
			},
			expectedError: false,
		},
		{
			name: "Team Not Found",
			setupMocks: func(mockTeamsRepo *mockcontract.MockTeamsRepository) {
				mockTeamsRepo.EXPECT().GetByID(
					mock.Anything,
					domain.TeamID(999),
				).Return(domain.Team{}, domain.ErrEntityNotFound)
			},
			teamID:        domain.TeamID(999),
			expectedTeam:  domain.Team{},
			expectedError: true,
			errorContains: "not found",
		},
		{
			name: "Database Error",
			setupMocks: func(mockTeamsRepo *mockcontract.MockTeamsRepository) {
				mockTeamsRepo.EXPECT().GetByID(
					mock.Anything,
					domain.TeamID(2),
				).Return(domain.Team{}, errors.New("database error"))
			},
			teamID:        domain.TeamID(2),
			expectedTeam:  domain.Team{},
			expectedError: true,
			errorContains: "database error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// Create mocks
			mockTxManager := mockdb.NewMockTxManager(t)
			mockTeamsRepo := mockcontract.NewMockTeamsRepository(t)
			mockUsersRepo := mockcontract.NewMockUsersRepository(t)
			mockProjectsRepo := mockcontract.NewMockProjectsRepository(t)

			// Setup mocks
			tt.setupMocks(mockTeamsRepo)

			// Create service
			mockUserNotificationsUseCase := mockcontract.NewMockUserNotificationsUseCase(t)
			service := New(mockTxManager, mockTeamsRepo, mockUsersRepo, mockUserNotificationsUseCase, mockProjectsRepo)

			// Call method
			team, err := service.GetTeamByID(context.Background(), tt.teamID)

			// Check result
			if tt.expectedError {
				require.Error(t, err)
				if tt.errorContains != "" {
					require.Contains(t, err.Error(), tt.errorContains)
				}
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.expectedTeam, team)
			}
		})
	}
}

func TestGetMembers(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name            string
		setupMocks      func(mockTeamsRepo *mockcontract.MockTeamsRepository)
		teamID          domain.TeamID
		expectedMembers []domain.TeamMember
		expectedError   bool
		errorContains   string
	}{
		{
			name: "Success",
			setupMocks: func(mockTeamsRepo *mockcontract.MockTeamsRepository) {
				mockTeamsRepo.EXPECT().GetMembers(
					mock.Anything,
					domain.TeamID(1),
				).Return([]domain.TeamMember{
					{
						UserID: 1,
						Role:   domain.RoleMember,
					},
					{
						UserID: 2,
						Role:   domain.RoleAdmin,
					},
					{
						UserID: 3,
						Role:   domain.RoleOwner,
					},
				}, nil)
			},
			teamID: domain.TeamID(1),
			expectedMembers: []domain.TeamMember{
				{
					UserID: 1,
					Role:   domain.RoleMember,
				},
				{
					UserID: 2,
					Role:   domain.RoleAdmin,
				},
				{
					UserID: 3,
					Role:   domain.RoleOwner,
				},
			},
			expectedError: false,
		},
		{
			name: "No Members",
			setupMocks: func(mockTeamsRepo *mockcontract.MockTeamsRepository) {
				mockTeamsRepo.EXPECT().GetMembers(
					mock.Anything,
					domain.TeamID(2),
				).Return([]domain.TeamMember{}, nil)
			},
			teamID:          domain.TeamID(2),
			expectedMembers: []domain.TeamMember{},
			expectedError:   false,
		},
		{
			name: "Team Not Found",
			setupMocks: func(mockTeamsRepo *mockcontract.MockTeamsRepository) {
				mockTeamsRepo.EXPECT().GetMembers(
					mock.Anything,
					domain.TeamID(999),
				).Return(nil, domain.ErrEntityNotFound)
			},
			teamID:          domain.TeamID(999),
			expectedMembers: nil,
			expectedError:   true,
			errorContains:   "not found",
		},
		{
			name: "Database Error",
			setupMocks: func(mockTeamsRepo *mockcontract.MockTeamsRepository) {
				mockTeamsRepo.EXPECT().GetMembers(
					mock.Anything,
					domain.TeamID(3),
				).Return(nil, errors.New("database error"))
			},
			teamID:          domain.TeamID(3),
			expectedMembers: nil,
			expectedError:   true,
			errorContains:   "database error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// Create mocks
			mockTxManager := mockdb.NewMockTxManager(t)
			mockTeamsRepo := mockcontract.NewMockTeamsRepository(t)
			mockUsersRepo := mockcontract.NewMockUsersRepository(t)
			mockProjectsRepo := mockcontract.NewMockProjectsRepository(t)

			// Setup mocks
			tt.setupMocks(mockTeamsRepo)

			// Create service
			mockUserNotificationsUseCase := mockcontract.NewMockUserNotificationsUseCase(t)
			service := New(mockTxManager, mockTeamsRepo, mockUsersRepo, mockUserNotificationsUseCase, mockProjectsRepo)

			// Call method
			members, err := service.GetMembers(context.Background(), tt.teamID)

			// Check result
			if tt.expectedError {
				require.Error(t, err)
				if tt.errorContains != "" {
					require.Contains(t, err.Error(), tt.errorContains)
				}
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.expectedMembers, members)
			}
		})
	}
}

func TestChangeMemberRole(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		setupMocks func(
			mockTeamsRepo *mockcontract.MockTeamsRepository,
			mockUsersRepo *mockcontract.MockUsersRepository,
			mockTxManager *mockdb.MockTxManager,
			mockUserNotificationsUseCase *mockcontract.MockUserNotificationsUseCase,
		)
		teamID        domain.TeamID
		userID        domain.UserID
		newRole       domain.Role
		expectedError bool
		errorContains string
	}{
		{
			name: "Owner promotes member to admin",
			setupMocks: func(
				mockTeamsRepo *mockcontract.MockTeamsRepository,
				mockUsersRepo *mockcontract.MockUsersRepository,
				mockTxManager *mockdb.MockTxManager,
				mockUserNotificationsUseCase *mockcontract.MockUserNotificationsUseCase,
			) {
				// Get current user
				mockUsersRepo.EXPECT().GetByID(
					mock.Anything,
					domain.UserID(1),
				).Return(domain.User{
					ID:          1,
					Username:    "owner",
					IsSuperuser: false,
				}, nil)

				// Get team
				mockTeamsRepo.EXPECT().GetByID(
					mock.Anything,
					domain.TeamID(1),
				).Return(domain.Team{
					ID:   1,
					Name: "Test Team",
					Members: []domain.TeamMember{
						{
							UserID: 1,
							Role:   domain.RoleOwner,
						},
						{
							UserID: 2,
							Role:   domain.RoleMember,
						},
					},
				}, nil)

				// Setup transaction
				mockTxManager.EXPECT().ReadCommitted(
					mock.Anything,
					mock.AnythingOfType("func(context.Context) error"),
				).RunAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
					return fn(ctx)
				})

				// Update member role
				mockTeamsRepo.EXPECT().UpdateMemberRole(
					mock.Anything,
					domain.TeamID(1),
					domain.UserID(2),
					domain.RoleAdmin,
				).Return(nil)

				// Create notification
				mockUserNotificationsUseCase.EXPECT().CreateNotification(
					mock.Anything,
					domain.UserID(2),
					domain.UserNotificationTypeRoleChanged,
					mock.AnythingOfType("domain.UserNotificationContent"),
				).Return(nil)
			},
			teamID:        domain.TeamID(1),
			userID:        domain.UserID(2),
			newRole:       domain.RoleAdmin,
			expectedError: false,
		},
		{
			name: "Member cannot change roles",
			setupMocks: func(
				mockTeamsRepo *mockcontract.MockTeamsRepository,
				mockUsersRepo *mockcontract.MockUsersRepository,
				mockTxManager *mockdb.MockTxManager,
				mockUserNotificationsUseCase *mockcontract.MockUserNotificationsUseCase,
			) {
				// Get current user
				mockUsersRepo.EXPECT().GetByID(
					mock.Anything,
					domain.UserID(3),
				).Return(domain.User{
					ID:          3,
					Username:    "member",
					IsSuperuser: false,
				}, nil)

				// Get team
				mockTeamsRepo.EXPECT().GetByID(
					mock.Anything,
					domain.TeamID(1),
				).Return(domain.Team{
					ID:   1,
					Name: "Test Team",
					Members: []domain.TeamMember{
						{
							UserID: 1,
							Role:   domain.RoleOwner,
						},
						{
							UserID: 2,
							Role:   domain.RoleAdmin,
						},
						{
							UserID: 3,
							Role:   domain.RoleMember,
						},
					},
				}, nil)
			},
			teamID:        domain.TeamID(1),
			userID:        domain.UserID(2),
			newRole:       domain.RoleMember,
			expectedError: true,
			errorContains: domain.ErrForbidden.Error(),
		},
		{
			name: "User not found in team",
			setupMocks: func(
				mockTeamsRepo *mockcontract.MockTeamsRepository,
				mockUsersRepo *mockcontract.MockUsersRepository,
				mockTxManager *mockdb.MockTxManager,
				mockUserNotificationsUseCase *mockcontract.MockUserNotificationsUseCase,
			) {
				// Get current user
				mockUsersRepo.EXPECT().GetByID(
					mock.Anything,
					domain.UserID(1),
				).Return(domain.User{
					ID:          1,
					Username:    "owner",
					IsSuperuser: false,
				}, nil)

				// Get team
				mockTeamsRepo.EXPECT().GetByID(
					mock.Anything,
					domain.TeamID(1),
				).Return(domain.Team{
					ID:   1,
					Name: "Test Team",
					Members: []domain.TeamMember{
						{
							UserID: 1,
							Role:   domain.RoleOwner,
						},
					},
				}, nil)
			},
			teamID:        domain.TeamID(1),
			userID:        domain.UserID(999),
			newRole:       domain.RoleAdmin,
			expectedError: true,
			errorContains: domain.ErrEntityNotFound.Error(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// Create mocks
			mockTxManager := mockdb.NewMockTxManager(t)
			mockTeamsRepo := mockcontract.NewMockTeamsRepository(t)
			mockUsersRepo := mockcontract.NewMockUsersRepository(t)
			mockUserNotificationsUseCase := mockcontract.NewMockUserNotificationsUseCase(t)
			mockProjectsRepo := mockcontract.NewMockProjectsRepository(t)

			// Setup mocks
			tt.setupMocks(mockTeamsRepo, mockUsersRepo, mockTxManager, mockUserNotificationsUseCase)

			// Create service
			service := New(mockTxManager, mockTeamsRepo, mockUsersRepo, mockUserNotificationsUseCase, mockProjectsRepo)

			// Create context with user ID - use the current user ID from the test setup
			var currentUserID domain.UserID
			if tt.name == "Member cannot change roles" {
				currentUserID = domain.UserID(3)
			} else {
				currentUserID = domain.UserID(1)
			}
			ctx := wardencontext.WithUserID(context.Background(), currentUserID)

			// Call method
			err := service.ChangeMemberRole(ctx, tt.teamID, tt.userID, tt.newRole)

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
