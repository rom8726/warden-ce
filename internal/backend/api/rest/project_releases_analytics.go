package rest

import (
	"context"
	"errors"

	"github.com/rom8726/warden/internal/backend/dto"
	"github.com/rom8726/warden/internal/domain"
	generatedapi "github.com/rom8726/warden/internal/generated/server"
)

func (r *RestAPI) GetProjectReleasesAnalytics(
	ctx context.Context,
	params generatedapi.GetProjectReleasesAnalyticsParams,
) (generatedapi.GetProjectReleasesAnalyticsRes, error) {
	projectID := domain.ProjectID(params.ProjectID)

	releases, statsMap, err := r.analyticsUseCase.ListReleases(ctx, projectID)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrEntityNotFound):
			return nil, &generatedapi.ErrorStatusCode{StatusCode: 404}
		default:
			return nil, &generatedapi.ErrorStatusCode{StatusCode: 500}
		}
	}

	result := make([]generatedapi.ReleaseAnalyticsSummary, 0, len(releases))
	for _, rel := range releases {
		stats, ok := statsMap[rel.Version]
		if !ok {
			continue
		}
		result = append(result, dto.ToReleaseAnalyticsSummary(rel, stats))
	}

	resp := generatedapi.GetProjectReleasesAnalyticsOKApplicationJSON(result)

	return &resp, nil
}
