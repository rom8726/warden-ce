package rest

import (
	"context"
	"log/slog"

	"github.com/pkg/errors"

	"github.com/rom8726/warden/internal/domain"
	generatedapi "github.com/rom8726/warden/internal/generated/server"
)

func (r *RestAPI) Verify2FA(
	ctx context.Context,
	req *generatedapi.TwoFAVerifyRequest,
) (generatedapi.Verify2FARes, error) {
	accessToken, refreshToken, expiresIn, err := r.usersUseCase.Verify2FA(ctx, req.Code, req.SessionID)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrInvalid2FACode):
			return &generatedapi.ErrorBadRequest{Error: generatedapi.ErrorBadRequestError{
				Message: generatedapi.NewOptString("invalid code"),
			}}, nil
		case errors.Is(err, domain.ErrInvalidToken), errors.Is(err, domain.ErrUserNotFound):
			return &generatedapi.ErrorUnauthorized{Error: generatedapi.ErrorUnauthorizedError{
				Message: generatedapi.NewOptString("unauthorized"),
			}}, nil
		case errors.Is(err, domain.ErrTooMany2FAAttempts):
			return &generatedapi.ErrorTooManyRequests{Error: generatedapi.ErrorTooManyRequestsError{
				Message: generatedapi.NewOptString("too many attempts. try again later"),
			}}, nil
		default:
			slog.Error("failed to verify 2FA", "error", err)

			return nil, r.NewError(ctx, err)
		}
	}
	resp := &generatedapi.TwoFAVerifyResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    expiresIn,
	}

	return resp, nil
}
