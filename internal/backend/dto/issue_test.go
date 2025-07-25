package dto

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/rom8726/warden/internal/domain"
	generatedapi "github.com/rom8726/warden/internal/generated/server"
)

var username = "testuser"

func TestDomainLevelToAPI(t *testing.T) {
	tests := []struct {
		name     string
		level    domain.IssueLevel
		expected generatedapi.IssueLevel
	}{
		{
			name:     "Exception level",
			level:    domain.IssueLevelException,
			expected: generatedapi.IssueLevelException,
		},
		{
			name:     "Error level",
			level:    domain.IssueLevelError,
			expected: generatedapi.IssueLevelError,
		},
		{
			name:     "Warning level",
			level:    domain.IssueLevelWarning,
			expected: generatedapi.IssueLevelWarning,
		},
		{
			name:     "Info level",
			level:    domain.IssueLevelInfo,
			expected: generatedapi.IssueLevelInfo,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := DomainLevelToAPI(tt.level)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestMakeIssuesListFilter(t *testing.T) {
	tests := []struct {
		name     string
		params   generatedapi.ListIssuesParams
		expected domain.ListIssuesFilter
	}{
		{
			name: "Basic params",
			params: generatedapi.ListIssuesParams{
				Page:    1,
				PerPage: 10,
			},
			expected: domain.ListIssuesFilter{
				PageNum:  1,
				PerPage:  10,
				OrderBy:  domain.OrderByFieldTotalEvents,
				OrderAsc: false,
			},
		},
		{
			name: "With project ID",
			params: generatedapi.ListIssuesParams{
				Page:      1,
				PerPage:   10,
				ProjectID: generatedapi.NewOptUint(123),
			},
			expected: domain.ListIssuesFilter{
				PageNum:  1,
				PerPage:  10,
				OrderBy:  domain.OrderByFieldTotalEvents,
				OrderAsc: false,
				ProjectID: func() *domain.ProjectID {
					id := domain.ProjectID(123)

					return &id
				}(),
			},
		},
		{
			name: "With level",
			params: generatedapi.ListIssuesParams{
				Page:    1,
				PerPage: 10,
				Level:   generatedapi.NewOptIssueLevel(generatedapi.IssueLevelError),
			},
			expected: domain.ListIssuesFilter{
				PageNum:  1,
				PerPage:  10,
				OrderBy:  domain.OrderByFieldTotalEvents,
				OrderAsc: false,
				Level: func() *domain.IssueLevel {
					level := domain.IssueLevelError

					return &level
				}(),
			},
		},
		{
			name: "With status",
			params: generatedapi.ListIssuesParams{
				Page:    1,
				PerPage: 10,
				Status:  generatedapi.NewOptIssueStatus(generatedapi.IssueStatusResolved),
			},
			expected: domain.ListIssuesFilter{
				PageNum:  1,
				PerPage:  10,
				OrderBy:  domain.OrderByFieldTotalEvents,
				OrderAsc: false,
				Status: func() *domain.IssueStatus {
					status := domain.IssueStatusResolved

					return &status
				}(),
			},
		},
		{
			name: "With all filters",
			params: generatedapi.ListIssuesParams{
				Page:      1,
				PerPage:   10,
				ProjectID: generatedapi.NewOptUint(123),
				Level:     generatedapi.NewOptIssueLevel(generatedapi.IssueLevelError),
				Status:    generatedapi.NewOptIssueStatus(generatedapi.IssueStatusResolved),
			},
			expected: domain.ListIssuesFilter{
				PageNum:  1,
				PerPage:  10,
				OrderBy:  domain.OrderByFieldTotalEvents,
				OrderAsc: false,
				ProjectID: func() *domain.ProjectID {
					id := domain.ProjectID(123)

					return &id
				}(),
				Level: func() *domain.IssueLevel {
					level := domain.IssueLevelError

					return &level
				}(),
				Status: func() *domain.IssueStatus {
					status := domain.IssueStatusResolved

					return &status
				}(),
			},
		},
		{
			name: "With sort_by total_events",
			params: generatedapi.ListIssuesParams{
				Page:    1,
				PerPage: 10,
				SortBy:  generatedapi.NewOptIssueSortColumn(generatedapi.IssueSortColumnTotalEvents),
			},
			expected: domain.ListIssuesFilter{
				PageNum:  1,
				PerPage:  10,
				OrderBy:  domain.OrderByFieldTotalEvents,
				OrderAsc: false,
			},
		},
		{
			name: "With sort_by first_seen",
			params: generatedapi.ListIssuesParams{
				Page:    1,
				PerPage: 10,
				SortBy:  generatedapi.NewOptIssueSortColumn(generatedapi.IssueSortColumnFirstSeen),
			},
			expected: domain.ListIssuesFilter{
				PageNum:  1,
				PerPage:  10,
				OrderBy:  domain.OrderByFieldFirstSeen,
				OrderAsc: false,
			},
		},
		{
			name: "With sort_by last_seen",
			params: generatedapi.ListIssuesParams{
				Page:    1,
				PerPage: 10,
				SortBy:  generatedapi.NewOptIssueSortColumn(generatedapi.IssueSortColumnLastSeen),
			},
			expected: domain.ListIssuesFilter{
				PageNum:  1,
				PerPage:  10,
				OrderBy:  domain.OrderByFieldLastSeen,
				OrderAsc: false,
			},
		},
		{
			name: "With sort_order asc",
			params: generatedapi.ListIssuesParams{
				Page:      1,
				PerPage:   10,
				SortOrder: generatedapi.NewOptSortOrder(generatedapi.SortOrderAsc),
			},
			expected: domain.ListIssuesFilter{
				PageNum:  1,
				PerPage:  10,
				OrderBy:  domain.OrderByFieldTotalEvents,
				OrderAsc: true,
			},
		},
		{
			name: "With sort_order desc",
			params: generatedapi.ListIssuesParams{
				Page:      1,
				PerPage:   10,
				SortOrder: generatedapi.NewOptSortOrder(generatedapi.SortOrderDesc),
			},
			expected: domain.ListIssuesFilter{
				PageNum:  1,
				PerPage:  10,
				OrderBy:  domain.OrderByFieldTotalEvents,
				OrderAsc: false,
			},
		},
		{
			name: "With sort_by and sort_order",
			params: generatedapi.ListIssuesParams{
				Page:      1,
				PerPage:   10,
				SortBy:    generatedapi.NewOptIssueSortColumn(generatedapi.IssueSortColumnFirstSeen),
				SortOrder: generatedapi.NewOptSortOrder(generatedapi.SortOrderAsc),
			},
			expected: domain.ListIssuesFilter{
				PageNum:  1,
				PerPage:  10,
				OrderBy:  domain.OrderByFieldFirstSeen,
				OrderAsc: true,
			},
		},
		{
			name: "With all parameters",
			params: generatedapi.ListIssuesParams{
				Page:      1,
				PerPage:   10,
				ProjectID: generatedapi.NewOptUint(123),
				Level:     generatedapi.NewOptIssueLevel(generatedapi.IssueLevelError),
				Status:    generatedapi.NewOptIssueStatus(generatedapi.IssueStatusResolved),
				SortBy:    generatedapi.NewOptIssueSortColumn(generatedapi.IssueSortColumnLastSeen),
				SortOrder: generatedapi.NewOptSortOrder(generatedapi.SortOrderAsc),
			},
			expected: domain.ListIssuesFilter{
				PageNum:  1,
				PerPage:  10,
				OrderBy:  domain.OrderByFieldLastSeen,
				OrderAsc: true,
				ProjectID: func() *domain.ProjectID {
					id := domain.ProjectID(123)

					return &id
				}(),
				Level: func() *domain.IssueLevel {
					level := domain.IssueLevelError

					return &level
				}(),
				Status: func() *domain.IssueStatus {
					status := domain.IssueStatusResolved

					return &status
				}(),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := MakeIssuesListFilter(tt.params)

			assert.Equal(t, tt.expected.PageNum, result.PageNum)
			assert.Equal(t, tt.expected.PerPage, result.PerPage)
			assert.Equal(t, tt.expected.OrderBy, result.OrderBy)
			assert.Equal(t, tt.expected.OrderAsc, result.OrderAsc)

			if tt.expected.ProjectID == nil {
				assert.Nil(t, result.ProjectID)
			} else {
				assert.NotNil(t, result.ProjectID)
				assert.Equal(t, *tt.expected.ProjectID, *result.ProjectID)
			}

			if tt.expected.Level == nil {
				assert.Nil(t, result.Level)
			} else {
				assert.NotNil(t, result.Level)
				assert.Equal(t, *tt.expected.Level, *result.Level)
			}

			if tt.expected.Status == nil {
				assert.Nil(t, result.Status)
			} else {
				assert.NotNil(t, result.Status)
				assert.Equal(t, *tt.expected.Status, *result.Status)
			}
		})
	}
}

func TestDomainIssueToAPI(t *testing.T) {
	now := time.Now()
	userID := domain.UserID(123)

	tests := []struct {
		name               string
		issue              domain.Issue
		projectName        string
		resolvedAt         *time.Time
		resolvedBy         *domain.UserID
		resolvedByUsername *string
		expected           generatedapi.Issue
	}{
		{
			name: "Basic issue without resolved info",
			issue: domain.Issue{
				ID:          domain.IssueID(1),
				ProjectID:   domain.ProjectID(100),
				Fingerprint: "fingerprint1",
				Source:      domain.SourceEvent,
				Status:      domain.IssueStatusUnresolved,
				Title:       "Test Issue",
				Level:       domain.IssueLevelError,
				Platform:    "go",
				FirstSeen:   now.Add(-24 * time.Hour),
				LastSeen:    now,
				TotalEvents: 5,
				CreatedAt:   now.Add(-24 * time.Hour),
				UpdatedAt:   now,
			},
			projectName:        "Test Project",
			resolvedAt:         nil,
			resolvedBy:         nil,
			resolvedByUsername: nil,
			expected: generatedapi.Issue{
				ID:          1,
				ProjectID:   100,
				Source:      generatedapi.IssueSourceEvent,
				Status:      generatedapi.IssueStatusUnresolved,
				ProjectName: "Test Project",
				Title:       "Test Issue",
				Message:     "Test Issue",
				Level:       generatedapi.IssueLevelError,
				Platform:    "go",
				Count:       5,
				FirstSeen:   now.Add(-24 * time.Hour),
				LastSeen:    now,
				ResolvedAt:  generatedapi.OptDateTime{},
				ResolvedBy:  generatedapi.OptString{},
			},
		},
		{
			name: "Issue with resolved info",
			issue: domain.Issue{
				ID:          domain.IssueID(2),
				ProjectID:   domain.ProjectID(200),
				Fingerprint: "fingerprint2",
				Source:      domain.SourceException,
				Status:      domain.IssueStatusResolved,
				Title:       "Resolved Issue",
				Level:       domain.IssueLevelWarning,
				Platform:    "python",
				FirstSeen:   now.Add(-48 * time.Hour),
				LastSeen:    now.Add(-24 * time.Hour),
				TotalEvents: 10,
				CreatedAt:   now.Add(-48 * time.Hour),
				UpdatedAt:   now,
			},
			projectName:        "Another Project",
			resolvedAt:         &now,
			resolvedBy:         &userID,
			resolvedByUsername: &username,
			expected: generatedapi.Issue{
				ID:          2,
				ProjectID:   200,
				Source:      generatedapi.IssueSourceException,
				Status:      generatedapi.IssueStatusResolved,
				ProjectName: "Another Project",
				Title:       "Resolved Issue",
				Message:     "Resolved Issue",
				Level:       generatedapi.IssueLevelWarning,
				Platform:    "python",
				Count:       10,
				FirstSeen:   now.Add(-48 * time.Hour),
				LastSeen:    now.Add(-24 * time.Hour),
				ResolvedAt: generatedapi.OptDateTime{
					Value: now,
					Set:   true,
				},
				ResolvedBy: generatedapi.OptString{
					Value: username,
					Set:   true,
				},
			},
		},
		{
			name: "Issue with resolved time but no user",
			issue: domain.Issue{
				ID:          domain.IssueID(3),
				ProjectID:   domain.ProjectID(300),
				Fingerprint: "fingerprint3",
				Source:      domain.SourceEvent,
				Status:      domain.IssueStatusResolved,
				Title:       "Auto-resolved Issue",
				Level:       domain.IssueLevelInfo,
				Platform:    "javascript",
				FirstSeen:   now.Add(-72 * time.Hour),
				LastSeen:    now.Add(-48 * time.Hour),
				TotalEvents: 3,
				CreatedAt:   now.Add(-72 * time.Hour),
				UpdatedAt:   now,
			},
			projectName:        "Third Project",
			resolvedAt:         &now,
			resolvedBy:         nil,
			resolvedByUsername: nil,
			expected: generatedapi.Issue{
				ID:          3,
				ProjectID:   300,
				Source:      generatedapi.IssueSourceEvent,
				Status:      generatedapi.IssueStatusResolved,
				ProjectName: "Third Project",
				Title:       "Auto-resolved Issue",
				Message:     "Auto-resolved Issue",
				Level:       generatedapi.IssueLevelInfo,
				Platform:    "javascript",
				Count:       3,
				FirstSeen:   now.Add(-72 * time.Hour),
				LastSeen:    now.Add(-48 * time.Hour),
				ResolvedAt: generatedapi.OptDateTime{
					Value: now,
					Set:   true,
				},
				ResolvedBy: generatedapi.OptString{},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := DomainIssueToAPI(
				tt.issue,
				tt.projectName,
				tt.resolvedAt,
				tt.resolvedBy,
				tt.resolvedByUsername,
			)

			assert.Equal(t, tt.expected.ID, result.ID)
			assert.Equal(t, tt.expected.ProjectID, result.ProjectID)
			assert.Equal(t, tt.expected.Source, result.Source)
			assert.Equal(t, tt.expected.Status, result.Status)
			assert.Equal(t, tt.expected.ProjectName, result.ProjectName)
			assert.Equal(t, tt.expected.Title, result.Title)
			assert.Equal(t, tt.expected.Message, result.Message)
			assert.Equal(t, tt.expected.Level, result.Level)
			assert.Equal(t, tt.expected.Platform, result.Platform)
			assert.Equal(t, tt.expected.Count, result.Count)
			assert.Equal(t, tt.expected.FirstSeen.Unix(), result.FirstSeen.Unix())
			assert.Equal(t, tt.expected.LastSeen.Unix(), result.LastSeen.Unix())
			assert.Equal(t, tt.expected.ResolvedAt.Set, result.ResolvedAt.Set)
			if tt.expected.ResolvedAt.Set {
				assert.Equal(t, tt.expected.ResolvedAt.Value.Unix(), result.ResolvedAt.Value.Unix())
			}
			assert.Equal(t, tt.expected.ResolvedBy.Set, result.ResolvedBy.Set)
			if tt.expected.ResolvedBy.Set {
				assert.Equal(t, tt.expected.ResolvedBy.Value, result.ResolvedBy.Value)
			}
		})
	}
}

func TestMakeIssueResponseWithEvent(t *testing.T) {
	now := time.Now()
	userID := domain.UserID(123)

	issue := domain.Issue{
		ID:          domain.IssueID(1),
		ProjectID:   domain.ProjectID(100),
		Fingerprint: "fingerprint1",
		Source:      domain.SourceEvent,
		Status:      domain.IssueStatusUnresolved,
		Title:       "Test Issue",
		Level:       domain.IssueLevelError,
		Platform:    "go",
		FirstSeen:   now.Add(-24 * time.Hour),
		LastSeen:    now,
		TotalEvents: 5,
		CreatedAt:   now.Add(-24 * time.Hour),
		UpdatedAt:   now,
	}

	events := []domain.Event{
		{
			ID:          "event1",
			ProjectID:   domain.ProjectID(100),
			Timestamp:   now.Add(-24 * time.Hour),
			Level:       "error",
			Platform:    "go",
			Message:     "Error message 1",
			GroupHash:   "hash1",
			Tags:        map[string]string{"key1": "value1"},
			ServerName:  "server1",
			Environment: "production",
		},
		{
			ID:          "event2",
			ProjectID:   domain.ProjectID(100),
			Timestamp:   now,
			Level:       "error",
			Platform:    "go",
			Message:     "Error message 2",
			GroupHash:   "hash1",
			Tags:        map[string]string{"key2": "value2"},
			ServerName:  "server2",
			Environment: "staging",
		},
	}

	issueExtended := domain.IssueExtendedWithChildren{
		Issue:              issue,
		ProjectName:        "Test Project",
		ResolvedBy:         &userID,
		ResolvedByUsername: &username,
		ResolvedAt:         &now,
		Events:             events,
	}

	t.Run("Issue with events", func(t *testing.T) {
		result := MakeIssueResponseWithEvent(issueExtended)

		// Check source
		assert.Equal(t, generatedapi.IssueSource(issueExtended.Source), result.Source)

		// Check issue fields
		assert.Equal(t, uint(1), result.Issue.ID)
		assert.Equal(t, uint(100), result.Issue.ProjectID)
		assert.Equal(t, generatedapi.IssueSourceEvent, result.Issue.Source)
		assert.Equal(t, generatedapi.IssueStatusUnresolved, result.Issue.Status)
		assert.Equal(t, "Test Project", result.Issue.ProjectName)
		assert.Equal(t, "Test Issue", result.Issue.Title)
		assert.Equal(t, "Test Issue", result.Issue.Message)
		assert.Equal(t, generatedapi.IssueLevelError, result.Issue.Level)
		assert.Equal(t, "go", result.Issue.Platform)
		assert.Equal(t, uint(5), result.Issue.Count)
		assert.Equal(t, now.Add(-24*time.Hour).Unix(), result.Issue.FirstSeen.Unix())
		assert.Equal(t, now.Unix(), result.Issue.LastSeen.Unix())
		assert.True(t, result.Issue.ResolvedAt.Set)
		assert.Equal(t, now.Unix(), result.Issue.ResolvedAt.Value.Unix())
		assert.True(t, result.Issue.ResolvedBy.Set)
		assert.Equal(t, username, result.Issue.ResolvedBy.Value)

		// Check events
		assert.Len(t, result.Events, 2)

		// Check first event
		assert.Equal(t, "event1", result.Events[0].EventID)
		assert.Equal(t, uint(100), result.Events[0].ProjectID)
		assert.Equal(t, "Error message 1", result.Events[0].Message)
		assert.Equal(t, generatedapi.IssueLevelError, result.Events[0].Level)
		assert.Equal(t, "go", result.Events[0].Platform)
		assert.Equal(t, now.Add(-24*time.Hour).Unix(), result.Events[0].Timestamp.Unix())
		assert.Equal(t, generatedapi.OptString{Value: "server1", Set: true}, result.Events[0].ServerName)
		assert.Equal(t, generatedapi.OptString{Value: "production", Set: true}, result.Events[0].Environment)
		assert.Equal(t, generatedapi.OptIssueEventTags{Value: generatedapi.IssueEventTags{"key1": "value1"}, Set: true}, result.Events[0].Tags)

		// Check the second event
		assert.Equal(t, "event2", result.Events[1].EventID)
		assert.Equal(t, uint(100), result.Events[1].ProjectID)
		assert.Equal(t, "Error message 2", result.Events[1].Message)
		assert.Equal(t, generatedapi.IssueLevelError, result.Events[1].Level)
		assert.Equal(t, "go", result.Events[1].Platform)
		assert.Equal(t, now.Unix(), result.Events[1].Timestamp.Unix())
		assert.Equal(t, generatedapi.OptString{Value: "server2", Set: true}, result.Events[1].ServerName)
		assert.Equal(t, generatedapi.OptString{Value: "staging", Set: true}, result.Events[1].Environment)
		assert.Equal(t, generatedapi.OptIssueEventTags{Value: generatedapi.IssueEventTags{"key2": "value2"}, Set: true}, result.Events[1].Tags)
	})
}
