package rest

import (
	"context"
	"errors"
	"log/slog"

	"github.com/rom8726/warden/internal/backend/dto"
	"github.com/rom8726/warden/internal/domain"
	generatedapi "github.com/rom8726/warden/internal/generated/server"
)

func (r *RestAPI) GetNotificationSetting(
	ctx context.Context,
	params generatedapi.GetNotificationSettingParams,
) (generatedapi.GetNotificationSettingRes, error) {
	settingID := domain.NotificationSettingID(params.SettingID)

	// Call the service
	setting, err := r.notificationsUseCase.GetNotificationSetting(ctx, settingID)
	if err != nil {
		slog.Error("get notification setting failed", "error", err, "setting_id", settingID)

		if errors.Is(err, domain.ErrEntityNotFound) {
			return &generatedapi.ErrorNotFound{
				Error: generatedapi.ErrorNotFoundError{
					Message: generatedapi.NewOptString(err.Error()),
				},
			}, nil
		}

		return nil, err
	}

	// Convert a domain model to an API model
	apiSetting := dto.DomainNotificationSettingToAPI(setting)

	return &apiSetting, nil
}
