package rest

import (
	"context"
	"errors"

	"github.com/rom8726/warden/internal/backend/dto"
	"github.com/rom8726/warden/internal/domain"
	generatedapi "github.com/rom8726/warden/internal/generated/server"
)

func (r *RestAPI) GetProjectReleaseAnalyticsDetails(
	ctx context.Context,
	params generatedapi.GetProjectReleaseAnalyticsDetailsParams,
) (generatedapi.GetProjectReleaseAnalyticsDetailsRes, error) {
	projectID := domain.ProjectID(params.ProjectID)
	version := params.Version

	details, err := r.analyticsUseCase.GetReleaseDetails(ctx, projectID, version, 10)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrEntityNotFound):
			return nil, &generatedapi.ErrorStatusCode{StatusCode: 404}
		default:
			return nil, &generatedapi.ErrorStatusCode{StatusCode: 500}
		}
	}

	result := dto.ToReleaseAnalyticsDetails(details)

	return &result, nil
}
