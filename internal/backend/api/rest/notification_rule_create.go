package rest

import (
	"context"
	"errors"
	"log/slog"

	"github.com/rom8726/warden/internal/backend/dto"
	"github.com/rom8726/warden/internal/domain"
	generatedapi "github.com/rom8726/warden/internal/generated/server"
)

func (r *RestAPI) CreateNotificationRule(
	ctx context.Context,
	req *generatedapi.CreateNotificationRuleRequest,
	params generatedapi.CreateNotificationRuleParams,
) (generatedapi.CreateNotificationRuleRes, error) {
	settingID := domain.NotificationSettingID(params.SettingID)

	// Convert request to domain DTO
	ruleDTO := dto.MakeNotificationRuleDTO(req, settingID)

	// Call the service
	rule, err := r.notificationsUseCase.CreateNotificationRule(ctx, ruleDTO)
	if err != nil {
		slog.Error("create notification rule failed", "error", err, "setting_id", settingID)

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

	// Convert a domain model to an API model
	apiRule := dto.DomainNotificationRuleToAPI(rule)

	return &apiRule, nil
}
