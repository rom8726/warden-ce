package rest

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"github.com/rom8726/warden/internal/domain"
	generatedapi "github.com/rom8726/warden/internal/generated/server"
)

func (r *RestAPI) RefreshToken(
	ctx context.Context,
	req *generatedapi.RefreshTokenRequest,
) (generatedapi.RefreshTokenRes, error) {
	accessToken, refreshToken, err := r.usersUseCase.LoginReissue(ctx, req.RefreshToken)
	if err != nil {
		slog.Error("login failed", "error", err)

		if errors.Is(err, domain.ErrInvalidToken) {
			return &generatedapi.ErrorUnauthorized{Error: generatedapi.ErrorUnauthorizedError{
				Message: generatedapi.NewOptString(domain.ErrInvalidToken.Error()),
			}}, nil
		}

		if errors.Is(err, domain.ErrEntityNotFound) || errors.Is(err, domain.ErrInactiveUser) {
			return &generatedapi.ErrorUnauthorized{Error: generatedapi.ErrorUnauthorizedError{
				Message: generatedapi.NewOptString(err.Error()),
			}}, nil
		}

		return nil, err
	}

	return &generatedapi.RefreshTokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    int(time.Now().Add(r.tokenizer.AccessTokenTTL()).Unix()),
	}, nil
}
