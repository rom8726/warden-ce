package event

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/rom8726/warden/internal/domain"
)

func TestParseEvent(t *testing.T) {
	projectID := domain.ProjectID(1)

	tests := []struct {
		name      string
		eventData map[string]any
		wantErr   bool
		checks    func(t *testing.T, event domain.Event)
	}{
		{
			name: "valid minimal event data",
			eventData: map[string]any{
				"event_id": "123",
			},
			wantErr: false,
			checks: func(t *testing.T, event domain.Event) {
				require.Equal(t, "123", string(event.ID))
			},
		},
		{
			name: "event with all fields",
			eventData: map[string]any{
				"event_id":    "456",
				"message":     "An error occurred",
				"level":       "critical",
				"platform":    "go",
				"timestamp":   "2025-06-16T10:00:00Z",
				"tags":        map[string]any{"key": "value"},
				"environment": "production",
				"server_name": "server-1",
			},
			wantErr: false,
			checks: func(t *testing.T, event domain.Event) {
				require.Equal(t, "456", string(event.ID))
				require.Equal(t, "An error occurred", event.Message)
				require.Equal(t, domain.IssueLevel("critical"), event.Level)
				require.Equal(t, "go", event.Platform)
				require.Equal(t, "value", event.Tags["key"])
				require.Equal(t, "production", event.Environment)
				require.Equal(t, "server-1", event.ServerName)
				require.Equal(t, time.Date(2025, 6, 16, 10, 0, 0, 0, time.UTC), event.Timestamp)
			},
		},
		{
			name: "invalid timestamp format",
			eventData: map[string]any{
				"event_id":  "789",
				"timestamp": "invalid-timestamp",
			},
			wantErr: false,
			checks: func(t *testing.T, event domain.Event) {
				require.Equal(t, "789", string(event.ID))
				require.WithinDuration(t, time.Now(), event.Timestamp, 2*time.Second)
			},
		},
		{
			name: "missing event_id",
			eventData: map[string]any{
				"message": "No event ID provided",
			},
			wantErr: true,
			checks:  nil,
		},
		{
			name: "exception event data",
			eventData: map[string]any{
				"event_id": "999",
				"exception": map[string]any{
					"type":  "panic",
					"value": "Something went wrong",
				},
			},
			wantErr: false,
			checks: func(t *testing.T, event domain.Event) {
				require.Equal(t, domain.SourceException, event.Source)
				require.Equal(t, domain.IssueLevelException, event.Level)
			},
		},
		{
			name: "fatal event data",
			eventData: map[string]any{
				"event_id": "999",
				"level":    "fatal",
				"exception": map[string]any{
					"type":  "panic",
					"value": "Something went wrong",
				},
			},
			wantErr: false,
			checks: func(t *testing.T, event domain.Event) {
				require.Equal(t, domain.SourceEvent, event.Source)
				require.Equal(t, domain.IssueLevelFatal, event.Level)
			},
		},
		{
			name: "tags with non-string values",
			eventData: map[string]any{
				"event_id": "555",
				"tags": map[string]any{
					"key1": "value1",
					"key2": 42, // integer
				},
			},
			wantErr: false,
			checks: func(t *testing.T, event domain.Event) {
				require.Equal(t, "value1", event.Tags["key1"])
				require.Equal(t, "42", event.Tags["key2"]) // should convert to string
			},
		},
		{
			name: "malformed tags field",
			eventData: map[string]any{
				"event_id": "666",
				"tags":     "not a map",
			},
			wantErr: false,
			checks: func(t *testing.T, event domain.Event) {
				require.Empty(t, event.Tags) // should handle gracefully
			},
		},
		{
			name: "missing optional fields",
			eventData: map[string]any{
				"event_id": "777",
			},
			wantErr: false,
			checks: func(t *testing.T, event domain.Event) {
				require.Equal(t, "unknown", event.Environment)
				require.Equal(t, "unknown", event.ServerName)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			event, err := ParseEvent(tt.eventData, projectID)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				if tt.checks != nil {
					tt.checks(t, event)
				}
			}
		})
	}
}

func TestParseEvent2(t *testing.T) {
	projectID := domain.ProjectID(42)

	eventDataStr := `{
  "event_id": "e1e2e3e4e5",
  "timestamp": "2025-06-07T11:59:59Z",
  "platform": "javascript",
  "level": "error",
  "environment": "production",
  "server_name": "web-frontend-01",
  "release": "frontend@2.1.0",
  "dist": "stable",
  "transaction": "/checkout",
  "request": {
    "url": "https://app.example.com/checkout?item=123",
    "method": "POST",
    "headers": {
      "User-Agent": "Mozilla/5.0 (Macintosh; Intel Mac OS X 13_4_1)",
      "Referer": "https://app.example.com/cart"
    },
    "data": "{\"foo\":\"bar\"}",
    "cookies": "sessionid=xyz123",
    "query_string": "item=123",
    "env": {
      "REMOTE_ADDR": "192.0.2.12"
    }
  },
  "tags": {
    "browser.name": "Chrome",
    "browser.version": "125.0.6422.112",
    "device.arch": "x86_64",
    "build_hash": "abcdef123",
    "env": "prod"
  },
  "user": {
    "id": "user_789",
    "email": "test@example.com",
    "ip_address": "192.0.2.10"
  },
  "contexts": {
    "browser": {
      "name": "Chrome",
      "version": "125.0.6422.112"
    },
    "os": {
      "name": "Mac OS X",
      "version": "13.4.1"
    },
    "runtime": {
      "name": "browser-js",
      "version": "1.90.0"
    },
    "device": {
      "arch": "x86_64"
    }
  },
  "exception": {
    "values": [
      {
        "type": "TypeError",
        "value": "Cannot read property 'foo' of undefined",
        "stacktrace": {
          "frames": [
            {
              "filename": "https://app.example.com/static/js/bundle.js",
              "function": "renderCheckout",
              "lineno": 108,
              "colno": 24
            },
            {
              "filename": "https://app.example.com/static/js/bundle.js",
              "function": "main",
              "lineno": 200,
              "colno": 5
            }
          ]
        }
      }
    ]
  },
  "message": "Uncaught TypeError: Cannot read property 'foo' of undefined"
}`

	var eventData map[string]any
	require.NoError(t, json.Unmarshal([]byte(eventDataStr), &eventData))

	tsStr := "2025-06-07T11:59:59Z"
	ts, _ := time.Parse(time.RFC3339, tsStr)

	event, err := ParseEvent(eventData, projectID)
	require.NoError(t, err)

	// Simple scalar fields
	require.Equal(t, "e1e2e3e4e5", string(event.ID))
	require.Equal(t, projectID, event.ProjectID)
	require.Equal(t, "Uncaught TypeError: Cannot read property 'foo' of undefined", event.Message)
	require.Equal(t, domain.IssueLevelException, event.Level)
	require.Equal(t, domain.SourceException, event.Source)
	require.Equal(t, "javascript", event.Platform)
	require.Equal(t, ts, event.Timestamp)
	require.Equal(t, "production", event.Environment)
	require.Equal(t, "web-frontend-01", event.ServerName)
	require.Equal(t, "frontend@2.1.0", event.Release)

	// Tags
	require.Len(t, event.Tags, 5)
	require.Equal(t, "Chrome", event.Tags["browser.name"])

	// Request context
	require.NotNil(t, event.RequestURL)
	require.Equal(t, "https://app.example.com/checkout?item=123", *event.RequestURL)
	require.NotNil(t, event.RequestMethod)
	require.Equal(t, "POST", *event.RequestMethod)
	require.NotNil(t, event.RequestQuery)
	require.Equal(t, "item=123", *event.RequestQuery)
	require.NotNil(t, event.RequestData)
	require.Equal(t, "{\"foo\":\"bar\"}", *event.RequestData)
	require.NotNil(t, event.RequestCookies)
	require.Equal(t, "sessionid=xyz123", *event.RequestCookies)
	require.NotNil(t, event.RequestHeaders)
	require.Equal(t, "https://app.example.com/cart", event.RequestHeaders["Referer"])
	require.NotNil(t, event.UserAgent)
	require.Equal(t, "Mozilla/5.0 (Macintosh; Intel Mac OS X 13_4_1)", *event.UserAgent)
	require.NotNil(t, event.RequestIP)
	require.Equal(t, "192.0.2.12", *event.RequestIP)

	// User data
	require.NotNil(t, event.UserID)
	require.Equal(t, "user_789", *event.UserID)
	require.NotNil(t, event.UserEmail)
	require.Equal(t, "test@example.com", *event.UserEmail)

	// Runtime context
	require.NotNil(t, event.RuntimeName)
	require.Equal(t, "browser-js", *event.RuntimeName)
	require.NotNil(t, event.RuntimeVersion)
	require.Equal(t, "1.90.0", *event.RuntimeVersion)
	require.NotNil(t, event.OSName)
	require.Equal(t, "Mac OS X", *event.OSName)
	require.NotNil(t, event.OSVersion)
	require.Equal(t, "13.4.1", *event.OSVersion)
	require.NotNil(t, event.BrowserName)
	require.Equal(t, "Chrome", *event.BrowserName)
	require.NotNil(t, event.BrowserVersion)
	require.Equal(t, "125.0.6422.112", *event.BrowserVersion)
	require.NotNil(t, event.DeviceArch)
	require.Equal(t, "x86_64", *event.DeviceArch)
}

func TestExtractExceptionValues(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		data          map[string]any
		expected      []ExceptionValue
		expectedError bool
		errorContains string
	}{
		{
			name: "No exception field",
			data: map[string]any{
				"message": "Test message",
			},
			expected:      []ExceptionValue{},
			expectedError: false,
		},
		{
			name: "Single exception as map",
			data: map[string]any{
				"exception": map[string]any{
					"type":  "ValueError",
					"value": "Invalid value",
				},
			},
			expected: []ExceptionValue{
				{
					Type:  "ValueError",
					Value: "Invalid value",
				},
			},
			expectedError: false,
		},
		{
			name: "Exception list",
			data: map[string]any{
				"exception": []any{
					map[string]any{
						"type":  "ValueError",
						"value": "Invalid value",
					},
					map[string]any{
						"type":  "TypeError",
						"value": "Invalid type",
					},
				},
			},
			expected: []ExceptionValue{
				{
					Type:  "ValueError",
					Value: "Invalid value",
				},
				{
					Type:  "TypeError",
					Value: "Invalid type",
				},
			},
			expectedError: false,
		},
		{
			name: "Exception with stacktrace",
			data: map[string]any{
				"exception": map[string]any{
					"type":  "ValueError",
					"value": "Invalid value",
					"stacktrace": map[string]any{
						"frames": []any{
							map[string]any{
								"filename": "test.py",
								"lineno":   10,
							},
						},
					},
				},
			},
			expected: []ExceptionValue{
				{
					Type:  "ValueError",
					Value: "Invalid value",
					Stacktrace: mustMarshalJSON(t, map[string]any{
						"frames": []any{
							map[string]any{
								"filename": "test.py",
								"lineno":   10,
							},
						},
					}),
				},
			},
			expectedError: false,
		},
		{
			name: "Missing type in exception",
			data: map[string]any{
				"exception": map[string]any{
					"value": "Invalid value",
				},
			},
			expected:      nil,
			expectedError: true,
			errorContains: "exception type is required",
		},
		{
			name: "Missing value in exception",
			data: map[string]any{
				"exception": map[string]any{
					"type": "ValueError",
				},
			},
			expected:      nil,
			expectedError: true,
			errorContains: "exception value is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// Call the function
			result, err := extractExceptionValues(tt.data)

			// Check results
			if tt.expectedError {
				require.Error(t, err)
				if tt.errorContains != "" {
					require.Contains(t, err.Error(), tt.errorContains)
				}
			} else {
				require.NoError(t, err)
				require.Equal(t, len(tt.expected), len(result))

				for i, expected := range tt.expected {
					require.Equal(t, expected.Type, result[i].Type)
					require.Equal(t, expected.Value, result[i].Value)
					if expected.Stacktrace != nil {
						require.NotNil(t, result[i].Stacktrace)
					}
				}
			}
		})
	}
}

func TestExtractSingleException(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		exMap         map[string]any
		expected      ExceptionValue
		expectedError bool
		errorContains string
	}{
		{
			name: "Valid exception",
			exMap: map[string]any{
				"type":  "ValueError",
				"value": "Invalid value",
			},
			expected: ExceptionValue{
				Type:  "ValueError",
				Value: "Invalid value",
			},
			expectedError: false,
		},
		{
			name: "Exception with stacktrace",
			exMap: map[string]any{
				"type":  "ValueError",
				"value": "Invalid value",
				"stacktrace": map[string]any{
					"frames": []any{
						map[string]any{
							"filename": "test.py",
							"lineno":   10,
						},
					},
				},
			},
			expected: ExceptionValue{
				Type:  "ValueError",
				Value: "Invalid value",
				Stacktrace: mustMarshalJSON(t, map[string]any{
					"frames": []any{
						map[string]any{
							"filename": "test.py",
							"lineno":   10,
						},
					},
				}),
			},
			expectedError: false,
		},
		{
			name: "Missing type",
			exMap: map[string]any{
				"value": "Invalid value",
			},
			expectedError: true,
			errorContains: "exception type is required",
		},
		{
			name: "Missing value",
			exMap: map[string]any{
				"type": "ValueError",
			},
			expectedError: true,
			errorContains: "exception value is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// Call the function
			result, err := extractSingleException(tt.exMap)

			// Check results
			if tt.expectedError {
				require.Error(t, err)
				if tt.errorContains != "" {
					require.Contains(t, err.Error(), tt.errorContains)
				}
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.expected.Type, result.Type)
				require.Equal(t, tt.expected.Value, result.Value)
				if tt.expected.Stacktrace != nil {
					require.NotNil(t, result.Stacktrace)
				}
			}
		})
	}
}

// Helper function to marshal JSON for tests.
func mustMarshalJSON(t *testing.T, v any) json.RawMessage {
	data, err := json.Marshal(v)
	require.NoError(t, err)

	return data
}
