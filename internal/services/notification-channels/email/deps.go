package email

import (
	"context"

	"github.com/rom8726/warden/internal/domain"
)

type TeamsRepository interface {
	GetUniqueUserIDsByTeamIDs(
		ctx context.Context,
		teamIDs []domain.TeamID,
	) ([]domain.UserID, error)
	GetTeamsByUserIDs(
		ctx context.Context,
		userIDs []domain.UserID,
	) (map[domain.UserID][]domain.Team, error)
	GetMembers(ctx context.Context, teamID domain.TeamID) ([]domain.TeamMember, error)
}

type UsersRepository interface {
	FetchByIDs(ctx context.Context, ids []domain.UserID) ([]domain.User, error)
}

type ProjectsRepository interface {
	GetByID(ctx context.Context, id domain.ProjectID) (domain.Project, error)
}
