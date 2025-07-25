package projects

import (
	"database/sql"
	"time"

	"github.com/rom8726/warden/internal/domain"
)

type projectModel struct {
	ID          uint           `db:"id"`
	Name        string         `db:"name"`
	Description sql.NullString `db:"description"`
	PublicKey   string         `db:"public_key"`
	TeamID      *uint          `db:"team_id"`
	CreatedAt   time.Time      `db:"created_at"`
	ArchivedAt  *time.Time     `db:"archived_at"`
}

func (m *projectModel) toDomain() domain.Project {
	var teamIDRef *domain.TeamID
	if m.TeamID != nil {
		teamIDDomain := domain.TeamID(*m.TeamID)
		teamIDRef = &teamIDDomain
	}

	return domain.Project{
		ID:          domain.ProjectID(m.ID),
		Name:        m.Name,
		Description: m.Description.String,
		PublicKey:   m.PublicKey,
		TeamID:      teamIDRef,
		CreatedAt:   m.CreatedAt,
		ArchivedAt:  m.ArchivedAt,
	}
}
