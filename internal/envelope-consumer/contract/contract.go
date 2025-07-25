package contract

import (
	"context"

	"github.com/rom8726/warden/internal/domain"
)

type EnvelopeUseCase interface {
	ProcessEnvelopeFromBytes(ctx context.Context, projectID domain.ProjectID, data []byte) error
}

type DataConsumer interface {
	Consume(ctx context.Context) <-chan []byte
	Close() error
}

type StoreEventUseCase interface {
	StoreEvent(
		ctx context.Context,
		projectID domain.ProjectID,
		req map[string]any,
	) (domain.EventID, error)
}

// EventUseCase handles event processing.
type EventUseCase interface {
	ProcessEvent(
		ctx context.Context,
		projectID domain.ProjectID,
		eventData map[string]any,
	) (domain.EventID, error)
}

type IssuesRepository interface {
	UpsertIssue(ctx context.Context, issue domain.IssueDTO) (domain.IssueUpsertResult, error)
}

type EventRepository interface {
	StoreWithFingerprints(ctx context.Context, event *domain.Event) error
}

type ReleaseRepository interface {
	GetByID(ctx context.Context, id domain.ReleaseID) (domain.Release, error)
	GetByProjectAndVersion(ctx context.Context, projectID domain.ProjectID, version string) (domain.Release, error)
	Create(ctx context.Context, release domain.ReleaseDTO) (domain.ReleaseID, error)
	ListByProject(ctx context.Context, projectID domain.ProjectID) ([]domain.Release, error)
}

type NotificationsQueueRepository interface {
	AddNotification(
		ctx context.Context,
		projectID domain.ProjectID,
		issueID domain.IssueID,
		level domain.IssueLevel,
		isNew, wasReactivated bool,
	) error
}

type IssueReleasesRepository interface {
	Create(ctx context.Context, issueID domain.IssueID, releaseID domain.ReleaseID, firstSeenIn bool) error
}

type EventProducer interface {
	SendStoreEvent(
		ctx context.Context,
		projectID domain.ProjectID,
		eventID domain.EventID,
		eventData map[string]any,
	) error
}
