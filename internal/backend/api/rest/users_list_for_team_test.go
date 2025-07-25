package rest

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/rom8726/warden/internal/domain"
	generatedapi "github.com/rom8726/warden/internal/generated/server"
	mockcontract "github.com/rom8726/warden/test_mocks/internal_/backend/contract"
)

func TestListUsersForTeam(t *testing.T) {
	tests := []struct {
		name        string
		setupMocks  func(*mockcontract.MockUsersUseCase)
		params      generatedapi.ListUsersForTeamParams
		expectedRes generatedapi.ListUsersForTeamRes
		expectedErr string
	}{
		{
			name: "success",
			setupMocks: func(mockUsersUseCase *mockcontract.MockUsersUseCase) {
				mockUsersUseCase.EXPECT().ListForTeamAdmin(context.Background(), domain.TeamID(1)).
					Return(
						[]domain.User{
							{ID: 1, Email: "user1@example.com"},
							{ID: 2, Email: "user2@example.com"},
						}, nil,
					)
			},
			params: generatedapi.ListUsersForTeamParams{TeamID: 1},
			expectedRes: &generatedapi.ListUsersResponse{
				{ID: 1, Email: "user1@example.com"},
				{ID: 2, Email: "user2@example.com"},
			},
			expectedErr: "",
		},
		{
			name: "forbidden error",
			setupMocks: func(mockUsersUseCase *mockcontract.MockUsersUseCase) {
				mockUsersUseCase.EXPECT().ListForTeamAdmin(context.Background(), domain.TeamID(1)).
					Return(
						nil, domain.ErrForbidden,
					)
			},
			params: generatedapi.ListUsersForTeamParams{TeamID: 1},
			expectedRes: &generatedapi.ErrorPermissionDenied{
				Error: generatedapi.ErrorPermissionDeniedError{
					Message: generatedapi.NewOptString("Only superusers and team admins\\owners can list users"),
				},
			},
			expectedErr: "",
		},
		{
			name: "unknown error",
			setupMocks: func(mockUsersUseCase *mockcontract.MockUsersUseCase) {
				mockUsersUseCase.EXPECT().ListForTeamAdmin(context.Background(), domain.TeamID(1)).
					Return(
						nil, errors.New("unknown error"),
					)
			},
			params:      generatedapi.ListUsersForTeamParams{TeamID: 1},
			expectedRes: nil,
			expectedErr: "unknown error",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockUsersUseCase := mockcontract.NewMockUsersUseCase(t)
			tc.setupMocks(mockUsersUseCase)

			api := &RestAPI{
				usersUseCase: mockUsersUseCase,
			}

			res, err := api.ListUsersForTeam(context.Background(), tc.params)

			if tc.expectedErr == "" {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedRes, res)
			} else {
				assert.Nil(t, res)
				assert.EqualError(t, err, tc.expectedErr)
			}

			mockUsersUseCase.AssertExpectations(t)
		})
	}
}
