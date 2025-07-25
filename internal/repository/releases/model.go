package releases

import (
	"database/sql"
	"time"

	"github.com/rom8726/warden/internal/domain"
)

type releaseModel struct {
	ID          uint           `db:"id"`
	ProjectID   uint           `db:"project_id"`
	Version     string         `db:"version"`
	Description sql.NullString `db:"description"`
	ReleasedAt  time.Time      `db:"released_at"`
	CreatedAt   time.Time      `db:"created_at"`
	UpdatedAt   time.Time      `db:"updated_at"`
}

func (m *releaseModel) toDomain() domain.Release {
	return domain.Release{
		ID:          domain.ReleaseID(m.ID),
		ProjectID:   domain.ProjectID(m.ProjectID),
		Version:     m.Version,
		Description: m.Description.String,
		ReleasedAt:  m.ReleasedAt,
		CreatedAt:   m.CreatedAt,
		UpdatedAt:   m.UpdatedAt,
	}
}
