package rest

import (
	"context"

	wardencontext "github.com/rom8726/warden/internal/context"
	generatedapi "github.com/rom8726/warden/internal/generated/server"
)

func (r *RestAPI) GetUnreadNotificationsCount(
	ctx context.Context,
) (generatedapi.GetUnreadNotificationsCountRes, error) {
	userID := wardencontext.UserID(ctx)

	count, err := r.userNotificationsUseCase.GetUnreadCount(ctx, userID)
	if err != nil {
		return nil, r.NewError(ctx, err)
	}

	return &generatedapi.UnreadCountResponse{
		Count: count,
	}, nil
}
