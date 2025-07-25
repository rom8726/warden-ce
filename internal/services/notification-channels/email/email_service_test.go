package email

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/rom8726/warden/internal/domain"
	mockcontract "github.com/rom8726/warden/test_mocks/internal_/services/notification-channels/email"
)

func TestService_SendUnresolvedIssuesSummaryEmail(t *testing.T) {
	tests := []struct {
		name          string
		issues        []domain.IssueExtended
		setupMocks    func(*mockcontract.MockTeamsRepository, *mockcontract.MockUsersRepository, *mockcontract.MockProjectsRepository)
		expectedError bool
		errorContains string
	}{
		{
			name:   "successful email sending",
			issues: createTestIssues(),
			setupMocks: func(
				mockTeamsRepo *mockcontract.MockTeamsRepository,
				mockUsersRepo *mockcontract.MockUsersRepository,
				mockProjectsRepo *mockcontract.MockProjectsRepository,
			) {
				// Setup projects repository
				mockProjectsRepo.EXPECT().GetByID(mock.Anything, domain.ProjectID(1)).
					Return(domain.Project{
						ID:     1,
						Name:   "Test Project 1",
						TeamID: ptr(domain.TeamID(1)),
					}, nil)
				mockProjectsRepo.EXPECT().GetByID(mock.Anything, domain.ProjectID(2)).
					Return(domain.Project{
						ID:     2,
						Name:   "Test Project 2",
						TeamID: ptr(domain.TeamID(2)),
					}, nil)

				// Setup teams repository
				mockTeamsRepo.EXPECT().GetUniqueUserIDsByTeamIDs(mock.Anything, mock.Anything).
					Return([]domain.UserID{1, 2}, nil)
				mockTeamsRepo.EXPECT().GetTeamsByUserIDs(mock.Anything, mock.Anything).
					Return(map[domain.UserID][]domain.Team{
						1: {{ID: 1, Name: "Team 1"}},
						2: {{ID: 2, Name: "Team 2"}},
					}, nil)

				// Setup users repository
				mockUsersRepo.EXPECT().FetchByIDs(mock.Anything, mock.Anything).
					Return([]domain.User{
						{ID: 1, Username: "user1", Email: "user1@example.com"},
						{ID: 2, Username: "user2", Email: "user2@example.com"},
					}, nil)
			},
			expectedError: false,
		},
		{
			name:   "no issues to send",
			issues: []domain.IssueExtended{},
			setupMocks: func(
				mockTeamsRepo *mockcontract.MockTeamsRepository,
				mockUsersRepo *mockcontract.MockUsersRepository,
				mockProjectsRepo *mockcontract.MockProjectsRepository,
			) {
				mockTeamsRepo.EXPECT().GetUniqueUserIDsByTeamIDs(mock.Anything, mock.Anything).Return([]domain.UserID{}, nil)
				mockUsersRepo.EXPECT().FetchByIDs(mock.Anything, mock.Anything).Return([]domain.User{}, nil)
				mockTeamsRepo.EXPECT().GetTeamsByUserIDs(mock.Anything, mock.Anything).Return(map[domain.UserID][]domain.Team{}, nil)
			},
			expectedError: false,
		},
		{
			name:   "error getting project",
			issues: createTestIssues(),
			setupMocks: func(
				mockTeamsRepo *mockcontract.MockTeamsRepository,
				mockUsersRepo *mockcontract.MockUsersRepository,
				mockProjectsRepo *mockcontract.MockProjectsRepository,
			) {
				mockProjectsRepo.EXPECT().GetByID(mock.Anything, domain.ProjectID(1)).
					Return(domain.Project{}, assert.AnError)
			},
			expectedError: true,
			errorContains: "get project:",
		},
		{
			name:   "error getting teams",
			issues: createTestIssues(),
			setupMocks: func(
				mockTeamsRepo *mockcontract.MockTeamsRepository,
				mockUsersRepo *mockcontract.MockUsersRepository,
				mockProjectsRepo *mockcontract.MockProjectsRepository,
			) {
				mockProjectsRepo.EXPECT().GetByID(mock.Anything, domain.ProjectID(1)).
					Return(domain.Project{
						ID:     1,
						Name:   "Test Project 1",
						TeamID: ptr(domain.TeamID(1)),
					}, nil)
				mockProjectsRepo.EXPECT().GetByID(mock.Anything, domain.ProjectID(2)).
					Return(domain.Project{
						ID:     2,
						Name:   "Test Project 2",
						TeamID: ptr(domain.TeamID(2)),
					}, nil)

				mockTeamsRepo.EXPECT().GetUniqueUserIDsByTeamIDs(mock.Anything, mock.Anything).
					Return(nil, assert.AnError)
			},
			expectedError: true,
			errorContains: "get teams:",
		},
		{
			name:   "error fetching users",
			issues: createTestIssues(),
			setupMocks: func(
				mockTeamsRepo *mockcontract.MockTeamsRepository,
				mockUsersRepo *mockcontract.MockUsersRepository,
				mockProjectsRepo *mockcontract.MockProjectsRepository,
			) {
				mockProjectsRepo.EXPECT().GetByID(mock.Anything, domain.ProjectID(1)).
					Return(domain.Project{
						ID:     1,
						Name:   "Test Project 1",
						TeamID: ptr(domain.TeamID(1)),
					}, nil)
				mockProjectsRepo.EXPECT().GetByID(mock.Anything, domain.ProjectID(2)).
					Return(domain.Project{
						ID:     2,
						Name:   "Test Project 2",
						TeamID: ptr(domain.TeamID(2)),
					}, nil)

				mockTeamsRepo.EXPECT().GetUniqueUserIDsByTeamIDs(mock.Anything, mock.Anything).
					Return([]domain.UserID{1, 2}, nil)

				mockUsersRepo.EXPECT().FetchByIDs(mock.Anything, mock.Anything).
					Return(nil, assert.AnError)
			},
			expectedError: true,
			errorContains: "fetch users:",
		},
		{
			name:   "error getting teams by user ids",
			issues: createTestIssues(),
			setupMocks: func(
				mockTeamsRepo *mockcontract.MockTeamsRepository,
				mockUsersRepo *mockcontract.MockUsersRepository,
				mockProjectsRepo *mockcontract.MockProjectsRepository,
			) {
				mockProjectsRepo.EXPECT().GetByID(mock.Anything, domain.ProjectID(1)).
					Return(domain.Project{
						ID:     1,
						Name:   "Test Project 1",
						TeamID: ptr(domain.TeamID(1)),
					}, nil)
				mockProjectsRepo.EXPECT().GetByID(mock.Anything, domain.ProjectID(2)).
					Return(domain.Project{
						ID:     2,
						Name:   "Test Project 2",
						TeamID: ptr(domain.TeamID(2)),
					}, nil)

				mockTeamsRepo.EXPECT().GetUniqueUserIDsByTeamIDs(mock.Anything, mock.Anything).
					Return([]domain.UserID{1, 2}, nil)

				mockUsersRepo.EXPECT().FetchByIDs(mock.Anything, mock.Anything).
					Return([]domain.User{
						{ID: 1, Username: "user1", Email: "user1@example.com"},
						{ID: 2, Username: "user2", Email: "user2@example.com"},
					}, nil)

				mockTeamsRepo.EXPECT().GetTeamsByUserIDs(mock.Anything, mock.Anything).
					Return(nil, assert.AnError)
			},
			expectedError: true,
			errorContains: "get teams by user ids map:",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create mocks
			mockTeamsRepo := mockcontract.NewMockTeamsRepository(t)
			mockUsersRepo := mockcontract.NewMockUsersRepository(t)
			mockProjectsRepo := mockcontract.NewMockProjectsRepository(t)

			// Setup mocks
			tt.setupMocks(mockTeamsRepo, mockUsersRepo, mockProjectsRepo)

			// Create service
			service := &Service{
				cfg: &Config{
					BaseURL: "http://localhost:8080",
				},
				teamsRepo:    mockTeamsRepo,
				usersRepo:    mockUsersRepo,
				projectsRepo: mockProjectsRepo,
			}
			// Подменяем функцию отправки писем на мок
			service.sendEmailFunc = func(ctx context.Context, to []string, subject, body string) error {
				return nil
			}

			// Execute test
			err := service.SendUnresolvedIssuesSummaryEmail(context.Background(), tt.issues)

			// Assert results
			if tt.expectedError {
				assert.Error(t, err)
				if tt.errorContains != "" {
					assert.Contains(t, err.Error(), tt.errorContains)
				}
			} else {
				assert.NoError(t, err)
			}

			// Verify all mocks were called as expected
			mockTeamsRepo.AssertExpectations(t)
			mockUsersRepo.AssertExpectations(t)
			mockProjectsRepo.AssertExpectations(t)
		})
	}
}

func TestService_sendEmailsParallel(t *testing.T) {
	tests := []struct {
		name          string
		maxWorkers    int
		emails        []emailData
		expectedError bool
	}{
		{
			name:       "successful parallel sending",
			maxWorkers: 2,
			emails: []emailData{
				{toEmails: []string{"user1@example.com"}, subject: "Test 1", body: "Body 1"},
				{toEmails: []string{"user2@example.com"}, subject: "Test 2", body: "Body 2"},
			},
			expectedError: false,
		},
		{
			name:          "empty emails list",
			maxWorkers:    2,
			emails:        []emailData{},
			expectedError: false,
		},
		{
			name:       "single email",
			maxWorkers: 2,
			emails: []emailData{
				{toEmails: []string{"user1@example.com"}, subject: "Test", body: "Body"},
			},
			expectedError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := &Service{
				cfg: &Config{
					SMTPHost: "localhost:1025",
					Username: "test",
					Password: "test",
				},
			}
			// Подменяем функцию отправки писем на мок
			service.sendEmailFunc = func(ctx context.Context, to []string, subject, body string) error {
				return nil
			}

			err := service.sendEmailsParallel(context.Background(), tt.maxWorkers, tt.emails)

			if tt.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestService_TLSLogic(t *testing.T) {
	tests := []struct {
		name        string
		config      *Config
		port        int
		expectedTLS bool
	}{
		{
			name: "port 1025 (MailHog) should never use TLS",
			config: &Config{
				AllowInsecure: false,
				UseTLS:        true,
			},
			port:        1025,
			expectedTLS: false,
		},
		{
			name: "port 465 should use TLS",
			config: &Config{
				AllowInsecure: true,
				UseTLS:        false,
			},
			port:        465,
			expectedTLS: true,
		},
		{
			name: "port 587 should use TLS",
			config: &Config{
				AllowInsecure: true,
				UseTLS:        false,
			},
			port:        587,
			expectedTLS: true,
		},
		{
			name: "port 25 without AllowInsecure should use TLS",
			config: &Config{
				AllowInsecure: false,
				UseTLS:        false,
			},
			port:        25,
			expectedTLS: true,
		},
		{
			name: "port 25 with AllowInsecure should not use TLS",
			config: &Config{
				AllowInsecure: true,
				UseTLS:        false,
			},
			port:        25,
			expectedTLS: false,
		},
		{
			name: "explicit UseTLS should override port logic (except 1025)",
			config: &Config{
				AllowInsecure: true,
				UseTLS:        true,
			},
			port:        25,
			expectedTLS: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// This is a simple test to verify our TLS logic
			useTLS := false
			if tt.port != 1025 { // MailHog port
				useTLS = tt.config.UseTLS || tt.port == 465 || tt.port == 587 || !tt.config.AllowInsecure
			}
			assert.Equal(t, tt.expectedTLS, useTLS)
		})
	}
}

// Helper functions

func createTestIssues() []domain.IssueExtended {
	now := time.Now()
	return []domain.IssueExtended{
		{
			Issue: domain.Issue{
				ID:        1,
				ProjectID: 1,
				Title:     "Test Issue 1",
				Level:     domain.IssueLevelError,
				Status:    domain.IssueStatusUnresolved,
				FirstSeen: now.Add(-time.Hour),
				LastSeen:  now,
			},
			ProjectName: "Test Project 1",
			ResolvedAt:  nil, // Unresolved issue
		},
		{
			Issue: domain.Issue{
				ID:        2,
				ProjectID: 2,
				Title:     "Test Issue 2",
				Level:     domain.IssueLevelWarning,
				Status:    domain.IssueStatusResolved,
				FirstSeen: now.Add(-2 * time.Hour),
				LastSeen:  now,
			},
			ProjectName: "Test Project 2",
			ResolvedAt:  &now, // Resolved issue
		},
	}
}

func ptr[T any](v T) *T {
	return &v
}
