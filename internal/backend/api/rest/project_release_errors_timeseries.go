package rest

import (
	"context"
	"errors"

	"github.com/rom8726/warden/internal/backend/dto"
	"github.com/rom8726/warden/internal/domain"
	generatedapi "github.com/rom8726/warden/internal/generated/server"
)

func (r *RestAPI) GetProjectReleaseErrorsTimeseries(
	ctx context.Context,
	params generatedapi.GetProjectReleaseErrorsTimeseriesParams,
) (generatedapi.GetProjectReleaseErrorsTimeseriesRes, error) {
	projectID := domain.ProjectID(params.ProjectID)
	release := params.Release

	interval := params.Interval
	granularity := params.Granularity

	period, err := dto.TimeseriesPeriodToDomainPeriod(interval, granularity)
	if err != nil {
		return nil, &generatedapi.ErrorStatusCode{StatusCode: 400}
	}

	var levels []domain.IssueLevel
	if params.Level.Set {
		levels = []domain.IssueLevel{domain.IssueLevel(params.Level.Value)}
	}
	groupBy := domain.EventTimeseriesGroupNone
	if params.GroupBy.Set {
		switch params.GroupBy.Value {
		case "level":
			groupBy = domain.EventTimeseriesGroupLevel
		}
	}

	series, err := r.analyticsUseCase.GetErrorsByTime(
		ctx,
		projectID,
		release,
		period.Interval,
		period.Granularity,
		levels,
		groupBy,
	)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrEntityNotFound):
			return nil, &generatedapi.ErrorStatusCode{StatusCode: 404}
		default:
			return nil, &generatedapi.ErrorStatusCode{StatusCode: 500}
		}
	}

	resp := generatedapi.TimeseriesResponse(dto.ToTimeseriesResponse(series, interval, granularity))

	return &resp, nil
}
