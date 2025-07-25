package users

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/rom8726/warden/internal/domain"
	mockusers "github.com/rom8726/warden/test_mocks/internal_/backend/usecases/users"
)

func TestNewAuthProviderChain(t *testing.T) {
	t.Parallel()

	// Create mock providers
	mockProvider1 := mockusers.NewMockAuthProvider(t)
	mockProvider2 := mockusers.NewMockAuthProvider(t)

	// Create a chain
	chain := NewAuthProviderChain(mockProvider1, mockProvider2)

	// Verify chain was created correctly
	require.NotNil(t, chain)
	require.Len(t, chain.providers, 2)
	require.Equal(t, mockProvider1, chain.providers[0])
	require.Equal(t, mockProvider2, chain.providers[1])
}

func TestAuthProviderChain_Authenticate(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		setupMocks    func(mockProvider1, mockProvider2 *mockusers.MockAuthProvider)
		username      string
		password      string
		expectedUser  *domain.User
		expectedError bool
		errorIs       error
		errorContains string
	}{
		{
			name: "First provider authenticates successfully",
			setupMocks: func(mockProvider1, mockProvider2 *mockusers.MockAuthProvider) {
				mockProvider1.EXPECT().CanHandle("user1").Return(true)
				mockProvider1.EXPECT().Authenticate(
					context.Background(),
					"user1",
					"password1",
				).Return(&domain.User{ID: 1, Username: "user1"}, nil)
			},
			username:      "user1",
			password:      "password1",
			expectedUser:  &domain.User{ID: 1, Username: "user1"},
			expectedError: false,
		},
		{
			name: "First provider can't handle, second authenticates successfully",
			setupMocks: func(mockProvider1, mockProvider2 *mockusers.MockAuthProvider) {
				mockProvider1.EXPECT().CanHandle("user2").Return(false)
				mockProvider2.EXPECT().CanHandle("user2").Return(true)
				mockProvider2.EXPECT().Authenticate(
					context.Background(),
					"user2",
					"password2",
				).Return(&domain.User{ID: 2, Username: "user2"}, nil)
			},
			username:      "user2",
			password:      "password2",
			expectedUser:  &domain.User{ID: 2, Username: "user2"},
			expectedError: false,
		},
		{
			name: "Both providers can't handle",
			setupMocks: func(mockProvider1, mockProvider2 *mockusers.MockAuthProvider) {
				mockProvider1.EXPECT().CanHandle("unknown").Return(false)
				mockProvider2.EXPECT().CanHandle("unknown").Return(false)
			},
			username:      "unknown",
			password:      "password",
			expectedUser:  nil,
			expectedError: true,
			errorIs:       domain.ErrInvalidCredentials,
		},
		{
			name: "Provider returns entity not found error",
			setupMocks: func(mockProvider1, mockProvider2 *mockusers.MockAuthProvider) {
				mockProvider1.EXPECT().CanHandle("user3").Return(true)
				mockProvider1.EXPECT().Authenticate(
					context.Background(),
					"user3",
					"password3",
				).Return(nil, domain.ErrEntityNotFound)

				mockProvider2.EXPECT().CanHandle("user3").Return(false)
			},
			username:      "user3",
			password:      "password3",
			expectedUser:  nil,
			expectedError: true,
			errorIs:       domain.ErrInvalidCredentials,
		},
		{
			name: "Provider returns invalid password error",
			setupMocks: func(mockProvider1, mockProvider2 *mockusers.MockAuthProvider) {
				mockProvider1.EXPECT().CanHandle("user4").Return(true)
				mockProvider1.EXPECT().Authenticate(
					context.Background(),
					"user4",
					"password4",
				).Return(nil, domain.ErrInvalidPassword)

				mockProvider2.EXPECT().CanHandle("user4").Return(false)
			},
			username:      "user4",
			password:      "password4",
			expectedUser:  nil,
			expectedError: true,
			errorIs:       domain.ErrInvalidCredentials,
		},
		{
			name: "Provider returns other error",
			setupMocks: func(mockProvider1, mockProvider2 *mockusers.MockAuthProvider) {
				mockProvider1.EXPECT().CanHandle("user5").Return(true)
				mockProvider1.EXPECT().Authenticate(
					context.Background(),
					"user5",
					"password5",
				).Return(nil, errors.New("database error"))

				mockProvider2.EXPECT().CanHandle("user5").Return(false)
			},
			username:      "user5",
			password:      "password5",
			expectedUser:  nil,
			expectedError: true,
			errorContains: "database error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// Create mock providers
			mockProvider1 := mockusers.NewMockAuthProvider(t)
			mockProvider2 := mockusers.NewMockAuthProvider(t)

			// Setup mocks
			tt.setupMocks(mockProvider1, mockProvider2)

			// Create a chain
			chain := NewAuthProviderChain(mockProvider1, mockProvider2)

			// Call method
			user, err := chain.Authenticate(context.Background(), tt.username, tt.password)

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

func TestAuthProviderChain_CanHandle(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		setupMocks     func(mockProvider1, mockProvider2 *mockusers.MockAuthProvider)
		username       string
		expectedResult bool
	}{
		{
			name: "First provider can handle",
			setupMocks: func(mockProvider1, mockProvider2 *mockusers.MockAuthProvider) {
				mockProvider1.EXPECT().CanHandle("user1").Return(true)
			},
			username:       "user1",
			expectedResult: true,
		},
		{
			name: "Second provider can handle",
			setupMocks: func(mockProvider1, mockProvider2 *mockusers.MockAuthProvider) {
				mockProvider1.EXPECT().CanHandle("user2").Return(false)
				mockProvider2.EXPECT().CanHandle("user2").Return(true)
			},
			username:       "user2",
			expectedResult: true,
		},
		{
			name: "No provider can handle",
			setupMocks: func(mockProvider1, mockProvider2 *mockusers.MockAuthProvider) {
				mockProvider1.EXPECT().CanHandle("unknown").Return(false)
				mockProvider2.EXPECT().CanHandle("unknown").Return(false)
			},
			username:       "unknown",
			expectedResult: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// Create mock providers
			mockProvider1 := mockusers.NewMockAuthProvider(t)
			mockProvider2 := mockusers.NewMockAuthProvider(t)

			// Setup mocks
			tt.setupMocks(mockProvider1, mockProvider2)

			// Create chain
			chain := NewAuthProviderChain(mockProvider1, mockProvider2)

			// Call method
			result := chain.CanHandle(tt.username)

			// Check result
			require.Equal(t, tt.expectedResult, result)
		})
	}
}
