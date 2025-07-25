package projects

import (
	"github.com/rom8726/warden/internal/domain"
)

type projectModelExtended struct {
	projectModel
	TeamName *string `db:"team_name"`
}

func (m *projectModelExtended) toDomain() domain.ProjectExtended {
	return domain.ProjectExtended{
		Project:  m.projectModel.toDomain(),
		TeamName: m.TeamName,
	}
}
