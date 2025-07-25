package kafka

import (
	"sync"

	"github.com/IBM/sarama"
)

var msgPool = sync.Pool{
	New: func() any { return new(sarama.ProducerMessage) },
}

func newMessage(topic string, p []byte) *sarama.ProducerMessage {
	value := make(sarama.ByteEncoder, len(p))
	copy(value, p)

	msg := msgPool.Get().(*sarama.ProducerMessage) //nolint:revive // never
	msg.Topic = topic
	msg.Value = value

	return msg
}

func freeMessage(msg *sarama.ProducerMessage) {
	msgPool.Put(msg)
}
