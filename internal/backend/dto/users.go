package dto

import (
	"github.com/rom8726/warden/internal/domain"
	generatedapi "github.com/rom8726/warden/internal/generated/server"
)

func DomainUsersToAPI(users []domain.User) []generatedapi.User {
	resp := make([]generatedapi.User, 0, len(users))

	for i := range users {
		user := users[i]
		var lastLogin generatedapi.OptDateTime
		if user.LastLogin != nil {
			lastLogin.Value = *user.LastLogin
			lastLogin.Set = true
		}

		resp = append(resp, generatedapi.User{
			ID:          uint(user.ID),
			Username:    user.Username,
			Email:       user.Email,
			IsSuperuser: user.IsSuperuser,
			IsActive:    user.IsActive,
			CreatedAt:   user.CreatedAt,
			LastLogin:   lastLogin,
		})
	}

	return resp
}
