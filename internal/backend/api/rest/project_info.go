package rest

import (
	"context"
	"errors"
	"log/slog"

	"github.com/rom8726/warden/internal/domain"
	generatedapi "github.com/rom8726/warden/internal/generated/server"
)

func (r *RestAPI) GetProject(
	ctx context.Context,
	params generatedapi.GetProjectParams,
) (generatedapi.GetProjectRes, error) {
	projectID := domain.ProjectID(params.ProjectID)

	// Check if the user has access to the project
	if err := r.permissionsService.CanAccessProject(ctx, projectID); err != nil {
		slog.Error("permission denied", "error", err, "project_id", projectID)

		if errors.Is(err, domain.ErrPermissionDenied) {
			return &generatedapi.ErrorPermissionDenied{Error: generatedapi.ErrorPermissionDeniedError{
				Message: generatedapi.NewOptString("permission denied"),
			}}, nil
		}

		if errors.Is(err, domain.ErrUserNotFound) {
			return &generatedapi.ErrorUnauthorized{Error: generatedapi.ErrorUnauthorizedError{
				Message: generatedapi.NewOptString("unauthorized"),
			}}, nil
		}

		return nil, err
	}

	project, err := r.projectsUseCase.GetProjectExtended(ctx, projectID)
	if err != nil {
		slog.Error("get project failed", "error", err)

		if errors.Is(err, domain.ErrEntityNotFound) {
			return &generatedapi.ErrorNotFound{Error: generatedapi.ErrorNotFoundError{
				Message: generatedapi.NewOptString(err.Error()),
			}}, nil
		}

		return nil, err
	}

	var (
		teamID   generatedapi.OptNilUint
		teamName generatedapi.OptNilString
	)

	if project.TeamID != nil {
		teamID.Value = uint(*project.TeamID)
		teamID.Set = true

		teamName.Value = *project.TeamName
		teamName.Set = true
	}

	return &generatedapi.ProjectResponse{
		Project: generatedapi.Project{
			ID:          uint(project.ID),
			Name:        project.Name,
			PublicKey:   project.PublicKey,
			Description: project.Description,
			TeamID:      teamID,
			TeamName:    teamName,
			CreatedAt:   project.CreatedAt,
		},
	}, nil
}
