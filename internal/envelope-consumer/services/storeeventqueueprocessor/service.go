package storeeventqueueprocessor

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
	eventUseCase contract.EventUseCase
	consumers    []contract.DataConsumer
	ctx          context.Context
	cancel       context.CancelFunc
}

// New creates new store event consumer.
func New(
	eventUseCase contract.EventUseCase,
	consumers []contract.DataConsumer,
) (*Service, error) {
	ctx, cancel := context.WithCancel(context.Background())

	return &Service{
		eventUseCase: eventUseCase,
		consumers:    consumers,
		ctx:          ctx,
		cancel:       cancel,
	}, nil
}

func (s *Service) IAmStoreEventConsumer() {}

// Start starts consuming messages from Kafka.
func (s *Service) Start(context.Context) error {
	slog.Info("Starting store event Kafka consumer")

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
	slog.Info("Stopping store event Kafka consumer")
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
	metrics.StoreEventMessagesReceived.WithLabelValues("received").Inc()

	// Parse store event message
	var storeEventMsg domain.StoreEventMessage
	if err := json.Unmarshal(data, &storeEventMsg); err != nil {
		slog.Error("Failed to unmarshal store event message", "error", err)
		metrics.StoreEventProcessingErrors.WithLabelValues("unmarshal_error").Inc()

		return
	}

	// Process store event through a use case
	eventID, err := s.eventUseCase.ProcessEvent(s.ctx, storeEventMsg.ProjectID, storeEventMsg.EventData)
	if err != nil {
		slog.Error("Failed to process store event",
			"event_id", storeEventMsg.EventID,
			"project_id", storeEventMsg.ProjectID,
			"error", err,
		)
		metrics.StoreEventProcessingErrors.WithLabelValues("processing_error").Inc()

		return
	}

	// Success processing metrics
	duration := time.Since(start)
	metrics.StoreEventProcessingDuration.WithLabelValues("process").Observe(duration.Seconds())

	slog.Debug("Store event processed successfully",
		"event_id", eventID,
		"project_id", storeEventMsg.ProjectID,
		"duration_ms", duration.Milliseconds(),
	)
}
