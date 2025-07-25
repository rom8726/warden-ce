package teams

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/rom8726/warden/internal/backend/contract"
	wardencontext "github.com/rom8726/warden/internal/context"
	"github.com/rom8726/warden/internal/domain"
	"github.com/rom8726/warden/pkg/db"
)

type TeamService struct {
	txManager                db.TxManager
	teamsRepo                contract.TeamsRepository
	usersRepo                contract.UsersRepository
	userNotificationsUseCase contract.UserNotificationsUseCase
	projectsRepo             contract.ProjectsRepository
}

func New(
	txManager db.TxManager,
	teamsRepo contract.TeamsRepository,
	usersRepo contract.UsersRepository,
	userNotificationsUseCase contract.UserNotificationsUseCase,
	projectsRepo contract.ProjectsRepository,
) *TeamService {
	return &TeamService{
		txManager:                txManager,
		teamsRepo:                teamsRepo,
		usersRepo:                usersRepo,
		userNotificationsUseCase: userNotificationsUseCase,
		projectsRepo:             projectsRepo,
	}
}

func (s *TeamService) Create(ctx context.Context, teamDTO domain.TeamDTO) (domain.Team, error) {
	// Get the current user from context
	currentUserID := wardencontext.UserID(ctx)
	ok, err := s.usersRepo.ExistsByID(ctx, currentUserID)
	if err != nil {
		return domain.Team{}, fmt.Errorf("get current user by id: %w", err)
	}

	if !ok {
		return domain.Team{}, domain.ErrUserNotFound
	}

	_, err = s.teamsRepo.GetByName(ctx, teamDTO.Name)
	if err == nil {
		return domain.Team{}, domain.ErrTeamNameAlreadyInUse
	}

	// Create the team
	team, err := s.teamsRepo.Create(ctx, teamDTO)
	if err != nil {
		return domain.Team{}, err
	}

	// Add the current user as an owner of the team
	err = s.teamsRepo.AddMember(ctx, team.ID, currentUserID, domain.RoleOwner)
	if err != nil {
		return domain.Team{}, fmt.Errorf("add current user as owner: %w", err)
	}

	// Refresh the team data to include the new member
	return s.teamsRepo.GetByID(ctx, team.ID)
}

func (s *TeamService) GetByID(ctx context.Context, id domain.TeamID) (domain.Team, error) {
	// Get the current user from context
	currentUserID := wardencontext.UserID(ctx)
	currentUser, err := s.usersRepo.GetByID(ctx, currentUserID)
	if err != nil {
		return domain.Team{}, fmt.Errorf("get current user by id: %w", err)
	}

	// Get the team
	team, err := s.teamsRepo.GetByID(ctx, id)
	if err != nil {
		return domain.Team{}, err
	}

	// Superusers can access any team
	if currentUser.IsSuperuser {
		return team, nil
	}

	// Check if the current user is a member of the team
	isMember := false
	for _, member := range team.Members {
		if member.UserID == currentUserID {
			isMember = true

			break
		}
	}
	if !isMember {
		return domain.Team{}, domain.ErrForbidden
	}

	return team, nil
}

func (s *TeamService) List(ctx context.Context) ([]domain.Team, error) {
	// Get the current user from context
	currentUserID := wardencontext.UserID(ctx)
	currentUser, err := s.usersRepo.GetByID(ctx, currentUserID)
	if err != nil {
		return nil, fmt.Errorf("get current user by id: %w", err)
	}

	// Superusers can see all teams
	if currentUser.IsSuperuser {
		return s.teamsRepo.List(ctx)
	}

	// Regular users can only see teams they are members of
	return s.teamsRepo.GetTeamsByUserID(ctx, currentUserID)
}

func (s *TeamService) Delete(ctx context.Context, id domain.TeamID) error {
	// Get the current user from context
	currentUserID := wardencontext.UserID(ctx)
	currentUser, err := s.usersRepo.GetByID(ctx, currentUserID)
	if err != nil {
		return fmt.Errorf("get current user by id: %w", err)
	}

	// Get the team
	team, err := s.teamsRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	// Check if the current user is a superuser or has owner role in the team
	if !currentUser.IsSuperuser {
		hasPermission := false
		for _, member := range team.Members {
			if member.UserID == currentUserID && member.Role == domain.RoleOwner {
				hasPermission = true

				break
			}
		}
		if !hasPermission {
			return domain.ErrForbidden
		}
	}

	projects, err := s.projectsRepo.List(ctx)
	if err != nil {
		return fmt.Errorf("check team projects: %w", err)
	}
	for _, project := range projects {
		if project.TeamID != nil && *project.TeamID == id {
			return domain.ErrTeamHasProjects
		}
	}

	return s.teamsRepo.Delete(ctx, id)
}

func (s *TeamService) AddMember(
	ctx context.Context,
	teamID domain.TeamID,
	userID domain.UserID,
	role domain.Role,
) error {
	// Get the current user from context
	currentUserID := wardencontext.UserID(ctx)
	currentUser, err := s.usersRepo.GetByID(ctx, currentUserID)
	if err != nil {
		return fmt.Errorf("get current user by id: %w", err)
	}

	// Check if a user to add exists
	_, err = s.usersRepo.GetByID(ctx, userID)
	if err != nil {
		return err
	}

	// Get the team
	team, err := s.teamsRepo.GetByID(ctx, teamID)
	if err != nil {
		return err
	}

	// Check if the current user is a superuser or has an admin/owner role in the team
	if !currentUser.IsSuperuser {
		hasPermission := false
		for _, member := range team.Members {
			if member.UserID == currentUserID && (member.Role == domain.RoleOwner || member.Role == domain.RoleAdmin) {
				hasPermission = true

				break
			}
		}
		if !hasPermission {
			return domain.ErrForbidden
		}
	}

	return s.txManager.ReadCommitted(ctx, func(ctx context.Context) error {
		// Add member to team
		err = s.teamsRepo.AddMember(ctx, teamID, userID, role)
		if err != nil {
			return err
		}

		// Create a notification for the added user
		content := domain.UserNotificationContent{
			TeamAdded: &domain.TeamAddedContent{
				TeamID:          uint(teamID),
				TeamName:        team.Name,
				Role:            string(role),
				AddedByUserID:   uint(currentUserID),
				AddedByUsername: currentUser.Username,
			},
		}

		err = s.userNotificationsUseCase.CreateNotification(ctx, userID, domain.UserNotificationTypeTeamAdded, content)
		if err != nil {
			// Log error but don't fail the operation
			return fmt.Errorf("create user notification: %w", err)
		}

		return nil
	})
}

// nolint:gocyclo // need refactoring
func (s *TeamService) RemoveMemberWithChecks(ctx context.Context, teamID domain.TeamID, userID domain.UserID) error {
	currentUserID := wardencontext.UserID(ctx)
	currentUser, err := s.getUserOrError(ctx, currentUserID, "get current user by id")
	if err != nil {
		return err
	}

	// Get the team
	team, err := s.teamsRepo.GetByID(ctx, teamID)
	if err != nil {
		return err
	}

	if currentUserID == userID {
		isOwner, ownerCount := s.isOwnerAndOwnerCount(team, currentUserID)
		if isOwner && ownerCount == 1 {
			return domain.ErrLastOwner
		}

		if currentUser.IsSuperuser {
			return domain.ErrForbidden
		}

		return s.txManager.ReadCommitted(ctx, func(ctx context.Context) error {
			if err := s.teamsRepo.RemoveMember(ctx, teamID, userID); err != nil {
				return err
			}

			return s.notifyTeamRemoved(ctx, team, userID, currentUserID, currentUser.Username)
		})
	}

	// Check if the current user is a superuser or has an admin/owner role in the team
	if !currentUser.IsSuperuser {
		hasPermission := false
		for _, member := range team.Members {
			if member.UserID == currentUserID && (member.Role == domain.RoleOwner || member.Role == domain.RoleAdmin) {
				hasPermission = true

				break
			}
		}
		if !hasPermission {
			return domain.ErrForbidden
		}
	}

	// Check if the user to remove is an owner and the current user is not an owner
	if !currentUser.IsSuperuser {
		isUserOwner := false
		isCurrentUserOwner := false
		for _, member := range team.Members {
			if member.UserID == userID && member.Role == domain.RoleOwner {
				isUserOwner = true
			}
			if member.UserID == currentUserID && member.Role == domain.RoleOwner {
				isCurrentUserOwner = true
			}
		}
		if isUserOwner && !isCurrentUserOwner {
			return domain.ErrForbidden
		}
	}

	user, err := s.getUserOrError(ctx, userID, "get user by id")
	if err != nil {
		return err
	}

	if user.IsSuperuser && user.ID == currentUserID {
		return domain.ErrForbidden
	}

	return s.txManager.ReadCommitted(ctx, func(ctx context.Context) error {
		// Remove a member from a team
		err = s.teamsRepo.RemoveMember(ctx, teamID, userID)
		if err != nil {
			return err
		}

		// Create a notification for the removed user
		return s.notifyTeamRemoved(ctx, team, userID, currentUserID, currentUser.Username)
	})
}

func (s *TeamService) GetTeamsByUserID(ctx context.Context, userID domain.UserID) ([]domain.Team, error) {
	return s.teamsRepo.GetTeamsByUserID(ctx, userID)
}

func (s *TeamService) GetTeamByID(ctx context.Context, id domain.TeamID) (domain.Team, error) {
	return s.teamsRepo.GetByID(ctx, id)
}

func (s *TeamService) GetByName(ctx context.Context, name string) (domain.Team, error) {
	return s.teamsRepo.GetByName(ctx, name)
}

func (s *TeamService) GetByProjectID(ctx context.Context, projectID domain.ProjectID) (domain.Team, error) {
	// Get the current user from context
	currentUserID := wardencontext.UserID(ctx)
	currentUser, err := s.usersRepo.GetByID(ctx, currentUserID)
	if err != nil {
		return domain.Team{}, fmt.Errorf("get current user by id: %w", err)
	}

	// Get the team by project ID
	team, err := s.teamsRepo.GetByProjectID(ctx, projectID)
	if err != nil {
		return domain.Team{}, err
	}

	// Superusers can access any team
	if currentUser.IsSuperuser {
		return team, nil
	}

	// Check if the current user is a member of the team
	isMember := false
	for _, member := range team.Members {
		if member.UserID == currentUserID {
			isMember = true

			break
		}
	}
	if !isMember {
		return domain.Team{}, domain.ErrForbidden
	}

	return team, nil
}

func (s *TeamService) GetMembers(ctx context.Context, teamID domain.TeamID) ([]domain.TeamMember, error) {
	return s.teamsRepo.GetMembers(ctx, teamID)
}

func (s *TeamService) isOwnerAndOwnerCount(team domain.Team, userID domain.UserID) (isOwner bool, ownerCount int) {
	for _, member := range team.Members {
		if member.Role == domain.RoleOwner {
			ownerCount++
		}
		if member.UserID == userID && member.Role == domain.RoleOwner {
			isOwner = true
		}
	}

	return
}

func (s *TeamService) notifyTeamRemoved(
	ctx context.Context,
	team domain.Team,
	userID, removedByUserID domain.UserID,
	removedByUsername string,
) error {
	content := domain.UserNotificationContent{
		TeamRemoved: &domain.TeamRemovedContent{
			TeamID:            uint(team.ID),
			TeamName:          team.Name,
			RemovedByUserID:   uint(removedByUserID),
			RemovedByUsername: removedByUsername,
		},
	}

	return s.userNotificationsUseCase.CreateNotification(ctx, userID, domain.UserNotificationTypeTeamRemoved, content)
}

func (s *TeamService) getUserOrError(ctx context.Context, userID domain.UserID, msg string) (domain.User, error) {
	user, err := s.usersRepo.GetByID(ctx, userID)
	if err != nil {
		return domain.User{}, fmt.Errorf("%s: %w", msg, err)
	}

	return user, nil
}

func (s *TeamService) ChangeMemberRole(
	ctx context.Context,
	teamID domain.TeamID,
	userID domain.UserID,
	newRole domain.Role,
) error {
	// Get the current user from context
	currentUserID := wardencontext.UserID(ctx)
	currentUser, err := s.usersRepo.GetByID(ctx, currentUserID)
	if err != nil {
		return fmt.Errorf("get current user by id: %w", err)
	}

	// Get the team
	team, err := s.teamsRepo.GetByID(ctx, teamID)
	if err != nil {
		return err
	}

	// Find the target member
	var targetMember *domain.TeamMember
	for i := range team.Members {
		if team.Members[i].UserID == userID {
			targetMember = &team.Members[i]

			break
		}
	}
	if targetMember == nil {
		return domain.ErrEntityNotFound
	}

	// Check permissions
	if err := s.validateRoleChangePermissions(currentUser, team, targetMember, newRole); err != nil {
		return err
	}

	// Check constraints
	if err := s.validateRoleChangeConstraints(team, targetMember, newRole); err != nil {
		return err
	}

	return s.txManager.ReadCommitted(ctx, func(ctx context.Context) error {
		// Update the role
		err = s.teamsRepo.UpdateMemberRole(ctx, teamID, userID, newRole)
		if err != nil {
			return err
		}

		// If this is ownership transfer, demote the old owner to admin
		if newRole == domain.RoleOwner && targetMember.Role != domain.RoleOwner {
			oldOwner := s.findOwner(team.Members)
			if oldOwner != nil && oldOwner.UserID != userID {
				err = s.teamsRepo.UpdateMemberRole(ctx, teamID, oldOwner.UserID, domain.RoleAdmin)
				if err != nil {
					return fmt.Errorf("update old owner role: %w", err)
				}
			}
		}

		// Create notification for the user
		content := domain.UserNotificationContent{
			RoleChanged: &domain.RoleChangedContent{
				TeamID:            uint(teamID),
				TeamName:          team.Name,
				OldRole:           string(targetMember.Role),
				NewRole:           string(newRole),
				ChangedByUserID:   uint(currentUserID),
				ChangedByUsername: currentUser.Username,
			},
		}

		err = s.userNotificationsUseCase.CreateNotification(
			ctx,
			userID,
			domain.UserNotificationTypeRoleChanged,
			content,
		)
		if err != nil {
			// Log error but don't fail the operation
			slog.Error("create user notification", "error", err)
		}

		return nil
	})
}

func (s *TeamService) validateRoleChangePermissions(
	currentUser domain.User,
	team domain.Team,
	targetMember *domain.TeamMember,
	newRole domain.Role,
) error {
	// Superusers can change any role
	if currentUser.IsSuperuser {
		return nil
	}

	// Find current user's role in the team
	var currentUserRole domain.Role
	for _, member := range team.Members {
		if member.UserID == currentUser.ID {
			currentUserRole = member.Role

			break
		}
	}

	// Only owners and admins can change roles
	if currentUserRole != domain.RoleOwner && currentUserRole != domain.RoleAdmin {
		return domain.ErrForbidden
	}

	// Admins cannot change owner roles
	if targetMember.Role == domain.RoleOwner && currentUserRole != domain.RoleOwner {
		return domain.ErrForbidden
	}

	// Admins cannot promote to owner
	if newRole == domain.RoleOwner && currentUserRole != domain.RoleOwner {
		return domain.ErrForbidden
	}

	return nil
}

func (s *TeamService) validateRoleChangeConstraints(
	team domain.Team,
	targetMember *domain.TeamMember,
	newRole domain.Role,
) error {
	// Cannot demote the only owner
	if targetMember.Role == domain.RoleOwner && newRole != domain.RoleOwner {
		ownerCount := 0
		for _, member := range team.Members {
			if member.Role == domain.RoleOwner {
				ownerCount++
			}
		}
		if ownerCount == 1 {
			return domain.ErrLastOwner
		}
	}

	return nil
}

func (s *TeamService) findOwner(members []domain.TeamMember) *domain.TeamMember {
	for i := range members {
		if members[i].Role == domain.RoleOwner {
			return &members[i]
		}
	}

	return nil
}
