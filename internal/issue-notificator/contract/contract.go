package contract

import (
	"context"

	"github.com/rom8726/warden/internal/domain"
)

type TeamsRepository interface {
	GetUniqueUserIDsByTeamIDs(
		ctx context.Context,
		teamIDs []domain.TeamID,
	) ([]domain.UserID, error)
	GetTeamsByUserIDs(
		ctx context.Context,
		userIDs []domain.UserID,
	) (map[domain.UserID][]domain.Team, error)
	GetMembers(ctx context.Context, teamID domain.TeamID) ([]domain.TeamMember, error)
}

type UsersRepository interface {
	FetchByIDs(ctx context.Context, ids []domain.UserID) ([]domain.User, error)
}

type ProjectsRepository interface {
	GetByID(ctx context.Context, id domain.ProjectID) (domain.Project, error)
}

type IssuesRepository interface {
	GetByID(ctx context.Context, id domain.IssueID) (domain.Issue, error)
	MarkAsNotified(ctx context.Context, issueID domain.IssueID) error
}

type NotificationsUseCase interface {
	GetNotificationSetting(
		ctx context.Context,
		id domain.NotificationSettingID,
	) (domain.NotificationSetting, error)

	TakePendingNotificationsWithSettings(
		ctx context.Context,
		limit uint,
	) ([]domain.NotificationWithSettings, error)
	MarkNotificationAsSent(ctx context.Context, id domain.NotificationID) error
	MarkNotificationAsFailed(ctx context.Context, id domain.NotificationID, reason string) error
	MarkNotificationAsSkipped(ctx context.Context, id domain.NotificationID, reason string) error
}

type NotificationSettingsRepository interface {
	GetSettingByID(
		ctx context.Context,
		id domain.NotificationSettingID,
	) (domain.NotificationSetting, error)
	ListSettings(
		ctx context.Context,
		projectID domain.ProjectID,
	) ([]domain.NotificationSetting, error)
}

type NotificationsQueueRepository interface {
	GetByID(ctx context.Context, id domain.NotificationID) (domain.Notification, error)
	TakePending(ctx context.Context, limit uint) ([]domain.Notification, error)
	MarkAsSent(ctx context.Context, id domain.NotificationID) error
	MarkAsFailed(ctx context.Context, id domain.NotificationID, reason string) error
	MarkAsSkipped(ctx context.Context, id domain.NotificationID, reason string) error
}
