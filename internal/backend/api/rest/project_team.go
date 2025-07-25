package rest

import (
	"context"

	"github.com/rom8726/warden/internal/domain"
	generatedapi "github.com/rom8726/warden/internal/generated/server"
)

func (r *RestAPI) GetProjectTeam(
	ctx context.Context,
	params generatedapi.GetProjectTeamParams,
) (generatedapi.GetProjectTeamRes, error) {
	team, err := r.teamsUseCase.GetByProjectID(ctx, domain.ProjectID(params.ProjectID))
	if err != nil {
		return &generatedapi.ErrorInternalServerError{
			Error: generatedapi.ErrorInternalServerErrorError{
				Message: generatedapi.NewOptString(err.Error()),
			},
		}, nil
	}

	// Convert members to API format
	apiMembers := make([]generatedapi.TeamMember, 0, len(team.Members))
	for _, member := range team.Members {
		apiMembers = append(apiMembers, generatedapi.TeamMember{
			UserID: uint(member.UserID),
			Role:   mapDomainRoleToAPI(member.Role),
		})
	}

	response := generatedapi.TeamResponse{
		Team: generatedapi.Team{
			ID:        uint(team.ID),
			Name:      team.Name,
			CreatedAt: team.CreatedAt,
			Members:   apiMembers,
		},
	}

	return &response, nil
}
