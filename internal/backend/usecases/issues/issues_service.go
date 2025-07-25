package issues

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/rom8726/warden/internal/backend/contract"
	wardencontext "github.com/rom8726/warden/internal/context"
	"github.com/rom8726/warden/internal/domain"
	"github.com/rom8726/warden/pkg/db"
)

const eventsLimitForIssue = 10

type Service struct {
	txManager                db.TxManager
	issuesRepo               contract.IssuesRepository
	projectsRepo             contract.ProjectsRepository
	eventsRepo               contract.EventRepository
	projectsService          contract.ProjectsUseCase
	resolutionsRepo          contract.ResolutionsRepository
	usersRepo                contract.UsersRepository
	teamsRepo                contract.TeamsRepository
	userNotificationsUseCase contract.UserNotificationsUseCase
}

func New(
	txManager db.TxManager,
	issuesRepo contract.IssuesRepository,
	projectsRepo contract.ProjectsRepository,
	eventsRepo contract.EventRepository,
	projectsService contract.ProjectsUseCase,
	resolutionsRepo contract.ResolutionsRepository,
	usersRepo contract.UsersRepository,
	teamsRepo contract.TeamsRepository,
	userNotificationsUseCase contract.UserNotificationsUseCase,
) *Service {
	return &Service{
		txManager:                txManager,
		issuesRepo:               issuesRepo,
		projectsRepo:             projectsRepo,
		eventsRepo:               eventsRepo,
		projectsService:          projectsService,
		resolutionsRepo:          resolutionsRepo,
		usersRepo:                usersRepo,
		teamsRepo:                teamsRepo,
		userNotificationsUseCase: userNotificationsUseCase,
	}
}

func (s *Service) List(ctx context.Context, filter *domain.ListIssuesFilter) ([]domain.IssueExtended, uint64, error) {
	return s.issuesRepo.ListExtended(ctx, filter)
}

func (s *Service) RecentIssues(ctx context.Context, limit uint) ([]domain.IssueExtended, error) {
	return s.issuesRepo.RecentIssues(ctx, limit)
}

func (s *Service) GetByIDWithChildren(
	ctx context.Context,
	id domain.IssueID,
) (domain.IssueExtendedWithChildren, error) {
	issue, err := s.issuesRepo.GetByID(ctx, id)
	if err != nil {
		return domain.IssueExtendedWithChildren{}, fmt.Errorf("get issue by ID: %w", err)
	}

	currentUserID := wardencontext.UserID(ctx)
	user, err := s.usersRepo.GetByID(ctx, currentUserID)
	if err != nil {
		return domain.IssueExtendedWithChildren{}, fmt.Errorf("get current user: %w", err)
	}

	userProjects, err := s.projectsService.GetProjectsByUserID(ctx, currentUserID, user.IsSuperuser)
	if err != nil {
		return domain.IssueExtendedWithChildren{}, fmt.Errorf("get user projects: %w", err)
	}

	allowed := false
	for _, p := range userProjects {
		if p.ID == issue.ProjectID {
			allowed = true

			break
		}
	}
	if !allowed {
		return domain.IssueExtendedWithChildren{}, errors.New("permission denied")
	}

	project, err := s.projectsRepo.GetByID(ctx, issue.ProjectID)
	if err != nil {
		return domain.IssueExtendedWithChildren{}, fmt.Errorf("get project by ID: %w", err)
	}

	events, err := s.eventsRepo.FetchForIssue(ctx, project.ID, issue.Fingerprint, eventsLimitForIssue)
	if err != nil {
		return domain.IssueExtendedWithChildren{}, fmt.Errorf("fetch events for issue: %w", err)
	}

	return domain.IssueExtendedWithChildren{
		Issue:       issue,
		ProjectName: project.Name,
		Events:      events,
	}, nil
}

func (s *Service) Timeseries(ctx context.Context, filter *domain.IssueTimeseriesFilter) ([]domain.Timeseries, error) {
	return s.issuesRepo.Timeseries(ctx, filter)
}

//nolint:gocyclo,nestif // need refactoring
func (s *Service) ChangeStatus(ctx context.Context, id domain.IssueID, status domain.IssueStatus) error {
	currentUserID := wardencontext.UserID(ctx)
	user, err := s.usersRepo.GetByID(ctx, currentUserID)
	if err != nil {
		return fmt.Errorf("get current user by ID: %w", err)
	}

	issue, err := s.issuesRepo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("get issue by ID: %w", err)
	}

	project, err := s.projectsRepo.GetByID(ctx, issue.ProjectID)
	if err != nil {
		return fmt.Errorf("get project by ID: %w", err)
	}

	// Get user's projects to check access
	userProjects, err := s.projectsService.GetProjectsByUserID(ctx, currentUserID, user.IsSuperuser)
	if err != nil {
		return fmt.Errorf("get user projects: %w", err)
	}

	// Check if a user has access to the issue's project
	var hasAccess bool
	for _, proj := range userProjects {
		if proj.ID == issue.ProjectID {
			hasAccess = true

			break
		}
	}

	if !hasAccess {
		return errors.New("user does not have access to this issue")
	}

	// Check if this is a regression (resolved -> unresolved)
	isRegression := issue.Status == domain.IssueStatusResolved && status == domain.IssueStatusUnresolved

	err = s.txManager.RepeatableRead(ctx, func(ctx context.Context) error {
		// Create a resolution record
		resolutionDTO := domain.ResolutionDTO{
			ProjectID:  issue.ProjectID,
			IssueID:    id,
			Status:     status,
			ResolvedBy: &currentUserID,
			Comment:    "", // Empty comment for now
		}

		_, err := s.resolutionsRepo.Create(ctx, resolutionDTO)
		if err != nil {
			return fmt.Errorf("create resolution: %w", err)
		}

		// Update issue status and resolved_by field
		err = s.issuesRepo.UpdateStatus(ctx, id, status)
		if err != nil {
			return fmt.Errorf("update issue status: %w", err)
		}

		// Create a notification for regression if applicable
		if isRegression {
			// Get team members for the project to notify them about regression
			teamMembers, err := s.getTeamMembersForProject(ctx, project)
			if err != nil {
				// Log error but don't fail the operation
				return fmt.Errorf("get team members for project: %w", err)
			}

			if len(teamMembers) == 0 {
				return nil
			}

			// Get the latest resolution to get resolved_at time
			resolutions, err := s.resolutionsRepo.GetByIssueID(ctx, id)
			if err != nil {
				return fmt.Errorf("get resolutions for issue: %w", err)
			}

			var resolvedAt string
			if len(resolutions) > 0 {
				// Get the most recent resolution
				latestResolution := resolutions[0]
				if latestResolution.ResolvedAt != nil {
					resolvedAt = latestResolution.ResolvedAt.Format(time.RFC3339)
				}
			}

			// Create notification content
			content := domain.UserNotificationContent{
				IssueRegression: &domain.IssueRegressionContent{
					IssueID:       uint(id),
					IssueTitle:    issue.Title,
					ProjectID:     uint(issue.ProjectID),
					ProjectName:   project.Name,
					ResolvedAt:    resolvedAt,
					ReactivatedAt: time.Now().Format(time.RFC3339),
				},
			}

			// Create notifications for all team members
			for _, member := range teamMembers {
				err = s.userNotificationsUseCase.CreateNotification(
					ctx,
					member.UserID,
					domain.UserNotificationTypeIssueRegression,
					content,
				)
				if err != nil {
					// Log error but don't fail the operation
					return fmt.Errorf("create user notification for user %d: %w", member.UserID, err)
				}
			}
		}

		return nil
	})
	if err != nil {
		return fmt.Errorf("change issue status: %w", err)
	}

	return nil
}

// getTeamMembersForProject returns all team members for a given project.
func (s *Service) getTeamMembersForProject(ctx context.Context, project domain.Project) ([]domain.TeamMember, error) {
	if project.TeamID == nil {
		// Project without team - return empty list
		return []domain.TeamMember{}, nil
	}

	// Get team members from the team repository
	teamMembers, err := s.teamsRepo.GetMembers(ctx, *project.TeamID)
	if err != nil {
		return nil, fmt.Errorf("get team members: %w", err)
	}

	return teamMembers, nil
}
