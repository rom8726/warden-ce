package domain

import (
	"encoding/json"
	"time"
)

type UserNotificationID uint

// UserNotificationType represents the type of user notification.
type UserNotificationType string

const (
	UserNotificationTypeTeamAdded       UserNotificationType = "team_added"
	UserNotificationTypeTeamRemoved     UserNotificationType = "team_removed"
	UserNotificationTypeRoleChanged     UserNotificationType = "role_changed"
	UserNotificationTypeIssueRegression UserNotificationType = "issue_regression"
)

// UserNotification represents a user notification.
type UserNotification struct {
	ID        UserNotificationID
	UserID    UserID
	Type      UserNotificationType
	Content   json.RawMessage
	IsRead    bool
	EmailSent bool
	CreatedAt time.Time
	UpdatedAt time.Time
}

// Uses embedded structs to avoid nesting and provide flat JSON structure.
type UserNotificationContent struct {
	TeamAdded       *TeamAddedContent       `json:"team_added,omitempty"`
	TeamRemoved     *TeamRemovedContent     `json:"team_removed,omitempty"`
	RoleChanged     *RoleChangedContent     `json:"role_changed,omitempty"`
	IssueRegression *IssueRegressionContent `json:"issue_regression,omitempty"`
}

// TeamAddedContent represents content for team added notifications.
type TeamAddedContent struct {
	TeamID          uint   `json:"team_id"`
	TeamName        string `json:"team_name"`
	Role            string `json:"role"`
	AddedByUserID   uint   `json:"added_by_user_id"`
	AddedByUsername string `json:"added_by_username"`
}

// TeamRemovedContent represents content for team removed notifications.
type TeamRemovedContent struct {
	TeamID            uint   `json:"team_id"`
	TeamName          string `json:"team_name"`
	RemovedByUserID   uint   `json:"removed_by_user_id"`
	RemovedByUsername string `json:"removed_by_username"`
}

// RoleChangedContent represents content for role changed notifications.
type RoleChangedContent struct {
	TeamID            uint   `json:"team_id"`
	TeamName          string `json:"team_name"`
	OldRole           string `json:"old_role"`
	NewRole           string `json:"new_role"`
	ChangedByUserID   uint   `json:"changed_by_user_id"`
	ChangedByUsername string `json:"changed_by_username"`
}

// IssueRegressionContent represents content for issue regression notifications.
type IssueRegressionContent struct {
	IssueID       uint   `json:"issue_id"`
	IssueTitle    string `json:"issue_title"`
	ProjectID     uint   `json:"project_id"`
	ProjectName   string `json:"project_name"`
	ResolvedAt    string `json:"resolved_at"`
	ReactivatedAt string `json:"reactivated_at"`
}
