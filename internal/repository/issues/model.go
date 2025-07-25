package issues

import (
	"time"

	"github.com/rom8726/warden/internal/domain"
)

type issueModel struct {
	ID                 uint       `db:"id"                   json:"id"`
	ProjectID          uint       `db:"project_id"           json:"project_id"`
	Fingerprint        string     `db:"fingerprint"          json:"fingerprint"`
	Source             string     `db:"source"               json:"source"`
	Status             string     `db:"status"               json:"status"`
	Title              string     `db:"title"                json:"title"`
	Level              string     `db:"level"                json:"level"`
	Platform           string     `db:"platform"             json:"platform"`
	FirstSeen          time.Time  `db:"first_seen"           json:"first_seen"`
	LastSeen           time.Time  `db:"last_seen"            json:"last_seen"`
	TotalEvents        uint       `db:"total_events"         json:"total_events"`
	CreatedAt          time.Time  `db:"created_at"           json:"created_at"`
	UpdatedAt          time.Time  `db:"updated_at"           json:"updated_at"`
	LastNotificationAt *time.Time `db:"last_notification_at" json:"last_notification_at"`
}

func (m *issueModel) toDomain() domain.Issue {
	return domain.Issue{
		ID:                 domain.IssueID(m.ID),
		ProjectID:          domain.ProjectID(m.ProjectID),
		Fingerprint:        m.Fingerprint,
		Source:             domain.IssueSource(m.Source),
		Status:             domain.IssueStatus(m.Status),
		Title:              m.Title,
		Level:              domain.IssueLevel(m.Level),
		Platform:           m.Platform,
		FirstSeen:          m.FirstSeen,
		LastSeen:           m.LastSeen,
		TotalEvents:        m.TotalEvents,
		CreatedAt:          m.CreatedAt,
		UpdatedAt:          m.UpdatedAt,
		LastNotificationAt: m.LastNotificationAt,
	}
}

type issueExtendedModel struct {
	issueModel
	ProjectName        string     `db:"project_name"         json:"project_name"`
	ResolvedBy         *uint      `db:"resolved_by"          json:"resolved_by"`
	ResolvedByUsername *string    `db:"resolved_by_username" json:"resolved_by_username"`
	ResolvedAt         *time.Time `db:"resolved_at"          json:"resolved_at"`
}

func (m *issueExtendedModel) toDomain() domain.IssueExtended {
	var resolvedBy *domain.UserID
	if m.ResolvedBy != nil {
		userID := domain.UserID(*m.ResolvedBy)
		resolvedBy = &userID
	}

	return domain.IssueExtended{
		Issue:              m.issueModel.toDomain(),
		ProjectName:        m.ProjectName,
		ResolvedBy:         resolvedBy,
		ResolvedByUsername: m.ResolvedByUsername,
		ResolvedAt:         m.ResolvedAt,
	}
}
