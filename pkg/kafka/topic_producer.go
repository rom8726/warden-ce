package kafka

import (
	"context"
	"log/slog"
)

type KafkaProducer interface {
	Produce(ctx context.Context, topic string, data []byte) error
}

type TopicProducer struct {
	producer KafkaProducer
	topic    string
}

func NewTopicProducer(producer KafkaProducer, topic string) *TopicProducer {
	return &TopicProducer{producer: producer, topic: topic}
}

func (p *TopicProducer) Produce(ctx context.Context, data []byte) error {
	return p.producer.Produce(ctx, p.topic, data)
}

type TopicProducerCreator struct {
	asyncProducer KafkaProducer
}

func NewTopicProducerCreator(asyncProducer *Producer) *TopicProducerCreator {
	return &TopicProducerCreator{
		asyncProducer: asyncProducer,
	}
}

type DataProducer interface {
	Produce(ctx context.Context, data []byte) error
}

func (c *TopicProducerCreator) Create(topic string) DataProducer {
	return NewTopicProducer(c.asyncProducer, topic)
}

type NopKafkaProducer struct{}

func NewNopKafkaProducer() *NopKafkaProducer {
	return &NopKafkaProducer{}
}

func (*NopKafkaProducer) Produce(_ context.Context, topic string, data []byte) error {
	slog.Warn("unexpected kafka produce message",
		slog.String("data", string(data)), slog.String("topic", topic))

	return nil
}

type NopTopicKafkaProducer struct{}

func NewNopTopicKafkaProducer() *NopTopicKafkaProducer {
	return &NopTopicKafkaProducer{}
}

func (*NopTopicKafkaProducer) Produce(_ context.Context, data []byte) error {
	slog.Warn("unexpected kafka produce message", slog.String("data", string(data)))

	return nil
}
