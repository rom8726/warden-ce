package projects

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/rom8726/warden/internal/domain"
	mockcontract "github.com/rom8726/warden/test_mocks/internal_/ingest-server/contract"
)

func TestValidateProjectKey(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		setupMocks     func(mockProjectRepo *mockcontract.MockProjectsRepository)
		projectID      domain.ProjectID
		key            string
		expectedResult bool
		expectedError  bool
		errorContains  string
	}{
		{
			name: "Valid key",
			setupMocks: func(mockProjectRepo *mockcontract.MockProjectsRepository) {
				mockProjectRepo.EXPECT().ValidateProjectKey(
					mock.Anything,
					domain.ProjectID(1),
					"valid-key",
				).Return(true, nil)
			},
			projectID:      domain.ProjectID(1),
			key:            "valid-key",
			expectedResult: true,
			expectedError:  false,
		},
		{
			name: "Invalid key",
			setupMocks: func(mockProjectRepo *mockcontract.MockProjectsRepository) {
				mockProjectRepo.EXPECT().ValidateProjectKey(
					mock.Anything,
					domain.ProjectID(1),
					"invalid-key",
				).Return(false, nil)
			},
			projectID:      domain.ProjectID(1),
			key:            "invalid-key",
			expectedResult: false,
			expectedError:  false,
		},
		{
			name: "Error validating key",
			setupMocks: func(mockProjectRepo *mockcontract.MockProjectsRepository) {
				mockProjectRepo.EXPECT().ValidateProjectKey(
					mock.Anything,
					domain.ProjectID(1),
					"error-key",
				).Return(false, errors.New("database error"))
			},
			projectID:      domain.ProjectID(1),
			key:            "error-key",
			expectedResult: false,
			expectedError:  true,
			errorContains:  "database error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// Create mocks
			mockProjectRepo := mockcontract.NewMockProjectsRepository(t)

			// Setup mocks
			tt.setupMocks(mockProjectRepo)

			// Create service
			service := New(mockProjectRepo)

			// Call method
			result, err := service.ValidateProjectKey(context.Background(), tt.projectID, tt.key)

			// Check result
			if tt.expectedError {
				require.Error(t, err)
				if tt.errorContains != "" {
					require.Contains(t, err.Error(), tt.errorContains)
				}
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.expectedResult, result)
			}
		})
	}
}

func TestValidateProjectKey_Caching(t *testing.T) {
	t.Parallel()

	// Create mocks
	mockProjectRepo := mockcontract.NewMockProjectsRepository(t)

	// Setup mocks - expect only one call for successful validation
	mockProjectRepo.EXPECT().ValidateProjectKey(
		mock.Anything,
		domain.ProjectID(1),
		"valid-key",
	).Return(true, nil).Once()

	// Create service
	service := New(mockProjectRepo)

	ctx := context.Background()
	projectID := domain.ProjectID(1)
	key := "valid-key"

	// First call - should hit repository
	result1, err := service.ValidateProjectKey(ctx, projectID, key)
	require.NoError(t, err)
	require.True(t, result1)

	// Second call - should hit cache, not repository
	result2, err := service.ValidateProjectKey(ctx, projectID, key)
	require.NoError(t, err)
	require.True(t, result2)

	// Verify that repository was called only once
	mockProjectRepo.AssertExpectations(t)
}

func TestValidateProjectKey_InvalidKeyNotCached(t *testing.T) {
	t.Parallel()

	// Create mocks
	mockProjectRepo := mockcontract.NewMockProjectsRepository(t)

	// Setup mocks - expect two calls for invalid key (not cached)
	mockProjectRepo.EXPECT().ValidateProjectKey(
		mock.Anything,
		domain.ProjectID(1),
		"invalid-key",
	).Return(false, nil).Times(2)

	// Create service
	service := New(mockProjectRepo)

	ctx := context.Background()
	projectID := domain.ProjectID(1)
	key := "invalid-key"

	// First call - should hit repository
	result1, err := service.ValidateProjectKey(ctx, projectID, key)
	require.NoError(t, err)
	require.False(t, result1)

	// Second call - should hit repository again (not cached)
	result2, err := service.ValidateProjectKey(ctx, projectID, key)
	require.NoError(t, err)
	require.False(t, result2)

	// Verify that repository was called twice
	mockProjectRepo.AssertExpectations(t)
}
