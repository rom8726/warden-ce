package rest

import (
	"context"
	"errors"
	"log/slog"

	"github.com/rom8726/warden/internal/backend/dto"
	"github.com/rom8726/warden/internal/domain"
	generatedapi "github.com/rom8726/warden/internal/generated/server"
)

func (r *RestAPI) CreateNotificationSetting(
	ctx context.Context,
	req *generatedapi.CreateNotificationSettingRequest,
	params generatedapi.CreateNotificationSettingParams,
) (generatedapi.CreateNotificationSettingRes, error) {
	projectID := domain.ProjectID(params.ProjectID)

	// Convert request to domain DTO
	settingDTO := dto.MakeNotificationSettingDTO(req, projectID)

	// Call the service
	setting, err := r.notificationsUseCase.CreateNotificationSetting(ctx, settingDTO)
	if err != nil {
		slog.Error("create notification setting failed", "error", err, "project_id", projectID)

		switch {
		case errors.Is(err, domain.ErrEntityNotFound):
			return &generatedapi.ErrorNotFound{
				Error: generatedapi.ErrorNotFoundError{
					Message: generatedapi.NewOptString(err.Error()),
				},
			}, nil
		}

		return nil, err
	}

	// Convert domain model to API model
	apiSetting := dto.DomainNotificationSettingToAPI(setting)

	return &apiSetting, nil
}
