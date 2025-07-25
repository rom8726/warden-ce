package kafka

import (
	"context"
	"log/slog"

	"github.com/IBM/sarama"

	"github.com/rom8726/warden/pkg/resilience"
)

type Producer struct {
	producer       sarama.AsyncProducer
	circuitBreaker resilience.CircuitBreaker
}

func NewProducer(addrs []string) (*Producer, error) {
	producerConfig := sarama.NewConfig()
	producerConfig.Version = sarama.MaxVersion
	producerConfig.Producer.RequiredAcks = sarama.WaitForLocal
	producerConfig.Producer.Return.Errors = true
	producerConfig.Producer.Return.Successes = true

	producer, err := sarama.NewAsyncProducer(addrs, producerConfig)
	if err != nil {
		return nil, err
	}

	// Create a circuit breaker for Kafka operations
	cb := resilience.NewKafkaCircuitBreaker()

	producerWithCB := &Producer{
		producer:       producer,
		circuitBreaker: cb,
	}

	go producerWithCB.dispatch()

	return producerWithCB, nil
}

func (p *Producer) Produce(ctx context.Context, topic string, data []byte) error {
	// Use circuit breaker and retry patterns
	return resilience.WithCircuitBreakerAndRetry(
		ctx,
		p.circuitBreaker,
		func(ctx context.Context) error {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case p.producer.Input() <- newMessage(topic, data):
				return nil
			}
		},
		resilience.KafkaRetryOptions()...,
	)
}

func (p *Producer) Close() error {
	return p.producer.Close()
}

func (p *Producer) dispatch() {
	for {
		select {
		case msg, ok := <-p.producer.Successes():
			if !ok {
				return
			}
			freeMessage(msg)
		case err, ok := <-p.producer.Errors():
			if !ok {
				return
			}
			freeMessage(err.Msg)
			slog.Error(err.Error())
		}
	}
}
