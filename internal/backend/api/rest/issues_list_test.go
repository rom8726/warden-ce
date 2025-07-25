package rest

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/rom8726/warden/internal/backend/dto"
	"github.com/rom8726/warden/internal/domain"
	generatedapi "github.com/rom8726/warden/internal/generated/server"
	mockcontract "github.com/rom8726/warden/test_mocks/internal_/backend/contract"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestRestAPI_ListIssues(t *testing.T) {
	t.Run("successful issues list", func(t *testing.T) {
		mockIssueUseCase := mockcontract.NewMockIssueUseCase(t)

		api := &RestAPI{
			issueUseCase: mockIssueUseCase,
		}

		params := generatedapi.ListIssuesParams{
			Page:    1,
			PerPage: 10,
		}

		expectedFilter := dto.MakeIssuesListFilter(params)
		expectedIssues := []domain.IssueExtended{
			{
				Issue: domain.Issue{
					ID:          domain.IssueID(1),
					ProjectID:   domain.ProjectID(1),
					Title:       "Test Issue 1",
					Fingerprint: "fp1",
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
			},
			{
				Issue: domain.Issue{
					ID:          domain.IssueID(2),
					ProjectID:   domain.ProjectID(1),
					Title:       "Test Issue 2",
					Fingerprint: "fp2",
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
				ResolvedAt:         timePtr(time.Now().Add(-12 * time.Hour)),
				ResolvedBy:         domainUserIDPtr(1),
				ResolvedByUsername: stringPtr("resolver"),
			},
		}
		expectedTotal := uint64(2)

		mockIssueUseCase.EXPECT().
			List(mock.Anything, &expectedFilter).
			Return(expectedIssues, expectedTotal, nil)

		resp, err := api.ListIssues(context.Background(), params)

		require.NoError(t, err)
		require.NotNil(t, resp)

		listResp, ok := resp.(*generatedapi.ListIssuesResponse)
		require.True(t, ok)
		assert.Len(t, listResp.Issues, 2)
		assert.Equal(t, uint(expectedTotal), listResp.Total)
		assert.Equal(t, params.Page, listResp.Page)
		assert.Equal(t, params.PerPage, listResp.PerPage)

		// Check first issue
		assert.Equal(t, uint(1), listResp.Issues[0].ID)
		assert.Equal(t, "Test Issue 1", listResp.Issues[0].Title)
		assert.Equal(t, "Test Project", listResp.Issues[0].ProjectName)
		assert.Equal(t, generatedapi.IssueStatusUnresolved, listResp.Issues[0].Status)
		assert.Equal(t, generatedapi.IssueLevelError, listResp.Issues[0].Level)

		// Check second issue (resolved)
		assert.Equal(t, uint(2), listResp.Issues[1].ID)
		assert.Equal(t, "Test Issue 2", listResp.Issues[1].Title)
		assert.Equal(t, "Test Project", listResp.Issues[1].ProjectName)
		assert.Equal(t, generatedapi.IssueStatusResolved, listResp.Issues[1].Status)
		assert.Equal(t, generatedapi.IssueLevelWarning, listResp.Issues[1].Level)
		assert.True(t, listResp.Issues[1].ResolvedAt.Set)
		assert.True(t, listResp.Issues[1].ResolvedBy.Set)
		assert.Equal(t, "resolver", listResp.Issues[1].ResolvedBy.Value)
	})

	t.Run("empty issues list", func(t *testing.T) {
		mockIssueUseCase := mockcontract.NewMockIssueUseCase(t)

		api := &RestAPI{
			issueUseCase: mockIssueUseCase,
		}

		params := generatedapi.ListIssuesParams{
			Page:    1,
			PerPage: 10,
		}

		expectedFilter := dto.MakeIssuesListFilter(params)
		expectedIssues := []domain.IssueExtended{}
		expectedTotal := uint64(0)

		mockIssueUseCase.EXPECT().
			List(mock.Anything, &expectedFilter).
			Return(expectedIssues, expectedTotal, nil)

		resp, err := api.ListIssues(context.Background(), params)

		require.NoError(t, err)
		require.NotNil(t, resp)

		listResp, ok := resp.(*generatedapi.ListIssuesResponse)
		require.True(t, ok)
		assert.Len(t, listResp.Issues, 0)
		assert.Equal(t, uint(expectedTotal), listResp.Total)
	})

	t.Run("list issues failed", func(t *testing.T) {
		mockIssueUseCase := mockcontract.NewMockIssueUseCase(t)

		api := &RestAPI{
			issueUseCase: mockIssueUseCase,
		}

		params := generatedapi.ListIssuesParams{
			Page:    1,
			PerPage: 10,
		}

		expectedFilter := dto.MakeIssuesListFilter(params)
		unexpectedErr := errors.New("database error")

		mockIssueUseCase.EXPECT().
			List(mock.Anything, &expectedFilter).
			Return(nil, uint64(0), unexpectedErr)

		resp, err := api.ListIssues(context.Background(), params)

		require.Error(t, err)
		require.Nil(t, resp)
		assert.Equal(t, unexpectedErr, err)
	})

	t.Run("with project filter", func(t *testing.T) {
		mockIssueUseCase := mockcontract.NewMockIssueUseCase(t)

		api := &RestAPI{
			issueUseCase: mockIssueUseCase,
		}

		projectID := uint(123)
		params := generatedapi.ListIssuesParams{
			Page:      1,
			PerPage:   10,
			ProjectID: generatedapi.OptUint{Value: projectID, Set: true},
		}

		expectedFilter := dto.MakeIssuesListFilter(params)
		expectedIssues := []domain.IssueExtended{}
		expectedTotal := uint64(0)

		mockIssueUseCase.EXPECT().
			List(mock.Anything, &expectedFilter).
			Return(expectedIssues, expectedTotal, nil)

		resp, err := api.ListIssues(context.Background(), params)

		require.NoError(t, err)
		require.NotNil(t, resp)

		_, ok := resp.(*generatedapi.ListIssuesResponse)
		require.True(t, ok)
		assert.Equal(t, projectID, expectedFilter.ProjectID.Uint())
	})

	t.Run("with status filter", func(t *testing.T) {
		mockIssueUseCase := mockcontract.NewMockIssueUseCase(t)

		api := &RestAPI{
			issueUseCase: mockIssueUseCase,
		}

		status := generatedapi.IssueStatusUnresolved
		params := generatedapi.ListIssuesParams{
			Page:    1,
			PerPage: 10,
			Status:  generatedapi.OptIssueStatus{Value: status, Set: true},
		}

		expectedFilter := dto.MakeIssuesListFilter(params)
		var expectedIssues []domain.IssueExtended
		expectedTotal := uint64(0)

		mockIssueUseCase.EXPECT().
			List(mock.Anything, &expectedFilter).
			Return(expectedIssues, expectedTotal, nil)

		resp, err := api.ListIssues(context.Background(), params)

		require.NoError(t, err)
		require.NotNil(t, resp)

		_, ok := resp.(*generatedapi.ListIssuesResponse)
		require.True(t, ok)
		assert.Equal(t, domain.IssueStatusUnresolved, *expectedFilter.Status)
	})
}
