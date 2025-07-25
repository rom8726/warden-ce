package projects

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log/slog"
	"time"

	"github.com/rom8726/warden/internal/backend/contract"
	wardencontext "github.com/rom8726/warden/internal/context"
	"github.com/rom8726/warden/internal/domain"
)

const (
	recentProjectsLimit = 5
)

type ProjectService struct {
	projectRepo      contract.ProjectsRepository
	issuesRepository contract.IssuesRepository
	teamsUseCase     contract.TeamsUseCase
}

func New(
	projectRepo contract.ProjectsRepository,
	issuesRepository contract.IssuesRepository,
	teamsUseCase contract.TeamsUseCase,
) *ProjectService {
	return &ProjectService{
		projectRepo:      projectRepo,
		issuesRepository: issuesRepository,
		teamsUseCase:     teamsUseCase,
	}
}

func (s *ProjectService) GetProject(ctx context.Context, id domain.ProjectID) (domain.Project, error) {
	return s.projectRepo.GetByID(ctx, id)
}

func (s *ProjectService) GetProjectExtended(ctx context.Context, id domain.ProjectID) (domain.ProjectExtended, error) {
	project, err := s.projectRepo.GetByID(ctx, id)
	if err != nil {
		return domain.ProjectExtended{}, fmt.Errorf("get project: %w", err)
	}

	projectExtended := domain.ProjectExtended{
		Project: project,
	}

	if project.TeamID != nil {
		team, err := s.teamsUseCase.GetTeamByID(ctx, *project.TeamID)
		if err != nil {
			return domain.ProjectExtended{}, fmt.Errorf("get team: %w", err)
		}

		projectExtended.TeamName = &team.Name
	}

	return projectExtended, nil
}

func (s *ProjectService) CreateProject(
	ctx context.Context,
	name, description string,
	teamID *domain.TeamID,
) (domain.Project, error) {
	publicKey, err := generateRandomKey(32)
	if err != nil {
		return domain.Project{}, fmt.Errorf("generate public key: %w", err)
	}

	if teamID != nil {
		_, err := s.teamsUseCase.GetTeamByID(ctx, *teamID)
		if err != nil {
			return domain.Project{}, fmt.Errorf("get team: %w", err)
		}
	}

	project := domain.ProjectDTO{
		Name:        name,
		Description: description,
		PublicKey:   publicKey,
		TeamID:      teamID,
	}

	id, err := s.projectRepo.Create(ctx, &project)
	if err != nil {
		return domain.Project{}, fmt.Errorf("create project: %w", err)
	}

	return domain.Project{
		ID:          id,
		Name:        name,
		Description: description,
		PublicKey:   publicKey,
		TeamID:      teamID,
		CreatedAt:   time.Now(),
	}, nil
}

func (s *ProjectService) List(ctx context.Context) ([]domain.ProjectExtended, error) {
	return s.projectRepo.List(ctx)
}

func (s *ProjectService) GetProjectsByUserID(
	ctx context.Context,
	userID domain.UserID,
	isSuperuser bool,
) ([]domain.ProjectExtended, error) {
	// Get all projects
	projects, err := s.projectRepo.List(ctx)
	if err != nil {
		return nil, fmt.Errorf("list projects: %w", err)
	}

	// If the user is a superuser, return all projects
	if isSuperuser {
		return projects, nil
	}

	// Get the teams that the user is a member of
	userTeams, err := s.teamsUseCase.GetTeamsByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("get user teams: %w", err)
	}

	// Create a map of team IDs for a quick lookup
	teamIDs := make(map[domain.TeamID]struct{}, len(userTeams))
	for _, team := range userTeams {
		teamIDs[team.ID] = struct{}{}
	}

	// Filter projects to only include those that belong to the user's teams
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

func (s *ProjectService) GeneralStats(
	ctx context.Context,
	id domain.ProjectID,
	period time.Duration,
) (domain.GeneralProjectStats, error) {
	_, err := s.projectRepo.GetByID(ctx, id)
	if err != nil {
		return domain.GeneralProjectStats{}, fmt.Errorf("get project: %w", err)
	}

	counters, err := s.issuesRepository.CountForLevels(ctx, id, period)
	if err != nil {
		return domain.GeneralProjectStats{}, fmt.Errorf("get issue counters: %w", err)
	}

	mostFrequestIssues, err := s.issuesRepository.MostFrequent(ctx, id, period, 6)
	if err != nil {
		return domain.GeneralProjectStats{}, fmt.Errorf("get most frequent issues: %w", err)
	}

	totalCnt := uint(0)
	for _, cnt := range counters {
		totalCnt += uint(cnt)
	}

	return domain.GeneralProjectStats{
		TotalIssues:        totalCnt,
		FatalIssues:        uint(counters[domain.IssueLevelFatal]),
		ErrorIssues:        uint(counters[domain.IssueLevelError]),
		WarningIssues:      uint(counters[domain.IssueLevelWarning]),
		InfoIssues:         uint(counters[domain.IssueLevelInfo]),
		DebugIssues:        uint(counters[domain.IssueLevelDebug]),
		ExceptionIssues:    uint(counters[domain.IssueLevelException]),
		MostFrequentIssues: mostFrequestIssues,
	}, nil
}

func (s *ProjectService) RecentProjects(ctx context.Context) ([]domain.ProjectExtended, error) {
	userID := wardencontext.UserID(ctx)

	return s.projectRepo.RecentProjects(ctx, userID, recentProjectsLimit)
}

func (s *ProjectService) UpdateInfo(
	ctx context.Context,
	id domain.ProjectID,
	name, description string,
) (domain.ProjectExtended, error) {
	// Check if the project exists
	project, err := s.projectRepo.GetByID(ctx, id)
	if err != nil {
		return domain.ProjectExtended{}, fmt.Errorf("failed to get project: %w", err)
	}

	// Update the project
	err = s.projectRepo.Update(ctx, id, name, description)
	if err != nil {
		return domain.ProjectExtended{}, fmt.Errorf("failed to update project: %w", err)
	}

	// Return the updated project with extended info
	project.Name = name
	project.Description = description

	projectExtended := domain.ProjectExtended{
		Project: project,
	}

	// If the project has a team, get the team name
	if project.TeamID != nil {
		team, err := s.teamsUseCase.GetTeamByID(ctx, *project.TeamID)
		if err != nil {
			// Log the error but don't fail the operation
			slog.Error("failed to get team name", "error", err, "team_id", *project.TeamID)
		} else {
			projectExtended.TeamName = &team.Name
		}
	}

	return projectExtended, nil
}

func (s *ProjectService) ArchiveProject(ctx context.Context, id domain.ProjectID) error {
	// Check if the project exists
	_, err := s.projectRepo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to get project: %w", err)
	}

	// Archive the project
	err = s.projectRepo.Archive(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to archive project: %w", err)
	}

	slog.Info("project archived", "project_id", id)

	return nil
}

func generateRandomKey(length int) (string, error) {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}

	return hex.EncodeToString(bytes), nil
}
