package rest

import (
	"context"
	"errors"
	"log/slog"

	"github.com/rom8726/warden/internal/domain"
	generatedapi "github.com/rom8726/warden/internal/generated/server"
)

func (r *RestAPI) DeleteNotificationRule(
	ctx context.Context,
	params generatedapi.DeleteNotificationRuleParams,
) (generatedapi.DeleteNotificationRuleRes, error) {
	ruleID := domain.NotificationRuleID(params.RuleID)

	// Call the service to delete the rule
	err := r.notificationsUseCase.DeleteNotificationRule(ctx, ruleID)
	if err != nil {
		slog.Error("delete notification rule failed", "error", err, "rule_id", ruleID)

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

	// Return a success response (204 No Content)
	return &generatedapi.DeleteNotificationRuleNoContent{}, nil
}
