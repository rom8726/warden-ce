package rest

import (
	"context"
	"errors"
	"log/slog"

	"github.com/rom8726/warden/internal/domain"
	generatedapi "github.com/rom8726/warden/internal/generated/server"
)

func (r *RestAPI) DeleteTeam(
	ctx context.Context,
	params generatedapi.DeleteTeamParams,
) (generatedapi.DeleteTeamRes, error) {
	teamID := domain.TeamID(params.TeamID)

	_, err := r.teamsUseCase.GetByID(ctx, teamID)
	if err != nil {
		slog.Error("get team failed", "error", err, "team_id", teamID)

		switch {
		case errors.Is(err, domain.ErrEntityNotFound):
			return &generatedapi.ErrorNotFound{
				Error: generatedapi.ErrorNotFoundError{
					Message: generatedapi.NewOptString(domain.ErrEntityNotFound.Error()),
				},
			}, nil
		}

		return nil, err
	}

	err = r.teamsUseCase.Delete(ctx, teamID)
	if err != nil {
		slog.Error("delete team failed", "error", err, "team_id", teamID)

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
		case errors.Is(err, domain.ErrTeamHasProjects):
			return &generatedapi.ErrorBadRequest{
				Error: generatedapi.ErrorBadRequestError{
					Message: generatedapi.NewOptString(domain.ErrTeamHasProjects.Error()),
				},
			}, nil
		}

		return nil, err
	}

	return &generatedapi.DeleteTeamNoContent{}, nil
}
