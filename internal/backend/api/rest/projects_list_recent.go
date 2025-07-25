package rest

import (
	"context"
	"log/slog"

	generatedapi "github.com/rom8726/warden/internal/generated/server"
)

func (r *RestAPI) RecentProjectsList(ctx context.Context) (generatedapi.RecentProjectsListRes, error) {
	list, err := r.projectsUseCase.RecentProjects(ctx)
	if err != nil {
		slog.Error("get recent projects failed", "error", err)

		return nil, err
	}

	items := make([]generatedapi.Project, 0, len(list))
	for i := range list {
		project := list[i]
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
