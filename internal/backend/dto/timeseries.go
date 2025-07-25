package dto

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/rom8726/warden/internal/domain"
	generatedapi "github.com/rom8726/warden/internal/generated/server"
)

// TimeseriesPeriodToDomainPeriod converts generatedapi period parameters to domain.Period.
func TimeseriesPeriodToDomainPeriod(interval, granularity string) (domain.Period, error) {
	intervalDuration, err := ParseHumanDuration(interval)
	if err != nil {
		return domain.Period{}, err
	}
	granularityDuration, err := ParseHumanDuration(granularity)
	if err != nil {
		return domain.Period{}, err
	}

	return domain.Period{
		Interval:    intervalDuration,
		Granularity: granularityDuration,
	}, nil
}

// ParseHumanDuration parses a human-readable duration string like "1h", "30m", "7d".
func ParseHumanDuration(str string) (time.Duration, error) {
	if len(str) < 2 {
		return 0, errors.New("duration too short")
	}

	unit := str[len(str)-1]
	val, err := strconv.Atoi(str[:len(str)-1])
	if err != nil {
		return 0, fmt.Errorf("invalid number: %w", err)
	}

	switch unit {
	case 'm':
		return time.Duration(val) * time.Minute, nil
	case 'h':
		return time.Duration(val) * time.Hour, nil
	case 'd':
		return time.Duration(val) * 24 * time.Hour, nil
	default:
		return 0, fmt.Errorf("unknown unit %q", unit)
	}
}

// DomainTimeseriesToAPI converts domain.Timeseries to generatedapi.TimeseriesData.
func DomainTimeseriesToAPI(data domain.Timeseries, interval, granularity string) generatedapi.TimeseriesData {
	return generatedapi.TimeseriesData{
		Period: generatedapi.Period{
			Interval:    interval,
			Granularity: granularity,
		},
		Name:        data.Name,
		Occurrences: data.Occurrences,
	}
}

// ToTimeseriesResponse конвертирует []domain.Timeseries в []generatedapi.TimeseriesData.
func ToTimeseriesResponse(series []domain.Timeseries, interval, granularity string) []generatedapi.TimeseriesData {
	result := make([]generatedapi.TimeseriesData, len(series))
	for i, s := range series {
		result[i] = DomainTimeseriesToAPI(s, interval, granularity)
	}

	return result
}
