package users

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/rom8726/warden/internal/domain"
	mockcontract "github.com/rom8726/warden/test_mocks/internal_/backend/contract"
)

func TestNewLocalAuthProvider(t *testing.T) {
	t.Parallel()

	// Create a mock repository
	mockRepo := mockcontract.NewMockUsersRepository(t)

	// Create provider
	provider := NewLocalAuthProvider(mockRepo)

	// Verify the provider was created correctly
	require.NotNil(t, provider)
	require.Equal(t, mockRepo, provider.repo)
}

func TestLocalAuthProvider_Authenticate(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		setupMocks    func(mockRepo *mockcontract.MockUsersRepository)
		username      string
		password      string
		expectedUser  *domain.User
		expectedError bool
		errorIs       error
		errorContains string
	}{
		{
			name: "Success with username",
			setupMocks: func(mockRepo *mockcontract.MockUsersRepository) {
				mockRepo.EXPECT().GetByUsername(
					mock.Anything,
					"user1",
				).Return(domain.User{
					ID:           1,
					Username:     "user1",
					PasswordHash: "$2a$10$55leG6UmKY/0JIc2EZYjB./Cl.aXAPG1.B1fJS8UofqEXRsWGfQuG", // hash for "password1"
					IsActive:     true,
				}, nil)
			},
			username: "user1",
			password: "password1",
			expectedUser: &domain.User{
				ID:           1,
				Username:     "user1",
				PasswordHash: "$2a$10$55leG6UmKY/0JIc2EZYjB./Cl.aXAPG1.B1fJS8UofqEXRsWGfQuG",
				IsActive:     true,
			},
			expectedError: false,
		},
		{
			name: "Success with email",
			setupMocks: func(mockRepo *mockcontract.MockUsersRepository) {
				mockRepo.EXPECT().GetByUsername(
					mock.Anything,
					"user2@example.com",
				).Return(domain.User{}, domain.ErrEntityNotFound)

				mockRepo.EXPECT().GetByEmail(
					mock.Anything,
					"user2@example.com",
				).Return(domain.User{
					ID:           2,
					Username:     "user2",
					Email:        "user2@example.com",
					PasswordHash: "$2a$10$55leG6UmKY/0JIc2EZYjB./Cl.aXAPG1.B1fJS8UofqEXRsWGfQuG", // hash for "password1"
					IsActive:     true,
				}, nil)
			},
			username: "user2@example.com",
			password: "password1",
			expectedUser: &domain.User{
				ID:           2,
				Username:     "user2",
				Email:        "user2@example.com",
				PasswordHash: "$2a$10$55leG6UmKY/0JIc2EZYjB./Cl.aXAPG1.B1fJS8UofqEXRsWGfQuG",
				IsActive:     true,
			},
			expectedError: false,
		},
		{
			name: "User not found by username or email",
			setupMocks: func(mockRepo *mockcontract.MockUsersRepository) {
				mockRepo.EXPECT().GetByUsername(
					mock.Anything,
					"unknown",
				).Return(domain.User{}, domain.ErrEntityNotFound)

				mockRepo.EXPECT().GetByEmail(
					mock.Anything,
					"unknown",
				).Return(domain.User{}, domain.ErrEntityNotFound)
			},
			username:      "unknown",
			password:      "password",
			expectedUser:  nil,
			expectedError: true,
			errorContains: "get user by email",
			errorIs:       domain.ErrEntityNotFound,
		},
		{
			name: "Error getting user by username",
			setupMocks: func(mockRepo *mockcontract.MockUsersRepository) {
				mockRepo.EXPECT().GetByUsername(
					mock.Anything,
					"user3",
				).Return(domain.User{}, errors.New("database error"))
			},
			username:      "user3",
			password:      "password3",
			expectedUser:  nil,
			expectedError: true,
			errorContains: "get user by username",
		},
		{
			name: "User is inactive",
			setupMocks: func(mockRepo *mockcontract.MockUsersRepository) {
				mockRepo.EXPECT().GetByUsername(
					mock.Anything,
					"inactive",
				).Return(domain.User{
					ID:           4,
					Username:     "inactive",
					PasswordHash: "$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy",
					IsActive:     false,
				}, nil)
			},
			username:      "inactive",
			password:      "password",
			expectedUser:  nil,
			expectedError: true,
			errorIs:       domain.ErrInactiveUser,
		},
		{
			name: "Invalid password",
			setupMocks: func(mockRepo *mockcontract.MockUsersRepository) {
				mockRepo.EXPECT().GetByUsername(
					mock.Anything,
					"user5",
				).Return(domain.User{
					ID:           5,
					Username:     "user5",
					PasswordHash: "$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy",
					IsActive:     true,
				}, nil)
			},
			username:      "user5",
			password:      "wrongpassword",
			expectedUser:  nil,
			expectedError: true,
			errorIs:       domain.ErrInvalidPassword,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// Create a mock repository
			mockRepo := mockcontract.NewMockUsersRepository(t)

			// Setup mocks
			tt.setupMocks(mockRepo)

			// Create provider
			provider := NewLocalAuthProvider(mockRepo)

			// Call method
			user, err := provider.Authenticate(context.Background(), tt.username, tt.password)

			// Check result
			if tt.expectedError {
				require.Error(t, err)
				if tt.errorIs != nil {
					require.ErrorIs(t, err, tt.errorIs)
				}
				if tt.errorContains != "" {
					require.Contains(t, err.Error(), tt.errorContains)
				}
			} else {
				require.NoError(t, err)
			}
			require.Equal(t, tt.expectedUser, user)
		})
	}
}

func TestLocalAuthProvider_CanHandle(t *testing.T) {
	t.Parallel()

	// Create a mock repository
	mockRepo := mockcontract.NewMockUsersRepository(t)

	// Create provider
	provider := NewLocalAuthProvider(mockRepo)

	// Call method with any username
	result := provider.CanHandle("anyusername")

	// Local provider should always return true
	require.True(t, result)
}
