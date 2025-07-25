package users

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	wardencontext "github.com/rom8726/warden/internal/context"
	"github.com/rom8726/warden/internal/domain"
	mockcontract "github.com/rom8726/warden/test_mocks/internal_/backend/contract"
	mockusers "github.com/rom8726/warden/test_mocks/internal_/backend/usecases/users"
)

func TestNew(t *testing.T) {
	t.Parallel()

	// Create mocks
	mockUsersRepo := mockcontract.NewMockUsersRepository(t)
	mockTeamsRepo := mockcontract.NewMockTeamsRepository(t)
	mockTokenizer := mockcontract.NewMockTokenizer(t)
	mockEmailer := mockcontract.NewMockEmailer(t)
	mockAuthProvider := mockusers.NewMockAuthProvider(t)
	mockRateLimiter := mockcontract.NewMockTwoFARateLimiter(t)

	// Create service
	service := New(
		mockUsersRepo,
		mockTeamsRepo,
		mockTokenizer,
		mockEmailer,
		mockRateLimiter,
		[]AuthProvider{mockAuthProvider},
	)

	// Verify service was created correctly
	require.NotNil(t, service)
	require.Equal(t, mockUsersRepo, service.usersRepo)
	require.Equal(t, mockTeamsRepo, service.teamsRepo)
	require.Equal(t, mockTokenizer, service.tokenizer)
	require.Equal(t, mockEmailer, service.emailer)

	// Verify that the authProvider is set
	require.NotNil(t, service.authProvider)
}

func TestLogin(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name       string
		setupMocks func(
			mockAuthProvider *mockusers.MockAuthProvider,
			mockTokenizer *mockcontract.MockTokenizer,
			mockUsersRepo *mockcontract.MockUsersRepository,
		)
		username             string
		password             string
		expectedAccessToken  string
		expectedRefreshToken string
		expectedTmpPasswd    bool
		expectedError        bool
		errorContains        string
	}{
		{
			name: "Success",
			setupMocks: func(
				mockAuthProvider *mockusers.MockAuthProvider,
				mockTokenizer *mockcontract.MockTokenizer,
				mockUsersRepo *mockcontract.MockUsersRepository,
			) {
				user := &domain.User{
					ID:            1,
					Username:      "user1",
					Email:         "user1@example.com",
					IsActive:      true,
					IsTmpPassword: false,
				}
				mockAuthProvider.EXPECT().CanHandle("user1").Return(true)
				mockAuthProvider.EXPECT().Authenticate(
					mock.Anything,
					"user1",
					"password1",
				).Return(user, nil)
				mockTokenizer.EXPECT().AccessToken(user).Return("access_token_1", nil)
				mockTokenizer.EXPECT().RefreshToken(user).Return("refresh_token_1", nil)
				mockUsersRepo.EXPECT().UpdateLastLogin(mock.Anything, domain.UserID(1)).Return(nil)
			},
			username:             "user1",
			password:             "password1",
			expectedAccessToken:  "access_token_1",
			expectedRefreshToken: "refresh_token_1",
			expectedTmpPasswd:    false,
			expectedError:        false,
		},
		{
			name: "Success with temporary password",
			setupMocks: func(
				mockAuthProvider *mockusers.MockAuthProvider,
				mockTokenizer *mockcontract.MockTokenizer,
				mockUsersRepo *mockcontract.MockUsersRepository,
			) {
				user := &domain.User{
					ID:            2,
					Username:      "user2",
					Email:         "user2@example.com",
					IsActive:      true,
					IsTmpPassword: true,
				}
				mockAuthProvider.EXPECT().CanHandle("user2").Return(true)
				mockAuthProvider.EXPECT().Authenticate(
					mock.Anything,
					"user2",
					"password2",
				).Return(user, nil)
				mockTokenizer.EXPECT().AccessToken(user).Return("access_token_2", nil)
				mockTokenizer.EXPECT().RefreshToken(user).Return("refresh_token_2", nil)
				mockUsersRepo.EXPECT().UpdateLastLogin(mock.Anything, domain.UserID(2)).Return(nil)
			},
			username:             "user2",
			password:             "password2",
			expectedAccessToken:  "access_token_2",
			expectedRefreshToken: "refresh_token_2",
			expectedTmpPasswd:    true,
			expectedError:        false,
		},
		{
			name: "Authentication failed",
			setupMocks: func(
				mockAuthProvider *mockusers.MockAuthProvider,
				mockTokenizer *mockcontract.MockTokenizer,
				mockUsersRepo *mockcontract.MockUsersRepository,
			) {
				// Mock the local auth provider behavior
				mockUsersRepo.EXPECT().GetByUsername(
					mock.Anything,
					"user3",
				).Return(domain.User{}, domain.ErrEntityNotFound)
				mockUsersRepo.EXPECT().GetByEmail(
					mock.Anything,
					"user3",
				).Return(domain.User{}, domain.ErrEntityNotFound)

				// This will be called first, but since we're returning false,
				// the local auth provider will be used instead
				mockAuthProvider.EXPECT().CanHandle("user3").Return(false)
			},
			username:             "user3",
			password:             "wrong_password",
			expectedAccessToken:  "",
			expectedRefreshToken: "",
			expectedTmpPasswd:    false,
			expectedError:        true,
			errorContains:        "authentication failed",
		},
		{
			name: "Inactive user",
			setupMocks: func(
				mockAuthProvider *mockusers.MockAuthProvider,
				mockTokenizer *mockcontract.MockTokenizer,
				mockUsersRepo *mockcontract.MockUsersRepository,
			) {
				user := &domain.User{
					ID:           4,
					Username:     "user4",
					Email:        "user4@example.com",
					IsActive:     false,
					TwoFAEnabled: false,
				}
				mockAuthProvider.EXPECT().CanHandle("user4").Return(true)
				mockAuthProvider.EXPECT().Authenticate(
					mock.Anything,
					"user4",
					"password4",
				).Return(user, nil)
			},
			username:             "user4",
			password:             "password4",
			expectedAccessToken:  "",
			expectedRefreshToken: "",
			expectedTmpPasswd:    false,
			expectedError:        true,
			errorContains:        "inactive user",
		},
		{
			name: "2FA required",
			setupMocks: func(
				mockAuthProvider *mockusers.MockAuthProvider,
				mockTokenizer *mockcontract.MockTokenizer,
				mockUsersRepo *mockcontract.MockUsersRepository,
			) {
				user := &domain.User{
					ID:           5,
					Username:     "user5",
					Email:        "user5@example.com",
					IsActive:     true,
					TwoFAEnabled: true,
				}
				mockAuthProvider.EXPECT().CanHandle("user5").Return(true)
				mockAuthProvider.EXPECT().Authenticate(
					mock.Anything,
					"user5",
					"password5",
				).Return(user, nil)
			},
			username:             "user5",
			password:             "password5",
			expectedAccessToken:  "",
			expectedRefreshToken: "",
			expectedTmpPasswd:    false,
			expectedError:        true,
			errorContains:        "2FA required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// Create mocks
			mockAuthProvider := mockusers.NewMockAuthProvider(t)
			mockTokenizer := mockcontract.NewMockTokenizer(t)
			mockUsersRepo := mockcontract.NewMockUsersRepository(t)
			mockTeamsRepo := mockcontract.NewMockTeamsRepository(t)
			mockEmailer := mockcontract.NewMockEmailer(t)
			mockRateLimiter := mockcontract.NewMockTwoFARateLimiter(t)

			// Setup mocks
			tt.setupMocks(mockAuthProvider, mockTokenizer, mockUsersRepo)

			// Create service
			service := New(
				mockUsersRepo,
				mockTeamsRepo,
				mockTokenizer,
				mockEmailer,
				mockRateLimiter,
				[]AuthProvider{mockAuthProvider},
			)

			// Call method
			accessToken, refreshToken, _, isTmpPasswd, err := service.Login(
				context.Background(),
				tt.username,
				tt.password,
			)

			// Check result
			if tt.expectedError {
				require.Error(t, err)
				if tt.errorContains != "" {
					require.Contains(t, err.Error(), tt.errorContains)
				}
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.expectedAccessToken, accessToken)
				require.Equal(t, tt.expectedRefreshToken, refreshToken)
				require.Equal(t, tt.expectedTmpPasswd, isTmpPasswd)
			}
		})
	}
}

func TestLoginReissue(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name       string
		setupMocks func(
			mockTokenizer *mockcontract.MockTokenizer,
			mockUsersRepo *mockcontract.MockUsersRepository,
		)
		refreshToken         string
		expectedAccessToken  string
		expectedRefreshToken string
		expectedError        bool
		errorContains        string
	}{
		{
			name: "Success",
			setupMocks: func(
				mockTokenizer *mockcontract.MockTokenizer,
				mockUsersRepo *mockcontract.MockUsersRepository,
			) {
				claims := &domain.TokenClaims{
					UserID: 1,
				}
				mockTokenizer.EXPECT().VerifyToken(
					"valid_refresh_token",
					domain.TokenTypeRefresh,
				).Return(claims, nil)
				user := domain.User{
					ID:       1,
					Username: "user1",
					Email:    "user1@example.com",
					IsActive: true,
				}
				mockUsersRepo.EXPECT().GetByID(
					mock.Anything,
					domain.UserID(1),
				).Return(user, nil)
				mockTokenizer.EXPECT().AccessToken(&user).Return("new_access_token", nil)
				mockTokenizer.EXPECT().RefreshToken(&user).Return("new_refresh_token", nil)
				mockUsersRepo.EXPECT().UpdateLastLogin(mock.Anything, domain.UserID(1)).Return(nil)
			},
			refreshToken:         "valid_refresh_token",
			expectedAccessToken:  "new_access_token",
			expectedRefreshToken: "new_refresh_token",
			expectedError:        false,
		},
		{
			name: "Invalid refresh token",
			setupMocks: func(
				mockTokenizer *mockcontract.MockTokenizer,
				mockUsersRepo *mockcontract.MockUsersRepository,
			) {
				mockTokenizer.EXPECT().VerifyToken(
					"invalid_refresh_token",
					domain.TokenTypeRefresh,
				).Return(nil, errors.New("invalid token"))
			},
			refreshToken:         "invalid_refresh_token",
			expectedAccessToken:  "",
			expectedRefreshToken: "",
			expectedError:        true,
			errorContains:        "verify refresh token",
		},
		{
			name: "User not found",
			setupMocks: func(
				mockTokenizer *mockcontract.MockTokenizer,
				mockUsersRepo *mockcontract.MockUsersRepository,
			) {
				claims := &domain.TokenClaims{
					UserID: 2,
				}
				mockTokenizer.EXPECT().VerifyToken(
					"valid_refresh_token_user_not_found",
					domain.TokenTypeRefresh,
				).Return(claims, nil)
				mockUsersRepo.EXPECT().GetByID(
					mock.Anything,
					domain.UserID(2),
				).Return(domain.User{}, domain.ErrEntityNotFound)
			},
			refreshToken:         "valid_refresh_token_user_not_found",
			expectedAccessToken:  "",
			expectedRefreshToken: "",
			expectedError:        true,
			errorContains:        "get user by uuid",
		},
		{
			name: "Inactive user",
			setupMocks: func(
				mockTokenizer *mockcontract.MockTokenizer,
				mockUsersRepo *mockcontract.MockUsersRepository,
			) {
				claims := &domain.TokenClaims{
					UserID: 3,
				}
				mockTokenizer.EXPECT().VerifyToken(
					"valid_refresh_token_inactive_user",
					domain.TokenTypeRefresh,
				).Return(claims, nil)
				user := domain.User{
					ID:       3,
					Username: "user3",
					Email:    "user3@example.com",
					IsActive: false,
				}
				mockUsersRepo.EXPECT().GetByID(
					mock.Anything,
					domain.UserID(3),
				).Return(user, nil)
			},
			refreshToken:         "valid_refresh_token_inactive_user",
			expectedAccessToken:  "",
			expectedRefreshToken: "",
			expectedError:        true,
			errorContains:        "inactive user",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// Create mocks
			mockTokenizer := mockcontract.NewMockTokenizer(t)
			mockUsersRepo := mockcontract.NewMockUsersRepository(t)
			mockTeamsRepo := mockcontract.NewMockTeamsRepository(t)
			mockEmailer := mockcontract.NewMockEmailer(t)
			mockAuthProvider := mockusers.NewMockAuthProvider(t)
			mockRateLimiter := mockcontract.NewMockTwoFARateLimiter(t)

			// Setup mocks
			tt.setupMocks(mockTokenizer, mockUsersRepo)

			// Create service
			service := New(
				mockUsersRepo,
				mockTeamsRepo,
				mockTokenizer,
				mockEmailer,
				mockRateLimiter,
				[]AuthProvider{mockAuthProvider},
			)

			// Call method
			accessToken, refreshToken, err := service.LoginReissue(
				context.Background(),
				tt.refreshToken,
			)

			// Check result
			if tt.expectedError {
				require.Error(t, err)
				if tt.errorContains != "" {
					require.Contains(t, err.Error(), tt.errorContains)
				}
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.expectedAccessToken, accessToken)
				require.Equal(t, tt.expectedRefreshToken, refreshToken)
			}
		})
	}
}

func TestGetByID(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name          string
		setupMocks    func(mockUsersRepo *mockcontract.MockUsersRepository)
		userID        domain.UserID
		expectedUser  domain.User
		expectedError bool
		errorContains string
	}{
		{
			name: "Success",
			setupMocks: func(mockUsersRepo *mockcontract.MockUsersRepository) {
				mockUsersRepo.EXPECT().GetByID(
					mock.Anything,
					domain.UserID(1),
				).Return(domain.User{
					ID:       1,
					Username: "user1",
					IsActive: true,
				}, nil)
			},
			userID: domain.UserID(1),
			expectedUser: domain.User{
				ID:       1,
				Username: "user1",
				IsActive: true,
			},
			expectedError: false,
		},
		{
			name: "User not found",
			setupMocks: func(mockUsersRepo *mockcontract.MockUsersRepository) {
				mockUsersRepo.EXPECT().GetByID(
					mock.Anything,
					domain.UserID(2),
				).Return(domain.User{}, domain.ErrEntityNotFound)
			},
			userID:        domain.UserID(2),
			expectedUser:  domain.User{},
			expectedError: true,
			errorContains: "not found",
		},
		{
			name: "Database error",
			setupMocks: func(mockUsersRepo *mockcontract.MockUsersRepository) {
				mockUsersRepo.EXPECT().GetByID(
					mock.Anything,
					domain.UserID(3),
				).Return(domain.User{}, errors.New("database error"))
			},
			userID:        domain.UserID(3),
			expectedUser:  domain.User{},
			expectedError: true,
			errorContains: "database error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// Create mocks
			mockUsersRepo := mockcontract.NewMockUsersRepository(t)
			mockTeamsRepo := mockcontract.NewMockTeamsRepository(t)
			mockTokenizer := mockcontract.NewMockTokenizer(t)
			mockEmailer := mockcontract.NewMockEmailer(t)
			mockAuthProvider := mockusers.NewMockAuthProvider(t)
			mockRateLimiter := mockcontract.NewMockTwoFARateLimiter(t)

			// Setup mocks
			tt.setupMocks(mockUsersRepo)

			// Create service
			service := New(
				mockUsersRepo,
				mockTeamsRepo,
				mockTokenizer,
				mockEmailer,
				mockRateLimiter,
				[]AuthProvider{mockAuthProvider},
			)

			// Call method
			user, err := service.GetByID(context.Background(), tt.userID)

			// Check result
			if tt.expectedError {
				require.Error(t, err)
				if tt.errorContains != "" {
					require.Contains(t, err.Error(), tt.errorContains)
				}
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.expectedUser, user)
			}
		})
	}
}

func TestCurrentUserInfo(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name             string
		setupMocks       func(mockUsersRepo *mockcontract.MockUsersRepository, mockTeamsRepo *mockcontract.MockTeamsRepository)
		userID           domain.UserID
		expectedUserInfo domain.UserInfo
		expectedError    bool
		errorContains    string
	}{
		{
			name: "Success",
			setupMocks: func(mockUsersRepo *mockcontract.MockUsersRepository, mockTeamsRepo *mockcontract.MockTeamsRepository) {
				mockUsersRepo.EXPECT().GetByID(
					mock.Anything,
					domain.UserID(1),
				).Return(domain.User{
					ID:       1,
					Username: "user1",
					IsActive: true,
				}, nil)
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
								Role:   domain.RoleOwner,
							},
						},
					},
					{
						ID:   2,
						Name: "Team 2",
						Members: []domain.TeamMember{
							{
								UserID: 1,
								Role:   domain.RoleMember,
							},
						},
					},
				}, nil)
			},
			userID: domain.UserID(1),
			expectedUserInfo: domain.UserInfo{
				User: domain.User{
					ID:       1,
					Username: "user1",
					IsActive: true,
				},
				Teams: []domain.UserTeamInfo{
					{
						ID:   1,
						Name: "Team 1",
						Role: domain.RoleOwner,
					},
					{
						ID:   2,
						Name: "Team 2",
						Role: domain.RoleMember,
					},
				},
			},
			expectedError: false,
		},
		{
			name: "User not found",
			setupMocks: func(mockUsersRepo *mockcontract.MockUsersRepository, mockTeamsRepo *mockcontract.MockTeamsRepository) {
				mockUsersRepo.EXPECT().GetByID(
					mock.Anything,
					domain.UserID(2),
				).Return(domain.User{}, domain.ErrEntityNotFound)
			},
			userID:           domain.UserID(2),
			expectedUserInfo: domain.UserInfo{},
			expectedError:    true,
			errorContains:    "get user by id",
		},
		{
			name: "Error getting teams",
			setupMocks: func(mockUsersRepo *mockcontract.MockUsersRepository, mockTeamsRepo *mockcontract.MockTeamsRepository) {
				mockUsersRepo.EXPECT().GetByID(
					mock.Anything,
					domain.UserID(3),
				).Return(domain.User{
					ID:       3,
					Username: "user3",
					IsActive: true,
				}, nil)
				mockTeamsRepo.EXPECT().GetTeamsByUserID(
					mock.Anything,
					domain.UserID(3),
				).Return(nil, errors.New("database error"))
			},
			userID:           domain.UserID(3),
			expectedUserInfo: domain.UserInfo{},
			expectedError:    true,
			errorContains:    "get teams by user id",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// Create mocks
			mockUsersRepo := mockcontract.NewMockUsersRepository(t)
			mockTeamsRepo := mockcontract.NewMockTeamsRepository(t)
			mockTokenizer := mockcontract.NewMockTokenizer(t)
			mockEmailer := mockcontract.NewMockEmailer(t)
			mockAuthProvider := mockusers.NewMockAuthProvider(t)
			mockRateLimiter := mockcontract.NewMockTwoFARateLimiter(t)

			// Setup mocks
			tt.setupMocks(mockUsersRepo, mockTeamsRepo)

			// Create service
			service := New(
				mockUsersRepo,
				mockTeamsRepo,
				mockTokenizer,
				mockEmailer,
				mockRateLimiter,
				[]AuthProvider{mockAuthProvider},
			)

			// Call method
			userInfo, err := service.CurrentUserInfo(context.Background(), tt.userID)

			// Check result
			if tt.expectedError {
				require.Error(t, err)
				if tt.errorContains != "" {
					require.Contains(t, err.Error(), tt.errorContains)
				}
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.expectedUserInfo, userInfo)
			}
		})
	}
}

func TestCreate(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name          string
		setupMocks    func(mockUsersRepo *mockcontract.MockUsersRepository)
		currentUser   domain.User
		username      string
		email         string
		password      string
		isSuperuser   bool
		expectedUser  domain.User
		expectedError bool
		errorContains string
	}{
		{
			name: "Success",
			setupMocks: func(mockUsersRepo *mockcontract.MockUsersRepository) {
				mockUsersRepo.EXPECT().GetByUsername(
					mock.Anything,
					"newuser",
				).Return(domain.User{}, domain.ErrEntityNotFound)
				mockUsersRepo.EXPECT().GetByEmail(
					mock.Anything,
					"newuser@example.com",
				).Return(domain.User{}, domain.ErrEntityNotFound)
				mockUsersRepo.EXPECT().Create(
					mock.Anything,
					mock.AnythingOfType("domain.UserDTO"),
				).Return(domain.User{
					ID:            1,
					Username:      "newuser",
					Email:         "newuser@example.com",
					IsActive:      true,
					IsSuperuser:   false,
					IsTmpPassword: true,
				}, nil)
			},
			currentUser: domain.User{
				ID:          999,
				Username:    "admin",
				IsSuperuser: true,
			},
			username:    "newuser",
			email:       "newuser@example.com",
			password:    "password123",
			isSuperuser: false,
			expectedUser: domain.User{
				ID:            1,
				Username:      "newuser",
				Email:         "newuser@example.com",
				IsActive:      true,
				IsSuperuser:   false,
				IsTmpPassword: true,
			},
			expectedError: false,
		},
		{
			name: "Not a superuser",
			setupMocks: func(mockUsersRepo *mockcontract.MockUsersRepository) {
				// No mocks needed
			},
			currentUser: domain.User{
				ID:          999,
				Username:    "regular",
				IsSuperuser: false,
			},
			username:      "newuser",
			email:         "newuser@example.com",
			password:      "password123",
			isSuperuser:   false,
			expectedUser:  domain.User{},
			expectedError: true,
			errorContains: "forbidden",
		},
		{
			name: "Username already in use",
			setupMocks: func(mockUsersRepo *mockcontract.MockUsersRepository) {
				mockUsersRepo.EXPECT().GetByUsername(
					mock.Anything,
					"existinguser",
				).Return(domain.User{
					ID:       2,
					Username: "existinguser",
				}, nil)
			},
			currentUser: domain.User{
				ID:          999,
				Username:    "admin",
				IsSuperuser: true,
			},
			username:      "existinguser",
			email:         "newuser@example.com",
			password:      "password123",
			isSuperuser:   false,
			expectedUser:  domain.User{},
			expectedError: true,
			errorContains: "username already in use",
		},
		{
			name: "Email already in use",
			setupMocks: func(mockUsersRepo *mockcontract.MockUsersRepository) {
				mockUsersRepo.EXPECT().GetByUsername(
					mock.Anything,
					"newuser",
				).Return(domain.User{}, domain.ErrEntityNotFound)
				mockUsersRepo.EXPECT().GetByEmail(
					mock.Anything,
					"existing@example.com",
				).Return(domain.User{
					ID:    3,
					Email: "existing@example.com",
				}, nil)
			},
			currentUser: domain.User{
				ID:          999,
				Username:    "admin",
				IsSuperuser: true,
			},
			username:      "newuser",
			email:         "existing@example.com",
			password:      "password123",
			isSuperuser:   false,
			expectedUser:  domain.User{},
			expectedError: true,
			errorContains: "email already in use",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// Create mocks
			mockUsersRepo := mockcontract.NewMockUsersRepository(t)
			mockTeamsRepo := mockcontract.NewMockTeamsRepository(t)
			mockTokenizer := mockcontract.NewMockTokenizer(t)
			mockEmailer := mockcontract.NewMockEmailer(t)
			mockAuthProvider := mockusers.NewMockAuthProvider(t)
			mockRateLimiter := mockcontract.NewMockTwoFARateLimiter(t)

			// Setup mocks
			tt.setupMocks(mockUsersRepo)

			// Create service
			service := New(
				mockUsersRepo,
				mockTeamsRepo,
				mockTokenizer,
				mockEmailer,
				mockRateLimiter,
				[]AuthProvider{mockAuthProvider},
			)

			// Call method
			user, err := service.Create(
				context.Background(),
				tt.currentUser,
				tt.username,
				tt.email,
				tt.password,
				tt.isSuperuser,
			)

			// Check result
			if tt.expectedError {
				require.Error(t, err)
				if tt.errorContains != "" {
					require.Contains(t, err.Error(), tt.errorContains)
				}
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.expectedUser, user)
			}
		})
	}
}

func TestListForTeamAdmin(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		setupMocks func(
			mockUsersRepo *mockcontract.MockUsersRepository,
			mockTeamsRepo *mockcontract.MockTeamsRepository,
		)
		isSuperUser   bool
		teamID        domain.TeamID
		expectedUsers []domain.User
		expectedError bool
		errorContains string
	}{
		{
			name: "Superuser can list all users",
			setupMocks: func(
				mockUsersRepo *mockcontract.MockUsersRepository,
				mockTeamsRepo *mockcontract.MockTeamsRepository,
			) {
				mockUsersRepo.EXPECT().List(mock.Anything).Return([]domain.User{
					{ID: 1, Username: "user1"},
					{ID: 2, Username: "user2"},
				}, nil)
			},
			isSuperUser:   true,
			teamID:        domain.TeamID(1),
			expectedUsers: []domain.User{{ID: 1, Username: "user1"}, {ID: 2, Username: "user2"}},
			expectedError: false,
		},
		{
			name: "Team admin can list team users",
			setupMocks: func(
				mockUsersRepo *mockcontract.MockUsersRepository,
				mockTeamsRepo *mockcontract.MockTeamsRepository,
			) {
				mockTeamsRepo.EXPECT().GetByID(mock.Anything, domain.TeamID(1)).Return(domain.Team{
					ID: 1,
					Members: []domain.TeamMember{
						{UserID: 999, Role: domain.RoleAdmin},
					},
				}, nil)
				mockUsersRepo.EXPECT().List(mock.Anything).Return([]domain.User{
					{ID: 1, Username: "user1"},
					{ID: 2, Username: "user2"},
				}, nil)
			},
			isSuperUser:   false,
			teamID:        domain.TeamID(1),
			expectedUsers: []domain.User{{ID: 1, Username: "user1"}, {ID: 2, Username: "user2"}},
			expectedError: false,
		},
		{
			name: "Team owner can list team users",
			setupMocks: func(
				mockUsersRepo *mockcontract.MockUsersRepository,
				mockTeamsRepo *mockcontract.MockTeamsRepository,
			) {
				mockTeamsRepo.EXPECT().GetByID(mock.Anything, domain.TeamID(2)).Return(domain.Team{
					ID: 2,
					Members: []domain.TeamMember{
						{UserID: 999, Role: domain.RoleOwner},
					},
				}, nil)
				mockUsersRepo.EXPECT().List(mock.Anything).Return([]domain.User{
					{ID: 3, Username: "user3"},
					{ID: 4, Username: "user4"},
				}, nil)
			},
			isSuperUser:   false,
			teamID:        domain.TeamID(2),
			expectedUsers: []domain.User{{ID: 3, Username: "user3"}, {ID: 4, Username: "user4"}},
			expectedError: false,
		},
		{
			name: "Non-admin team member is forbidden",
			setupMocks: func(
				mockUsersRepo *mockcontract.MockUsersRepository,
				mockTeamsRepo *mockcontract.MockTeamsRepository,
			) {
				mockTeamsRepo.EXPECT().GetByID(mock.Anything, domain.TeamID(3)).Return(domain.Team{
					ID: 3,
					Members: []domain.TeamMember{
						{UserID: 999, Role: domain.RoleMember},
					},
				}, nil)
			},
			isSuperUser:   false,
			teamID:        domain.TeamID(3),
			expectedUsers: nil,
			expectedError: true,
			errorContains: "forbidden",
		},
		{
			name: "Error fetching team details",
			setupMocks: func(
				mockUsersRepo *mockcontract.MockUsersRepository,
				mockTeamsRepo *mockcontract.MockTeamsRepository,
			) {
				mockTeamsRepo.EXPECT().GetByID(mock.Anything, domain.TeamID(4)).Return(domain.Team{}, errors.New("database error"))
			},
			isSuperUser:   false,
			teamID:        domain.TeamID(4),
			expectedUsers: nil,
			expectedError: true,
			errorContains: "get team by id",
		},
		{
			name: "Error fetching users",
			setupMocks: func(
				mockUsersRepo *mockcontract.MockUsersRepository,
				mockTeamsRepo *mockcontract.MockTeamsRepository,
			) {
				mockTeamsRepo.EXPECT().GetByID(mock.Anything, domain.TeamID(5)).Return(domain.Team{
					ID: 5,
					Members: []domain.TeamMember{
						{UserID: 999, Role: domain.RoleAdmin},
					},
				}, nil)
				mockUsersRepo.EXPECT().List(mock.Anything).Return(nil, errors.New("database error"))
			},
			isSuperUser:   false,
			teamID:        domain.TeamID(5),
			expectedUsers: nil,
			expectedError: true,
			errorContains: "database error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// Create mocks
			mockUsersRepo := mockcontract.NewMockUsersRepository(t)
			mockTeamsRepo := mockcontract.NewMockTeamsRepository(t)
			mockTokenizer := mockcontract.NewMockTokenizer(t)
			mockEmailer := mockcontract.NewMockEmailer(t)
			mockAuthProvider := mockusers.NewMockAuthProvider(t)
			mockRateLimiter := mockcontract.NewMockTwoFARateLimiter(t)

			// Setup mocks
			tt.setupMocks(mockUsersRepo, mockTeamsRepo)

			// Create service
			service := New(
				mockUsersRepo,
				mockTeamsRepo,
				mockTokenizer,
				mockEmailer,
				mockRateLimiter,
				[]AuthProvider{mockAuthProvider},
			)

			// Setup context
			ctx := context.Background()
			ctx = wardencontext.WithUserID(ctx, domain.UserID(999))
			if tt.isSuperUser {
				ctx = wardencontext.WithIsSuper(ctx, true)
			}

			// Call method
			users, err := service.ListForTeamAdmin(ctx, tt.teamID)

			// Check results
			if tt.expectedError {
				require.Error(t, err)
				if tt.errorContains != "" {
					require.Contains(t, err.Error(), tt.errorContains)
				}
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.expectedUsers, users)
			}
		})
	}
}
