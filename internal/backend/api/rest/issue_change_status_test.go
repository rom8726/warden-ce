package rest

import (
	"context"
	"errors"
	"testing"

	"github.com/rom8726/warden/internal/domain"
	generatedapi "github.com/rom8726/warden/internal/generated/server"
	mockcontract "github.com/rom8726/warden/test_mocks/internal_/backend/contract"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestRestAPI_ChangeIssueStatus(t *testing.T) {
	t.Run("successful status change to resolved", func(t *testing.T) {
		mockIssueUseCase := mockcontract.NewMockIssueUseCase(t)
		mockPermissionsService := mockcontract.NewMockPermissionsService(t)

		api := &RestAPI{
			issueUseCase:       mockIssueUseCase,
			permissionsService: mockPermissionsService,
		}

		params := generatedapi.ChangeIssueStatusParams{
			IssueID: 123,
		}

		req := &generatedapi.ChangeIssueStatusReq{
			Status: generatedapi.IssueStatusResolved,
		}

		mockPermissionsService.EXPECT().
			CanManageIssue(mock.Anything, domain.IssueID(123)).
			Return(nil)

		mockIssueUseCase.EXPECT().
			ChangeStatus(mock.Anything, domain.IssueID(123), domain.IssueStatusResolved).
			Return(nil)

		resp, err := api.ChangeIssueStatus(context.Background(), req, params)

		require.NoError(t, err)
		require.NotNil(t, resp)

		_, ok := resp.(*generatedapi.ChangeIssueStatusNoContent)
		require.True(t, ok)
	})

	t.Run("successful status change to unresolved", func(t *testing.T) {
		mockIssueUseCase := mockcontract.NewMockIssueUseCase(t)
		mockPermissionsService := mockcontract.NewMockPermissionsService(t)

		api := &RestAPI{
			issueUseCase:       mockIssueUseCase,
			permissionsService: mockPermissionsService,
		}

		params := generatedapi.ChangeIssueStatusParams{
			IssueID: 123,
		}

		req := &generatedapi.ChangeIssueStatusReq{
			Status: generatedapi.IssueStatusUnresolved,
		}

		mockPermissionsService.EXPECT().
			CanManageIssue(mock.Anything, domain.IssueID(123)).
			Return(nil)

		mockIssueUseCase.EXPECT().
			ChangeStatus(mock.Anything, domain.IssueID(123), domain.IssueStatusUnresolved).
			Return(nil)

		resp, err := api.ChangeIssueStatus(context.Background(), req, params)

		require.NoError(t, err)
		require.NotNil(t, resp)

		_, ok := resp.(*generatedapi.ChangeIssueStatusNoContent)
		require.True(t, ok)
	})

	t.Run("permission denied", func(t *testing.T) {
		mockPermissionsService := mockcontract.NewMockPermissionsService(t)

		api := &RestAPI{
			permissionsService: mockPermissionsService,
		}

		params := generatedapi.ChangeIssueStatusParams{
			IssueID: 123,
		}

		req := &generatedapi.ChangeIssueStatusReq{
			Status: generatedapi.IssueStatusResolved,
		}

		mockPermissionsService.EXPECT().
			CanManageIssue(mock.Anything, domain.IssueID(123)).
			Return(domain.ErrPermissionDenied)

		resp, err := api.ChangeIssueStatus(context.Background(), req, params)

		require.NoError(t, err)
		require.NotNil(t, resp)

		errorResp, ok := resp.(*generatedapi.ErrorNotFound)
		require.True(t, ok)
		assert.Equal(t, "issue not found", errorResp.Error.Message.Value)
	})

	t.Run("unauthorized", func(t *testing.T) {
		mockPermissionsService := mockcontract.NewMockPermissionsService(t)

		api := &RestAPI{
			permissionsService: mockPermissionsService,
		}

		params := generatedapi.ChangeIssueStatusParams{
			IssueID: 123,
		}

		req := &generatedapi.ChangeIssueStatusReq{
			Status: generatedapi.IssueStatusResolved,
		}

		mockPermissionsService.EXPECT().
			CanManageIssue(mock.Anything, domain.IssueID(123)).
			Return(domain.ErrUserNotFound)

		resp, err := api.ChangeIssueStatus(context.Background(), req, params)

		require.NoError(t, err)
		require.NotNil(t, resp)

		errorResp, ok := resp.(*generatedapi.ErrorUnauthorized)
		require.True(t, ok)
		assert.Equal(t, "unauthorized", errorResp.Error.Message.Value)
	})

	t.Run("issue not found", func(t *testing.T) {
		mockIssueUseCase := mockcontract.NewMockIssueUseCase(t)
		mockPermissionsService := mockcontract.NewMockPermissionsService(t)

		api := &RestAPI{
			issueUseCase:       mockIssueUseCase,
			permissionsService: mockPermissionsService,
		}

		params := generatedapi.ChangeIssueStatusParams{
			IssueID: 123,
		}

		req := &generatedapi.ChangeIssueStatusReq{
			Status: generatedapi.IssueStatusResolved,
		}

		mockPermissionsService.EXPECT().
			CanManageIssue(mock.Anything, domain.IssueID(123)).
			Return(nil)

		mockIssueUseCase.EXPECT().
			ChangeStatus(mock.Anything, domain.IssueID(123), domain.IssueStatusResolved).
			Return(domain.ErrEntityNotFound)

		resp, err := api.ChangeIssueStatus(context.Background(), req, params)

		require.NoError(t, err)
		require.NotNil(t, resp)

		errorResp, ok := resp.(*generatedapi.ErrorNotFound)
		require.True(t, ok)
		assert.Equal(t, domain.ErrEntityNotFound.Error(), errorResp.Error.Message.Value)
	})

	t.Run("permission check failed with unexpected error", func(t *testing.T) {
		mockPermissionsService := mockcontract.NewMockPermissionsService(t)

		api := &RestAPI{
			permissionsService: mockPermissionsService,
		}

		params := generatedapi.ChangeIssueStatusParams{
			IssueID: 123,
		}

		req := &generatedapi.ChangeIssueStatusReq{
			Status: generatedapi.IssueStatusResolved,
		}

		unexpectedErr := errors.New("database error")
		mockPermissionsService.EXPECT().
			CanManageIssue(mock.Anything, domain.IssueID(123)).
			Return(unexpectedErr)

		resp, err := api.ChangeIssueStatus(context.Background(), req, params)

		require.Error(t, err)
		require.Nil(t, resp)
		assert.Equal(t, unexpectedErr, err)
	})

	t.Run("change status failed with unexpected error", func(t *testing.T) {
		mockIssueUseCase := mockcontract.NewMockIssueUseCase(t)
		mockPermissionsService := mockcontract.NewMockPermissionsService(t)

		api := &RestAPI{
			issueUseCase:       mockIssueUseCase,
			permissionsService: mockPermissionsService,
		}

		params := generatedapi.ChangeIssueStatusParams{
			IssueID: 123,
		}

		req := &generatedapi.ChangeIssueStatusReq{
			Status: generatedapi.IssueStatusResolved,
		}

		unexpectedErr := errors.New("database error")
		mockPermissionsService.EXPECT().
			CanManageIssue(mock.Anything, domain.IssueID(123)).
			Return(nil)

		mockIssueUseCase.EXPECT().
			ChangeStatus(mock.Anything, domain.IssueID(123), domain.IssueStatusResolved).
			Return(unexpectedErr)

		resp, err := api.ChangeIssueStatus(context.Background(), req, params)

		require.Error(t, err)
		require.Nil(t, resp)
		assert.Equal(t, unexpectedErr, err)
	})
}
