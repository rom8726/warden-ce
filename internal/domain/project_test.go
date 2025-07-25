package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestProjectID_String(t *testing.T) {
	tests := []struct {
		name      string
		projectID ProjectID
		expected  string
	}{
		{
			name:      "zero",
			projectID: 0,
			expected:  "0",
		},
		{
			name:      "positive number",
			projectID: 123,
			expected:  "123",
		},
		{
			name:      "large number",
			projectID: 999999,
			expected:  "999999",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.projectID.String()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestProjectID_Uint(t *testing.T) {
	tests := []struct {
		name      string
		projectID ProjectID
		expected  uint
	}{
		{
			name:      "zero",
			projectID: 0,
			expected:  0,
		},
		{
			name:      "positive number",
			projectID: 123,
			expected:  123,
		},
		{
			name:      "large number",
			projectID: 999999,
			expected:  999999,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.projectID.Uint()
			assert.Equal(t, tt.expected, result)
		})
	}
}
