package dto

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/rom8726/warden/internal/domain"
	generatedapi "github.com/rom8726/warden/internal/generated/server"
)

func TestParseHumanDuration(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected time.Duration
		hasError bool
	}{
		{
			name:     "Minutes",
			input:    "30m",
			expected: 30 * time.Minute,
			hasError: false,
		},
		{
			name:     "Hours",
			input:    "24h",
			expected: 24 * time.Hour,
			hasError: false,
		},
		{
			name:     "Days",
			input:    "7d",
			expected: 7 * 24 * time.Hour,
			hasError: false,
		},
		{
			name:     "Single digit",
			input:    "1h",
			expected: time.Hour,
			hasError: false,
		},
		{
			name:     "Multiple digits",
			input:    "123m",
			expected: 123 * time.Minute,
			hasError: false,
		},
		{
			name:     "Too short",
			input:    "m",
			expected: 0,
			hasError: true,
		},
		{
			name:     "Invalid number",
			input:    "abc",
			expected: 0,
			hasError: true,
		},
		{
			name:     "Invalid unit",
			input:    "10x",
			expected: 0,
			hasError: true,
		},
		{
			name:     "Empty string",
			input:    "",
			expected: 0,
			hasError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ParseHumanDuration(tt.input)

			if tt.hasError {
				assert.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}

func TestTimeseriesPeriodToDomainPeriod(t *testing.T) {
	tests := []struct {
		name        string
		interval    string
		granularity string
		expected    domain.Period
		hasError    bool
	}{
		{
			name:        "Valid period - hours and minutes",
			interval:    "24h",
			granularity: "30m",
			expected: domain.Period{
				Interval:    24 * time.Hour,
				Granularity: 30 * time.Minute,
			},
			hasError: false,
		},
		{
			name:        "Valid period - days and hours",
			interval:    "7d",
			granularity: "1h",
			expected: domain.Period{
				Interval:    7 * 24 * time.Hour,
				Granularity: time.Hour,
			},
			hasError: false,
		},
		{
			name:        "Invalid interval",
			interval:    "invalid",
			granularity: "1h",
			expected:    domain.Period{},
			hasError:    true,
		},
		{
			name:        "Invalid granularity",
			interval:    "24h",
			granularity: "invalid",
			expected:    domain.Period{},
			hasError:    true,
		},
		{
			name:        "Both invalid",
			interval:    "invalid",
			granularity: "also-invalid",
			expected:    domain.Period{},
			hasError:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := TimeseriesPeriodToDomainPeriod(tt.interval, tt.granularity)

			if tt.hasError {
				assert.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.expected.Interval, result.Interval)
				assert.Equal(t, tt.expected.Granularity, result.Granularity)
			}
		})
	}
}

func TestDomainTimeseriesToAPI(t *testing.T) {
	tests := []struct {
		name        string
		timeseries  domain.Timeseries
		interval    string
		granularity string
		expected    generatedapi.TimeseriesData
	}{
		{
			name: "Basic timeseries",
			timeseries: domain.Timeseries{
				Name: "error",
				Period: domain.Period{
					Interval:    24 * time.Hour,
					Granularity: time.Hour,
				},
				Occurrences: []uint{1, 2, 3, 4, 5},
			},
			interval:    "24h",
			granularity: "1h",
			expected: generatedapi.TimeseriesData{
				Period: generatedapi.Period{
					Interval:    "24h",
					Granularity: "1h",
				},
				Name:        "error",
				Occurrences: []uint{1, 2, 3, 4, 5},
			},
		},
		{
			name: "Timeseries with different period",
			timeseries: domain.Timeseries{
				Name: "warning",
				Period: domain.Period{
					Interval:    7 * 24 * time.Hour,
					Granularity: 12 * time.Hour,
				},
				Occurrences: []uint{10, 20, 30, 40, 50, 60},
			},
			interval:    "7d",
			granularity: "12h",
			expected: generatedapi.TimeseriesData{
				Period: generatedapi.Period{
					Interval:    "7d",
					Granularity: "12h",
				},
				Name:        "warning",
				Occurrences: []uint{10, 20, 30, 40, 50, 60},
			},
		},
		{
			name: "Timeseries with no occurrences",
			timeseries: domain.Timeseries{
				Name: "info",
				Period: domain.Period{
					Interval:    30 * 24 * time.Hour,
					Granularity: 24 * time.Hour,
				},
				Occurrences: []uint{},
			},
			interval:    "30d",
			granularity: "1d",
			expected: generatedapi.TimeseriesData{
				Period: generatedapi.Period{
					Interval:    "30d",
					Granularity: "1d",
				},
				Name:        "info",
				Occurrences: []uint{},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := DomainTimeseriesToAPI(tt.timeseries, tt.interval, tt.granularity)

			assert.Equal(t, tt.expected.Period.Interval, result.Period.Interval)
			assert.Equal(t, tt.expected.Period.Granularity, result.Period.Granularity)
			assert.Equal(t, tt.expected.Name, result.Name)
			assert.Equal(t, tt.expected.Occurrences, result.Occurrences)
		})
	}
}
