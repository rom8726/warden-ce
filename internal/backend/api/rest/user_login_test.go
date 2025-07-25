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

func TestRestAPI_Login(t *testing.T) {
	t.Run("successful login", func(t *testing.T) {
		mockUsersUseCase := mockcontract.NewMockUsersUseCase(t)
		mockTokenizer := mockcontract.NewMockTokenizer(t)

		api := &RestAPI{
			usersUseCase: mockUsersUseCase,
			tokenizer:    mockTokenizer,
		}

		req := &generatedapi.LoginRequest{
			Username: "testuser",
			Password: "testpass",
		}

		expectedAccessToken := "access_token"
		expectedRefreshToken := "refresh_token"
		expectedIsTmpPassword := false
		expectedTTL := 3 * time.Hour

		mockUsersUseCase.EXPECT().
			Login(mock.Anything, "testuser", "testpass").
			Return(expectedAccessToken, expectedRefreshToken, "", expectedIsTmpPassword, nil)

		mockTokenizer.EXPECT().
			AccessTokenTTL().
			Return(expectedTTL)

		resp, err := api.Login(context.Background(), req)

		require.NoError(t, err)
		require.NotNil(t, resp)

		loginResp, ok := resp.(*generatedapi.LoginResponse)
		require.True(t, ok)
		assert.Equal(t, expectedAccessToken, loginResp.AccessToken)
		assert.Equal(t, expectedRefreshToken, loginResp.RefreshToken)
		assert.Equal(t, expectedIsTmpPassword, loginResp.IsTmpPassword)
	})

	t.Run("invalid credentials", func(t *testing.T) {
		mockUsersUseCase := mockcontract.NewMockUsersUseCase(t)

		api := &RestAPI{
			usersUseCase: mockUsersUseCase,
		}

		req := &generatedapi.LoginRequest{
			Username: "testuser",
			Password: "wrongpass",
		}

		mockUsersUseCase.EXPECT().
			Login(mock.Anything, "testuser", "wrongpass").
			Return("", "", "", false, domain.ErrInvalidCredentials)

		resp, err := api.Login(context.Background(), req)

		require.NoError(t, err)
		require.NotNil(t, resp)

		errorResp, ok := resp.(*generatedapi.ErrorInvalidCredentials)
		require.True(t, ok)
		assert.Equal(t, domain.ErrInvalidCredentials.Error(), errorResp.Error.Message.Value)
	})

	t.Run("inactive user", func(t *testing.T) {
		mockUsersUseCase := mockcontract.NewMockUsersUseCase(t)

		api := &RestAPI{
			usersUseCase: mockUsersUseCase,
		}

		req := &generatedapi.LoginRequest{
			Username: "inactiveuser",
			Password: "testpass",
		}

		mockUsersUseCase.EXPECT().
			Login(mock.Anything, "inactiveuser", "testpass").
			Return("", "", "", false, domain.ErrInactiveUser)

		resp, err := api.Login(context.Background(), req)

		require.NoError(t, err)
		require.NotNil(t, resp)

		errorResp, ok := resp.(*generatedapi.ErrorInvalidCredentials)
		require.True(t, ok)
		assert.Equal(t, domain.ErrInactiveUser.Error(), errorResp.Error.Message.Value)
	})

	t.Run("unexpected error", func(t *testing.T) {
		mockUsersUseCase := mockcontract.NewMockUsersUseCase(t)

		api := &RestAPI{
			usersUseCase: mockUsersUseCase,
		}

		req := &generatedapi.LoginRequest{
			Username: "testuser",
			Password: "testpass",
		}

		unexpectedErr := errors.New("database error")
		mockUsersUseCase.EXPECT().
			Login(mock.Anything, "testuser", "testpass").
			Return("", "", "", false, unexpectedErr)

		resp, err := api.Login(context.Background(), req)

		require.Error(t, err)
		require.Nil(t, resp)
		assert.Equal(t, unexpectedErr, err)
	})
}
