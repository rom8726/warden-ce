package rest

import (
	"context"
	"log/slog"

	"github.com/rom8726/warden/internal/domain"
	generatedapi "github.com/rom8726/warden/internal/generated/server"
)

func (r *RestAPI) SendTestNotification(
	ctx context.Context,
	params generatedapi.SendTestNotificationParams,
) (generatedapi.SendTestNotificationRes, error) {
	err := r.notificationsUseCase.SendTestNotification(
		ctx,
		domain.ProjectID(params.ProjectID),
		domain.NotificationSettingID(params.SettingID),
	)
	if err != nil {
		slog.Error("failed to send test notification", "error", err)

		return nil, err
	}

	return &generatedapi.SendTestNotificationNoContent{}, nil
}
