package rest

import (
	"context"
	"errors"
	"log/slog"

	"github.com/rom8726/warden/internal/domain"
	generatedapi "github.com/rom8726/warden/internal/generated/server"
)

func (r *RestAPI) AddProject(
	ctx context.Context,
	req *generatedapi.AddProjectRequest,
) (generatedapi.AddProjectRes, error) {
	// Convert OptNilUint to *domain.TeamID
	var teamID *domain.TeamID
	if teamIDValue, ok := req.GetTeamID().Get(); ok {
		// TeamID is set and not null
		id := domain.TeamID(teamIDValue)
		teamID = &id
	}

	_, err := r.projectsUseCase.CreateProject(ctx, req.Name, req.Description, teamID)
	if err != nil {
		slog.Error("add project failed", "error", err)

		if errors.Is(err, domain.ErrEntityNotFound) {
			return &generatedapi.ErrorNotFound{Error: generatedapi.ErrorNotFoundError{
				Message: generatedapi.NewOptString(err.Error()),
			}}, nil
		}

		return nil, err
	}

	return &generatedapi.AddProjectCreated{}, nil
}
