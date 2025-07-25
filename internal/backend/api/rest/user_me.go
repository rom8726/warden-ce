package rest

import (
	"context"
	"log/slog"

	wardencontext "github.com/rom8726/warden/internal/context"
	generatedapi "github.com/rom8726/warden/internal/generated/server"
)

func (r *RestAPI) GetCurrentUser(ctx context.Context) (generatedapi.GetCurrentUserRes, error) {
	userInfo, err := r.usersUseCase.CurrentUserInfo(ctx, wardencontext.UserID(ctx))
	if err != nil {
		slog.Error("get current user info failed", "error", err)

		return nil, err
	}

	var lastLogin generatedapi.OptDateTime
	if userInfo.User.LastLogin != nil {
		lastLogin.Value = *userInfo.User.LastLogin
		lastLogin.Set = true
	}

	userTeams := make([]generatedapi.UserTeam, 0, len(userInfo.Teams))
	for _, teamInfo := range userInfo.Teams {
		// Convert boolean to OptBool
		var canLeave generatedapi.OptBool
		canLeave.SetTo(teamInfo.CanLeave)

		userTeams = append(userTeams, generatedapi.UserTeam{
			ID:       uint(teamInfo.ID),
			Name:     teamInfo.Name,
			Role:     generatedapi.UserTeamRole(teamInfo.Role),
			CanLeave: canLeave,
		})
	}

	return &generatedapi.User{
		ID:            uint(userInfo.User.ID),
		Username:      userInfo.User.Username,
		Email:         userInfo.User.Email,
		IsSuperuser:   userInfo.User.IsSuperuser,
		IsActive:      userInfo.User.IsActive,
		IsTmpPassword: userInfo.User.IsTmpPassword,
		TwoFaEnabled:  userInfo.User.TwoFAEnabled,
		CreatedAt:     userInfo.User.CreatedAt,
		LastLogin:     lastLogin,
		Teams:         userTeams,
	}, nil
}
