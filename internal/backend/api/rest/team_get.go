package rest

import (
	"context"
	"errors"
	"log/slog"

	"github.com/rom8726/warden/internal/domain"
	generatedapi "github.com/rom8726/warden/internal/generated/server"
)

func (r *RestAPI) GetTeam(ctx context.Context, params generatedapi.GetTeamParams) (generatedapi.GetTeamRes, error) {
	teamID := domain.TeamID(params.TeamID)

	team, err := r.teamsUseCase.GetByID(ctx, teamID)
	if err != nil {
		if errors.Is(err, domain.ErrEntityNotFound) {
			return &generatedapi.ErrorNotFound{
				Error: generatedapi.ErrorNotFoundError{
					Message: generatedapi.NewOptString("Team not found"),
				},
			}, nil
		}

		slog.Error("get team failed", "error", err, "team_id", teamID)

		return nil, err
	}

	members := make([]generatedapi.TeamMember, 0, len(team.Members))
	for j := range team.Members {
		member := team.Members[j]
		members = append(members, generatedapi.TeamMember{
			UserID: uint(member.UserID),
			Role:   mapDomainRoleToAPI(member.Role),
		})
	}

	apiTeam := generatedapi.Team{
		ID:        uint(team.ID),
		Name:      team.Name,
		CreatedAt: team.CreatedAt,
		Members:   members,
	}

	return &apiTeam, nil
}
