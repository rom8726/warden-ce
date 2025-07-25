package rest

import (
	"context"
	"errors"
	"log/slog"

	"github.com/rom8726/warden/internal/domain"
	generatedapi "github.com/rom8726/warden/internal/generated/server"
)

func (r *RestAPI) RemoveTeamMember(
	ctx context.Context,
	params generatedapi.RemoveTeamMemberParams,
) (generatedapi.RemoveTeamMemberRes, error) {
	teamID := domain.TeamID(params.TeamID)
	userID := domain.UserID(params.UserID)

	err := r.teamsUseCase.RemoveMemberWithChecks(ctx, teamID, userID)
	if err != nil {
		slog.Error("remove team member failed", "error", err, "team_id", teamID, "user_id", userID)

		switch {
		case errors.Is(err, domain.ErrEntityNotFound):
			return &generatedapi.ErrorNotFound{
				Error: generatedapi.ErrorNotFoundError{
					Message: generatedapi.NewOptString(domain.ErrEntityNotFound.Error()),
				},
			}, nil
		case errors.Is(err, domain.ErrForbidden):
			return &generatedapi.ErrorPermissionDenied{
				Error: generatedapi.ErrorPermissionDeniedError{
					Message: generatedapi.NewOptString(domain.ErrForbidden.Error()),
				},
			}, nil
		case errors.Is(err, domain.ErrLastOwner):
			return &generatedapi.ErrorPermissionDenied{
				Error: generatedapi.ErrorPermissionDeniedError{
					Message: generatedapi.NewOptString(domain.ErrLastOwner.Error()),
				},
			}, nil
		}

		return nil, err
	}

	return &generatedapi.RemoveTeamMemberNoContent{}, nil
}
