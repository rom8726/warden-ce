package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIssueID_Uint(t *testing.T) {
	tests := []struct {
		name     string
		issueID  IssueID
		expected uint
	}{
		{
			name:     "zero",
			issueID:  0,
			expected: 0,
		},
		{
			name:     "positive number",
			issueID:  123,
			expected: 123,
		},
		{
			name:     "large number",
			issueID:  999999,
			expected: 999999,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.issueID.Uint()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestIsNotifiableLevel(t *testing.T) {
	tests := []struct {
		name     string
		level    IssueLevel
		expected bool
	}{
		{
			name:     "debug level should not be notifiable",
			level:    IssueLevelDebug,
			expected: false,
		},
		{
			name:     "info level should not be notifiable",
			level:    IssueLevelInfo,
			expected: false,
		},
		{
			name:     "warning level should be notifiable",
			level:    IssueLevelWarning,
			expected: true,
		},
		{
			name:     "error level should be notifiable",
			level:    IssueLevelError,
			expected: true,
		},
		{
			name:     "fatal level should be notifiable",
			level:    IssueLevelFatal,
			expected: true,
		},
		{
			name:     "exception level should be notifiable",
			level:    IssueLevelException,
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsNotifiableLevel(tt.level)
			assert.Equal(t, tt.expected, result)
		})
	}
}
