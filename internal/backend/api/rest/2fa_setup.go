package rest

import (
	"context"
	"errors"
	"log/slog"

	wardencontext "github.com/rom8726/warden/internal/context"
	"github.com/rom8726/warden/internal/domain"
	generatedapi "github.com/rom8726/warden/internal/generated/server"
)

func (r *RestAPI) Setup2FA(ctx context.Context) (generatedapi.Setup2FARes, error) {
	userID := wardencontext.UserID(ctx)
	secret, qrURL, qrImage, err := r.usersUseCase.Setup2FA(ctx, userID)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrUserNotFound):
			return &generatedapi.ErrorUnauthorized{Error: generatedapi.ErrorUnauthorizedError{
				Message: generatedapi.NewOptString("unauthorized"),
			}}, nil
		default:
			slog.Error("failed to setup 2FA", "error", err)

			return nil, r.NewError(ctx, err)
		}
	}
	resp := &generatedapi.TwoFASetupResponse{
		Secret:  secret,
		QrURL:   qrURL,
		QrImage: qrImage, // base64 PNG
	}

	return resp, nil
}
