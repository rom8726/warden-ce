package rest

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	wardencontext "github.com/rom8726/warden/internal/context"
	"github.com/rom8726/warden/internal/domain"
	generatedapi "github.com/rom8726/warden/internal/generated/ingestserver"
	mockcontract "github.com/rom8726/warden/test_mocks/internal_/ingest-server/contract"
)

func TestSecurityHandler_HandleSentryAuth(t *testing.T) {
	t.Run("store event operation with valid key", func(t *testing.T) {
		mockProjectService := mockcontract.NewMockProjectsUseCase(t)

		handler := &SecurityHandler{
			projectService: mockProjectService,
		}

		ctx := wardencontext.WithProjectID(context.Background(), 123)
		tokenHolder := generatedapi.SentryAuth{
			APIKey: `sentry_key="valid_key"`,
		}

		mockProjectService.EXPECT().
			ValidateProjectKey(mock.Anything, domain.ProjectID(123), "valid_key").
			Return(true, nil)

		resultCtx, err := handler.HandleSentryAuth(ctx, generatedapi.StoreEventOperation, tokenHolder)

		require.NoError(t, err)
		assert.Equal(t, ctx, resultCtx)
	})

	t.Run("receive envelope operation with valid key", func(t *testing.T) {
		mockProjectService := mockcontract.NewMockProjectsUseCase(t)

		handler := &SecurityHandler{
			projectService: mockProjectService,
		}

		ctx := wardencontext.WithProjectID(context.Background(), 456)
		tokenHolder := generatedapi.SentryAuth{
			APIKey: `sentry_key="valid_key"`,
		}

		mockProjectService.EXPECT().
			ValidateProjectKey(mock.Anything, domain.ProjectID(456), "valid_key").
			Return(true, nil)

		resultCtx, err := handler.HandleSentryAuth(ctx, generatedapi.ReceiveEnvelopeOperation, tokenHolder)

		require.NoError(t, err)
		assert.Equal(t, ctx, resultCtx)
	})

	t.Run("invalid key", func(t *testing.T) {
		mockProjectService := mockcontract.NewMockProjectsUseCase(t)

		handler := &SecurityHandler{
			projectService: mockProjectService,
		}

		ctx := wardencontext.WithProjectID(context.Background(), 123)
		tokenHolder := generatedapi.SentryAuth{
			APIKey: `sentry_key="invalid_key"`,
		}

		mockProjectService.EXPECT().
			ValidateProjectKey(mock.Anything, domain.ProjectID(123), "invalid_key").
			Return(false, nil)

		resultCtx, err := handler.HandleSentryAuth(ctx, generatedapi.StoreEventOperation, tokenHolder)

		require.Error(t, err)
		assert.Nil(t, resultCtx)
		assert.Equal(t, "invalid or unauthorized key", err.Error())
	})

	t.Run("validation error", func(t *testing.T) {
		mockProjectService := mockcontract.NewMockProjectsUseCase(t)

		handler := &SecurityHandler{
			projectService: mockProjectService,
		}

		ctx := wardencontext.WithProjectID(context.Background(), 123)
		tokenHolder := generatedapi.SentryAuth{
			APIKey: `sentry_key="valid_key"`,
		}

		expectedErr := errors.New("validation error")
		mockProjectService.EXPECT().
			ValidateProjectKey(mock.Anything, domain.ProjectID(123), "valid_key").
			Return(false, expectedErr)

		resultCtx, err := handler.HandleSentryAuth(ctx, generatedapi.StoreEventOperation, tokenHolder)

		require.Error(t, err)
		assert.Nil(t, resultCtx)
		assert.Equal(t, expectedErr, err)
	})
}

func TestParseSentryKeyFromAuth(t *testing.T) {
	t.Run("valid sentry key", func(t *testing.T) {
		auth := `sentry_key="test_key_123"`
		result := parseSentryKeyFromAuth(auth)
		assert.Equal(t, "test_key_123", result)
	})

	t.Run("multiple parts with sentry key", func(t *testing.T) {
		auth := `sentry_key="test_key",other_param="value"`
		result := parseSentryKeyFromAuth(auth)
		assert.Equal(t, "test_key", result)
	})

	t.Run("no sentry key", func(t *testing.T) {
		auth := `other_param="value"`
		result := parseSentryKeyFromAuth(auth)
		assert.Equal(t, "", result)
	})

	t.Run("empty auth", func(t *testing.T) {
		auth := ""
		result := parseSentryKeyFromAuth(auth)
		assert.Equal(t, "", result)
	})

	t.Run("sentry key with spaces", func(t *testing.T) {
		auth := `sentry_key=" test_key "`
		result := parseSentryKeyFromAuth(auth)
		assert.Equal(t, "test_key", result)
	})
}
