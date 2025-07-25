package rest

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/rom8726/warden/internal/domain"
	generatedapi "github.com/rom8726/warden/internal/generated/server"
	mockcontract "github.com/rom8726/warden/test_mocks/internal_/backend/contract"
)

func TestRestAPI_GetRecentIssues(t *testing.T) {
	t.Run("successful recent issues list", func(t *testing.T) {
		mockIssueUseCase := mockcontract.NewMockIssueUseCase(t)

		api := &RestAPI{
			issueUseCase: mockIssueUseCase,
		}

		params := generatedapi.GetRecentIssuesParams{
			Limit: 10,
		}

		expectedIssues := []domain.IssueExtended{
			{
				Issue: domain.Issue{
					ID:          domain.IssueID(1),
					ProjectID:   domain.ProjectID(1),
					Title:       "Recent Issue 1",
					Fingerprint: "fp1",
					Status:      domain.IssueStatusUnresolved,
					Level:       domain.IssueLevelError,
					Platform:    "python",
					FirstSeen:   time.Now().Add(-2 * time.Hour),
					LastSeen:    time.Now().Add(-1 * time.Hour),
					TotalEvents: 3,
					CreatedAt:   time.Now().Add(-2 * time.Hour),
					UpdatedAt:   time.Now().Add(-1 * time.Hour),
				},
				ProjectName: "Test Project 1",
			},
			{
				Issue: domain.Issue{
					ID:          domain.IssueID(2),
					ProjectID:   domain.ProjectID(2),
					Title:       "Recent Issue 2",
					Fingerprint: "fp2",
					Status:      domain.IssueStatusUnresolved,
					Level:       domain.IssueLevelWarning,
					Platform:    "javascript",
					FirstSeen:   time.Now().Add(-3 * time.Hour),
					LastSeen:    time.Now().Add(-30 * time.Minute),
					TotalEvents: 2,
					CreatedAt:   time.Now().Add(-3 * time.Hour),
					UpdatedAt:   time.Now().Add(-30 * time.Minute),
				},
				ProjectName: "Test Project 2",
			},
		}

		mockIssueUseCase.EXPECT().
			RecentIssues(mock.Anything, params.Limit).
			Return(expectedIssues, nil)

		resp, err := api.GetRecentIssues(context.Background(), params)

		require.NoError(t, err)
		require.NotNil(t, resp)

		listResp, ok := resp.(*generatedapi.ListIssueSummariesResponse)
		require.True(t, ok)
		assert.Len(t, listResp.Issues, 2)

		// Check first issue
		assert.Equal(t, uint(1), listResp.Issues[0].ID)
		assert.Equal(t, "Recent Issue 1", listResp.Issues[0].Title)
		assert.Equal(t, uint(1), listResp.Issues[0].ProjectID)
		assert.Equal(t, generatedapi.IssueLevelError, listResp.Issues[0].Level)
		assert.Equal(t, uint(3), listResp.Issues[0].Count)

		// Check second issue
		assert.Equal(t, uint(2), listResp.Issues[1].ID)
		assert.Equal(t, "Recent Issue 2", listResp.Issues[1].Title)
		assert.Equal(t, uint(2), listResp.Issues[1].ProjectID)
		assert.Equal(t, generatedapi.IssueLevelWarning, listResp.Issues[1].Level)
		assert.Equal(t, uint(2), listResp.Issues[1].Count)
	})

	t.Run("empty recent issues list", func(t *testing.T) {
		mockIssueUseCase := mockcontract.NewMockIssueUseCase(t)

		api := &RestAPI{
			issueUseCase: mockIssueUseCase,
		}

		params := generatedapi.GetRecentIssuesParams{
			Limit: 10,
		}

		expectedIssues := []domain.IssueExtended{}

		mockIssueUseCase.EXPECT().
			RecentIssues(mock.Anything, params.Limit).
			Return(expectedIssues, nil)

		resp, err := api.GetRecentIssues(context.Background(), params)

		require.NoError(t, err)
		require.NotNil(t, resp)

		listResp, ok := resp.(*generatedapi.ListIssueSummariesResponse)
		require.True(t, ok)
		assert.Len(t, listResp.Issues, 0)
	})

	t.Run("get recent issues failed", func(t *testing.T) {
		mockIssueUseCase := mockcontract.NewMockIssueUseCase(t)

		api := &RestAPI{
			issueUseCase: mockIssueUseCase,
		}

		params := generatedapi.GetRecentIssuesParams{
			Limit: 10,
		}

		unexpectedErr := errors.New("database error")

		mockIssueUseCase.EXPECT().
			RecentIssues(mock.Anything, params.Limit).
			Return(nil, unexpectedErr)

		resp, err := api.GetRecentIssues(context.Background(), params)

		require.Error(t, err)
		require.Nil(t, resp)
		assert.Equal(t, unexpectedErr, err)
	})

	t.Run("with custom limit", func(t *testing.T) {
		mockIssueUseCase := mockcontract.NewMockIssueUseCase(t)

		api := &RestAPI{
			issueUseCase: mockIssueUseCase,
		}

		params := generatedapi.GetRecentIssuesParams{
			Limit: 5,
		}

		expectedIssues := []domain.IssueExtended{}

		mockIssueUseCase.EXPECT().
			RecentIssues(mock.Anything, params.Limit).
			Return(expectedIssues, nil)

		resp, err := api.GetRecentIssues(context.Background(), params)

		require.NoError(t, err)
		require.NotNil(t, resp)

		listResp, ok := resp.(*generatedapi.ListIssueSummariesResponse)
		require.True(t, ok)
		assert.Len(t, listResp.Issues, 0)
	})
}
