package rest

import (
	"context"
	"log/slog"

	"github.com/rom8726/warden/internal/backend/dto"
	"github.com/rom8726/warden/internal/domain"
	generatedapi "github.com/rom8726/warden/internal/generated/server"
)

func (r *RestAPI) GetProjectIssueEventsTimeseries(
	ctx context.Context,
	params generatedapi.GetProjectIssueEventsTimeseriesParams,
) (generatedapi.GetProjectIssueEventsTimeseriesRes, error) {
	period, err := dto.TimeseriesPeriodToDomainPeriod(params.Interval, params.Granularity)
	if err != nil {
		slog.Error("invalid period", "error", err)

		return nil, err
	}

	projectID := domain.ProjectID(params.ProjectID)
	issueID := domain.IssueID(params.IssueID)

	filter := domain.IssueEventsTimeseriesFilter{
		Period:    period,
		ProjectID: projectID,
		IssueID:   issueID,
		Levels:    nil,
		GroupBy:   domain.EventTimeseriesGroupLevel,
	}

	list, err := r.eventUseCase.IssueTimeseries(ctx, &filter)
	if err != nil {
		slog.Error("get project stats failed", "error", err)

		return nil, err
	}

	items := make([]generatedapi.TimeseriesData, 0, len(list))
	for i := range list {
		elem := list[i]
		items = append(items, dto.DomainTimeseriesToAPI(elem, params.Interval, params.Granularity))
	}

	resp := generatedapi.TimeseriesResponse(items)

	return &resp, nil
}
