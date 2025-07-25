package rest

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/rom8726/warden/internal/domain"
	generatedapi "github.com/rom8726/warden/internal/generated/server"
	mockcontract "github.com/rom8726/warden/test_mocks/internal_/backend/contract"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestRestAPI_ListUsers(t *testing.T) {
	t.Run("successful users list", func(t *testing.T) {
		mockUsersUseCase := mockcontract.NewMockUsersUseCase(t)

		api := &RestAPI{
			usersUseCase: mockUsersUseCase,
		}

		lastLogin1 := time.Now().Add(-24 * time.Hour)
		lastLogin2 := time.Now().Add(-12 * time.Hour)

		expectedUsers := []domain.User{
			{
				ID:          domain.UserID(1),
				Username:    "user1",
				Email:       "user1@example.com",
				IsSuperuser: false,
				IsActive:    true,
				LastLogin:   &lastLogin1,
				CreatedAt:   time.Now().Add(-7 * 24 * time.Hour),
			},
			{
				ID:          domain.UserID(2),
				Username:    "user2",
				Email:       "user2@example.com",
				IsSuperuser: true,
				IsActive:    true,
				LastLogin:   &lastLogin2,
				CreatedAt:   time.Now().Add(-3 * 24 * time.Hour),
			},
			{
				ID:          domain.UserID(3),
				Username:    "user3",
				Email:       "user3@example.com",
				IsSuperuser: false,
				IsActive:    false,
				LastLogin:   nil,
				CreatedAt:   time.Now().Add(-1 * 24 * time.Hour),
			},
		}

		mockUsersUseCase.EXPECT().
			List(mock.Anything).
			Return(expectedUsers, nil)

		resp, err := api.ListUsers(context.Background())

		require.NoError(t, err)
		require.NotNil(t, resp)

		listResp, ok := resp.(*generatedapi.ListUsersResponse)
		require.True(t, ok)
		assert.Len(t, *listResp, 3)

		// Check first user
		assert.Equal(t, uint(1), (*listResp)[0].ID)
		assert.Equal(t, "user1", (*listResp)[0].Username)
		assert.Equal(t, "user1@example.com", (*listResp)[0].Email)
		assert.False(t, (*listResp)[0].IsSuperuser)
		assert.True(t, (*listResp)[0].IsActive)
		assert.True(t, (*listResp)[0].LastLogin.Set)
		assert.Equal(t, lastLogin1, (*listResp)[0].LastLogin.Value)

		// Check second user (superuser)
		assert.Equal(t, uint(2), (*listResp)[1].ID)
		assert.Equal(t, "user2", (*listResp)[1].Username)
		assert.Equal(t, "user2@example.com", (*listResp)[1].Email)
		assert.True(t, (*listResp)[1].IsSuperuser)
		assert.True(t, (*listResp)[1].IsActive)
		assert.True(t, (*listResp)[1].LastLogin.Set)
		assert.Equal(t, lastLogin2, (*listResp)[1].LastLogin.Value)

		// Check third user (inactive, no last login)
		assert.Equal(t, uint(3), (*listResp)[2].ID)
		assert.Equal(t, "user3", (*listResp)[2].Username)
		assert.Equal(t, "user3@example.com", (*listResp)[2].Email)
		assert.False(t, (*listResp)[2].IsSuperuser)
		assert.False(t, (*listResp)[2].IsActive)
		assert.False(t, (*listResp)[2].LastLogin.Set)
	})

	t.Run("empty users list", func(t *testing.T) {
		mockUsersUseCase := mockcontract.NewMockUsersUseCase(t)

		api := &RestAPI{
			usersUseCase: mockUsersUseCase,
		}

		expectedUsers := []domain.User{}

		mockUsersUseCase.EXPECT().
			List(mock.Anything).
			Return(expectedUsers, nil)

		resp, err := api.ListUsers(context.Background())

		require.NoError(t, err)
		require.NotNil(t, resp)

		listResp, ok := resp.(*generatedapi.ListUsersResponse)
		require.True(t, ok)
		assert.Len(t, *listResp, 0)
	})

	t.Run("forbidden error", func(t *testing.T) {
		mockUsersUseCase := mockcontract.NewMockUsersUseCase(t)

		api := &RestAPI{
			usersUseCase: mockUsersUseCase,
		}

		mockUsersUseCase.EXPECT().
			List(mock.Anything).
			Return(nil, domain.ErrForbidden)

		resp, err := api.ListUsers(context.Background())

		require.NoError(t, err)
		require.NotNil(t, resp)

		errorResp, ok := resp.(*generatedapi.ErrorPermissionDenied)
		require.True(t, ok)
		assert.Equal(t, "Only superusers can list users", errorResp.Error.Message.Value)
	})

	t.Run("unexpected error", func(t *testing.T) {
		mockUsersUseCase := mockcontract.NewMockUsersUseCase(t)

		api := &RestAPI{
			usersUseCase: mockUsersUseCase,
		}

		unexpectedErr := errors.New("database error")
		mockUsersUseCase.EXPECT().
			List(mock.Anything).
			Return(nil, unexpectedErr)

		resp, err := api.ListUsers(context.Background())

		require.Error(t, err)
		require.Nil(t, resp)
		assert.Equal(t, unexpectedErr, err)
	})
}
