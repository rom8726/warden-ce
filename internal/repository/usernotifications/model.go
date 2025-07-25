package usernotifications

import (
	"encoding/json"
	"time"

	"github.com/rom8726/warden/internal/domain"
)

type userNotificationModel struct {
	ID        uint            `db:"id"`
	UserID    uint            `db:"user_id"`
	Type      string          `db:"type"`
	Content   json.RawMessage `db:"content"`
	IsRead    bool            `db:"is_read"`
	EmailSent bool            `db:"email_sent"`
	CreatedAt time.Time       `db:"created_at"`
	UpdatedAt time.Time       `db:"updated_at"`
}

func (m *userNotificationModel) toDomain() domain.UserNotification {
	return domain.UserNotification{
		ID:        domain.UserNotificationID(m.ID),
		UserID:    domain.UserID(m.UserID),
		Type:      domain.UserNotificationType(m.Type),
		Content:   m.Content,
		IsRead:    m.IsRead,
		EmailSent: m.EmailSent,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
	}
}
