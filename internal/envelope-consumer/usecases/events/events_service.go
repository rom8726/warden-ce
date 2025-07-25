package events

import (
	"context"
	"fmt"
	"strconv"
	"time"

	eventcommon "github.com/rom8726/warden/internal/common/event"
	"github.com/rom8726/warden/internal/domain"
	"github.com/rom8726/warden/internal/envelope-consumer/contract"
	"github.com/rom8726/warden/pkg/db"
	"github.com/rom8726/warden/pkg/metrics"
)

type EventService struct {
	txManager              db.TxManager
	issueRepo              contract.IssuesRepository
	eventRepo              contract.EventRepository
	releaseRepo            contract.ReleaseRepository
	notificationsQueueRepo contract.NotificationsQueueRepository
	issueReleasesRepo      contract.IssueReleasesRepository
	cacheService           contract.CacheService
}

func New(
	txManager db.TxManager,
	issueRepo contract.IssuesRepository,
	eventRepo contract.EventRepository,
	releaseRepo contract.ReleaseRepository,
	notificationsQueueRepo contract.NotificationsQueueRepository,
	issueReleasesRepo contract.IssueReleasesRepository,
	cacheService contract.CacheService,
) *EventService {
	return &EventService{
		txManager:              txManager,
		issueRepo:              issueRepo,
		eventRepo:              eventRepo,
		releaseRepo:            releaseRepo,
		notificationsQueueRepo: notificationsQueueRepo,
		issueReleasesRepo:      issueReleasesRepo,
		cacheService:           cacheService,
	}
}

// ProcessEvent processes an event from the Sentry SDK.
//
//nolint:gocyclo // need refactoring
func (s *EventService) ProcessEvent(
	ctx context.Context,
	projectID domain.ProjectID,
	eventData map[string]any,
) (domain.EventID, error) {
	// Start timing for overall event processing
	start := time.Now()
	projectIDStr := strconv.FormatUint(uint64(projectID), 10)

	event, err := eventcommon.ParseEvent(eventData, projectID)
	if err != nil {
		return "", fmt.Errorf("parse event: %w", err)
	}

	err = s.txManager.ReadCommitted(ctx, func(ctx context.Context) error {
		// Store the event with fingerprints
		if err := s.eventRepo.StoreWithFingerprints(ctx, &event); err != nil {
			return fmt.Errorf("store event: %w", err)
		}

		issue := domain.IssueDTO{
			ProjectID:   projectID,
			Fingerprint: event.GroupHash,
			Source:      event.Source,
			Status:      domain.IssueStatusUnresolved,
			Title:       event.Message,
			Level:       event.Level,
			Platform:    event.Platform,
		}

		// Start timing for UpsertIssue operation
		upsertStart := time.Now()
		upsertRes, err := s.issueRepo.UpsertIssue(ctx, issue)

		// Record metrics for UpsertIssue operation
		metrics.UpsertIssueDuration.WithLabelValues(projectIDStr).Observe(time.Since(upsertStart).Seconds())

		if err != nil {
			metrics.UpsertIssueErrors.WithLabelValues(projectIDStr).Inc()

			return fmt.Errorf("upsert issue: %w", err)
		}

		// Record successful UpsertIssue operation
		status := "updated"
		if upsertRes.IsNew {
			status = "created"
		}
		metrics.UpsertIssueTotal.WithLabelValues(projectIDStr, status).Inc()

		// Create or get a release using cache
		releaseID, err := s.cacheService.GetOrCreateRelease(ctx, projectID, event.Release, s.releaseRepo)
		if err != nil {
			return fmt.Errorf("get or create release: %w", err)
		}

		// Create or get issue_release using cache
		err = s.cacheService.GetOrCreateIssueRelease(
			ctx,
			upsertRes.ID,
			releaseID,
			upsertRes.IsNew,
			s.issueReleasesRepo,
		)
		if err != nil {
			return fmt.Errorf("get or create issue release: %w", err)
		}

		if domain.IsNotifiableLevel(issue.Level) && (upsertRes.IsNew || upsertRes.WasReactivated) {
			err := s.notificationsQueueRepo.AddNotification(
				ctx,
				projectID,
				upsertRes.ID,
				issue.Level,
				upsertRes.IsNew, upsertRes.WasReactivated,
			)
			if err != nil {
				return fmt.Errorf("add notification: %w", err)
			}
		}

		return nil
	})
	if err != nil {
		return "", fmt.Errorf("store event and issue: %w", err)
	}

	// Record overall processing time
	metrics.ProcessingTime.WithLabelValues("event").Observe(time.Since(start).Seconds())
	metrics.EventsProcessed.WithLabelValues(projectIDStr).Inc()

	return event.ID, nil
}
