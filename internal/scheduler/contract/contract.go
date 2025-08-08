package contract

import (
	"context"
	"encoding/json"
	"time"

	"github.com/rom8726/warden/internal/domain"
)

type IssuesRepository interface {
	ListUnresolved(ctx context.Context) ([]domain.IssueExtended, error)
	NewIssuesForRelease(ctx context.Context, projectID domain.ProjectID, release string) ([]string, error)
	RegressionsForRelease(ctx context.Context, projectID domain.ProjectID, release string) ([]string, error)
	ResolvedInRelease(ctx context.Context, projectID domain.ProjectID, release string) ([]string, error)
	FixTimesForRelease(
		ctx context.Context,
		projectID domain.ProjectID,
		release string,
	) (map[string]time.Duration, error)
	DeleteOld(ctx context.Context, maxAge time.Duration, limit uint) (uint, error)
}

type NotificationsQueueRepository interface {
	AddNotification(
		ctx context.Context,
		projectID domain.ProjectID,
		issueID domain.IssueID,
		level domain.IssueLevel,
		isNew, wasReactivated bool,
	) error
	DeleteOld(ctx context.Context, maxAge time.Duration, limit uint) (uint, error)
}

type Emailer interface {
	SendUnresolvedIssuesSummaryEmail(ctx context.Context, issues []domain.IssueExtended) error
}

type ReleaseRepository interface {
	ListByProject(ctx context.Context, projectID domain.ProjectID) ([]domain.Release, error)
}

type ReleaseStatsRepository interface {
	Create(ctx context.Context, stats domain.ReleaseStats) error
}

type EventRepository interface {
	AggregateBySegment(
		ctx context.Context,
		projectID domain.ProjectID,
		release string,
		segment domain.SegmentName,
	) (map[string]uint, error)
}

type ProjectsRepository interface {
	List(ctx context.Context) ([]domain.ProjectExtended, error)
	GetByID(ctx context.Context, id domain.ProjectID) (domain.Project, error)
}

type AnalyticsUseCase interface {
	RecalculateReleaseStatsForAllProjects(ctx context.Context) error
}

type UserNotificationsUseCase interface {
	DeleteOldNotifications(ctx context.Context, maxAge time.Duration, limit uint) (uint, error)
}

type UserNotificationsRepository interface {
	Create(
		ctx context.Context,
		userID domain.UserID,
		notificationType domain.UserNotificationType,
		content json.RawMessage,
	) (domain.UserNotification, error)
	GetByID(ctx context.Context, id domain.UserNotificationID) (domain.UserNotification, error)
	GetByUserID(ctx context.Context, userID domain.UserID, limit, offset uint) ([]domain.UserNotification, error)
	GetUnreadCount(ctx context.Context, userID domain.UserID) (uint, error)
	MarkAsRead(ctx context.Context, id domain.UserNotificationID) error
	MarkAllAsRead(ctx context.Context, userID domain.UserID) error
	DeleteOld(ctx context.Context, maxAge time.Duration, limit uint) (uint, error)
}

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
