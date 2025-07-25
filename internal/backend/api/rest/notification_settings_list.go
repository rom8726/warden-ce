package rest

import (
	"context"
	"errors"
	"log/slog"

	"github.com/rom8726/warden/internal/backend/dto"
	"github.com/rom8726/warden/internal/domain"
	generatedapi "github.com/rom8726/warden/internal/generated/server"
)

func (r *RestAPI) ListNotificationSettings(
	ctx context.Context,
	params generatedapi.ListNotificationSettingsParams,
) (generatedapi.ListNotificationSettingsRes, error) {
	projectID := domain.ProjectID(params.ProjectID)

	// Call the service to list settings
	settings, err := r.notificationsUseCase.ListNotificationSettings(ctx, projectID)
	if err != nil {
		slog.Error("list notification settings failed", "error", err, "project_id", projectID)

		switch {
		case errors.Is(err, domain.ErrEntityNotFound):
			return &generatedapi.ErrorNotFound{
				Error: generatedapi.ErrorNotFoundError{
					Message: generatedapi.NewOptString(err.Error()),
				},
			}, nil
		default:
			return nil, err
		}
	}

	// Convert domain models to API response
	response := dto.MakeListNotificationSettingsResponse(settings)

	return &response, nil
}
