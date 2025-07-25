package events

import (
	"context"
	"fmt"

	"github.com/rom8726/warden/internal/backend/contract"
	"github.com/rom8726/warden/internal/domain"
)

type EventService struct {
	issueRepo contract.IssuesRepository
	eventRepo contract.EventRepository
}

func New(
	issueRepo contract.IssuesRepository,
	eventRepo contract.EventRepository,
) *EventService {
	return &EventService{
		issueRepo: issueRepo,
		eventRepo: eventRepo,
	}
}

func (s *EventService) Timeseries(
	ctx context.Context,
	filter *domain.EventTimeseriesFilter,
) ([]domain.Timeseries, error) {
	return s.eventRepo.Timeseries(ctx, filter)
}

func (s *EventService) IssueTimeseries(
	ctx context.Context,
	filter *domain.IssueEventsTimeseriesFilter,
) ([]domain.Timeseries, error) {
	issue, err := s.issueRepo.GetByID(ctx, filter.IssueID)
	if err != nil {
		return nil, fmt.Errorf("get issue: %w", err)
	}

	return s.eventRepo.IssueTimeseries(ctx, issue.Fingerprint, filter)
}
