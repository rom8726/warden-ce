package envelope

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"time"

	"github.com/rom8726/warden/internal/domain"
	"github.com/rom8726/warden/internal/ingest-server/contract"
	"github.com/rom8726/warden/pkg/metrics"
)

type EnvelopeService struct {
	envelopProducer contract.EnvelopProducer
}

func New(
	envelopProducer contract.EnvelopProducer,
) *EnvelopeService {
	return &EnvelopeService{
		envelopProducer: envelopProducer,
	}
}

// ReceiveEnvelope uses Fire-and-Forget through Kafka for async processing.
func (s *EnvelopeService) ReceiveEnvelope(ctx context.Context, projectID domain.ProjectID, data io.Reader) error {
	start := time.Now()

	// Read all data from the reader
	dataBytes, err := io.ReadAll(data)
	if err != nil {
		slog.Error("Error reading envelope data", "error", err)

		return fmt.Errorf("error reading envelope data: %w", err)
	}
	if len(dataBytes) == 0 {
		slog.Error("Empty envelope")

		return errors.New("empty envelope")
	}

	// Send an envelope to Kafka (Fire-and-Forget) - no need to parse header here
	if err := s.envelopProducer.SendEnvelope(ctx, projectID, dataBytes); err != nil {
		slog.Error("Failed to send envelope to Kafka", "error", err)
		metrics.EnvelopeProcessingErrors.WithLabelValues("kafka_send_failed").Inc()

		return fmt.Errorf("failed to send envelope to Kafka: %w", err)
	}

	// Metrics
	duration := time.Since(start)
	metrics.EnvelopeProcessingDuration.WithLabelValues("receive").Observe(duration.Seconds())

	slog.Debug("Envelope sent to Kafka successfully",
		"project_id", projectID,
		"size_bytes", len(dataBytes),
		"duration_ms", duration.Milliseconds(),
	)

	return nil
}
