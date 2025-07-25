package rest

import (
	"context"
	"errors"
	"log/slog"

	"github.com/rom8726/warden/internal/domain"
	generatedapi "github.com/rom8726/warden/internal/generated/server"
)

func (r *RestAPI) DeleteUser(
	ctx context.Context,
	params generatedapi.DeleteUserParams,
) (generatedapi.DeleteUserRes, error) {
	userID := domain.UserID(params.UserID)
	err := r.usersUseCase.Delete(ctx, userID)
	if err != nil {
		slog.Error("delete user failed", "error", err, "user_id", userID)

		switch {
		case errors.Is(err, domain.ErrEntityNotFound):
			return &generatedapi.ErrorNotFound{
				Error: generatedapi.ErrorNotFoundError{
					Message: generatedapi.NewOptString(err.Error()),
				},
			}, nil
		case errors.Is(err, domain.ErrForbidden):
			return &generatedapi.ErrorPermissionDenied{
				Error: generatedapi.ErrorPermissionDeniedError{
					Message: generatedapi.NewOptString("Only superusers can create new users"),
				},
			}, nil
		}

		return nil, err
	}

	return &generatedapi.DeleteUserNoContent{}, nil
}
