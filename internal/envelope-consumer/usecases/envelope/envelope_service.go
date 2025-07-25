package envelope

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/rom8726/warden/internal/domain"
	"github.com/rom8726/warden/internal/envelope-consumer/contract"
	"github.com/rom8726/warden/pkg/metrics"
)

type EnvelopeService struct {
	eventUseCase contract.StoreEventUseCase
}

func New(
	eventUseCase contract.StoreEventUseCase,
) *EnvelopeService {
	return &EnvelopeService{
		eventUseCase: eventUseCase,
	}
}

// ProcessEnvelopeFromBytes processes envelope from byte data (used by Kafka consumer).
func (s *EnvelopeService) ProcessEnvelopeFromBytes(ctx context.Context, projectID domain.ProjectID, data []byte) error {
	start := time.Now()
	projectIDStr := projectID.String()

	scanner := bufio.NewScanner(strings.NewReader(string(data)))

	// The first line is the envelope header
	if !scanner.Scan() {
		if err := scanner.Err(); err != nil {
			slog.Error("Error reading envelope header", "error", err)

			return fmt.Errorf("error reading envelope header: %w", err)
		}

		slog.Error("Empty envelope")

		return domain.ErrNoEnvelope
	}

	headerLine := scanner.Text()
	var envelopeHeader map[string]any
	if err := json.Unmarshal([]byte(headerLine), &envelopeHeader); err != nil {
		slog.Error("Failed to parse envelope header", "error", err)

		return domain.ErrInvalidEnvelopeHeader
	}

	slog.Debug("Processing envelope from Kafka", "project_id", projectID, "header", envelopeHeader)

	// Process each item in the envelope
	for {
		// Read item header
		if !scanner.Scan() {
			break // End of envelope
		}
		itemHeaderLine := scanner.Text()

		var itemHeader map[string]any
		if err := json.Unmarshal([]byte(itemHeaderLine), &itemHeader); err != nil {
			slog.Error("Failed to parse item header", "error", err)

			continue // Skip this item
		}

		// Get item type
		itemType, ok := itemHeader["type"].(string)
		if !ok {
			slog.Error("Item header missing type field")

			continue // Skip this item
		}

		// Get item length
		itemLengthRaw, ok := itemHeader["length"].(float64)
		if !ok {
			slog.Error("Item header missing length field")

			continue // Skip this item
		}
		itemLength := int(itemLengthRaw)

		// Read item payload
		var payloadBuilder strings.Builder
		bytesRead := 0
		for bytesRead < itemLength && scanner.Scan() {
			line := scanner.Text()
			payloadBuilder.WriteString(line)
			payloadBuilder.WriteString("\n")
			bytesRead += len(line) + 1 // +1 for a newline
		}

		payload := payloadBuilder.String()

		// Process the item based on its type
		switch itemType {
		case "event":
			// Track event received
			metrics.EventsReceived.WithLabelValues(projectIDStr).Inc()

			// Parse event data
			var eventData map[string]any
			if err := json.Unmarshal([]byte(payload), &eventData); err != nil {
				slog.Error("Failed to parse event data", "error", err)
				metrics.ValidationErrors.WithLabelValues("invalid_json").Inc()

				continue
			}

			// Process the event
			eventID, err := s.eventUseCase.StoreEvent(ctx, projectID, eventData)
			if err != nil {
				slog.Error("Failed to process event", "error", err)
				metrics.ValidationErrors.WithLabelValues("process_event").Inc()

				continue
			}

			metrics.EventsProcessed.WithLabelValues(projectIDStr).Inc()
			slog.Debug("Event processed successfully", "event_id", eventID)

		default:
			slog.Info("Skipping unsupported item type", "type", itemType)
		}
	}

	if err := scanner.Err(); err != nil {
		slog.Error("Error reading envelope", "error", err)

		return fmt.Errorf("error reading envelope: %w", err)
	}

	// Track processing time
	metrics.ProcessingTime.WithLabelValues("envelope").Observe(time.Since(start).Seconds())

	return nil
}
