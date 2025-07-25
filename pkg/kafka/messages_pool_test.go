package kafka

import (
	"testing"

	"github.com/IBM/sarama"
	"github.com/stretchr/testify/assert"
)

func TestNewMessage(t *testing.T) {
	t.Run("create message with data", func(t *testing.T) {
		topic := "test-topic"
		data := []byte("test message")

		msg := newMessage(topic, data)

		assert.NotNil(t, msg)
		assert.Equal(t, topic, msg.Topic)
		assert.Equal(t, sarama.ByteEncoder(data), msg.Value)
	})

	t.Run("create message with empty data", func(t *testing.T) {
		topic := "test-topic"
		data := []byte{}

		msg := newMessage(topic, data)

		assert.NotNil(t, msg)
		assert.Equal(t, topic, msg.Topic)
		assert.Equal(t, sarama.ByteEncoder(data), msg.Value)
	})

	t.Run("create message with nil data", func(t *testing.T) {
		topic := "test-topic"
		var data []byte

		msg := newMessage(topic, data)

		assert.NotNil(t, msg)
		assert.Equal(t, topic, msg.Topic)
		assert.Equal(t, sarama.ByteEncoder{}, msg.Value)
	})
}

func TestFreeMessage(t *testing.T) {
	t.Run("free message", func(t *testing.T) {
		msg := &sarama.ProducerMessage{
			Topic: "test-topic",
			Value: sarama.ByteEncoder("test message"),
		}

		// Это не должно паниковать
		freeMessage(msg)
	})

	t.Run("free nil message", func(t *testing.T) {
		// Это не должно паниковать
		freeMessage(nil)
	})
}
