package rest

import (
	"context"
	"log/slog"

	"github.com/pkg/errors"

	wardencontext "github.com/rom8726/warden/internal/context"
	"github.com/rom8726/warden/internal/domain"
	generatedapi "github.com/rom8726/warden/internal/generated/server"
)

func (r *RestAPI) Confirm2FA(
	ctx context.Context,
	req *generatedapi.TwoFAConfirmRequest,
) (generatedapi.Confirm2FARes, error) {
	userID := wardencontext.UserID(ctx)
	err := r.usersUseCase.Confirm2FA(ctx, userID, req.Code)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrInvalid2FACode):
			return &generatedapi.ErrorBadRequest{Error: generatedapi.ErrorBadRequestError{
				Message: generatedapi.NewOptString("invalid code"),
			}}, nil
		case errors.Is(err, domain.ErrUserNotFound):
			return &generatedapi.ErrorUnauthorized{Error: generatedapi.ErrorUnauthorizedError{
				Message: generatedapi.NewOptString("unauthorized"),
			}}, nil
		case errors.Is(err, domain.ErrTooMany2FAAttempts):
			return &generatedapi.ErrorTooManyRequests{Error: generatedapi.ErrorTooManyRequestsError{
				Message: generatedapi.NewOptString("too many attempts. try again later"),
			}}, nil
		default:
			slog.Error("failed to confirm 2FA", "error", err)

			return nil, r.NewError(ctx, err)
		}
	}

	return &generatedapi.Confirm2FANoContent{}, nil
}
