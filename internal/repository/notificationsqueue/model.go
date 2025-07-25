package notificationsqueue

import (
	"time"

	"github.com/rom8726/warden/internal/domain"
)

type notificationModel struct {
	ID             uint       `db:"id"`
	IssueID        uint       `db:"issue_id"`
	ProjectID      uint       `db:"project_id"`
	Level          string     `db:"level"`
	IsNew          bool       `db:"is_new"`
	WasReactivated bool       `db:"was_reactivated"`
	SentAt         *time.Time `db:"sent_at"`
	Status         string     `db:"status"`
	FailReason     *string    `db:"fail_reason"`
	CreatedAt      time.Time  `db:"created_at"`
	UpdatedAt      time.Time  `db:"updated_at"`
}

func (m *notificationModel) toDomain() domain.Notification {
	return domain.Notification{
		ID:             domain.NotificationID(m.ID),
		ProjectID:      domain.ProjectID(m.ProjectID),
		IssueID:        domain.IssueID(m.IssueID),
		Level:          domain.IssueLevel(m.Level),
		IsNew:          m.IsNew,
		WasReactivated: m.WasReactivated,
		SentAt:         m.SentAt,
		Status:         domain.NotificationStatus(m.Status),
		FailReason:     m.FailReason,
		CreatedAt:      m.CreatedAt,
	}
}
