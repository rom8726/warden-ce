package rest

import (
	"context"
	"errors"
	"log/slog"

	"github.com/rom8726/warden/internal/domain"
	generatedapi "github.com/rom8726/warden/internal/generated/server"
)

func (r *RestAPI) CheckTeamExists(
	ctx context.Context,
	params generatedapi.CheckTeamExistsParams,
) (generatedapi.CheckTeamExistsRes, error) {
	_, err := r.teamsUseCase.GetByName(ctx, params.TeamName)
	if err != nil {
		if errors.Is(err, domain.ErrEntityNotFound) {
			return &generatedapi.TeamExistsResponse{
				Exists: false,
			}, nil
		}

		slog.Error("check team exists failed", "error", err, "team_name", params.TeamName)

		return nil, err
	}

	return &generatedapi.TeamExistsResponse{
		Exists: true,
	}, nil
}
