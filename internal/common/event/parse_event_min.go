package event

import (
	"fmt"

	"github.com/rom8726/warden/internal/domain"
)

func ParseEventMinimal(eventData map[string]any) (domain.EventMinimal, error) {
	// Extract message
	message := extractMessage(eventData)

	// Extract level (default to "error" if not present)
	level := extractLevel(eventData)

	// Extract platform (default to "unknown" if not present)
	platform := extractPlatform(eventData)

	source := domain.SourceEvent
	var exceptionData domain.ExceptionData
	if _, ok := eventData["exception"]; ok { //nolint:nestif // need refactoring
		if level != domain.IssueLevelFatal {
			source = domain.SourceException
			level = domain.IssueLevelException
		}

		exValues, err := extractExceptionValues(eventData)
		if err != nil {
			return domain.EventMinimal{}, fmt.Errorf("extract exception values: %w", err)
		}
		if len(exValues) > 0 {
			value := exValues[0]
			exceptionData.ExceptionStacktrace = value.Stacktrace
			exceptionData.ExceptionType = &value.Type
			exceptionData.ExceptionValue = &value.Value
			if message == "" {
				message = value.Value
			}
		}
	}

	return domain.EventMinimal{
		Level:         level,
		Source:        source,
		Platform:      platform,
		Message:       message,
		ExceptionData: exceptionData,
	}, nil
}
