package rest

import (
	"context"
	"errors"
	"log/slog"

	"github.com/rom8726/warden/internal/backend/dto"
	"github.com/rom8726/warden/internal/domain"
	generatedapi "github.com/rom8726/warden/internal/generated/server"
)

func (r *RestAPI) UpdateNotificationRule(
	ctx context.Context,
	req *generatedapi.UpdateNotificationRuleRequest,
	params generatedapi.UpdateNotificationRuleParams,
) (generatedapi.UpdateNotificationRuleRes, error) {
	ruleID := domain.NotificationRuleID(params.RuleID)

	// Get the existing rule
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

	// Update the rule with the values from the request
	updatedRule := dto.UpdateNotificationRuleFromRequest(rule, req)

	// Call the service to update the rule
	err = r.notificationsUseCase.UpdateNotificationRule(ctx, updatedRule)
	if err != nil {
		slog.Error("update notification rule failed", "error", err, "rule_id", ruleID)

		return nil, err
	}

	// Get the updated rule
	rule, err = r.notificationsUseCase.GetNotificationRule(ctx, ruleID)
	if err != nil {
		slog.Error("get updated notification rule failed", "error", err, "rule_id", ruleID)

		return nil, err
	}

	// Convert domain model to API model
	apiRule := dto.DomainNotificationRuleToAPI(rule)

	return &apiRule, nil
}
