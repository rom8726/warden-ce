package rest

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"github.com/rom8726/warden/internal/domain"
	generatedapi "github.com/rom8726/warden/internal/generated/server"
)

func (r *RestAPI) Login(ctx context.Context, req *generatedapi.LoginRequest) (generatedapi.LoginRes, error) {
	accessToken, refreshToken, sessionID, isTmpPwd, err := r.usersUseCase.Login(ctx, req.Username, req.Password)
	if err != nil {
		if errors.Is(err, domain.ErrInvalidCredentials) || errors.Is(err, domain.ErrInactiveUser) {
			return &generatedapi.ErrorInvalidCredentials{Error: generatedapi.ErrorInvalidCredentialsError{
				Message: generatedapi.NewOptString(err.Error()),
			}}, nil
		}

		if errors.Is(err, domain.ErrTwoFARequired) {
			return &generatedapi.Error2FARequired{Error: generatedapi.Error2FARequiredError{
				Code:      "2fa_required",
				SessionID: sessionID,
				Message:   "2FA required",
			}}, nil
		}

		slog.Error("login failed", "error", err)

		return nil, err
	}

	return &generatedapi.LoginResponse{
		AccessToken:   accessToken,
		RefreshToken:  refreshToken,
		ExpiresIn:     int(time.Now().Add(r.tokenizer.AccessTokenTTL()).Unix()),
		IsTmpPassword: isTmpPwd,
	}, nil
}
