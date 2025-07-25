package envelopequeueprocessor

import (
	"context"
	"encoding/json"
	"log/slog"
	"time"

	"github.com/rom8726/di"

	"github.com/rom8726/warden/internal/domain"
	"github.com/rom8726/warden/internal/envelope-consumer/contract"
	"github.com/rom8726/warden/pkg/metrics"
)

var _ di.Servicer = (*Service)(nil)

type Service struct {
	envelopeUseCase contract.EnvelopeUseCase
	consumers       []contract.DataConsumer
	ctx             context.Context
	cancel          context.CancelFunc
}

// New creates new envelope consumer.
func New(
	envelopeUseCase contract.EnvelopeUseCase,
	consumers []contract.DataConsumer,
) (*Service, error) {
	ctx, cancel := context.WithCancel(context.Background())

	return &Service{
		envelopeUseCase: envelopeUseCase,
		consumers:       consumers,
		ctx:             ctx,
		cancel:          cancel,
	}, nil
}

// Start starts consuming messages from Kafka.
func (s *Service) Start(context.Context) error {
	slog.Info("Starting envelope Kafka consumer")

	worker := func(consumer contract.DataConsumer) {
		messages := consumer.Consume(s.ctx)
		for {
			select {
			case <-s.ctx.Done():
				return
			case data, ok := <-messages:
				if !ok {
					return
				}
				s.processMessage(data)
			}
		}
	}

	for i := range s.consumers {
		consumer := s.consumers[i]
		go worker(consumer)
	}

	return nil
}

// Stop stops the consumer.
func (s *Service) Stop(context.Context) error {
	slog.Info("Stopping envelope Kafka consumer")
	s.cancel()

	var lastErr error
	for _, consumer := range s.consumers {
		if err := consumer.Close(); err != nil && lastErr == nil {
			lastErr = err
		}
	}

	return lastErr
}

// processMessage processes individual message from Kafka.
func (s *Service) processMessage(data []byte) {
	start := time.Now()

	// Increment received messages counter
	metrics.EnvelopeMessagesReceived.WithLabelValues("received").Inc()

	// Parse envelope message
	var envelopeMsg domain.EnvelopeMessage
	if err := json.Unmarshal(data, &envelopeMsg); err != nil {
		slog.Error("Failed to unmarshal envelope message", "error", err)
		metrics.EnvelopeProcessingErrors.WithLabelValues("unmarshal_error").Inc()

		return
	}

	// Check deadline
	if envelopeMsg.IsExpired() {
		slog.Warn("Envelope message expired",
			"envelope_id", envelopeMsg.EnvelopeID,
			"project_id", envelopeMsg.ProjectID,
		)
		metrics.EnvelopeProcessingErrors.WithLabelValues("expired").Inc()

		return
	}

	// Process envelope through a use case
	if err := s.envelopeUseCase.ProcessEnvelopeFromBytes(s.ctx, envelopeMsg.ProjectID, envelopeMsg.Data); err != nil {
		slog.Error("Failed to process envelope",
			"envelope_id", envelopeMsg.EnvelopeID,
			"project_id", envelopeMsg.ProjectID,
			"error", err,
		)
		metrics.EnvelopeProcessingErrors.WithLabelValues("processing_error").Inc()

		return
	}

	// Success processing metrics
	duration := time.Since(start)
	metrics.EnvelopeProcessingDuration.WithLabelValues("process").Observe(duration.Seconds())

	slog.Debug("Envelope processed successfully",
		"envelope_id", envelopeMsg.EnvelopeID,
		"project_id", envelopeMsg.ProjectID,
		"duration_ms", duration.Milliseconds(),
		"size_bytes", len(envelopeMsg.Data),
	)
}
