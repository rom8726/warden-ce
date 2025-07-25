package rest

import (
	"context"
	"errors"
	"testing"
	"time"

	wardencontext "github.com/rom8726/warden/internal/context"
	"github.com/rom8726/warden/internal/domain"
	generatedapi "github.com/rom8726/warden/internal/generated/server"
	mockcontract "github.com/rom8726/warden/test_mocks/internal_/backend/contract"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestRestAPI_CreateUser(t *testing.T) {
	t.Run("successful user creation", func(t *testing.T) {
		mockUsersUseCase := mockcontract.NewMockUsersUseCase(t)

		api := &RestAPI{
			usersUseCase: mockUsersUseCase,
		}

		ctx := wardencontext.WithUserID(context.Background(), domain.UserID(1))
		req := &generatedapi.CreateUserRequest{
			Username: "newuser",
			Email:    "newuser@example.com",
			Password: "password123",
		}

		expectedUser := domain.User{
			ID:            domain.UserID(2),
			Username:      "newuser",
			Email:         "newuser@example.com",
			IsSuperuser:   false,
			IsActive:      true,
			IsTmpPassword: false,
			CreatedAt:     time.Now(),
		}

		mockUsersUseCase.EXPECT().
			GetByID(mock.Anything, domain.UserID(1)).
			Return(expectedUser, nil)

		mockUsersUseCase.EXPECT().
			Create(mock.Anything, expectedUser, "newuser", "newuser@example.com", "password123", false).
			Return(expectedUser, nil)

		resp, err := api.CreateUser(ctx, req)

		require.NoError(t, err)
		require.NotNil(t, resp)

		createResp, ok := resp.(*generatedapi.CreateUserResponse)
		require.True(t, ok)
		assert.Equal(t, uint(expectedUser.ID), createResp.User.ID)
		assert.Equal(t, expectedUser.Username, createResp.User.Username)
		assert.Equal(t, expectedUser.Email, createResp.User.Email)
		assert.Equal(t, expectedUser.IsSuperuser, createResp.User.IsSuperuser)
		assert.Equal(t, expectedUser.IsActive, createResp.User.IsActive)
		assert.Equal(t, expectedUser.IsTmpPassword, createResp.User.IsTmpPassword)
	})

	t.Run("create superuser", func(t *testing.T) {
		mockUsersUseCase := mockcontract.NewMockUsersUseCase(t)

		api := &RestAPI{
			usersUseCase: mockUsersUseCase,
		}

		ctx := wardencontext.WithUserID(context.Background(), domain.UserID(1))
		req := &generatedapi.CreateUserRequest{
			Username:    "superuser",
			Email:       "superuser@example.com",
			Password:    "password123",
			IsSuperuser: generatedapi.NewOptBool(true),
		}

		currentUser := domain.User{
			ID:          domain.UserID(1),
			IsSuperuser: true,
		}

		expectedUser := domain.User{
			ID:            domain.UserID(2),
			Username:      "superuser",
			Email:         "superuser@example.com",
			IsSuperuser:   true,
			IsActive:      true,
			IsTmpPassword: false,
			CreatedAt:     time.Now(),
		}

		mockUsersUseCase.EXPECT().
			GetByID(mock.Anything, domain.UserID(1)).
			Return(currentUser, nil)

		mockUsersUseCase.EXPECT().
			Create(mock.Anything, currentUser, "superuser", "superuser@example.com", "password123", true).
			Return(expectedUser, nil)

		resp, err := api.CreateUser(ctx, req)

		require.NoError(t, err)
		require.NotNil(t, resp)

		createResp, ok := resp.(*generatedapi.CreateUserResponse)
		require.True(t, ok)
		assert.True(t, createResp.User.IsSuperuser)
	})

	t.Run("permission denied", func(t *testing.T) {
		mockUsersUseCase := mockcontract.NewMockUsersUseCase(t)

		api := &RestAPI{
			usersUseCase: mockUsersUseCase,
		}

		ctx := wardencontext.WithUserID(context.Background(), domain.UserID(1))
		req := &generatedapi.CreateUserRequest{
			Username: "newuser",
			Email:    "newuser@example.com",
			Password: "password123",
		}

		currentUser := domain.User{
			ID:          domain.UserID(1),
			IsSuperuser: false,
		}

		mockUsersUseCase.EXPECT().
			GetByID(mock.Anything, domain.UserID(1)).
			Return(currentUser, nil)

		mockUsersUseCase.EXPECT().
			Create(mock.Anything, currentUser, "newuser", "newuser@example.com", "password123", false).
			Return(domain.User{}, domain.ErrForbidden)

		resp, err := api.CreateUser(ctx, req)

		require.NoError(t, err)
		require.NotNil(t, resp)

		errorResp, ok := resp.(*generatedapi.ErrorPermissionDenied)
		require.True(t, ok)
		assert.Equal(t, "Only superusers can create new users", errorResp.Error.Message.Value)
	})

	t.Run("username already in use", func(t *testing.T) {
		mockUsersUseCase := mockcontract.NewMockUsersUseCase(t)

		api := &RestAPI{
			usersUseCase: mockUsersUseCase,
		}

		ctx := wardencontext.WithUserID(context.Background(), domain.UserID(1))
		req := &generatedapi.CreateUserRequest{
			Username: "existinguser",
			Email:    "newuser@example.com",
			Password: "password123",
		}

		currentUser := domain.User{
			ID:          domain.UserID(1),
			IsSuperuser: true,
		}

		mockUsersUseCase.EXPECT().
			GetByID(mock.Anything, domain.UserID(1)).
			Return(currentUser, nil)

		mockUsersUseCase.EXPECT().
			Create(mock.Anything, currentUser, "existinguser", "newuser@example.com", "password123", false).
			Return(domain.User{}, domain.ErrUsernameAlreadyInUse)

		resp, err := api.CreateUser(ctx, req)

		require.NoError(t, err)
		require.NotNil(t, resp)

		errorResp, ok := resp.(*generatedapi.ErrorBadRequest)
		require.True(t, ok)
		assert.Equal(t, "username already in use", errorResp.Error.Message.Value)
	})

	t.Run("email already in use", func(t *testing.T) {
		mockUsersUseCase := mockcontract.NewMockUsersUseCase(t)

		api := &RestAPI{
			usersUseCase: mockUsersUseCase,
		}

		ctx := wardencontext.WithUserID(context.Background(), domain.UserID(1))
		req := &generatedapi.CreateUserRequest{
			Username: "newuser",
			Email:    "existing@example.com",
			Password: "password123",
		}

		currentUser := domain.User{
			ID:          domain.UserID(1),
			IsSuperuser: true,
		}

		mockUsersUseCase.EXPECT().
			GetByID(mock.Anything, domain.UserID(1)).
			Return(currentUser, nil)

		mockUsersUseCase.EXPECT().
			Create(mock.Anything, currentUser, "newuser", "existing@example.com", "password123", false).
			Return(domain.User{}, domain.ErrEmailAlreadyInUse)

		resp, err := api.CreateUser(ctx, req)

		require.NoError(t, err)
		require.NotNil(t, resp)

		errorResp, ok := resp.(*generatedapi.ErrorBadRequest)
		require.True(t, ok)
		assert.Equal(t, "email already in use", errorResp.Error.Message.Value)
	})

	t.Run("get current user failed", func(t *testing.T) {
		mockUsersUseCase := mockcontract.NewMockUsersUseCase(t)

		api := &RestAPI{
			usersUseCase: mockUsersUseCase,
		}

		ctx := wardencontext.WithUserID(context.Background(), domain.UserID(1))
		req := &generatedapi.CreateUserRequest{
			Username: "newuser",
			Email:    "newuser@example.com",
			Password: "password123",
		}

		unexpectedErr := errors.New("database error")
		mockUsersUseCase.EXPECT().
			GetByID(mock.Anything, domain.UserID(1)).
			Return(domain.User{}, unexpectedErr)

		resp, err := api.CreateUser(ctx, req)

		require.Error(t, err)
		require.Nil(t, resp)
		assert.Equal(t, unexpectedErr, err)
	})

	t.Run("unexpected error", func(t *testing.T) {
		mockUsersUseCase := mockcontract.NewMockUsersUseCase(t)

		api := &RestAPI{
			usersUseCase: mockUsersUseCase,
		}

		ctx := wardencontext.WithUserID(context.Background(), domain.UserID(1))
		req := &generatedapi.CreateUserRequest{
			Username: "newuser",
			Email:    "newuser@example.com",
			Password: "password123",
		}

		currentUser := domain.User{
			ID:          domain.UserID(1),
			IsSuperuser: true,
		}

		unexpectedErr := errors.New("database error")
		mockUsersUseCase.EXPECT().
			GetByID(mock.Anything, domain.UserID(1)).
			Return(currentUser, nil)

		mockUsersUseCase.EXPECT().
			Create(mock.Anything, currentUser, "newuser", "newuser@example.com", "password123", false).
			Return(domain.User{}, unexpectedErr)

		resp, err := api.CreateUser(ctx, req)

		require.Error(t, err)
		require.Nil(t, resp)
		assert.Equal(t, unexpectedErr, err)
	})
}
