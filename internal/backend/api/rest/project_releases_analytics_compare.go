package rest

import (
	"context"
	"errors"

	"github.com/rom8726/warden/internal/backend/dto"
	"github.com/rom8726/warden/internal/domain"
	generatedapi "github.com/rom8726/warden/internal/generated/server"
)

func (r *RestAPI) CompareProjectReleasesAnalytics(
	ctx context.Context,
	req *generatedapi.CompareProjectReleasesAnalyticsReq,
	params generatedapi.CompareProjectReleasesAnalyticsParams,
) (generatedapi.CompareProjectReleasesAnalyticsRes, error) {
	projectID := domain.ProjectID(params.ProjectID)
	base := req.BaseVersion
	target := req.TargetVersion

	comp, err := r.analyticsUseCase.CompareReleases(ctx, projectID, base, target)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrEntityNotFound):
			return nil, &generatedapi.ErrorStatusCode{StatusCode: 404}
		default:
			return nil, &generatedapi.ErrorStatusCode{StatusCode: 500}
		}
	}

	result := dto.ToReleaseComparison(comp)

	return &result, nil
}
