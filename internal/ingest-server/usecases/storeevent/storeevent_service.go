package storeevent

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/rom8726/warden/internal/domain"
	"github.com/rom8726/warden/internal/ingest-server/contract"
	"github.com/rom8726/warden/pkg/metrics"
)

type StoreEventService struct {
	storeEventProducer contract.EventProducer
}

func New(storeEventProducer contract.EventProducer) *StoreEventService {
	return &StoreEventService{
		storeEventProducer: storeEventProducer,
	}
}

// StoreEvent uses Fire-and-Forget through Kafka for async processing.
func (s *StoreEventService) StoreEvent(
	ctx context.Context,
	projectID domain.ProjectID,
	req map[string]any,
) (domain.EventID, error) {
	start := time.Now()
	projectIDStr := projectID.String()

	// Track event received
	metrics.EventsReceived.WithLabelValues(projectIDStr).Inc()

	// Extract event_id from request
	eventID, err := s.extractEventID(req)
	if err != nil {
		slog.Error("Failed to extract event_id", "error", err, "project_id", projectID)
		metrics.ValidationErrors.WithLabelValues("missing_event_id").Inc()

		return "", fmt.Errorf("extract event_id: %w", err)
	}

	// Send store event to Kafka (Fire-and-Forget)
	if err := s.storeEventProducer.SendStoreEvent(ctx, projectID, eventID, req); err != nil {
		slog.Error("Failed to send store event to Kafka", "error", err, "project_id", projectID, "event_id", eventID)
		metrics.StoreEventProcessingErrors.WithLabelValues("kafka_send_failed").Inc()

		return "", fmt.Errorf("failed to send store event to Kafka: %w", err)
	}

	// Metrics
	duration := time.Since(start)
	metrics.StoreEventProcessingDuration.WithLabelValues("store").Observe(duration.Seconds())

	slog.Debug("Store event sent to Kafka successfully",
		"project_id", projectID,
		"event_id", eventID,
		"duration_ms", duration.Milliseconds(),
	)

	return eventID, nil
}

// extractEventID extracts event_id from the request data.
func (s *StoreEventService) extractEventID(req map[string]any) (domain.EventID, error) {
	eventIDRaw, ok := req["event_id"]
	if !ok {
		return "", errors.New("event_id is required")
	}

	eventID := fmt.Sprint(eventIDRaw)

	if eventID == "" {
		return "", errors.New("event_id cannot be empty")
	}

	return domain.EventID(eventID), nil
}
