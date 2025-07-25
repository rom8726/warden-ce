//nolint:interfacebloat // our way
package contract

import (
	"context"
	"encoding/json"
	"time"

	"github.com/rom8726/warden/internal/domain"
	"github.com/rom8726/warden/pkg/kafka"
)

type UsersUseCase interface {
	Login(
		ctx context.Context,
		username, password string,
	) (accessToken, refreshToken, sessionID string, isTmpPassword bool, err error)
	LoginReissue(
		ctx context.Context,
		currRefreshToken string,
	) (accessToken, refreshToken string, err error)
	List(ctx context.Context) ([]domain.User, error)
	ListForTeamAdmin(ctx context.Context, teamID domain.TeamID) ([]domain.User, error)
	GetByID(ctx context.Context, id domain.UserID) (domain.User, error)
	CurrentUserInfo(ctx context.Context, id domain.UserID) (domain.UserInfo, error)
	Create(
		ctx context.Context,
		currentUser domain.User,
		username, email, password string,
		isSuperuser bool,
	) (domain.User, error)
	SetSuperuserStatus(ctx context.Context, id domain.UserID, isSuperuser bool) (domain.User, error)
	SetActiveStatus(ctx context.Context, id domain.UserID, isActive bool) (domain.User, error)
	Delete(ctx context.Context, id domain.UserID) error
	UpdatePassword(ctx context.Context, id domain.UserID, oldPassword, newPassword string) error
	ForgotPassword(ctx context.Context, email string) error
	ResetPassword(ctx context.Context, token, newPassword string) error
	Setup2FA(ctx context.Context, userID domain.UserID) (secret, qrURL, qrImage string, err error)
	Confirm2FA(ctx context.Context, userID domain.UserID, code string) error
	Send2FACode(ctx context.Context, userID domain.UserID, action string) error
	Disable2FA(ctx context.Context, userID domain.UserID, emailCode string) error
	Reset2FA(ctx context.Context, userID domain.UserID, emailCode string) (secret, qrURL, qrImage string, err error)
	Verify2FA(ctx context.Context, code, sessionID string) (accessToken, refreshToken string, expiresIn int, err error)
}

type UsersRepository interface {
	FetchByIDs(ctx context.Context, ids []domain.UserID) ([]domain.User, error)
	Create(ctx context.Context, user domain.UserDTO) (domain.User, error)
	GetByID(ctx context.Context, id domain.UserID) (domain.User, error)
	GetByUsername(ctx context.Context, username string) (domain.User, error)
	GetByEmail(ctx context.Context, email string) (domain.User, error)
	ExistsByID(ctx context.Context, id domain.UserID) (bool, error)
	Update(ctx context.Context, user *domain.User) error
	Delete(ctx context.Context, id domain.UserID) error
	List(ctx context.Context) ([]domain.User, error)
	UpdateLastLogin(ctx context.Context, id domain.UserID) error
	UpdatePassword(ctx context.Context, id domain.UserID, passwordHash string) error
	Update2FA(ctx context.Context, id domain.UserID, enabled bool, secret string, confirmedAt *time.Time) error
}

type Tokenizer interface {
	AccessToken(user *domain.User) (string, error)
	RefreshToken(user *domain.User) (string, error)
	VerifyToken(token string, tokenType domain.TokenType) (*domain.TokenClaims, error)
	ResetPasswordToken(user *domain.User) (string, time.Duration, error)
	AccessTokenTTL() time.Duration
	SecretKey() string
}

// EventUseCase handles event processing.
type EventUseCase interface {
	Timeseries(
		ctx context.Context,
		filter *domain.EventTimeseriesFilter,
	) ([]domain.Timeseries, error)
	IssueTimeseries(
		ctx context.Context,
		filter *domain.IssueEventsTimeseriesFilter,
	) ([]domain.Timeseries, error)
}

type EventRepository interface {
	Timeseries(
		ctx context.Context,
		filter *domain.EventTimeseriesFilter,
	) ([]domain.Timeseries, error)
	IssueTimeseries(
		ctx context.Context,
		fingerprint string, // from issue
		filter *domain.IssueEventsTimeseriesFilter,
	) ([]domain.Timeseries, error)
	FetchForIssue(
		ctx context.Context,
		projectID domain.ProjectID,
		groupHash string,
		limit uint,
	) ([]domain.Event, error)
	EventsByRelease(
		ctx context.Context,
		projectID domain.ProjectID,
		release string,
		limit uint,
	) ([]domain.Event, error)
	TopIssuesByRelease(
		ctx context.Context,
		projectID domain.ProjectID,
		release string,
		limit uint,
	) ([]string, error)
	AggregateBySegment(
		ctx context.Context,
		projectID domain.ProjectID,
		release string,
		segment domain.SegmentName,
	) (map[string]uint, error)
}

type ResolutionsRepository interface {
	Create(ctx context.Context, resolutionDTO domain.ResolutionDTO) (domain.Resolution, error)
	GetByIssueID(ctx context.Context, issueID domain.IssueID) ([]domain.Resolution, error)
}

type ProjectsUseCase interface {
	CreateProject(ctx context.Context, name, description string, teamID *domain.TeamID) (domain.Project, error)
	GetProjectExtended(ctx context.Context, id domain.ProjectID) (domain.ProjectExtended, error)
	List(ctx context.Context) ([]domain.ProjectExtended, error)
	GeneralStats(
		ctx context.Context,
		id domain.ProjectID,
		period time.Duration,
	) (domain.GeneralProjectStats, error)
	GetProjectsByUserID(ctx context.Context, userID domain.UserID, isSuperuser bool) ([]domain.ProjectExtended, error)
	RecentProjects(ctx context.Context) ([]domain.ProjectExtended, error)
	UpdateInfo(ctx context.Context, id domain.ProjectID, name, description string) (domain.ProjectExtended, error)
	ArchiveProject(ctx context.Context, id domain.ProjectID) error
}

type ProjectsRepository interface {
	GetByID(ctx context.Context, id domain.ProjectID) (domain.Project, error)
	Create(ctx context.Context, project *domain.ProjectDTO) (domain.ProjectID, error)
	List(ctx context.Context) ([]domain.ProjectExtended, error)
	RecentProjects(ctx context.Context, userID domain.UserID, limit uint) ([]domain.ProjectExtended, error)
	Update(ctx context.Context, id domain.ProjectID, name, description string) error
	Archive(ctx context.Context, id domain.ProjectID) error
}

type IssueUseCase interface {
	GetByIDWithChildren(
		ctx context.Context,
		id domain.IssueID,
	) (domain.IssueExtendedWithChildren, error)
	List(ctx context.Context, filter *domain.ListIssuesFilter) ([]domain.IssueExtended, uint64, error)
	RecentIssues(ctx context.Context, limit uint) ([]domain.IssueExtended, error)
	Timeseries(ctx context.Context, filter *domain.IssueTimeseriesFilter) ([]domain.Timeseries, error)
	ChangeStatus(ctx context.Context, id domain.IssueID, status domain.IssueStatus) error
}

type TeamsUseCase interface {
	Create(ctx context.Context, teamDTO domain.TeamDTO) (domain.Team, error)
	GetByID(ctx context.Context, id domain.TeamID) (domain.Team, error)
	GetByName(ctx context.Context, name string) (domain.Team, error)
	GetByProjectID(ctx context.Context, projectID domain.ProjectID) (domain.Team, error)
	List(ctx context.Context) ([]domain.Team, error)
	Delete(ctx context.Context, id domain.TeamID) error
	AddMember(
		ctx context.Context,
		teamID domain.TeamID,
		userID domain.UserID,
		role domain.Role,
	) error
	RemoveMemberWithChecks(ctx context.Context, teamID domain.TeamID, userID domain.UserID) error
	ChangeMemberRole(ctx context.Context, teamID domain.TeamID, userID domain.UserID, newRole domain.Role) error
	GetTeamsByUserID(ctx context.Context, userID domain.UserID) ([]domain.Team, error)
	GetTeamByID(ctx context.Context, id domain.TeamID) (domain.Team, error)
	GetMembers(ctx context.Context, teamID domain.TeamID) ([]domain.TeamMember, error)
}

type TeamsRepository interface {
	Create(ctx context.Context, teamDTO domain.TeamDTO) (domain.Team, error)
	GetByID(ctx context.Context, id domain.TeamID) (domain.Team, error)
	GetByName(ctx context.Context, name string) (domain.Team, error)
	GetByProjectID(ctx context.Context, projectID domain.ProjectID) (domain.Team, error)
	List(ctx context.Context) ([]domain.Team, error)
	Delete(ctx context.Context, id domain.TeamID) error
	AddMember(
		ctx context.Context,
		teamID domain.TeamID,
		userID domain.UserID,
		role domain.Role,
	) error
	RemoveMember(ctx context.Context, teamID domain.TeamID, userID domain.UserID) error
	UpdateMemberRole(ctx context.Context, teamID domain.TeamID, userID domain.UserID, newRole domain.Role) error
	GetMembers(ctx context.Context, teamID domain.TeamID) ([]domain.TeamMember, error)
	GetTeamsByUserID(ctx context.Context, userID domain.UserID) ([]domain.Team, error)
}

type NotificationsUseCase interface {
	// Notification Settings
	CreateNotificationSetting(
		ctx context.Context,
		settingDTO domain.NotificationSettingDTO,
	) (domain.NotificationSetting, error)
	GetNotificationSetting(
		ctx context.Context,
		id domain.NotificationSettingID,
	) (domain.NotificationSetting, error)
	UpdateNotificationSetting(
		ctx context.Context,
		setting domain.NotificationSetting,
	) error
	DeleteNotificationSetting(
		ctx context.Context,
		id domain.NotificationSettingID,
	) error
	ListNotificationSettings(
		ctx context.Context,
		projectID domain.ProjectID,
	) ([]domain.NotificationSetting, error)

	// Notification Rules
	CreateNotificationRule(
		ctx context.Context,
		ruleDTO domain.NotificationRuleDTO,
	) (domain.NotificationRule, error)
	GetNotificationRule(
		ctx context.Context,
		id domain.NotificationRuleID,
	) (domain.NotificationRule, error)
	UpdateNotificationRule(
		ctx context.Context,
		rule domain.NotificationRule,
	) error
	DeleteNotificationRule(
		ctx context.Context,
		id domain.NotificationRuleID,
	) error
	ListNotificationRules(
		ctx context.Context,
		settingID domain.NotificationSettingID,
	) ([]domain.NotificationRule, error)

	SendTestNotification(
		ctx context.Context,
		projectID domain.ProjectID,
		notificationSettingID domain.NotificationSettingID,
	) error
}

type NotificationSettingsRepository interface {
	CreateSetting(
		ctx context.Context,
		settingDTO domain.NotificationSettingDTO,
	) (domain.NotificationSetting, error)
	GetSettingByID(
		ctx context.Context,
		id domain.NotificationSettingID,
	) (domain.NotificationSetting, error)
	UpdateSetting(ctx context.Context, setting domain.NotificationSetting) error
	DeleteSetting(ctx context.Context, id domain.NotificationSettingID) error
	ListSettings(
		ctx context.Context,
		projectID domain.ProjectID,
	) ([]domain.NotificationSetting, error)
	ListAllSettings(
		ctx context.Context,
	) ([]domain.NotificationSetting, error)
}

type NotificationRulesRepository interface {
	CreateRule(
		ctx context.Context,
		ruleDTO domain.NotificationRuleDTO,
	) (domain.NotificationRule, error)
	GetRuleByID(ctx context.Context, id domain.NotificationRuleID) (domain.NotificationRule, error)
	UpdateRule(ctx context.Context, rule domain.NotificationRule) error
	DeleteRule(ctx context.Context, id domain.NotificationRuleID) error
	ListRules(
		ctx context.Context,
		settingID domain.NotificationSettingID,
	) ([]domain.NotificationRule, error)
}

type NotificationsQueueRepository interface {
	GetByID(ctx context.Context, id domain.NotificationID) (domain.Notification, error)
	AddNotification(
		ctx context.Context,
		projectID domain.ProjectID,
		issueID domain.IssueID,
		level domain.IssueLevel,
		isNew, wasReactivated bool,
	) error
}

type Notificator interface {
	SendTestNotification(
		ctx context.Context,
		notificationSettingID domain.NotificationSettingID,
	) error
}

type UserNotificationsUseCase interface {
	CreateNotification(
		ctx context.Context,
		userID domain.UserID,
		notificationType domain.UserNotificationType,
		content domain.UserNotificationContent,
	) error
	GetUserNotifications(
		ctx context.Context,
		userID domain.UserID,
		limit, offset uint,
	) ([]domain.UserNotification, error)
	GetUnreadCount(ctx context.Context, userID domain.UserID) (uint, error)
	MarkAsRead(ctx context.Context, notificationID domain.UserNotificationID) error
	MarkAllAsRead(ctx context.Context, userID domain.UserID) error
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
	MarkEmailAsSent(ctx context.Context, id domain.UserNotificationID) error
	MarkEmailAsFailed(ctx context.Context, id domain.UserNotificationID, reason string) error
}

type PermissionsService interface {
	CanAccessProject(ctx context.Context, projectID domain.ProjectID) error
	CanAccessIssue(ctx context.Context, issueID domain.IssueID) error
	CanManageProject(ctx context.Context, projectID domain.ProjectID, isIssueManagement bool) error
	CanManageIssue(ctx context.Context, issueID domain.IssueID) error
	GetAccessibleProjects(ctx context.Context, projects []domain.ProjectExtended) ([]domain.ProjectExtended, error)
}

type Emailer interface {
	Send(
		ctx context.Context,
		issue *domain.Issue,
		project *domain.Project,
		configData json.RawMessage,
		isRegress bool,
	) error
	SendResetPasswordEmail(ctx context.Context, email, token string) error
	Send2FACodeEmail(ctx context.Context, email, code, action string) error
}

type TwoFARateLimiter interface {
	Inc(userID domain.UserID) (attempts int, blocked bool)
	Reset(userID domain.UserID)
	IsBlocked(userID domain.UserID) bool
}

type ReleaseRepository interface {
	Create(ctx context.Context, release domain.ReleaseDTO) (domain.ReleaseID, error)
	GetByID(ctx context.Context, id domain.ReleaseID) (domain.Release, error)
	GetByProjectAndVersion(ctx context.Context, projectID domain.ProjectID, version string) (domain.Release, error)
	ListByProject(ctx context.Context, projectID domain.ProjectID) ([]domain.Release, error)
}

type ReleaseStatsRepository interface {
	GetByProjectAndRelease(
		ctx context.Context,
		projectID domain.ProjectID,
		release string,
	) (domain.ReleaseStats, error)
	Create(ctx context.Context, stats domain.ReleaseStats) error
}

type AnalyticsUseCase interface {
	ListReleases(
		ctx context.Context,
		projectID domain.ProjectID,
	) ([]domain.Release, map[string]domain.ReleaseStats, error)
	GetReleaseDetails(
		ctx context.Context,
		projectID domain.ProjectID,
		version string,
		topLimit uint,
	) (domain.ReleaseAnalyticsDetails, error)
	CompareReleases(
		ctx context.Context,
		projectID domain.ProjectID,
		baseVersion, targetVersion string,
	) (domain.ReleaseComparison, error)
	GetErrorsByTime(
		ctx context.Context,
		projectID domain.ProjectID,
		release string,
		period, granularity time.Duration,
		levels []domain.IssueLevel,
		groupBy domain.EventTimeseriesGroup,
	) ([]domain.Timeseries, error)
	GetUserSegments(
		ctx context.Context,
		projectID domain.ProjectID,
		release string,
	) (domain.UserSegmentsAnalytics, error)
}

type TopicProducerCreator interface {
	Create(topic string) kafka.DataProducer
}

type SettingsUseCase interface {
	GetSetting(ctx context.Context, name string) (*domain.Setting, error)
	SetSetting(ctx context.Context, name string, value any, description string) error
	DeleteSetting(ctx context.Context, name string) error
	ListSettings(ctx context.Context) ([]*domain.Setting, error)
}

// SettingRepository defines the interface for settings operations.
type SettingRepository interface {
	GetByName(ctx context.Context, name string) (*domain.Setting, error)
	SetByName(ctx context.Context, name string, value any, description string) error
	DeleteByName(ctx context.Context, name string) error
	List(ctx context.Context) ([]*domain.Setting, error)
}

type IssuesRepository interface {
	UpsertIssue(ctx context.Context, issue domain.IssueDTO) (domain.IssueUpsertResult, error)
	GetByID(ctx context.Context, id domain.IssueID) (domain.Issue, error)
	ListByFingerprints(ctx context.Context, fingerprints []string) ([]domain.Issue, error)
	CountForLevels(
		ctx context.Context,
		projectID domain.ProjectID,
		period time.Duration,
	) (map[domain.IssueLevel]uint64, error)
	MostFrequent(
		ctx context.Context,
		projectID domain.ProjectID,
		period time.Duration,
		limit uint,
	) ([]domain.IssueExtended, error)
	ListExtended(ctx context.Context, filter *domain.ListIssuesFilter) ([]domain.IssueExtended, uint64, error)
	RecentIssues(ctx context.Context, limit uint) ([]domain.IssueExtended, error)
	Timeseries(
		ctx context.Context,
		filter *domain.IssueTimeseriesFilter,
	) ([]domain.Timeseries, error)
	UpdateStatus(ctx context.Context, issueID domain.IssueID, status domain.IssueStatus) error
	MarkAsNotified(ctx context.Context, issueID domain.IssueID) error
}

type NotificationChannel interface {
	Type() domain.NotificationType
	Send(
		ctx context.Context,
		issue *domain.Issue,
		project *domain.Project,
		config json.RawMessage,
		isRegress bool,
	) error
}

// ComponentVersion represents version information for a system component.
type ComponentVersion struct {
	Name      string
	Version   string
	BuildTime string
	Status    string
}

// VersionsUseCase handles a version collection from system components.
type VersionsUseCase interface {
	GetVersions(ctx context.Context) ([]ComponentVersion, error)
}
