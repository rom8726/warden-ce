package rest

import (
	"context"
	"errors"

	"github.com/rom8726/warden/internal/domain"
	generatedapi "github.com/rom8726/warden/internal/generated/server"
)

func (r *RestAPI) MarkNotificationAsRead(
	ctx context.Context,
	params generatedapi.MarkNotificationAsReadParams,
) (generatedapi.MarkNotificationAsReadRes, error) {
	notificationID := domain.UserNotificationID(params.NotificationID)

	err := r.userNotificationsUseCase.MarkAsRead(ctx, notificationID)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrEntityNotFound):
			return &generatedapi.ErrorNotFound{Error: generatedapi.ErrorNotFoundError{
				Message: generatedapi.NewOptString(err.Error()),
			}}, nil
		default:
			return nil, r.NewError(ctx, err)
		}
	}

	return &generatedapi.MarkNotificationAsReadNoContent{}, nil
}
