package rest

import (
	"context"
	"errors"
	"log/slog"

	"github.com/rom8726/warden/internal/domain"
	generatedapi "github.com/rom8726/warden/internal/generated/server"
)

func (r *RestAPI) ChangeTeamMemberRole(
	ctx context.Context,
	req *generatedapi.ChangeTeamMemberRoleRequest,
	params generatedapi.ChangeTeamMemberRoleParams,
) (generatedapi.ChangeTeamMemberRoleRes, error) {
	teamID := domain.TeamID(params.TeamID)
	userID := domain.UserID(params.UserID)

	var role domain.Role
	switch req.GetRole() {
	case generatedapi.ChangeTeamMemberRoleRequestRoleOwner:
		role = domain.RoleOwner
	case generatedapi.ChangeTeamMemberRoleRequestRoleAdmin:
		role = domain.RoleAdmin
	case generatedapi.ChangeTeamMemberRoleRequestRoleMember:
		role = domain.RoleMember
	default:
		return &generatedapi.ErrorBadRequest{
			Error: generatedapi.ErrorBadRequestError{
				Message: generatedapi.NewOptString("Invalid role"),
			},
		}, nil
	}

	err := r.teamsUseCase.ChangeMemberRole(ctx, teamID, userID, role)
	if err != nil {
		slog.Error("change team member role failed", "error", err, "team_id", teamID, "user_id", userID)

		switch {
		case errors.Is(err, domain.ErrEntityNotFound):
			return &generatedapi.ErrorNotFound{
				Error: generatedapi.ErrorNotFoundError{
					Message: generatedapi.NewOptString("Team or user not found"),
				},
			}, nil
		case errors.Is(err, domain.ErrForbidden):
			return &generatedapi.ErrorPermissionDenied{
				Error: generatedapi.ErrorPermissionDeniedError{
					Message: generatedapi.NewOptString("Insufficient permissions to change role"),
				},
			}, nil
		case errors.Is(err, domain.ErrLastOwner):
			return &generatedapi.ErrorBadRequest{
				Error: generatedapi.ErrorBadRequestError{
					Message: generatedapi.NewOptString("Cannot demote the only owner of the team"),
				},
			}, nil
		}

		return nil, err
	}

	return &generatedapi.ChangeTeamMemberRoleOK{}, nil
}
