package rest

import (
	"context"
	"log/slog"

	wardencontext "github.com/rom8726/warden/internal/context"
	"github.com/rom8726/warden/internal/domain"
	generatedapi "github.com/rom8726/warden/internal/generated/server"
)

func (r *RestAPI) ListProjects(ctx context.Context) (generatedapi.ListProjectsRes, error) {
	userID := wardencontext.UserID(ctx)

	// Get all projects
	allProjects, err := r.projectsUseCase.List(ctx)
	if err != nil {
		slog.Error("get all projects failed", "error", err)

		return nil, err
	}

	// Filter projects based on user permissions
	projects, err := r.permissionsService.GetAccessibleProjects(ctx, allProjects)
	if err != nil {
		slog.Error("filter projects failed", "error", err, "user_id", userID)

		return nil, err
	}

	items := make([]generatedapi.Project, 0, len(projects))
	for i := range projects {
		project := projects[i]
		items = append(items, generatedapi.Project{
			ID:          uint(project.ID),
			Name:        project.Name,
			PublicKey:   project.PublicKey,
			Description: project.Description,
			TeamID:      convertTeamID(project.TeamID),
			TeamName:    convertTeamName(project.TeamName),
			CreatedAt:   project.CreatedAt,
		})
	}

	resp := generatedapi.ListProjectsResponse(items)

	return &resp, nil
}

func convertTeamID(teamID *domain.TeamID) generatedapi.OptNilUint {
	if teamID == nil {
		return generatedapi.OptNilUint{}
	}

	return generatedapi.NewOptNilUint(uint(*teamID))
}

func convertTeamName(teamName *string) generatedapi.OptNilString {
	if teamName == nil {
		return generatedapi.OptNilString{}
	}

	return generatedapi.NewOptNilString(*teamName)
}
