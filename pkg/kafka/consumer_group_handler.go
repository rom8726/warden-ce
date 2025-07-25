package kafka

import (
	"github.com/IBM/sarama"
)

type ConsumerGroupHandler interface {
	sarama.ConsumerGroupHandler

	Close()

	Reset()

	WaitReady()
}

type syncConsumerGroup struct {
	out   chan<- []byte
	ready chan struct{}
}

// NewSyncConsumerGroup constructor.
//
//nolint:ireturn // required by sarama package.
func NewSyncConsumerGroup(out chan<- []byte) ConsumerGroupHandler {
	return &syncConsumerGroup{
		out:   out,
		ready: make(chan struct{}),
	}
}

func (*syncConsumerGroup) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

func (c *syncConsumerGroup) Close() {
	close(c.out)
}

func (c *syncConsumerGroup) Reset() {
	c.ready = make(chan struct{})
}

func (c *syncConsumerGroup) Setup(sarama.ConsumerGroupSession) error {
	close(c.ready)

	return nil
}

func (c *syncConsumerGroup) WaitReady() {
	<-c.ready
}

func (c *syncConsumerGroup) ConsumeClaim(
	session sarama.ConsumerGroupSession,
	claim sarama.ConsumerGroupClaim,
) error {
	ctx := session.Context()
	messages := claim.Messages()

loop:
	for {
		select {
		case <-ctx.Done():
			break loop
		case message, ok := <-messages:
			if !ok {
				break loop
			}

			value := make([]byte, len(message.Value))
			copy(value, message.Value)

			select {
			case <-ctx.Done():
				break loop
			case c.out <- value:
				session.MarkMessage(message, "")
			}
		}
	}

	return nil
}
