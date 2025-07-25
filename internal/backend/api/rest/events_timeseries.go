package rest

import (
	"context"
	"log/slog"

	"github.com/rom8726/warden/internal/backend/dto"
	"github.com/rom8726/warden/internal/domain"
	generatedapi "github.com/rom8726/warden/internal/generated/server"
)

func (r *RestAPI) GetEventsTimeseries(
	ctx context.Context,
	params generatedapi.GetEventsTimeseriesParams,
) (generatedapi.GetEventsTimeseriesRes, error) {
	period, err := dto.TimeseriesPeriodToDomainPeriod(params.Interval, params.Granularity)
	if err != nil {
		slog.Error("invalid period", "error", err)

		return nil, err
	}

	var projectID *domain.ProjectID
	if params.ProjectID.Set {
		domainProjectID := domain.ProjectID(params.ProjectID.Value)
		projectID = &domainProjectID
	}

	filter := domain.EventTimeseriesFilter{
		Period:    period,
		ProjectID: projectID,
		Levels:    nil,
		GroupBy:   domain.EventTimeseriesGroupLevel,
	}

	list, err := r.eventUseCase.Timeseries(ctx, &filter)
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
