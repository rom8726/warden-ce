package rest

import (
	"context"
	"errors"
	"log/slog"

	"github.com/rom8726/warden/internal/backend/dto"
	"github.com/rom8726/warden/internal/domain"
	generatedapi "github.com/rom8726/warden/internal/generated/server"
)

func (r *RestAPI) ListUsers(ctx context.Context) (generatedapi.ListUsersRes, error) {
	users, err := r.usersUseCase.List(ctx)
	if err != nil {
		slog.Error("list users failed", "error", err)

		if errors.Is(err, domain.ErrForbidden) {
			return &generatedapi.ErrorPermissionDenied{
				Error: generatedapi.ErrorPermissionDeniedError{
					Message: generatedapi.NewOptString("Only superusers can list users"),
				},
			}, nil
		}

		return nil, err
	}

	resp := dto.DomainUsersToAPI(users)
	listResp := generatedapi.ListUsersResponse(resp)

	return &listResp, nil
}
