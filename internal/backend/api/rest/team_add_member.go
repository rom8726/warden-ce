package rest

import (
	"context"
	"errors"
	"log/slog"

	"github.com/rom8726/warden/internal/domain"
	generatedapi "github.com/rom8726/warden/internal/generated/server"
)

func (r *RestAPI) AddTeamMember(
	ctx context.Context,
	req *generatedapi.AddTeamMemberRequest,
	params generatedapi.AddTeamMemberParams,
) (generatedapi.AddTeamMemberRes, error) {
	teamID := domain.TeamID(params.TeamID)
	userID := domain.UserID(req.GetUserID())

	var role domain.Role
	switch req.GetRole() {
	case generatedapi.AddTeamMemberRequestRoleOwner:
		role = domain.RoleOwner
	case generatedapi.AddTeamMemberRequestRoleAdmin:
		role = domain.RoleAdmin
	case generatedapi.AddTeamMemberRequestRoleMember:
		role = domain.RoleMember
	default:
		role = domain.RoleMember // Default to member if unknown
	}

	err := r.teamsUseCase.AddMember(ctx, teamID, userID, role)
	if err != nil {
		slog.Error("add team member failed", "error", err, "team_id", teamID, "user_id", userID)

		switch {
		case errors.Is(err, domain.ErrEntityNotFound):
			return &generatedapi.ErrorNotFound{
				Error: generatedapi.ErrorNotFoundError{
					Message: generatedapi.NewOptString(err.Error()),
				},
			}, nil
		case errors.Is(err, domain.ErrForbidden):
			return &generatedapi.ErrorPermissionDenied{
				Error: generatedapi.ErrorPermissionDeniedError{
					Message: generatedapi.NewOptString(domain.ErrForbidden.Error()),
				},
			}, nil
		}

		return nil, err
	}

	return &generatedapi.AddTeamMemberCreated{}, nil
}
