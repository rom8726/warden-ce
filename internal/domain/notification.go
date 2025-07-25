package domain

import (
	"encoding/json"
	"time"
)

type (
	NotificationID        uint
	NotificationSettingID uint
	NotificationRuleID    uint
)

// NotificationType represents the type of notification.
type NotificationType string

const (
	NotificationTypeEmail      NotificationType = "email"
	NotificationTypeTelegram   NotificationType = "telegram"
	NotificationTypeSlack      NotificationType = "slack"
	NotificationTypeMattermost NotificationType = "mattermost"
	NotificationTypeWebhook    NotificationType = "webhook"
	NotificationTypePachca     NotificationType = "pachca"
)

type NotificationStatus string

const (
	NotificationStatusPending NotificationStatus = "pending"
	NotificationStatusSent    NotificationStatus = "sent"
	NotificationStatusFailed  NotificationStatus = "failed"
)

// NotificationSetting represents a notification setting for a project.
type NotificationSetting struct {
	ID        NotificationSettingID
	ProjectID ProjectID
	Type      NotificationType
	Config    json.RawMessage
	Enabled   bool
	CreatedAt time.Time
	UpdatedAt time.Time
	Rules     []NotificationRule
}

// NotificationRule represents a rule for when to send notifications.
type NotificationRule struct {
	ID                  NotificationRuleID
	NotificationSetting NotificationSettingID
	EventLevel          IssueLevel
	Fingerprint         *string
	IsNewError          *bool
	IsRegression        *bool
	CreatedAt           time.Time
}

type NotificationSettingDTO struct {
	ProjectID ProjectID
	Type      NotificationType
	Config    json.RawMessage
	Enabled   bool
}

type NotificationRuleDTO struct {
	NotificationSetting NotificationSettingID
	EventLevel          string
	Fingerprint         *string
	IsNewError          *bool
	IsRegression        *bool
}

type Notification struct {
	ID             NotificationID
	ProjectID      ProjectID
	IssueID        IssueID
	Level          IssueLevel
	IsNew          bool
	WasReactivated bool
	SentAt         *time.Time
	Status         NotificationStatus
	FailReason     *string
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

type NotificationWithSettings struct {
	Notification
	Settings []NotificationSetting
}
