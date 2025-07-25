package resolutions

import (
	"time"

	"github.com/rom8726/warden/internal/domain"
)

type resolutionModel struct {
	ID         uint       `db:"id"`
	ProjectID  uint       `db:"project_id"`
	IssueID    uint       `db:"issue_id"`
	Status     string     `db:"status"`
	ResolvedBy *uint      `db:"resolved_by"`
	ResolvedAt *time.Time `db:"resolved_at"`
	Comment    string     `db:"comment"`
	CreatedAt  time.Time  `db:"created_at"`
	UpdatedAt  time.Time  `db:"updated_at"`
}

func (m *resolutionModel) toDomain() domain.Resolution {
	var resolvedBy *domain.UserID
	if m.ResolvedBy != nil {
		userID := domain.UserID(*m.ResolvedBy)
		resolvedBy = &userID
	}

	return domain.Resolution{
		ID:         domain.ResolutionID(m.ID),
		ProjectID:  domain.ProjectID(m.ProjectID),
		IssueID:    domain.IssueID(m.IssueID),
		Status:     domain.IssueStatus(m.Status),
		ResolvedBy: resolvedBy,
		ResolvedAt: m.ResolvedAt,
		Comment:    m.Comment,
		CreatedAt:  m.CreatedAt,
		UpdatedAt:  m.UpdatedAt,
	}
}

func fromDomain(resolution *domain.Resolution) resolutionModel {
	var resolvedBy *uint
	if resolution.ResolvedBy != nil {
		id := uint(*resolution.ResolvedBy)
		resolvedBy = &id
	}

	return resolutionModel{
		ID:         uint(resolution.ID),
		ProjectID:  uint(resolution.ProjectID),
		IssueID:    uint(resolution.IssueID),
		Status:     string(resolution.Status),
		ResolvedBy: resolvedBy,
		ResolvedAt: resolution.ResolvedAt,
		Comment:    resolution.Comment,
		CreatedAt:  resolution.CreatedAt,
		UpdatedAt:  resolution.UpdatedAt,
	}
}

func fromDTO(dto domain.ResolutionDTO) resolutionModel {
	var resolvedBy *uint
	if dto.ResolvedBy != nil {
		id := uint(*dto.ResolvedBy)
		resolvedBy = &id
	}

	now := time.Now()

	return resolutionModel{
		ProjectID:  uint(dto.ProjectID),
		IssueID:    uint(dto.IssueID),
		Status:     string(dto.Status),
		ResolvedBy: resolvedBy,
		ResolvedAt: nil, // Will be set if status is resolved
		Comment:    dto.Comment,
		CreatedAt:  now,
		UpdatedAt:  now,
	}
}
