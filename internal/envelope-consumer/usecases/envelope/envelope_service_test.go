package envelope

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/rom8726/warden/internal/domain"
	mockcontract "github.com/rom8726/warden/test_mocks/internal_/envelope-consumer/contract"
)

func TestNew(t *testing.T) {
	t.Parallel()

	// Create mocks
	mockEventUseCase := mockcontract.NewMockStoreEventUseCase(t)

	// Create a new EnvelopeService
	service := New(mockEventUseCase)

	// Verify the service was created correctly
	require.NotNil(t, service)
	require.Equal(t, mockEventUseCase, service.eventUseCase)
}

func TestProcessEnvelopeFromBytes_Success(t *testing.T) {
	t.Parallel()

	// Create mocks
	mockEventUseCase := mockcontract.NewMockStoreEventUseCase(t)

	// Set up the mock to return success when StoreEvent is called
	mockEventUseCase.EXPECT().StoreEvent(mock.Anything, domain.ProjectID(1), mock.Anything).
		Return(domain.EventID("event-123"), nil)

	// Create a new EnvelopeService
	service := New(mockEventUseCase)

	// Create envelope data
	envelopeData := `{"version": "1.0"}
{"type": "event", "length": 13}
{"data": "test"}`

	// Call ProcessEnvelopeFromBytes
	err := service.ProcessEnvelopeFromBytes(context.Background(), 1, []byte(envelopeData))

	// Verify no error
	require.NoError(t, err)

	// Verify the mock was called
	mockEventUseCase.AssertExpectations(t)
}

func TestProcessEnvelopeFromBytes_EmptyEnvelope(t *testing.T) {
	t.Parallel()

	// Create mocks
	mockEventUseCase := mockcontract.NewMockStoreEventUseCase(t)

	// Create a new EnvelopeService
	service := New(mockEventUseCase)

	// Call ProcessEnvelopeFromBytes with empty data
	err := service.ProcessEnvelopeFromBytes(context.Background(), 1, []byte{})

	// Verify the error
	require.Error(t, err)
	require.Contains(t, err.Error(), "empty envelope")

	// Verify the mock was not called
	mockEventUseCase.AssertNotCalled(t, "StoreEvent")
}

func TestProcessEnvelopeFromBytes_InvalidHeader(t *testing.T) {
	t.Parallel()

	// Create mocks
	mockEventUseCase := mockcontract.NewMockStoreEventUseCase(t)

	// Create a new EnvelopeService
	service := New(mockEventUseCase)

	// Create envelope data with an invalid header
	envelopeData := `invalid json
{"type": "event", "length": 13}
{"data": "test"}`

	// Call ProcessEnvelopeFromBytes
	err := service.ProcessEnvelopeFromBytes(context.Background(), 1, []byte(envelopeData))

	// Verify the error
	require.Error(t, err)
	require.Contains(t, err.Error(), "invalid envelope header")

	// Verify the mock was not called
	mockEventUseCase.AssertNotCalled(t, "StoreEvent")
}

func TestProcessEnvelopeFromBytes_InvalidItemHeader(t *testing.T) {
	t.Parallel()

	// Create mocks
	mockEventUseCase := mockcontract.NewMockStoreEventUseCase(t)

	// Create a new EnvelopeService
	service := New(mockEventUseCase)

	// Create envelope data with an invalid item header
	envelopeData := `{"version": "1.0"}
invalid json
{"data": "test"}`

	// Call ProcessEnvelopeFromBytes
	err := service.ProcessEnvelopeFromBytes(context.Background(), 1, []byte(envelopeData))

	// Verify no error (invalid item headers are skipped)
	require.NoError(t, err)

	// Verify the mock was not called
	mockEventUseCase.AssertNotCalled(t, "StoreEvent")
}

func TestProcessEnvelopeFromBytes_MissingTypeField(t *testing.T) {
	t.Parallel()

	// Create mocks
	mockEventUseCase := mockcontract.NewMockStoreEventUseCase(t)

	// Create a new EnvelopeService
	service := New(mockEventUseCase)

	// Create envelope data with a missing type field
	envelopeData := `{"version": "1.0"}
{"length": 13}
{"data": "test"}`

	// Call ProcessEnvelopeFromBytes
	err := service.ProcessEnvelopeFromBytes(context.Background(), 1, []byte(envelopeData))

	// Verify no error (items with missing type are skipped)
	require.NoError(t, err)

	// Verify the mock was not called
	mockEventUseCase.AssertNotCalled(t, "StoreEvent")
}

func TestProcessEnvelopeFromBytes_MissingLengthField(t *testing.T) {
	t.Parallel()

	// Create mocks
	mockEventUseCase := mockcontract.NewMockStoreEventUseCase(t)

	// Create a new EnvelopeService
	service := New(mockEventUseCase)

	// Create envelope data with missing length field
	envelopeData := `{"version": "1.0"}
{"type": "event"}
{"data": "test"}`

	// Call ProcessEnvelopeFromBytes
	err := service.ProcessEnvelopeFromBytes(context.Background(), 1, []byte(envelopeData))

	// Verify no error (items with missing length are skipped)
	require.NoError(t, err)

	// Verify the mock was not called
	mockEventUseCase.AssertNotCalled(t, "StoreEvent")
}

func TestProcessEnvelopeFromBytes_InvalidEventData(t *testing.T) {
	t.Parallel()

	// Create mocks
	mockEventUseCase := mockcontract.NewMockStoreEventUseCase(t)

	// Create a new EnvelopeService
	service := New(mockEventUseCase)

	// Create envelope data with invalid event data
	envelopeData := `{"version": "1.0"}
{"type": "event", "length": 13}
invalid json`

	// Call ProcessEnvelopeFromBytes
	err := service.ProcessEnvelopeFromBytes(context.Background(), 1, []byte(envelopeData))

	// Verify no error (invalid event data is skipped)
	require.NoError(t, err)

	// Verify the mock was not called
	mockEventUseCase.AssertNotCalled(t, "StoreEvent")
}

func TestProcessEnvelopeFromBytes_ProcessEventError(t *testing.T) {
	t.Parallel()

	// Create mocks
	mockEventUseCase := mockcontract.NewMockStoreEventUseCase(t)

	// Set up the mock to return an error when StoreEvent is called
	expectedError := errors.New("store event error")
	mockEventUseCase.EXPECT().StoreEvent(mock.Anything, domain.ProjectID(1), mock.Anything).
		Return(domain.EventID(""), expectedError)

	// Create a new EnvelopeService
	service := New(mockEventUseCase)

	// Create envelope data
	envelopeData := `{"version": "1.0"}
{"type": "event", "length": 13}
{"data": "test"}`

	// Call ProcessEnvelopeFromBytes
	err := service.ProcessEnvelopeFromBytes(context.Background(), 1, []byte(envelopeData))

	// Verify no error (process event errors are logged but don't fail the envelope)
	require.NoError(t, err)

	// Verify the mock was called
	mockEventUseCase.AssertExpectations(t)
}

func TestProcessEnvelopeFromBytes_UnsupportedItemType(t *testing.T) {
	t.Parallel()

	// Create mocks
	mockEventUseCase := mockcontract.NewMockStoreEventUseCase(t)

	// Create a new EnvelopeService
	service := New(mockEventUseCase)

	// Create envelope data with an unsupported item type
	envelopeData := `{"version": "1.0"}
{"type": "unsupported", "length": 13}
{"data": "test"}`

	// Call ProcessEnvelopeFromBytes
	err := service.ProcessEnvelopeFromBytes(context.Background(), 1, []byte(envelopeData))

	// Verify no error (unsupported types are skipped)
	require.NoError(t, err)

	// Verify the mock was not called
	mockEventUseCase.AssertNotCalled(t, "StoreEvent")
}

func TestProcessEnvelopeFromBytes_MultipleItems(t *testing.T) {
	t.Parallel()

	// Create mocks
	mockEventUseCase := mockcontract.NewMockStoreEventUseCase(t)

	// Set up the mock to return success for both StoreEvent calls
	mockEventUseCase.EXPECT().StoreEvent(mock.Anything, domain.ProjectID(1), mock.Anything).
		Return(domain.EventID("event-1"), nil).Once()
	mockEventUseCase.EXPECT().StoreEvent(mock.Anything, domain.ProjectID(1), mock.Anything).
		Return(domain.EventID("event-2"), nil).Once()

	// Create a new EnvelopeService
	service := New(mockEventUseCase)

	// Create envelope data with multiple events
	envelopeData := `{"version": "1.0"}
{"type": "event", "length": 13}
{"data": "test1"}
{"type": "event", "length": 13}
{"data": "test2"}`

	// Call ProcessEnvelopeFromBytes
	err := service.ProcessEnvelopeFromBytes(context.Background(), 1, []byte(envelopeData))

	// Verify no error
	require.NoError(t, err)

	// Verify the mock was called twice
	mockEventUseCase.AssertExpectations(t)
	mockEventUseCase.AssertNumberOfCalls(t, "StoreEvent", 2)
}
