package rest

import (
	"context"
	"log/slog"

	"github.com/rom8726/warden/internal/domain"
	generatedapi "github.com/rom8726/warden/internal/generated/server"
)

func (r *RestAPI) ListTeams(ctx context.Context) (generatedapi.ListTeamsRes, error) {
	teams, err := r.teamsUseCase.List(ctx)
	if err != nil {
		slog.Error("list teams failed", "error", err)

		return nil, err
	}

	items := make([]generatedapi.Team, 0, len(teams))
	for i := range teams {
		team := teams[i]

		members := make([]generatedapi.TeamMember, 0, len(team.Members))
		for j := range team.Members {
			member := team.Members[j]
			members = append(members, generatedapi.TeamMember{
				UserID: uint(member.UserID),
				Role:   mapDomainRoleToAPI(member.Role),
			})
		}

		items = append(items, generatedapi.Team{
			ID:        uint(team.ID),
			Name:      team.Name,
			CreatedAt: team.CreatedAt,
			Members:   members,
		})
	}

	resp := generatedapi.ListTeamsResponse(items)

	return &resp, nil
}

func mapDomainRoleToAPI(role domain.Role) generatedapi.TeamMemberRole {
	switch role {
	case domain.RoleOwner:
		return generatedapi.TeamMemberRoleOwner
	case domain.RoleAdmin:
		return generatedapi.TeamMemberRoleAdmin
	case domain.RoleMember:
		return generatedapi.TeamMemberRoleMember
	default:
		return generatedapi.TeamMemberRoleMember // Default to member if unknown
	}
}
