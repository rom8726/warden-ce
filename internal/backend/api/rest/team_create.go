package rest

import (
	"context"
	"errors"
	"log/slog"

	"github.com/rom8726/warden/internal/domain"
	generatedapi "github.com/rom8726/warden/internal/generated/server"
)

func (r *RestAPI) CreateTeam(
	ctx context.Context,
	req *generatedapi.CreateTeamRequest,
) (generatedapi.CreateTeamRes, error) {
	teamDTO := domain.TeamDTO{
		Name: req.Name,
	}

	team, err := r.teamsUseCase.Create(ctx, teamDTO)
	if err != nil {
		slog.Error("create team failed", "error", err)

		if errors.Is(err, domain.ErrTeamNameAlreadyInUse) {
			return &generatedapi.ErrorBadRequest{Error: generatedapi.ErrorBadRequestError{
				Message: generatedapi.NewOptString(err.Error()),
			}}, nil
		}

		return &generatedapi.ErrorInternalServerError{
			Error: generatedapi.ErrorInternalServerErrorError{
				Message: generatedapi.NewOptString("Failed to create team"),
			},
		}, err
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

	return &generatedapi.CreateTeamResponse{
		Team: apiTeam,
	}, nil
}
