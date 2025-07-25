package kafka

import (
	"context"
	"errors"
	"log/slog"

	"github.com/IBM/sarama"
)

var (
	errEmptyTopic   = errors.New("empty topic")
	errEmptyGroupID = errors.New("empty group ID")
)

type Consumer struct {
	consumer sarama.ConsumerGroup
	topics   []string
}

func NewConsumer(
	addrs []string,
	topic, groupID string,
) (*Consumer, error) {
	if topic == "" {
		return nil, errEmptyTopic
	}
	if groupID == "" {
		return nil, errEmptyGroupID
	}

	consumerConfig := sarama.NewConfig()
	consumerConfig.Version = sarama.MaxVersion
	consumer, err := sarama.NewConsumerGroup(addrs, groupID, consumerConfig)
	if err != nil {
		return nil, err
	}

	return &Consumer{
		consumer: consumer,
		topics:   []string{topic},
	}, nil
}

func (c *Consumer) Consume(ctx context.Context) <-chan []byte {
	out := make(chan []byte, 1)

	go func() {
		handler := NewSyncConsumerGroup(out)
		defer handler.Close()

		for {
			err := c.consumer.Consume(ctx, c.topics, handler)
			if errors.Is(err, sarama.ErrClosedConsumerGroup) || ctx.Err() != nil {
				break
			}

			if err != nil {
				slog.Error(err.Error())
			}

			handler.Reset()
		}
	}()

	return out
}

func (c *Consumer) Close() error {
	return c.consumer.Close()
}
