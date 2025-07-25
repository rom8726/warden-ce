package envelope

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/rom8726/warden/internal/domain"
	mockcontract "github.com/rom8726/warden/test_mocks/internal_/ingest-server/contract"
)

func TestReceiveEnvelope_Success(t *testing.T) {
	t.Parallel()

	// Create mocks
	mockEnvelopProducer := mockcontract.NewMockEnvelopProducer(t)

	// Set up the mock to return success when SendEnvelope is called
	mockEnvelopProducer.EXPECT().SendEnvelope(mock.Anything, domain.ProjectID(1), mock.Anything).
		Return(nil)

	// Create a new EnvelopeService
	service := New(mockEnvelopProducer)

	// Create a reader with valid envelope data
	envelopeData := `{"version": "1.0"}
{"type": "event", "length": 13}
{"data": "test"}`
	reader := strings.NewReader(envelopeData)

	// Call ReceiveEnvelope
	err := service.ReceiveEnvelope(context.Background(), 1, reader)

	// Verify no error
	require.NoError(t, err)

	// Verify the mock was called with correct data
	mockEnvelopProducer.AssertExpectations(t)
	mockEnvelopProducer.AssertCalled(t, "SendEnvelope", mock.Anything, domain.ProjectID(1), []byte(envelopeData))
}

func TestReceiveEnvelope_ProducerError(t *testing.T) {
	t.Parallel()

	// Create mocks
	mockEnvelopProducer := mockcontract.NewMockEnvelopProducer(t)

	// Set up the mock to return an error when SendEnvelope is called
	expectedError := errors.New("kafka send error")
	mockEnvelopProducer.EXPECT().SendEnvelope(mock.Anything, domain.ProjectID(1), mock.Anything).
		Return(expectedError)

	// Create a new EnvelopeService
	service := New(mockEnvelopProducer)

	// Create a reader with valid envelope data
	envelopeData := `{"version": "1.0"}
{"type": "event", "length": 13}
{"data": "test"}`
	reader := strings.NewReader(envelopeData)

	// Call ReceiveEnvelope
	err := service.ReceiveEnvelope(context.Background(), 1, reader)

	// Verify the error
	require.Error(t, err)
	require.Contains(t, err.Error(), "failed to send envelope to Kafka")

	// Verify the mock was called
	mockEnvelopProducer.AssertExpectations(t)
}

func TestReceiveEnvelope_ReadError(t *testing.T) {
	t.Parallel()

	// Create mocks
	mockEnvelopProducer := mockcontract.NewMockEnvelopProducer(t)

	// Create a new EnvelopeService
	service := New(mockEnvelopProducer)

	// Create a reader that will cause an error (closed reader)
	reader := strings.NewReader("")
	// Simulate read error by using a context that's already cancelled
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	// Call ReceiveEnvelope
	err := service.ReceiveEnvelope(ctx, 1, reader)

	// Verify the error
	require.Error(t, err)
	require.Contains(t, err.Error(), "empty envelope")

	// Verify the mock was not called
	mockEnvelopProducer.AssertNotCalled(t, "SendEnvelope")
}
