package rest

import (
	"context"
	"log/slog"

	wardencontext "github.com/rom8726/warden/internal/context"
	generatedapi "github.com/rom8726/warden/internal/generated/server"
)

func (r *RestAPI) Send2FACode(ctx context.Context) (generatedapi.Send2FACodeRes, error) {
	err := r.usersUseCase.Send2FACode(ctx, wardencontext.UserID(ctx), "disable")
	if err != nil {
		slog.Error("send 2fa code failed", "error", err)

		return nil, err
	}

	return &generatedapi.Send2FACodeNoContent{}, nil
}
