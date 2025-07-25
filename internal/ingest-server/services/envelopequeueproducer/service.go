package envelopequeueproducer

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"time"

	"github.com/rom8726/warden/internal/domain"
	"github.com/rom8726/warden/internal/ingest-server/contract"
	"github.com/rom8726/warden/pkg/kafka"
	"github.com/rom8726/warden/pkg/metrics"
)

type Service struct {
	envelopeProducerHigh   kafka.DataProducer
	envelopeProducerNormal kafka.DataProducer
	envelopeProducerLow    kafka.DataProducer
}

func New(producerCreator contract.TopicProducerCreator) *Service {
	return &Service{
		envelopeProducerHigh:   producerCreator.Create(domain.EnvelopeTopicHigh),
		envelopeProducerNormal: producerCreator.Create(domain.EnvelopeTopicNormal),
		envelopeProducerLow:    producerCreator.Create(domain.EnvelopeTopicLow),
	}
}

// SendEnvelope asynchronously sends envelope to Kafka (Fire-and-Forget).
func (p *Service) SendEnvelope(ctx context.Context, projectID domain.ProjectID, data []byte) error {
	start := time.Now()

	// Create an envelope message
	message := domain.NewEnvelopeMessage(projectID, data)

	// Determine priority based on data size
	priority := calculatePriority(data)
	message.SetPriority(priority)

	// Select the appropriate producer based on priority
	var producer kafka.DataProducer
	switch {
	case priority >= 8:
		producer = p.envelopeProducerHigh
	case priority >= 4:
		producer = p.envelopeProducerNormal
	default:
		producer = p.envelopeProducerLow
	}

	// Serialize message
	messageData, err := json.Marshal(message)
	if err != nil {
		metrics.EnvelopeProcessingErrors.WithLabelValues("serialization_error").Inc()

		return fmt.Errorf("marshal envelope message: %w", err)
	}

	// Send it to Kafka (Fire-and-Forget)
	if err := producer.Produce(ctx, messageData); err != nil {
		metrics.EnvelopeProcessingErrors.WithLabelValues("kafka_send_error").Inc()

		return fmt.Errorf("send envelope to kafka: %w", err)
	}

	// Metrics
	duration := time.Since(start)
	metrics.EnvelopeProcessingDuration.WithLabelValues("send").Observe(duration.Seconds())
	metrics.EnvelopeMessagesSent.WithLabelValues(fmt.Sprintf("priority_%d", priority)).Inc()

	slog.Debug("Envelope sent to Kafka",
		"project_id", projectID,
		"envelope_id", message.EnvelopeID,
		"priority", priority,
		"size_bytes", len(data),
		"duration_ms", duration.Milliseconds(),
	)

	return nil
}

// calculatePriority calculates processing priority based on data size and other factors.
func calculatePriority(data []byte) int {
	// Base priority based on size
	size := len(data)

	switch {
	case size < 1024: // < 1KB
		return 8 // High priority for small messages
	case size < 10240: // < 10KB
		return 6 // Medium priority
	case size < 102400: // < 100KB
		return 4 // Low priority
	default:
		return 2 // Very low priority for large messages
	}
}
