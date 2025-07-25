package notifications

import (
	"encoding/json"
	"time"

	"github.com/rom8726/warden/internal/domain"
)

type notificationSettingModel struct {
	ID        uint            `db:"id"`
	ProjectID uint            `db:"project_id"`
	Type      string          `db:"type"`
	Config    json.RawMessage `db:"config"`
	Enabled   bool            `db:"enabled"`
	CreatedAt time.Time       `db:"created_at"`
	UpdatedAt time.Time       `db:"updated_at"`
}

type notificationRuleModel struct {
	ID                  uint      `db:"id"`
	NotificationSetting uint      `db:"notification_setting_id"`
	EventLevel          string    `db:"event_level"`
	Fingerprint         *string   `db:"fingerprint"`
	IsNewError          *bool     `db:"is_new_error"`
	IsRegression        *bool     `db:"is_regression"`
	CreatedAt           time.Time `db:"created_at"`
}

func (m *notificationSettingModel) toDomain() domain.NotificationSetting {
	return domain.NotificationSetting{
		ID:        domain.NotificationSettingID(m.ID),
		ProjectID: domain.ProjectID(m.ProjectID),
		Type:      domain.NotificationType(m.Type),
		Config:    m.Config,
		Enabled:   m.Enabled,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
		Rules:     []domain.NotificationRule{}, // Will be populated separately
	}
}

func (m *notificationRuleModel) toDomain() domain.NotificationRule {
	return domain.NotificationRule{
		ID:                  domain.NotificationRuleID(m.ID),
		NotificationSetting: domain.NotificationSettingID(m.NotificationSetting),
		EventLevel:          domain.IssueLevel(m.EventLevel),
		Fingerprint:         m.Fingerprint,
		IsNewError:          m.IsNewError,
		IsRegression:        m.IsRegression,
		CreatedAt:           m.CreatedAt,
	}
}

func settingFromDomain(setting domain.NotificationSetting) notificationSettingModel {
	return notificationSettingModel{
		ID:        uint(setting.ID),
		ProjectID: uint(setting.ProjectID),
		Type:      string(setting.Type),
		Config:    setting.Config,
		Enabled:   setting.Enabled,
		CreatedAt: setting.CreatedAt,
		UpdatedAt: setting.UpdatedAt,
	}
}

func ruleFromDomain(rule domain.NotificationRule) notificationRuleModel {
	return notificationRuleModel{
		ID:                  uint(rule.ID),
		NotificationSetting: uint(rule.NotificationSetting),
		EventLevel:          string(rule.EventLevel),
		Fingerprint:         rule.Fingerprint,
		IsNewError:          rule.IsNewError,
		IsRegression:        rule.IsRegression,
		CreatedAt:           rule.CreatedAt,
	}
}

func settingFromDTO(dto domain.NotificationSettingDTO) notificationSettingModel {
	now := time.Now()

	if dto.Config == nil {
		dto.Config = json.RawMessage("{}")
	}

	return notificationSettingModel{
		ProjectID: uint(dto.ProjectID),
		Type:      string(dto.Type),
		Config:    dto.Config,
		Enabled:   dto.Enabled,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

func ruleFromDTO(dto domain.NotificationRuleDTO) notificationRuleModel {
	return notificationRuleModel{
		NotificationSetting: uint(dto.NotificationSetting),
		EventLevel:          dto.EventLevel,
		Fingerprint:         dto.Fingerprint,
		IsNewError:          dto.IsNewError,
		IsRegression:        dto.IsRegression,
		CreatedAt:           time.Now(),
	}
}
