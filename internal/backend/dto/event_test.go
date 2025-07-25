package dto

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/rom8726/warden/internal/domain"
	generatedapi "github.com/rom8726/warden/internal/generated/server"
)

func TestDomainIssueEventToAPI(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name     string
		event    domain.Event
		expected generatedapi.IssueEvent
	}{
		{
			name: "Basic event",
			event: domain.Event{
				ID:          "event1",
				ProjectID:   domain.ProjectID(100),
				Timestamp:   now,
				Level:       "error",
				Platform:    "go",
				Message:     "Error message",
				GroupHash:   "hash1",
				Tags:        map[string]string{"key1": "value1"},
				ServerName:  "server1",
				Environment: "production",
			},
			expected: generatedapi.IssueEvent{
				EventID:     "event1",
				ProjectID:   100,
				Message:     "Error message",
				Level:       generatedapi.IssueLevelError,
				Platform:    "go",
				Timestamp:   now,
				ServerName:  generatedapi.NewOptString("server1"),
				Environment: generatedapi.NewOptString("production"),
				Tags: generatedapi.OptIssueEventTags{
					Value: generatedapi.IssueEventTags{"key1": "value1"},
					Set:   true,
				},
			},
		},
		{
			name: "Event with multiple tags",
			event: domain.Event{
				ID:        "event2",
				ProjectID: domain.ProjectID(200),
				Timestamp: now.Add(-1 * time.Hour),
				Level:     "warning",
				Platform:  "python",
				Message:   "Warning message",
				GroupHash: "hash2",
				Tags: map[string]string{
					"key1": "value1",
					"key2": "value2",
					"key3": "value3",
				},
				ServerName:  "server2",
				Environment: "staging",
			},
			expected: generatedapi.IssueEvent{
				EventID:     "event2",
				ProjectID:   200,
				Message:     "Warning message",
				Level:       generatedapi.IssueLevelWarning,
				Platform:    "python",
				Timestamp:   now.Add(-1 * time.Hour),
				ServerName:  generatedapi.NewOptString("server2"),
				Environment: generatedapi.NewOptString("staging"),
				Tags: generatedapi.OptIssueEventTags{
					Value: generatedapi.IssueEventTags{
						"key1": "value1",
						"key2": "value2",
						"key3": "value3",
					},
					Set: true,
				},
			},
		},
		{
			name: "Event with no tags",
			event: domain.Event{
				ID:          "event3",
				ProjectID:   domain.ProjectID(300),
				Timestamp:   now.Add(-2 * time.Hour),
				Level:       "info",
				Platform:    "javascript",
				Message:     "Info message",
				GroupHash:   "hash3",
				Tags:        nil,
				ServerName:  "server3",
				Environment: "development",
			},
			expected: generatedapi.IssueEvent{
				EventID:     "event3",
				ProjectID:   300,
				Message:     "Info message",
				Level:       generatedapi.IssueLevelInfo,
				Platform:    "javascript",
				Timestamp:   now.Add(-2 * time.Hour),
				ServerName:  generatedapi.NewOptString("server3"),
				Environment: generatedapi.NewOptString("development"),
				Tags: generatedapi.OptIssueEventTags{
					Set: false,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := DomainIssueEventToAPI(tt.event)

			assert.Equal(t, tt.expected.EventID, result.EventID)
			assert.Equal(t, tt.expected.ProjectID, result.ProjectID)
			assert.Equal(t, tt.expected.Message, result.Message)
			assert.Equal(t, tt.expected.Level, result.Level)
			assert.Equal(t, tt.expected.Platform, result.Platform)
			assert.Equal(t, tt.expected.Timestamp.Unix(), result.Timestamp.Unix())
			assert.Equal(t, tt.expected.ServerName, result.ServerName)
			assert.Equal(t, tt.expected.Environment, result.Environment)
			assert.Equal(t, tt.expected.Tags, result.Tags)
		})
	}
}
