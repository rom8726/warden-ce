package storeeventqueueproducer

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"time"

	"github.com/rom8726/warden/internal/domain"
	"github.com/rom8726/warden/pkg/kafka"
	"github.com/rom8726/warden/pkg/metrics"
)

type Service struct {
	storeEventProducerHigh   kafka.DataProducer
	storeEventProducerNormal kafka.DataProducer
	storeEventProducerLow    kafka.DataProducer
}

func New(producerCreator *kafka.TopicProducerCreator) *Service {
	return &Service{
		storeEventProducerHigh:   producerCreator.Create(domain.StoreEventTopicHigh),
		storeEventProducerNormal: producerCreator.Create(domain.StoreEventTopicNormal),
		storeEventProducerLow:    producerCreator.Create(domain.StoreEventTopicLow),
	}
}

// SendStoreEvent asynchronously sends store event to Kafka (Fire-and-Forget).
func (p *Service) SendStoreEvent(
	ctx context.Context,
	projectID domain.ProjectID,
	eventID domain.EventID,
	eventData map[string]any,
) error {
	start := time.Now()

	// Create a store event message
	message := domain.NewStoreEventMessage(projectID, eventID, eventData)

	// Determine priority based on event data size and type
	priority := calculatePriority(eventData)
	message.SetPriority(priority)

	// Select the appropriate producer based on priority
	var producer kafka.DataProducer
	switch {
	case priority >= 8:
		producer = p.storeEventProducerHigh
	case priority >= 4:
		producer = p.storeEventProducerNormal
	default:
		producer = p.storeEventProducerLow
	}

	// Serialize message
	messageData, err := json.Marshal(message)
	if err != nil {
		metrics.StoreEventProcessingErrors.WithLabelValues("serialization_error").Inc()

		return fmt.Errorf("marshal store event message: %w", err)
	}

	// Send it to Kafka (Fire-and-Forget)
	if err := producer.Produce(ctx, messageData); err != nil {
		metrics.StoreEventProcessingErrors.WithLabelValues("kafka_send_error").Inc()

		return fmt.Errorf("send store event to kafka: %w", err)
	}

	// Metrics
	duration := time.Since(start)
	metrics.StoreEventProcessingDuration.WithLabelValues("send").Observe(duration.Seconds())
	metrics.StoreEventMessagesSent.WithLabelValues(fmt.Sprintf("priority_%d", priority)).Inc()

	slog.Debug("Store event sent to Kafka",
		"project_id", projectID,
		"event_id", eventID,
		"priority", priority,
		"duration_ms", duration.Milliseconds(),
	)

	return nil
}

// calculatePriority calculates processing priority based on event data and other factors.
func calculatePriority(eventData map[string]any) int {
	// Check if it's an exception (high priority)
	if _, ok := eventData["exception"]; ok {
		return 9 // High priority for exceptions
	}

	// Check level
	if level, ok := eventData["level"].(string); ok {
		switch level {
		case "fatal", "error":
			return 8 // High priority for errors
		case "warning":
			return 6 // Medium priority for warnings
		case "info", "debug":
			return 4 // Low priority for info/debug
		}
	}

	// Default priority
	return 5
}
