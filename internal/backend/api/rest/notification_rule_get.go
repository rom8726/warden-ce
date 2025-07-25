package rest

import (
	"context"
	"errors"
	"log/slog"

	"github.com/rom8726/warden/internal/backend/dto"
	"github.com/rom8726/warden/internal/domain"
	generatedapi "github.com/rom8726/warden/internal/generated/server"
)

func (r *RestAPI) GetNotificationRule(
	ctx context.Context,
	params generatedapi.GetNotificationRuleParams,
) (generatedapi.GetNotificationRuleRes, error) {
	ruleID := domain.NotificationRuleID(params.RuleID)

	// Call the service
	rule, err := r.notificationsUseCase.GetNotificationRule(ctx, ruleID)
	if err != nil {
		slog.Error("get notification rule failed", "error", err, "rule_id", ruleID)

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
