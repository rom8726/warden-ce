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

func TestRestAPI_RefreshToken(t *testing.T) {
	t.Run("successful refresh", func(t *testing.T) {
		mockUsersUseCase := mockcontract.NewMockUsersUseCase(t)
		mockTokenizer := mockcontract.NewMockTokenizer(t)

		api := &RestAPI{
			usersUseCase: mockUsersUseCase,
			tokenizer:    mockTokenizer,
		}

		req := &generatedapi.RefreshTokenRequest{
			RefreshToken: "refresh_token",
		}

		expectedAccessToken := "new_access_token"
		expectedRefreshToken := "new_refresh_token"
		expectedTTL := 3 * time.Hour

		mockUsersUseCase.EXPECT().
			LoginReissue(mock.Anything, "refresh_token").
			Return(expectedAccessToken, expectedRefreshToken, nil)

		mockTokenizer.EXPECT().
			AccessTokenTTL().
			Return(expectedTTL)

		resp, err := api.RefreshToken(context.Background(), req)

		require.NoError(t, err)
		require.NotNil(t, resp)

		refreshResp, ok := resp.(*generatedapi.RefreshTokenResponse)
		require.True(t, ok)
		assert.Equal(t, expectedAccessToken, refreshResp.AccessToken)
		assert.Equal(t, expectedRefreshToken, refreshResp.RefreshToken)
	})

	t.Run("invalid token", func(t *testing.T) {
		mockUsersUseCase := mockcontract.NewMockUsersUseCase(t)

		api := &RestAPI{
			usersUseCase: mockUsersUseCase,
		}

		req := &generatedapi.RefreshTokenRequest{
			RefreshToken: "invalid_token",
		}

		mockUsersUseCase.EXPECT().
			LoginReissue(mock.Anything, "invalid_token").
			Return("", "", domain.ErrInvalidToken)

		resp, err := api.RefreshToken(context.Background(), req)

		require.NoError(t, err)
		require.NotNil(t, resp)

		errorResp, ok := resp.(*generatedapi.ErrorUnauthorized)
		require.True(t, ok)
		assert.Equal(t, domain.ErrInvalidToken.Error(), errorResp.Error.Message.Value)
	})

	t.Run("entity not found", func(t *testing.T) {
		mockUsersUseCase := mockcontract.NewMockUsersUseCase(t)

		api := &RestAPI{
			usersUseCase: mockUsersUseCase,
		}

		req := &generatedapi.RefreshTokenRequest{
			RefreshToken: "refresh_token",
		}

		mockUsersUseCase.EXPECT().
			LoginReissue(mock.Anything, "refresh_token").
			Return("", "", domain.ErrEntityNotFound)

		resp, err := api.RefreshToken(context.Background(), req)

		require.NoError(t, err)
		require.NotNil(t, resp)

		errorResp, ok := resp.(*generatedapi.ErrorUnauthorized)
		require.True(t, ok)
		assert.Equal(t, domain.ErrEntityNotFound.Error(), errorResp.Error.Message.Value)
	})

	t.Run("inactive user", func(t *testing.T) {
		mockUsersUseCase := mockcontract.NewMockUsersUseCase(t)

		api := &RestAPI{
			usersUseCase: mockUsersUseCase,
		}

		req := &generatedapi.RefreshTokenRequest{
			RefreshToken: "refresh_token",
		}

		mockUsersUseCase.EXPECT().
			LoginReissue(mock.Anything, "refresh_token").
			Return("", "", domain.ErrInactiveUser)

		resp, err := api.RefreshToken(context.Background(), req)

		require.NoError(t, err)
		require.NotNil(t, resp)

		errorResp, ok := resp.(*generatedapi.ErrorUnauthorized)
		require.True(t, ok)
		assert.Equal(t, domain.ErrInactiveUser.Error(), errorResp.Error.Message.Value)
	})

	t.Run("unexpected error", func(t *testing.T) {
		mockUsersUseCase := mockcontract.NewMockUsersUseCase(t)

		api := &RestAPI{
			usersUseCase: mockUsersUseCase,
		}

		req := &generatedapi.RefreshTokenRequest{
			RefreshToken: "refresh_token",
		}

		unexpectedErr := errors.New("database error")
		mockUsersUseCase.EXPECT().
			LoginReissue(mock.Anything, "refresh_token").
			Return("", "", unexpectedErr)

		resp, err := api.RefreshToken(context.Background(), req)

		require.Error(t, err)
		require.Nil(t, resp)
		assert.Equal(t, unexpectedErr, err)
	})
}
