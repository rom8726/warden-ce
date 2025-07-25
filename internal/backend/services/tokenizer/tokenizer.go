// Package tokenizer creates and verifies JWT tokens
package tokenizer

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/golang-jwt/jwt"

	"github.com/rom8726/warden/internal/backend/contract"
	"github.com/rom8726/warden/internal/domain"
)

var (
	errUnexpectedMethod = errors.New("unexpected token signing method")
	errInvalidType      = errors.New("invalid token type")
	errInvalidToken     = errors.New("invalid token")
)

type Service struct {
	secretKey        []byte
	accessTTL        time.Duration
	refreshTTL       time.Duration
	resetPasswordTTL time.Duration
	permissionsSvc   contract.PermissionsService
	teamsUseCase     contract.TeamsUseCase
	projectsRepo     contract.ProjectsRepository
}

type ServiceParams struct {
	SecretKey                               []byte
	AccessTTL, RefreshTTL, ResetPasswordTTL time.Duration
}

func New(
	params *ServiceParams,
	permissionsSvc contract.PermissionsService,
	teamsUseCase contract.TeamsUseCase,
	projectsRepo contract.ProjectsRepository,
) *Service {
	return &Service{
		secretKey:        params.SecretKey,
		accessTTL:        params.AccessTTL,
		refreshTTL:       params.RefreshTTL,
		resetPasswordTTL: params.ResetPasswordTTL,
		permissionsSvc:   permissionsSvc,
		teamsUseCase:     teamsUseCase,
		projectsRepo:     projectsRepo,
	}
}

func (s *Service) SecretKey() string {
	return string(s.secretKey)
}

func (s *Service) AccessToken(user *domain.User) (string, error) {
	return s.generateToken(user, domain.TokenTypeAccess, s.accessTTL)
}

func (s *Service) RefreshToken(user *domain.User) (string, error) {
	return s.generateToken(user, domain.TokenTypeRefresh, s.refreshTTL)
}

func (s *Service) ResetPasswordToken(user *domain.User) (string, time.Duration, error) {
	token, err := s.generateToken(user, domain.TokenTypeResetPassword, s.resetPasswordTTL)
	if err != nil {
		return "", 0, err
	}

	return token, s.resetPasswordTTL, nil
}

func (s *Service) AccessTokenTTL() time.Duration {
	return s.accessTTL
}

func (s *Service) VerifyToken(token string, tokenType domain.TokenType) (*domain.TokenClaims, error) {
	claims, err := s.verifyToken(token, tokenType)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", domain.ErrInvalidToken, err) //nolint:errorlint // ok
	}

	return claims, nil
}

func (s *Service) verifyToken(token string, tokenType domain.TokenType) (*domain.TokenClaims, error) {
	parsed, err := jwt.ParseWithClaims(
		token,
		&domain.TokenClaims{},
		func(token *jwt.Token) (any, error) {
			_, ok := token.Method.(*jwt.SigningMethodHMAC)
			if !ok {
				return nil, errUnexpectedMethod
			}

			return s.secretKey, nil
		},
	)
	if err != nil {
		return nil, err
	}

	claims, ok := parsed.Claims.(*domain.TokenClaims)
	if !ok {
		return nil, errInvalidToken
	}

	if claims.TokenType != tokenType {
		return nil, errInvalidType
	}

	return claims, nil
}

//nolint:gocyclo // need refactor
func (s *Service) generateUserPermissions(ctx context.Context, user *domain.User) (domain.UserPermissions, error) {
	permissions := domain.UserPermissions{
		ProjectPermissions: make(map[domain.ProjectID]domain.ProjectPermission),
		TeamRoles:          make(map[domain.TeamID]domain.Role),
		CanCreateProjects:  user.IsSuperuser,
		CanCreateTeams:     user.IsSuperuser,
		CanManageUsers:     user.IsSuperuser,
	}

	if user.IsSuperuser {
		allProjects, err := s.projectsRepo.List(ctx)
		if err != nil {
			return permissions, fmt.Errorf("failed to get all projects for superuser: %w", err)
		}

		for _, project := range allProjects {
			projectPerm := domain.ProjectPermission{
				CanRead:   true,
				CanWrite:  true,
				CanDelete: true,
				CanManage: true,
				TeamRole:  domain.RoleOwner,
			}
			permissions.ProjectPermissions[project.ID] = projectPerm
		}

		allTeams, err := s.teamsUseCase.List(ctx)
		if err == nil {
			for _, team := range allTeams {
				permissions.TeamRoles[team.ID] = domain.RoleOwner
			}
		}

		return permissions, nil
	}

	userTeams, err := s.teamsUseCase.GetTeamsByUserID(ctx, user.ID)
	if err != nil {
		return permissions, fmt.Errorf("failed to get user teams: %w", err)
	}

	for _, team := range userTeams {
		members, err := s.teamsUseCase.GetMembers(ctx, team.ID)
		if err != nil {
			continue
		}

		var userRole domain.Role
		for _, member := range members {
			if member.UserID == user.ID {
				userRole = member.Role
				permissions.TeamRoles[team.ID] = member.Role

				break
			}
		}

		switch userRole {
		case domain.RoleOwner, domain.RoleAdmin:
			permissions.CanCreateProjects = true
			permissions.CanCreateTeams = true
		}
	}

	allProjects, err := s.projectsRepo.List(ctx)
	if err != nil {
		return permissions, fmt.Errorf("failed to get all projects: %w", err)
	}

	accessibleProjects, err := s.permissionsSvc.GetAccessibleProjects(ctx, allProjects)
	if err != nil {
		return permissions, fmt.Errorf("failed to get accessible projects: %w", err)
	}

	for _, project := range accessibleProjects {
		if project.TeamID != nil {
			userRole, hasRole := permissions.TeamRoles[*project.TeamID]
			if hasRole {
				projectPerm := domain.ProjectPermission{
					CanRead:   true,
					CanWrite:  userRole == domain.RoleOwner || userRole == domain.RoleAdmin,
					CanDelete: userRole == domain.RoleOwner || userRole == domain.RoleAdmin,
					CanManage: userRole == domain.RoleOwner || userRole == domain.RoleAdmin,
					TeamRole:  userRole,
				}
				permissions.ProjectPermissions[project.ID] = projectPerm
			}
		} else {
			projectPerm := domain.ProjectPermission{
				CanRead:   true,
				CanWrite:  false,
				CanDelete: false,
				CanManage: false,
			}
			permissions.ProjectPermissions[project.ID] = projectPerm
		}
	}

	return permissions, nil
}

func (s *Service) generateToken(user *domain.User, tokenType domain.TokenType, ttl time.Duration) (string, error) {
	now := time.Now().UTC()

	var permissions domain.UserPermissions
	if tokenType == domain.TokenTypeAccess {
		ctx := context.Background()
		var err error
		permissions, err = s.generateUserPermissions(ctx, user)
		if err != nil {
			slog.Error("failed to generate user permissions", "user_id", user.ID, "error", err)
		}
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &domain.TokenClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: now.Add(ttl).Unix(),
			IssuedAt:  now.Unix(),
		},
		TokenType:   tokenType,
		UserID:      uint(user.ID),
		Username:    user.Username,
		IsSuperuser: user.IsSuperuser,
		Permissions: permissions,
	})

	return token.SignedString(s.secretKey)
}
