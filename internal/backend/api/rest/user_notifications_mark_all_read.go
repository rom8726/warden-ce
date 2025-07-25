package rest

import (
	"context"

	wardencontext "github.com/rom8726/warden/internal/context"
	generatedapi "github.com/rom8726/warden/internal/generated/server"
)

func (r *RestAPI) MarkAllNotificationsAsRead(
	ctx context.Context,
) (generatedapi.MarkAllNotificationsAsReadRes, error) {
	userID := wardencontext.UserID(ctx)

	err := r.userNotificationsUseCase.MarkAllAsRead(ctx, userID)
	if err != nil {
		return nil, r.NewError(ctx, err)
	}

	return &generatedapi.MarkAllNotificationsAsReadNoContent{}, nil
}
