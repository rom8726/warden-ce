package rest

import (
	"context"
	"errors"
	"log/slog"

	"github.com/rom8726/warden/internal/backend/dto"
	"github.com/rom8726/warden/internal/domain"
	generatedapi "github.com/rom8726/warden/internal/generated/server"
)

func (r *RestAPI) ListNotificationRules(
	ctx context.Context,
	params generatedapi.ListNotificationRulesParams,
) (generatedapi.ListNotificationRulesRes, error) {
	settingID := domain.NotificationSettingID(params.SettingID)

	// Call the service to list rules
	rules, err := r.notificationsUseCase.ListNotificationRules(ctx, settingID)
	if err != nil {
		slog.Error("list notification rules failed", "error", err, "setting_id", settingID)

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
	response := dto.MakeListNotificationRulesResponse(rules)

	return &response, nil
}
