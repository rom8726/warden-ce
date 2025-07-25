package projects

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
	mockProjectRepo := mockcontract.NewMockProjectsRepository(t)
	mockIssuesRepo := mockcontract.NewMockIssuesRepository(t)
	mockTeamsUseCase := mockcontract.NewMockTeamsUseCase(t)

	// Create service
	service := New(mockProjectRepo, mockIssuesRepo, mockTeamsUseCase)

	// Verify service was created correctly
	require.NotNil(t, service)
	require.Equal(t, mockProjectRepo, service.projectRepo)
	require.Equal(t, mockIssuesRepo, service.issuesRepository)
	require.Equal(t, mockTeamsUseCase, service.teamsUseCase)
}

func TestGetProject(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name            string
		setupMocks      func(mockProjectRepo *mockcontract.MockProjectsRepository)
		projectID       domain.ProjectID
		expectedProject domain.Project
		expectedError   bool
		errorContains   string
	}{
		{
			name: "Success",
			setupMocks: func(mockProjectRepo *mockcontract.MockProjectsRepository) {
				mockProjectRepo.EXPECT().GetByID(
					mock.Anything,
					domain.ProjectID(1),
				).Return(domain.Project{
					ID:        1,
					Name:      "Test Project",
					PublicKey: "test-public-key",
					CreatedAt: time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
				}, nil)
			},
			projectID: domain.ProjectID(1),
			expectedProject: domain.Project{
				ID:        1,
				Name:      "Test Project",
				PublicKey: "test-public-key",
				CreatedAt: time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
			},
			expectedError: false,
		},
		{
			name: "Project not found",
			setupMocks: func(mockProjectRepo *mockcontract.MockProjectsRepository) {
				mockProjectRepo.EXPECT().GetByID(
					mock.Anything,
					domain.ProjectID(999),
				).Return(domain.Project{}, domain.ErrEntityNotFound)
			},
			projectID:       domain.ProjectID(999),
			expectedProject: domain.Project{},
			expectedError:   true,
			errorContains:   "not found",
		},
		{
			name: "Database error",
			setupMocks: func(mockProjectRepo *mockcontract.MockProjectsRepository) {
				mockProjectRepo.EXPECT().GetByID(
					mock.Anything,
					domain.ProjectID(1),
				).Return(domain.Project{}, errors.New("database error"))
			},
			projectID:       domain.ProjectID(1),
			expectedProject: domain.Project{},
			expectedError:   true,
			errorContains:   "database error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// Create mocks
			mockProjectRepo := mockcontract.NewMockProjectsRepository(t)
			mockIssuesRepo := mockcontract.NewMockIssuesRepository(t)
			mockTeamsUseCase := mockcontract.NewMockTeamsUseCase(t)

			// Setup mocks
			tt.setupMocks(mockProjectRepo)

			// Create service
			service := New(mockProjectRepo, mockIssuesRepo, mockTeamsUseCase)

			// Call method
			project, err := service.GetProject(context.Background(), tt.projectID)

			// Check result
			if tt.expectedError {
				require.Error(t, err)
				if tt.errorContains != "" {
					require.Contains(t, err.Error(), tt.errorContains)
				}
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.expectedProject, project)
			}
		})
	}
}

func TestGetProjectExtended(t *testing.T) {
	t.Parallel()

	teamID := domain.TeamID(1)

	tests := []struct {
		name       string
		setupMocks func(
			mockProjectRepo *mockcontract.MockProjectsRepository,
			mockTeamsUseCase *mockcontract.MockTeamsUseCase,
		)
		projectID               domain.ProjectID
		expectedProjectExtended domain.ProjectExtended
		expectedError           bool
		errorContains           string
	}{
		{
			name: "Success - Project with team",
			setupMocks: func(
				mockProjectRepo *mockcontract.MockProjectsRepository,
				mockTeamsUseCase *mockcontract.MockTeamsUseCase,
			) {
				mockProjectRepo.EXPECT().GetByID(
					mock.Anything,
					domain.ProjectID(1),
				).Return(domain.Project{
					ID:        1,
					Name:      "Test Project",
					PublicKey: "test-public-key",
					TeamID:    &teamID,
					CreatedAt: time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
				}, nil)

				mockTeamsUseCase.EXPECT().GetTeamByID(
					mock.Anything,
					domain.TeamID(1),
				).Return(domain.Team{
					ID:   1,
					Name: "Test Team",
				}, nil)
			},
			projectID: domain.ProjectID(1),
			expectedProjectExtended: domain.ProjectExtended{
				Project: domain.Project{
					ID:        1,
					Name:      "Test Project",
					PublicKey: "test-public-key",
					TeamID:    &teamID,
					CreatedAt: time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
				},
				TeamName: func() *string {
					s := "Test Team"

					return &s
				}(),
			},
			expectedError: false,
		},
		{
			name: "Success - Project without team",
			setupMocks: func(
				mockProjectRepo *mockcontract.MockProjectsRepository,
				mockTeamsUseCase *mockcontract.MockTeamsUseCase,
			) {
				mockProjectRepo.EXPECT().GetByID(
					mock.Anything,
					domain.ProjectID(2),
				).Return(domain.Project{
					ID:        2,
					Name:      "Personal Project",
					PublicKey: "personal-public-key",
					TeamID:    nil,
					CreatedAt: time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
				}, nil)
			},
			projectID: domain.ProjectID(2),
			expectedProjectExtended: domain.ProjectExtended{
				Project: domain.Project{
					ID:        2,
					Name:      "Personal Project",
					PublicKey: "personal-public-key",
					TeamID:    nil,
					CreatedAt: time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
				},
				TeamName: nil,
			},
			expectedError: false,
		},
		{
			name: "Project not found",
			setupMocks: func(
				mockProjectRepo *mockcontract.MockProjectsRepository,
				mockTeamsUseCase *mockcontract.MockTeamsUseCase,
			) {
				mockProjectRepo.EXPECT().GetByID(
					mock.Anything,
					domain.ProjectID(999),
				).Return(domain.Project{}, domain.ErrEntityNotFound)
			},
			projectID:               domain.ProjectID(999),
			expectedProjectExtended: domain.ProjectExtended{},
			expectedError:           true,
			errorContains:           "get project",
		},
		{
			name: "Team not found",
			setupMocks: func(
				mockProjectRepo *mockcontract.MockProjectsRepository,
				mockTeamsUseCase *mockcontract.MockTeamsUseCase,
			) {
				mockProjectRepo.EXPECT().GetByID(
					mock.Anything,
					domain.ProjectID(3),
				).Return(domain.Project{
					ID:        3,
					Name:      "Test Project",
					PublicKey: "test-public-key",
					TeamID:    &teamID,
					CreatedAt: time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
				}, nil)

				mockTeamsUseCase.EXPECT().GetTeamByID(
					mock.Anything,
					domain.TeamID(1),
				).Return(domain.Team{}, domain.ErrEntityNotFound)
			},
			projectID:               domain.ProjectID(3),
			expectedProjectExtended: domain.ProjectExtended{},
			expectedError:           true,
			errorContains:           "get team",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// Create mocks
			mockProjectRepo := mockcontract.NewMockProjectsRepository(t)
			mockIssuesRepo := mockcontract.NewMockIssuesRepository(t)
			mockTeamsUseCase := mockcontract.NewMockTeamsUseCase(t)

			// Setup mocks
			tt.setupMocks(mockProjectRepo, mockTeamsUseCase)

			// Create service
			service := New(mockProjectRepo, mockIssuesRepo, mockTeamsUseCase)

			// Call method
			projectExtended, err := service.GetProjectExtended(context.Background(), tt.projectID)

			// Check result
			if tt.expectedError {
				require.Error(t, err)
				if tt.errorContains != "" {
					require.Contains(t, err.Error(), tt.errorContains)
				}
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.expectedProjectExtended, projectExtended)
			}
		})
	}
}

func TestCreateProject(t *testing.T) {
	t.Parallel()

	teamID := domain.TeamID(1)

	tests := []struct {
		name            string
		setupMocks      func(mockProjectRepo *mockcontract.MockProjectsRepository, mockTeamsUseCase *mockcontract.MockTeamsUseCase)
		projectName     string
		description     string
		teamID          *domain.TeamID
		expectedProject domain.Project
		expectedError   bool
		errorContains   string
	}{
		{
			name: "Success - Project with team",
			setupMocks: func(mockProjectRepo *mockcontract.MockProjectsRepository, mockTeamsUseCase *mockcontract.MockTeamsUseCase) {
				mockTeamsUseCase.EXPECT().GetTeamByID(
					mock.Anything,
					domain.TeamID(1),
				).Return(
					domain.Team{
						ID:   1,
						Name: "Test Team",
					},
					nil,
				)

				mockProjectRepo.EXPECT().Create(
					mock.Anything,
					mock.AnythingOfType("*domain.ProjectDTO"),
				).Run(func(ctx context.Context, projectDTO *domain.ProjectDTO) {
					require.Equal(t, "Test Project", projectDTO.Name)
					require.Equal(t, "Some description", projectDTO.Description)
					require.NotEmpty(t, projectDTO.PublicKey)
					require.Equal(t, &teamID, projectDTO.TeamID)
				}).Return(domain.ProjectID(1), nil)
			},
			projectName: "Test Project",
			description: "Some description",
			teamID:      &teamID,
			expectedProject: domain.Project{
				ID:          1,
				Name:        "Test Project",
				Description: "Some description",
				TeamID:      &teamID,
			},
			expectedError: false,
		},
		{
			name: "Success - Project without team",
			setupMocks: func(mockProjectRepo *mockcontract.MockProjectsRepository, mockTeamsUseCase *mockcontract.MockTeamsUseCase) {
				mockProjectRepo.EXPECT().Create(
					mock.Anything,
					mock.AnythingOfType("*domain.ProjectDTO"),
				).Run(func(ctx context.Context, projectDTO *domain.ProjectDTO) {
					require.Equal(t, "Personal Project", projectDTO.Name)
					require.Equal(t, "Some description", projectDTO.Description)
					require.NotEmpty(t, projectDTO.PublicKey)
					require.Nil(t, projectDTO.TeamID)
				}).Return(domain.ProjectID(2), nil)
			},
			projectName: "Personal Project",
			description: "Some description",
			teamID:      nil,
			expectedProject: domain.Project{
				ID:          2,
				Name:        "Personal Project",
				Description: "Some description",
				TeamID:      nil,
			},
			expectedError: false,
		},
		{
			name: "Error creating project",
			setupMocks: func(mockProjectRepo *mockcontract.MockProjectsRepository, mockTeamsUseCase *mockcontract.MockTeamsUseCase) {
				mockTeamsUseCase.EXPECT().GetTeamByID(
					mock.Anything,
					domain.TeamID(1),
				).Return(
					domain.Team{
						ID:   1,
						Name: "Test Team",
					},
					nil,
				)

				mockProjectRepo.EXPECT().Create(
					mock.Anything,
					mock.AnythingOfType("*domain.ProjectDTO"),
				).Return(domain.ProjectID(0), errors.New("database error"))
			},
			projectName:     "Test Project",
			description:     "Some description",
			teamID:          &teamID,
			expectedProject: domain.Project{},
			expectedError:   true,
			errorContains:   "create project",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// Create mocks
			mockProjectRepo := mockcontract.NewMockProjectsRepository(t)
			mockIssuesRepo := mockcontract.NewMockIssuesRepository(t)
			mockTeamsUseCase := mockcontract.NewMockTeamsUseCase(t)

			// Setup mocks
			tt.setupMocks(mockProjectRepo, mockTeamsUseCase)

			// Create service
			service := New(mockProjectRepo, mockIssuesRepo, mockTeamsUseCase)

			// Call method
			project, err := service.CreateProject(context.Background(), tt.projectName, tt.description, tt.teamID)

			// Check result
			if tt.expectedError {
				require.Error(t, err)
				if tt.errorContains != "" {
					require.Contains(t, err.Error(), tt.errorContains)
				}
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.expectedProject.ID, project.ID)
				require.Equal(t, tt.expectedProject.Name, project.Name)
				require.Equal(t, tt.expectedProject.Description, project.Description)
				require.Equal(t, tt.expectedProject.TeamID, project.TeamID)
				require.NotEmpty(t, project.PublicKey)
				require.WithinDuration(t, time.Now(), project.CreatedAt, 2*time.Second)
			}
		})
	}
}

func TestList(t *testing.T) {
	t.Parallel()

	teamID := domain.TeamID(1)
	teamName := "Test Team"

	tests := []struct {
		name             string
		setupMocks       func(mockProjectRepo *mockcontract.MockProjectsRepository)
		expectedProjects []domain.ProjectExtended
		expectedError    bool
		errorContains    string
	}{
		{
			name: "Success - Multiple projects",
			setupMocks: func(mockProjectRepo *mockcontract.MockProjectsRepository) {
				mockProjectRepo.EXPECT().List(
					mock.Anything,
				).Return([]domain.ProjectExtended{
					{
						Project: domain.Project{
							ID:        1,
							Name:      "Project 1",
							PublicKey: "public-key-1",
							TeamID:    &teamID,
							CreatedAt: time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
						},
						TeamName: &teamName,
					},
					{
						Project: domain.Project{
							ID:        2,
							Name:      "Project 2",
							PublicKey: "public-key-2",
							TeamID:    nil,
							CreatedAt: time.Date(2023, 1, 2, 0, 0, 0, 0, time.UTC),
						},
						TeamName: nil,
					},
				}, nil)
			},
			expectedProjects: []domain.ProjectExtended{
				{
					Project: domain.Project{
						ID:        1,
						Name:      "Project 1",
						PublicKey: "public-key-1",
						TeamID:    &teamID,
						CreatedAt: time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
					},
					TeamName: &teamName,
				},
				{
					Project: domain.Project{
						ID:        2,
						Name:      "Project 2",
						PublicKey: "public-key-2",
						TeamID:    nil,
						CreatedAt: time.Date(2023, 1, 2, 0, 0, 0, 0, time.UTC),
					},
					TeamName: nil,
				},
			},
			expectedError: false,
		},
		{
			name: "Success - No projects",
			setupMocks: func(mockProjectRepo *mockcontract.MockProjectsRepository) {
				mockProjectRepo.EXPECT().List(
					mock.Anything,
				).Return([]domain.ProjectExtended{}, nil)
			},
			expectedProjects: []domain.ProjectExtended{},
			expectedError:    false,
		},
		{
			name: "Error listing projects",
			setupMocks: func(mockProjectRepo *mockcontract.MockProjectsRepository) {
				mockProjectRepo.EXPECT().List(
					mock.Anything,
				).Return(nil, errors.New("database error"))
			},
			expectedProjects: nil,
			expectedError:    true,
			errorContains:    "database error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// Create mocks
			mockProjectRepo := mockcontract.NewMockProjectsRepository(t)
			mockIssuesRepo := mockcontract.NewMockIssuesRepository(t)
			mockTeamsUseCase := mockcontract.NewMockTeamsUseCase(t)

			// Setup mocks
			tt.setupMocks(mockProjectRepo)

			// Create service
			service := New(mockProjectRepo, mockIssuesRepo, mockTeamsUseCase)

			// Call method
			projects, err := service.List(context.Background())

			// Check result
			if tt.expectedError {
				require.Error(t, err)
				if tt.errorContains != "" {
					require.Contains(t, err.Error(), tt.errorContains)
				}
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.expectedProjects, projects)
			}
		})
	}
}

func TestGetProjectsByUserID(t *testing.T) {
	t.Parallel()

	teamID1 := domain.TeamID(1)
	teamID2 := domain.TeamID(2)
	teamName1 := "Team 1"
	teamName2 := "Team 2"

	tests := []struct {
		name       string
		setupMocks func(
			mockProjectRepo *mockcontract.MockProjectsRepository,
			mockTeamsUseCase *mockcontract.MockTeamsUseCase,
		)
		userID           domain.UserID
		isSuperuser      bool
		expectedProjects []domain.ProjectExtended
		expectedError    bool
		errorContains    string
	}{
		{
			name: "Success - Superuser gets all projects",
			setupMocks: func(
				mockProjectRepo *mockcontract.MockProjectsRepository,
				mockTeamsUseCase *mockcontract.MockTeamsUseCase,
			) {
				mockProjectRepo.EXPECT().List(
					mock.Anything,
				).Return([]domain.ProjectExtended{
					{
						Project: domain.Project{
							ID:        1,
							Name:      "Project 1",
							PublicKey: "public-key-1",
							TeamID:    &teamID1,
							CreatedAt: time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
						},
						TeamName: &teamName1,
					},
					{
						Project: domain.Project{
							ID:        2,
							Name:      "Project 2",
							PublicKey: "public-key-2",
							TeamID:    &teamID2,
							CreatedAt: time.Date(2023, 1, 2, 0, 0, 0, 0, time.UTC),
						},
						TeamName: &teamName2,
					},
					{
						Project: domain.Project{
							ID:        3,
							Name:      "Personal Project",
							PublicKey: "public-key-3",
							TeamID:    nil,
							CreatedAt: time.Date(2023, 1, 3, 0, 0, 0, 0, time.UTC),
						},
						TeamName: nil,
					},
				}, nil)
			},
			userID:      domain.UserID(1),
			isSuperuser: true,
			expectedProjects: []domain.ProjectExtended{
				{
					Project: domain.Project{
						ID:        1,
						Name:      "Project 1",
						PublicKey: "public-key-1",
						TeamID:    &teamID1,
						CreatedAt: time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
					},
					TeamName: &teamName1,
				},
				{
					Project: domain.Project{
						ID:        2,
						Name:      "Project 2",
						PublicKey: "public-key-2",
						TeamID:    &teamID2,
						CreatedAt: time.Date(2023, 1, 2, 0, 0, 0, 0, time.UTC),
					},
					TeamName: &teamName2,
				},
				{
					Project: domain.Project{
						ID:        3,
						Name:      "Personal Project",
						PublicKey: "public-key-3",
						TeamID:    nil,
						CreatedAt: time.Date(2023, 1, 3, 0, 0, 0, 0, time.UTC),
					},
					TeamName: nil,
				},
			},
			expectedError: false,
		},
		{
			name: "Success - Regular user gets only their team projects and personal projects",
			setupMocks: func(
				mockProjectRepo *mockcontract.MockProjectsRepository,
				mockTeamsUseCase *mockcontract.MockTeamsUseCase,
			) {
				mockProjectRepo.EXPECT().List(
					mock.Anything,
				).Return([]domain.ProjectExtended{
					{
						Project: domain.Project{
							ID:        1,
							Name:      "Project 1",
							PublicKey: "public-key-1",
							TeamID:    &teamID1,
							CreatedAt: time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
						},
						TeamName: &teamName1,
					},
					{
						Project: domain.Project{
							ID:        2,
							Name:      "Project 2",
							PublicKey: "public-key-2",
							TeamID:    &teamID2,
							CreatedAt: time.Date(2023, 1, 2, 0, 0, 0, 0, time.UTC),
						},
						TeamName: &teamName2,
					},
					{
						Project: domain.Project{
							ID:        3,
							Name:      "Personal Project",
							PublicKey: "public-key-3",
							TeamID:    nil,
							CreatedAt: time.Date(2023, 1, 3, 0, 0, 0, 0, time.UTC),
						},
						TeamName: nil,
					},
				}, nil)

				mockTeamsUseCase.EXPECT().GetTeamsByUserID(
					mock.Anything,
					domain.UserID(1),
				).Return([]domain.Team{
					{
						ID:   1,
						Name: "Team 1",
					},
				}, nil)
			},
			userID:      domain.UserID(1),
			isSuperuser: false,
			expectedProjects: []domain.ProjectExtended{
				{
					Project: domain.Project{
						ID:        3,
						Name:      "Personal Project",
						PublicKey: "public-key-3",
						TeamID:    nil,
						CreatedAt: time.Date(2023, 1, 3, 0, 0, 0, 0, time.UTC),
					},
					TeamName: nil,
				},
				{
					Project: domain.Project{
						ID:        1,
						Name:      "Project 1",
						PublicKey: "public-key-1",
						TeamID:    &teamID1,
						CreatedAt: time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
					},
					TeamName: &teamName1,
				},
			},
			expectedError: false,
		},
		{
			name: "Error listing projects",
			setupMocks: func(
				mockProjectRepo *mockcontract.MockProjectsRepository,
				mockTeamsUseCase *mockcontract.MockTeamsUseCase,
			) {
				mockProjectRepo.EXPECT().List(
					mock.Anything,
				).Return(nil, errors.New("database error"))
			},
			userID:           domain.UserID(1),
			isSuperuser:      true,
			expectedProjects: nil,
			expectedError:    true,
			errorContains:    "list projects",
		},
		{
			name: "Error getting user teams",
			setupMocks: func(
				mockProjectRepo *mockcontract.MockProjectsRepository,
				mockTeamsUseCase *mockcontract.MockTeamsUseCase,
			) {
				mockProjectRepo.EXPECT().List(
					mock.Anything,
				).Return([]domain.ProjectExtended{
					{
						Project: domain.Project{
							ID:        1,
							Name:      "Project 1",
							PublicKey: "public-key-1",
							TeamID:    &teamID1,
							CreatedAt: time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
						},
						TeamName: &teamName1,
					},
				}, nil)

				mockTeamsUseCase.EXPECT().GetTeamsByUserID(
					mock.Anything,
					domain.UserID(1),
				).Return(nil, errors.New("database error"))
			},
			userID:           domain.UserID(1),
			isSuperuser:      false,
			expectedProjects: nil,
			expectedError:    true,
			errorContains:    "get user teams",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// Create mocks
			mockProjectRepo := mockcontract.NewMockProjectsRepository(t)
			mockIssuesRepo := mockcontract.NewMockIssuesRepository(t)
			mockTeamsUseCase := mockcontract.NewMockTeamsUseCase(t)

			// Setup mocks
			tt.setupMocks(mockProjectRepo, mockTeamsUseCase)

			// Create service
			service := New(mockProjectRepo, mockIssuesRepo, mockTeamsUseCase)

			// Call method
			projects, err := service.GetProjectsByUserID(context.Background(), tt.userID, tt.isSuperuser)

			// Check result
			if tt.expectedError {
				require.Error(t, err)
				if tt.errorContains != "" {
					require.Contains(t, err.Error(), tt.errorContains)
				}
			} else {
				require.NoError(t, err)
				// Since the order of projects in the result might be different, we need to check that all expected projects are in the result
				require.Equal(t, len(tt.expectedProjects), len(projects))
				for _, expectedProject := range tt.expectedProjects {
					found := false
					for _, project := range projects {
						if project.ID == expectedProject.ID {
							require.Equal(t, expectedProject, project)
							found = true

							break
						}
					}
					require.True(t, found, "Expected project with ID %d not found in result", expectedProject.ID)
				}
			}
		})
	}
}

func TestUpdateInfo(t *testing.T) {
	t.Parallel()

	teamID := domain.TeamID(1)
	teamName := "Test Team"

	tests := []struct {
		name                    string
		setupMocks              func(mockProjectRepo *mockcontract.MockProjectsRepository, mockTeamsUseCase *mockcontract.MockTeamsUseCase)
		projectID               domain.ProjectID
		newName                 string
		newDescription          string
		expectedProjectExtended domain.ProjectExtended
		expectedError           bool
		errorContains           string
	}{
		{
			name: "Success - Project with team",
			setupMocks: func(mockProjectRepo *mockcontract.MockProjectsRepository, mockTeamsUseCase *mockcontract.MockTeamsUseCase) {
				// Mock getting the project
				mockProjectRepo.EXPECT().GetByID(
					mock.Anything,
					domain.ProjectID(1),
				).Return(domain.Project{
					ID:          1,
					Name:        "Old Name",
					Description: "Old Description",
					PublicKey:   "test-public-key",
					TeamID:      &teamID,
					CreatedAt:   time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
				}, nil)

				// Mock updating the project
				mockProjectRepo.EXPECT().Update(
					mock.Anything,
					domain.ProjectID(1),
					"New Name",
					"New Description",
				).Return(nil)

				// Mock getting the team name
				mockTeamsUseCase.EXPECT().GetTeamByID(
					mock.Anything,
					domain.TeamID(1),
				).Return(domain.Team{
					ID:   1,
					Name: "Test Team",
				}, nil)
			},
			projectID:      domain.ProjectID(1),
			newName:        "New Name",
			newDescription: "New Description",
			expectedProjectExtended: domain.ProjectExtended{
				Project: domain.Project{
					ID:          1,
					Name:        "New Name",
					Description: "New Description",
					PublicKey:   "test-public-key",
					TeamID:      &teamID,
					CreatedAt:   time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
				},
				TeamName: &teamName,
			},
			expectedError: false,
		},
		{
			name: "Success - Project without team",
			setupMocks: func(mockProjectRepo *mockcontract.MockProjectsRepository, mockTeamsUseCase *mockcontract.MockTeamsUseCase) {
				// Mock getting the project
				mockProjectRepo.EXPECT().GetByID(
					mock.Anything,
					domain.ProjectID(2),
				).Return(domain.Project{
					ID:          2,
					Name:        "Old Personal Project",
					Description: "Old Personal Description",
					PublicKey:   "personal-public-key",
					TeamID:      nil,
					CreatedAt:   time.Date(2023, 1, 2, 0, 0, 0, 0, time.UTC),
				}, nil)

				// Mock updating the project
				mockProjectRepo.EXPECT().Update(
					mock.Anything,
					domain.ProjectID(2),
					"New Personal Project",
					"New Personal Description",
				).Return(nil)
			},
			projectID:      domain.ProjectID(2),
			newName:        "New Personal Project",
			newDescription: "New Personal Description",
			expectedProjectExtended: domain.ProjectExtended{
				Project: domain.Project{
					ID:          2,
					Name:        "New Personal Project",
					Description: "New Personal Description",
					PublicKey:   "personal-public-key",
					TeamID:      nil,
					CreatedAt:   time.Date(2023, 1, 2, 0, 0, 0, 0, time.UTC),
				},
				TeamName: nil,
			},
			expectedError: false,
		},
		{
			name: "Error - Project not found",
			setupMocks: func(mockProjectRepo *mockcontract.MockProjectsRepository, mockTeamsUseCase *mockcontract.MockTeamsUseCase) {
				// Mock getting the project
				mockProjectRepo.EXPECT().GetByID(
					mock.Anything,
					domain.ProjectID(999),
				).Return(domain.Project{}, domain.ErrEntityNotFound)
			},
			projectID:               domain.ProjectID(999),
			newName:                 "New Name",
			newDescription:          "New Description",
			expectedProjectExtended: domain.ProjectExtended{},
			expectedError:           true,
			errorContains:           "failed to get project",
		},
		{
			name: "Error - Update fails",
			setupMocks: func(mockProjectRepo *mockcontract.MockProjectsRepository, mockTeamsUseCase *mockcontract.MockTeamsUseCase) {
				// Mock getting the project
				mockProjectRepo.EXPECT().GetByID(
					mock.Anything,
					domain.ProjectID(3),
				).Return(domain.Project{
					ID:          3,
					Name:        "Old Name",
					Description: "Old Description",
					PublicKey:   "test-public-key",
					TeamID:      nil,
					CreatedAt:   time.Date(2023, 1, 3, 0, 0, 0, 0, time.UTC),
				}, nil)

				// Mock updating the project with error
				mockProjectRepo.EXPECT().Update(
					mock.Anything,
					domain.ProjectID(3),
					"New Name",
					"New Description",
				).Return(errors.New("database error"))
			},
			projectID:               domain.ProjectID(3),
			newName:                 "New Name",
			newDescription:          "New Description",
			expectedProjectExtended: domain.ProjectExtended{},
			expectedError:           true,
			errorContains:           "failed to update project",
		},
		{
			name: "Error - Team not found",
			setupMocks: func(mockProjectRepo *mockcontract.MockProjectsRepository, mockTeamsUseCase *mockcontract.MockTeamsUseCase) {
				// Mock getting the project
				mockProjectRepo.EXPECT().GetByID(
					mock.Anything,
					domain.ProjectID(4),
				).Return(domain.Project{
					ID:          4,
					Name:        "Old Name",
					Description: "Old Description",
					PublicKey:   "test-public-key",
					TeamID:      &teamID,
					CreatedAt:   time.Date(2023, 1, 4, 0, 0, 0, 0, time.UTC),
				}, nil)

				// Mock updating the project
				mockProjectRepo.EXPECT().Update(
					mock.Anything,
					domain.ProjectID(4),
					"New Name",
					"New Description",
				).Return(nil)

				// Mock getting the team name with error
				mockTeamsUseCase.EXPECT().GetTeamByID(
					mock.Anything,
					domain.TeamID(1),
				).Return(domain.Team{}, errors.New("team not found"))
			},
			projectID:      domain.ProjectID(4),
			newName:        "New Name",
			newDescription: "New Description",
			expectedProjectExtended: domain.ProjectExtended{
				Project: domain.Project{
					ID:          4,
					Name:        "New Name",
					Description: "New Description",
					PublicKey:   "test-public-key",
					TeamID:      &teamID,
					CreatedAt:   time.Date(2023, 1, 4, 0, 0, 0, 0, time.UTC),
				},
				TeamName: nil,
			},
			expectedError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// Create mocks
			mockProjectRepo := mockcontract.NewMockProjectsRepository(t)
			mockIssuesRepo := mockcontract.NewMockIssuesRepository(t)
			mockTeamsUseCase := mockcontract.NewMockTeamsUseCase(t)

			// Setup mocks
			tt.setupMocks(mockProjectRepo, mockTeamsUseCase)

			// Create service
			service := New(mockProjectRepo, mockIssuesRepo, mockTeamsUseCase)

			// Call method
			projectExtended, err := service.UpdateInfo(context.Background(), tt.projectID, tt.newName, tt.newDescription)

			// Check result
			if tt.expectedError {
				require.Error(t, err)
				if tt.errorContains != "" {
					require.Contains(t, err.Error(), tt.errorContains)
				}
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.expectedProjectExtended, projectExtended)
			}
		})
	}
}

func TestGeneralStats(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		setupMocks    func(mockIssuesRepo *mockcontract.MockIssuesRepository, mockProjectsRepo *mockcontract.MockProjectsRepository)
		projectID     domain.ProjectID
		period        time.Duration
		expectedStats domain.GeneralProjectStats
		expectedError bool
		errorContains string
	}{
		{
			name: "Success",
			setupMocks: func(mockIssuesRepo *mockcontract.MockIssuesRepository, mockProjectsRepo *mockcontract.MockProjectsRepository) {
				mockProjectsRepo.EXPECT().GetByID(
					mock.Anything,
					domain.ProjectID(1),
				).Return(domain.Project{}, nil)

				mockIssuesRepo.EXPECT().CountForLevels(
					mock.Anything,
					domain.ProjectID(1),
					24*time.Hour,
				).Return(map[domain.IssueLevel]uint64{
					domain.IssueLevelError:     10,
					domain.IssueLevelWarning:   5,
					domain.IssueLevelInfo:      3,
					domain.IssueLevelException: 2,
				}, nil)

				mockIssuesRepo.EXPECT().MostFrequent(
					mock.Anything,
					domain.ProjectID(1),
					24*time.Hour,
					uint(6),
				).Return([]domain.IssueExtended{
					{
						Issue: domain.Issue{
							ID:          1,
							ProjectID:   1,
							Title:       "Issue 1",
							Level:       domain.IssueLevelError,
							TotalEvents: 10,
						},
						ProjectName: "Test Project",
					},
					{
						Issue: domain.Issue{
							ID:          2,
							ProjectID:   1,
							Title:       "Issue 2",
							Level:       domain.IssueLevelWarning,
							TotalEvents: 5,
						},
						ProjectName: "Test Project",
					},
				}, nil)
			},
			projectID: domain.ProjectID(1),
			period:    24 * time.Hour,
			expectedStats: domain.GeneralProjectStats{
				TotalIssues:     20,
				ErrorIssues:     10,
				WarningIssues:   5,
				InfoIssues:      3,
				ExceptionIssues: 2,
				MostFrequentIssues: []domain.IssueExtended{
					{
						Issue: domain.Issue{
							ID:          1,
							ProjectID:   1,
							Title:       "Issue 1",
							Level:       domain.IssueLevelError,
							TotalEvents: 10,
						},
						ProjectName: "Test Project",
					},
					{
						Issue: domain.Issue{
							ID:          2,
							ProjectID:   1,
							Title:       "Issue 2",
							Level:       domain.IssueLevelWarning,
							TotalEvents: 5,
						},
						ProjectName: "Test Project",
					},
				},
			},
			expectedError: false,
		},
		{
			name: "Error getting issue counters",
			setupMocks: func(mockIssuesRepo *mockcontract.MockIssuesRepository, mockProjectsRepo *mockcontract.MockProjectsRepository) {
				mockProjectsRepo.EXPECT().GetByID(
					mock.Anything,
					domain.ProjectID(1),
				).Return(domain.Project{}, nil)

				mockIssuesRepo.EXPECT().CountForLevels(
					mock.Anything,
					domain.ProjectID(1),
					24*time.Hour,
				).Return(nil, errors.New("database error"))
			},
			projectID:     domain.ProjectID(1),
			period:        24 * time.Hour,
			expectedStats: domain.GeneralProjectStats{},
			expectedError: true,
			errorContains: "get issue counters",
		},
		{
			name: "Error getting most frequent issues",
			setupMocks: func(mockIssuesRepo *mockcontract.MockIssuesRepository, mockProjectsRepo *mockcontract.MockProjectsRepository) {
				mockProjectsRepo.EXPECT().GetByID(
					mock.Anything,
					domain.ProjectID(1),
				).Return(domain.Project{}, nil)

				mockIssuesRepo.EXPECT().CountForLevels(
					mock.Anything,
					domain.ProjectID(1),
					24*time.Hour,
				).Return(map[domain.IssueLevel]uint64{
					domain.IssueLevelError:     10,
					domain.IssueLevelWarning:   5,
					domain.IssueLevelInfo:      3,
					domain.IssueLevelException: 2,
				}, nil)

				mockIssuesRepo.EXPECT().MostFrequent(
					mock.Anything,
					domain.ProjectID(1),
					24*time.Hour,
					uint(6),
				).Return(nil, errors.New("database error"))
			},
			projectID:     domain.ProjectID(1),
			period:        24 * time.Hour,
			expectedStats: domain.GeneralProjectStats{},
			expectedError: true,
			errorContains: "get most frequent issues",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// Create mocks
			mockProjectRepo := mockcontract.NewMockProjectsRepository(t)
			mockIssuesRepo := mockcontract.NewMockIssuesRepository(t)
			mockTeamsUseCase := mockcontract.NewMockTeamsUseCase(t)

			// Setup mocks
			tt.setupMocks(mockIssuesRepo, mockProjectRepo)

			// Create service
			service := New(mockProjectRepo, mockIssuesRepo, mockTeamsUseCase)

			// Call method
			stats, err := service.GeneralStats(context.Background(), tt.projectID, tt.period)

			// Check result
			if tt.expectedError {
				require.Error(t, err)
				if tt.errorContains != "" {
					require.Contains(t, err.Error(), tt.errorContains)
				}
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.expectedStats, stats)
			}
		})
	}
}
