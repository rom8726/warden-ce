package contract

import (
	"context"
	"io"

	"github.com/rom8726/warden/internal/domain"
	"github.com/rom8726/warden/pkg/kafka"
)

type EnvelopeUseCase interface {
	ReceiveEnvelope(ctx context.Context, projectID domain.ProjectID, data io.Reader) error
}

type StoreEventUseCase interface {
	StoreEvent(
		ctx context.Context,
		projectID domain.ProjectID,
		req map[string]any,
	) (domain.EventID, error)
}

type ProjectsUseCase interface {
	ValidateProjectKey(
		ctx context.Context,
		projectID domain.ProjectID,
		key string,
	) (bool, error)
}

type EnvelopProducer interface {
	SendEnvelope(ctx context.Context, projectID domain.ProjectID, data []byte) error
}

type EventProducer interface {
	SendStoreEvent(
		ctx context.Context,
		projectID domain.ProjectID,
		eventID domain.EventID,
		eventData map[string]any,
	) error
}

type ProjectsRepository interface {
	ValidateProjectKey(ctx context.Context, projectID domain.ProjectID, key string) (bool, error)
	GetProjectIDs(ctx context.Context) ([]domain.ProjectID, error)
}

type TopicProducerCreator interface {
	Create(topic string) kafka.DataProducer
}
