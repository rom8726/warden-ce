package kafka

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestNewTopicProducer(t *testing.T) {
	t.Run("create topic producer", func(t *testing.T) {
		mockProd := &mockProducer{}
		topic := "test-topic"

		topicProducer := NewTopicProducer(mockProd, topic)

		assert.NotNil(t, topicProducer)
		assert.Equal(t, mockProd, topicProducer.producer)
		assert.Equal(t, topic, topicProducer.topic)
	})
}

func TestTopicProducer_Produce(t *testing.T) {
	t.Run("produce message successfully", func(t *testing.T) {
		mockProd := &mockProducer{}
		topic := "test-topic"
		data := []byte("test message")

		mockProd.On("Produce", mock.Anything, topic, data).Return(nil)

		topicProducer := NewTopicProducer(mockProd, topic)

		err := topicProducer.Produce(context.Background(), data)

		assert.NoError(t, err)
		mockProd.AssertExpectations(t)
	})

	t.Run("produce message with error", func(t *testing.T) {
		mockProd := &mockProducer{}
		topic := "test-topic"
		data := []byte("test message")
		expectedErr := assert.AnError

		mockProd.On("Produce", mock.Anything, topic, data).Return(expectedErr)

		topicProducer := NewTopicProducer(mockProd, topic)

		err := topicProducer.Produce(context.Background(), data)

		assert.Equal(t, expectedErr, err)
		mockProd.AssertExpectations(t)
	})

	t.Run("produce empty message", func(t *testing.T) {
		mockProd := &mockProducer{}
		topic := "test-topic"
		data := []byte{}

		mockProd.On("Produce", mock.Anything, topic, data).Return(nil)

		topicProducer := NewTopicProducer(mockProd, topic)

		err := topicProducer.Produce(context.Background(), data)

		assert.NoError(t, err)
		mockProd.AssertExpectations(t)
	})
}

// mockProducer is a mock implementation of Producer interface for testing
type mockProducer struct {
	mock.Mock
}

func (m *mockProducer) Produce(ctx context.Context, topic string, data []byte) error {
	args := m.Called(ctx, topic, data)
	return args.Error(0)
}
