package rest

import (
	"context"
	"errors"
	"log/slog"

	"github.com/rom8726/warden/internal/domain"
	generatedapi "github.com/rom8726/warden/internal/generated/server"
)

func (r *RestAPI) SetUserActiveStatus(
	ctx context.Context,
	req *generatedapi.SetUserActiveStatusRequest,
	params generatedapi.SetUserActiveStatusParams,
) (generatedapi.SetUserActiveStatusRes, error) {
	userID := domain.UserID(params.UserID)
	user, err := r.usersUseCase.SetActiveStatus(ctx, userID, req.IsActive)
	if err != nil {
		slog.Error("set user active status failed", "error", err, "user_id", userID, "is_active", req.IsActive)

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

	var lastLogin generatedapi.OptDateTime
	if user.LastLogin != nil {
		lastLogin.Value = *user.LastLogin
		lastLogin.Set = true
	}

	teams, err := r.teamsUseCase.GetTeamsByUserID(ctx, userID)
	if err != nil {
		slog.Error("get user teams failed", "error", err, "user_id", userID)

		return nil, err
	}

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
		ID:            uint(user.ID),
		Username:      user.Username,
		Email:         user.Email,
		IsSuperuser:   user.IsSuperuser,
		IsActive:      user.IsActive,
		IsTmpPassword: user.IsTmpPassword,
		CreatedAt:     user.CreatedAt,
		LastLogin:     lastLogin,
		Teams:         userTeams,
	}, nil
}
