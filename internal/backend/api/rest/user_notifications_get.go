package rest

import (
	"context"

	"github.com/rom8726/warden/internal/backend/dto"
	wardencontext "github.com/rom8726/warden/internal/context"
	generatedapi "github.com/rom8726/warden/internal/generated/server"
)

func (r *RestAPI) GetUserNotifications(
	ctx context.Context,
	params generatedapi.GetUserNotificationsParams,
) (generatedapi.GetUserNotificationsRes, error) {
	userID := wardencontext.UserID(ctx)

	limit := params.Limit.Or(50)
	offset := params.Offset.Or(0)

	notifications, err := r.userNotificationsUseCase.GetUserNotifications(ctx, userID, limit, offset)
	if err != nil {
		return nil, r.NewError(ctx, err)
	}

	dtoNotifications := make([]generatedapi.UserNotification, 0, len(notifications))
	for _, notification := range notifications {
		notif, err := dto.UserNotificationToDTO(notification)
		if err != nil {
			return nil, err
		}

		dtoNotifications = append(dtoNotifications, notif)
	}

	return &generatedapi.UserNotificationsResponse{
		Notifications: dtoNotifications,
		Total:         len(dtoNotifications),
	}, nil
}
