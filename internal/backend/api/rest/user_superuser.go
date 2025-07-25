package rest

import (
	"context"
	"errors"
	"log/slog"

	"github.com/rom8726/warden/internal/domain"
	generatedapi "github.com/rom8726/warden/internal/generated/server"
)

func (r *RestAPI) SetSuperuserStatus(
	ctx context.Context,
	req *generatedapi.SetSuperuserStatusRequest,
	params generatedapi.SetSuperuserStatusParams,
) (generatedapi.SetSuperuserStatusRes, error) {
	userID := domain.UserID(params.UserID)
	user, err := r.usersUseCase.SetSuperuserStatus(ctx, userID, req.IsSuperuser)
	if err != nil {
		slog.Error("set superuser status failed",
			"error", err, "user_id", userID, "is_superuser", req.IsSuperuser)

		switch {
		case errors.Is(err, domain.ErrEntityNotFound):
			return &generatedapi.ErrorNotFound{Error: generatedapi.ErrorNotFoundError{
				Message: generatedapi.NewOptString(err.Error()),
			}}, nil
		case errors.Is(err, domain.ErrForbidden):
			return &generatedapi.ErrorPermissionDenied{Error: generatedapi.ErrorPermissionDeniedError{
				Message: generatedapi.NewOptString(err.Error()),
			}}, nil
		}

		return nil, err
	}

	// Convert domain.User to API response format
	var lastLogin generatedapi.OptDateTime
	if user.LastLogin != nil {
		lastLogin.Value = *user.LastLogin
		lastLogin.Set = true
	}

	// Get user's teams
	teams, err := r.teamsUseCase.GetTeamsByUserID(ctx, userID)
	if err != nil {
		slog.Error("get user teams failed", "error", err, "user_id", userID)

		return nil, err
	}

	// Convert teams to UserTeam objects
	userTeams := make([]generatedapi.UserTeam, 0, len(teams))
	for _, team := range teams {
		// Find the user's role in this team
		var role string
		for _, member := range team.Members {
			if member.UserID == userID {
				role = string(member.Role)

				break
			}
		}

		userTeams = append(userTeams, generatedapi.UserTeam{
			ID:   uint(team.ID),
			Name: team.Name,
			Role: generatedapi.UserTeamRole(role),
		})
	}

	return &generatedapi.User{
		ID:          uint(user.ID),
		Username:    user.Username,
		Email:       user.Email,
		IsSuperuser: user.IsSuperuser,
		IsActive:    user.IsActive,
		CreatedAt:   user.CreatedAt,
		LastLogin:   lastLogin,
		Teams:       userTeams,
	}, nil
}
