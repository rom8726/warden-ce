package permissions

import (
	"context"

	"github.com/rom8726/warden/internal/backend/contract"
	wardencontext "github.com/rom8726/warden/internal/context"
	"github.com/rom8726/warden/internal/domain"
)

// Service handles permission checks for various operations.
type Service struct {
	teamsUseCase contract.TeamsUseCase
	projectRepo  contract.ProjectsRepository
	issueRepo    contract.IssuesRepository
}

// New creates a new permissions service.
func New(
	teamsUseCase contract.TeamsUseCase,
	projectRepo contract.ProjectsRepository,
	issueRepo contract.IssuesRepository,
) *Service {
	return &Service{
		teamsUseCase: teamsUseCase,
		projectRepo:  projectRepo,
		issueRepo:    issueRepo,
	}
}

// CanAccessProject checks if a user can access a project.
func (s *Service) CanAccessProject(ctx context.Context, projectID domain.ProjectID) error {
	// Get the project to check its team
	project, err := s.projectRepo.GetByID(ctx, projectID)
	if err != nil {
		return err
	}

	isSuper := wardencontext.IsSuper(ctx)
	if isSuper {
		return nil
	}

	userID := wardencontext.UserID(ctx)
	if userID == 0 {
		return domain.ErrUserNotFound
	}

	// If the project has no team, it's accessible to all users
	if project.TeamID == nil {
		return nil
	}

	// Get the teams that the user is a member of
	userTeams, err := s.teamsUseCase.GetTeamsByUserID(ctx, userID)
	if err != nil {
		return err
	}

	// Check if the user is a member of the project's team
	for _, team := range userTeams {
		if team.ID == *project.TeamID {
			return nil
		}
	}

	return domain.ErrPermissionDenied
}

// CanAccessIssue checks if a user can access an issue.
func (s *Service) CanAccessIssue(ctx context.Context, issueID domain.IssueID) error {
	isSuper := wardencontext.IsSuper(ctx)
	if isSuper {
		return nil
	}

	// Get the issue to find its project
	issue, err := s.issueRepo.GetByID(ctx, issueID)
	if err != nil {
		return err
	}

	// Check if the user can access the project that the issue belongs to
	return s.CanAccessProject(ctx, issue.ProjectID)
}

// CanManageProject checks if a user can manage a project (create, update, delete).
func (s *Service) CanManageProject(ctx context.Context, projectID domain.ProjectID, isIssueManagement bool) error {
	// Get the project to check its team
	project, err := s.projectRepo.GetByID(ctx, projectID)
	if err != nil {
		return err
	}

	isSuper := wardencontext.IsSuper(ctx)
	if isSuper {
		return nil
	}

	userID := wardencontext.UserID(ctx)
	if userID == 0 {
		return domain.ErrUserNotFound
	}

	// If the project has no team, it can't be managed by regular users
	if project.TeamID == nil {
		if isIssueManagement {
			return nil
		}

		return domain.ErrPermissionDenied
	}

	// Get the team members to check the user's role
	members, err := s.teamsUseCase.GetMembers(ctx, *project.TeamID)
	if err != nil {
		return err
	}

	// Check if the user is an owner or admin of the team
	for _, member := range members {
		if member.UserID == userID && isIssueManagement {
			return nil
		}

		if member.UserID == userID && (member.Role == domain.RoleOwner || member.Role == domain.RoleAdmin) {
			return nil
		}
	}

	return domain.ErrPermissionDenied
}

// CanManageIssue checks if a user can manage an issue (update, delete).
func (s *Service) CanManageIssue(ctx context.Context, issueID domain.IssueID) error {
	// Get the issue to find its project
	issue, err := s.issueRepo.GetByID(ctx, issueID)
	if err != nil {
		return err
	}

	isSuper := wardencontext.IsSuper(ctx)
	if isSuper {
		return nil
	}

	// Check if the user can manage the project that the issue belongs to
	return s.CanManageProject(ctx, issue.ProjectID, true)
}

// GetAccessibleProjects returns all projects that a user can access.
func (s *Service) GetAccessibleProjects(
	ctx context.Context,
	projects []domain.ProjectExtended,
) ([]domain.ProjectExtended, error) {
	isSuper := wardencontext.IsSuper(ctx)
	if isSuper {
		return projects, nil
	}

	userID := wardencontext.UserID(ctx)
	if userID == 0 {
		return nil, domain.ErrUserNotFound
		// If there's no user ID in the context, it means the request hasn't gone through
		// the authentication middleware yet. In this case, we'll return all projects
		// and let the authentication middleware handle it later.
	}

	// Get the teams that the user is a member of
	userTeams, err := s.teamsUseCase.GetTeamsByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Create a map of team IDs for a quick lookup
	teamIDs := make(map[domain.TeamID]struct{}, len(userTeams))
	for _, team := range userTeams {
		teamIDs[team.ID] = struct{}{}
	}

	// Filter projects to only include those that belong to the user's teams or have no team
	filteredProjects := make([]domain.ProjectExtended, 0, len(projects))
	for _, project := range projects {
		// Include projects without a team (personal projects)
		if project.TeamID == nil {
			filteredProjects = append(filteredProjects, project)

			continue
		}

		// Include projects that belong to the user's teams
		if _, ok := teamIDs[*project.TeamID]; ok {
			filteredProjects = append(filteredProjects, project)
		}
	}

	return filteredProjects, nil
}
