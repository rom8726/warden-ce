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

func TestRestAPI_GetIssue(t *testing.T) {
	t.Run("successful issue info", func(t *testing.T) {
		mockIssueUseCase := mockcontract.NewMockIssueUseCase(t)
		mockPermissionsService := mockcontract.NewMockPermissionsService(t)

		api := &RestAPI{
			issueUseCase:       mockIssueUseCase,
			permissionsService: mockPermissionsService,
		}

		params := generatedapi.GetIssueParams{
			ProjectID: 1,
			IssueID:   123,
		}

		expectedIssue := &domain.IssueExtendedWithChildren{
			Issue: domain.Issue{
				ID:          domain.IssueID(123),
				ProjectID:   domain.ProjectID(1),
				Title:       "Test Issue",
				Fingerprint: "fp123",
				Status:      domain.IssueStatusUnresolved,
				Level:       domain.IssueLevelError,
				Platform:    "python",
				FirstSeen:   time.Now().Add(-24 * time.Hour),
				LastSeen:    time.Now(),
				TotalEvents: 5,
				CreatedAt:   time.Now().Add(-24 * time.Hour),
				UpdatedAt:   time.Now(),
			},
			ProjectName: "Test Project",
		}

		mockPermissionsService.EXPECT().
			CanAccessIssue(mock.Anything, domain.IssueID(123)).
			Return(nil)

		mockIssueUseCase.EXPECT().
			GetByIDWithChildren(mock.Anything, domain.IssueID(123)).
			Return(*expectedIssue, nil)

		resp, err := api.GetIssue(context.Background(), params)

		require.NoError(t, err)
		require.NotNil(t, resp)

		issueResp, ok := resp.(*generatedapi.IssueResponse)
		require.True(t, ok)
		assert.Equal(t, uint(123), issueResp.Issue.ID)
		assert.Equal(t, "Test Issue", issueResp.Issue.Title)
		assert.Equal(t, "Test Project", issueResp.Issue.ProjectName)
		assert.Equal(t, generatedapi.IssueStatusUnresolved, issueResp.Issue.Status)
		assert.Equal(t, generatedapi.IssueLevelError, issueResp.Issue.Level)
		assert.Equal(t, "python", issueResp.Issue.Platform)
		assert.Equal(t, uint(5), issueResp.Issue.Count)
	})

	t.Run("resolved issue", func(t *testing.T) {
		mockIssueUseCase := mockcontract.NewMockIssueUseCase(t)
		mockPermissionsService := mockcontract.NewMockPermissionsService(t)

		api := &RestAPI{
			issueUseCase:       mockIssueUseCase,
			permissionsService: mockPermissionsService,
		}

		params := generatedapi.GetIssueParams{
			ProjectID: 1,
			IssueID:   123,
		}

		resolvedAt := time.Now().Add(-12 * time.Hour)
		resolvedBy := domain.UserID(1)
		resolvedByUsername := "resolver"

		expectedIssue := &domain.IssueExtendedWithChildren{
			Issue: domain.Issue{
				ID:          domain.IssueID(123),
				ProjectID:   domain.ProjectID(1),
				Title:       "Resolved Issue",
				Fingerprint: "fp123",
				Status:      domain.IssueStatusResolved,
				Level:       domain.IssueLevelWarning,
				Platform:    "javascript",
				FirstSeen:   time.Now().Add(-48 * time.Hour),
				LastSeen:    time.Now().Add(-12 * time.Hour),
				TotalEvents: 3,
				CreatedAt:   time.Now().Add(-48 * time.Hour),
				UpdatedAt:   time.Now().Add(-12 * time.Hour),
			},
			ProjectName:        "Test Project",
			ResolvedAt:         &resolvedAt,
			ResolvedBy:         &resolvedBy,
			ResolvedByUsername: &resolvedByUsername,
		}

		mockPermissionsService.EXPECT().
			CanAccessIssue(mock.Anything, domain.IssueID(123)).
			Return(nil)

		mockIssueUseCase.EXPECT().
			GetByIDWithChildren(mock.Anything, domain.IssueID(123)).
			Return(*expectedIssue, nil)

		resp, err := api.GetIssue(context.Background(), params)

		require.NoError(t, err)
		require.NotNil(t, resp)

		issueResp, ok := resp.(*generatedapi.IssueResponse)
		require.True(t, ok)
		assert.Equal(t, uint(123), issueResp.Issue.ID)
		assert.Equal(t, "Resolved Issue", issueResp.Issue.Title)
		assert.Equal(t, generatedapi.IssueStatusResolved, issueResp.Issue.Status)
		assert.True(t, issueResp.Issue.ResolvedAt.Set)
		assert.True(t, issueResp.Issue.ResolvedBy.Set)
		assert.Equal(t, "resolver", issueResp.Issue.ResolvedBy.Value)
	})

	t.Run("permission denied", func(t *testing.T) {
		mockPermissionsService := mockcontract.NewMockPermissionsService(t)

		api := &RestAPI{
			permissionsService: mockPermissionsService,
		}

		params := generatedapi.GetIssueParams{
			ProjectID: 1,
			IssueID:   123,
		}

		mockPermissionsService.EXPECT().
			CanAccessIssue(mock.Anything, domain.IssueID(123)).
			Return(domain.ErrPermissionDenied)

		resp, err := api.GetIssue(context.Background(), params)

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

		params := generatedapi.GetIssueParams{
			ProjectID: 1,
			IssueID:   123,
		}

		mockPermissionsService.EXPECT().
			CanAccessIssue(mock.Anything, domain.IssueID(123)).
			Return(domain.ErrUserNotFound)

		resp, err := api.GetIssue(context.Background(), params)

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

		params := generatedapi.GetIssueParams{
			ProjectID: 1,
			IssueID:   123,
		}

		mockPermissionsService.EXPECT().
			CanAccessIssue(mock.Anything, domain.IssueID(123)).
			Return(nil)

		mockIssueUseCase.EXPECT().
			GetByIDWithChildren(mock.Anything, domain.IssueID(123)).
			Return(domain.IssueExtendedWithChildren{}, domain.ErrEntityNotFound)

		resp, err := api.GetIssue(context.Background(), params)

		require.Error(t, err)
		require.Nil(t, resp)
		assert.Equal(t, domain.ErrEntityNotFound, err)
	})

	t.Run("permission check failed with unexpected error", func(t *testing.T) {
		mockPermissionsService := mockcontract.NewMockPermissionsService(t)

		api := &RestAPI{
			permissionsService: mockPermissionsService,
		}

		params := generatedapi.GetIssueParams{
			ProjectID: 1,
			IssueID:   123,
		}

		unexpectedErr := errors.New("database error")
		mockPermissionsService.EXPECT().
			CanAccessIssue(mock.Anything, domain.IssueID(123)).
			Return(unexpectedErr)

		resp, err := api.GetIssue(context.Background(), params)

		require.Error(t, err)
		require.Nil(t, resp)
		assert.Equal(t, unexpectedErr, err)
	})

	t.Run("get issue failed with unexpected error", func(t *testing.T) {
		mockIssueUseCase := mockcontract.NewMockIssueUseCase(t)
		mockPermissionsService := mockcontract.NewMockPermissionsService(t)

		api := &RestAPI{
			issueUseCase:       mockIssueUseCase,
			permissionsService: mockPermissionsService,
		}

		params := generatedapi.GetIssueParams{
			ProjectID: 1,
			IssueID:   123,
		}

		unexpectedErr := errors.New("database error")
		mockPermissionsService.EXPECT().
			CanAccessIssue(mock.Anything, domain.IssueID(123)).
			Return(nil)

		mockIssueUseCase.EXPECT().
			GetByIDWithChildren(mock.Anything, domain.IssueID(123)).
			Return(domain.IssueExtendedWithChildren{}, unexpectedErr)

		resp, err := api.GetIssue(context.Background(), params)

		require.Error(t, err)
		require.Nil(t, resp)
		assert.Equal(t, unexpectedErr, err)
	})
}
