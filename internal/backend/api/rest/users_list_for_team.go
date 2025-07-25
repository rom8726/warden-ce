package rest

import (
	"context"
	"errors"
	"log/slog"

	"github.com/rom8726/warden/internal/backend/dto"
	"github.com/rom8726/warden/internal/domain"
	generatedapi "github.com/rom8726/warden/internal/generated/server"
)

func (r *RestAPI) ListUsersForTeam(
	ctx context.Context,
	params generatedapi.ListUsersForTeamParams,
) (generatedapi.ListUsersForTeamRes, error) {
	users, err := r.usersUseCase.ListForTeamAdmin(ctx, domain.TeamID(params.TeamID))
	if err != nil {
		slog.Error("list users for team failed", "error", err)

		if errors.Is(err, domain.ErrForbidden) {
			return &generatedapi.ErrorPermissionDenied{
				Error: generatedapi.ErrorPermissionDeniedError{
					Message: generatedapi.NewOptString("Only superusers and team admins\\owners can list users"),
				},
			}, nil
		}

		return nil, err
	}

	resp := dto.DomainUsersToAPI(users)
	listResp := generatedapi.ListUsersResponse(resp)

	return &listResp, nil
}
