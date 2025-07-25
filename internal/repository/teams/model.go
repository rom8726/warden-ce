package teams

import (
	"time"

	"github.com/rom8726/warden/internal/domain"
)

type teamModel struct {
	ID        uint      `db:"id"`
	Name      string    `db:"name"`
	CreatedAt time.Time `db:"created_at"`
}

type teamMemberModel struct {
	TeamID uint   `db:"team_id"`
	UserID uint   `db:"user_id"`
	Role   string `db:"role"`
}

func (m *teamModel) toDomain() domain.Team {
	return domain.Team{
		ID:        domain.TeamID(m.ID),
		Name:      m.Name,
		CreatedAt: m.CreatedAt,
		Members:   []domain.TeamMember{}, // Will be populated separately
	}
}

func (m *teamMemberModel) toDomain() domain.TeamMember {
	return domain.TeamMember{
		TeamID: domain.TeamID(m.TeamID),
		UserID: domain.UserID(m.UserID),
		Role:   domain.Role(m.Role),
	}
}
