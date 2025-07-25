package tokenizer

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/golang-jwt/jwt/test"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/rom8726/warden/internal/domain"
	mockcontract "github.com/rom8726/warden/test_mocks/internal_/backend/contract"
)

func TestService_AccessToken(t *testing.T) {
	mockPermissionsSvc := mockcontract.NewMockPermissionsService(t)
	mockTeamsUseCase := mockcontract.NewMockTeamsUseCase(t)
	mockProjectsRepo := mockcontract.NewMockProjectsRepository(t)

	mockProjectsRepo.EXPECT().List(mock.Anything).Return([]domain.ProjectExtended{}, nil).Maybe()
	mockTeamsUseCase.EXPECT().GetTeamsByUserID(mock.Anything, mock.Anything).Return([]domain.Team{}, nil).Maybe()
	mockPermissionsSvc.EXPECT().GetAccessibleProjects(mock.Anything, mock.Anything).Return([]domain.ProjectExtended{}, nil).Maybe()

	srv := New(&ServiceParams{
		SecretKey:        []byte("secret"),
		AccessTTL:        time.Minute,
		RefreshTTL:       time.Minute,
		ResetPasswordTTL: time.Minute,
	}, mockPermissionsSvc, mockTeamsUseCase, mockProjectsRepo)
	user := domain.User{
		ID: 123,
	}
	token, err := srv.AccessToken(&user)
	require.NoError(t, err)
	require.NotEmpty(t, token)
}

func TestService_RefreshToken(t *testing.T) {
	mockPermissionsSvc := mockcontract.NewMockPermissionsService(t)
	mockTeamsUseCase := mockcontract.NewMockTeamsUseCase(t)
	mockProjectsRepo := mockcontract.NewMockProjectsRepository(t)

	mockProjectsRepo.EXPECT().List(mock.Anything).Return([]domain.ProjectExtended{}, nil).Maybe()
	mockTeamsUseCase.EXPECT().GetTeamsByUserID(mock.Anything, mock.Anything).Return([]domain.Team{}, nil).Maybe()
	mockPermissionsSvc.EXPECT().GetAccessibleProjects(mock.Anything, mock.Anything).Return([]domain.ProjectExtended{}, nil).Maybe()

	srv := New(&ServiceParams{
		SecretKey:        []byte("secret"),
		AccessTTL:        time.Minute,
		RefreshTTL:       time.Minute,
		ResetPasswordTTL: time.Minute,
	}, mockPermissionsSvc, mockTeamsUseCase, mockProjectsRepo)
	user := domain.User{
		ID: 123,
	}
	token, err := srv.RefreshToken(&user)
	require.NoError(t, err)
	require.NotEmpty(t, token)
}

func TestService_ResetPasswordToken(t *testing.T) {
	mockPermissionsSvc := mockcontract.NewMockPermissionsService(t)
	mockTeamsUseCase := mockcontract.NewMockTeamsUseCase(t)
	mockProjectsRepo := mockcontract.NewMockProjectsRepository(t)

	mockProjectsRepo.EXPECT().List(mock.Anything).Return([]domain.ProjectExtended{}, nil).Maybe()
	mockTeamsUseCase.EXPECT().GetTeamsByUserID(mock.Anything, mock.Anything).Return([]domain.Team{}, nil).Maybe()
	mockPermissionsSvc.EXPECT().GetAccessibleProjects(mock.Anything, mock.Anything).Return([]domain.ProjectExtended{}, nil).Maybe()

	srv := New(&ServiceParams{
		SecretKey:        []byte("secret"),
		AccessTTL:        time.Minute,
		RefreshTTL:       time.Minute,
		ResetPasswordTTL: time.Minute,
	}, mockPermissionsSvc, mockTeamsUseCase, mockProjectsRepo)
	user := domain.User{
		ID: 123,
	}
	token, _, err := srv.ResetPasswordToken(&user)
	require.NoError(t, err)
	require.NotEmpty(t, token)
}

func TestService_VerifyToken(t *testing.T) {
	mockPermissionsSvc := mockcontract.NewMockPermissionsService(t)
	mockTeamsUseCase := mockcontract.NewMockTeamsUseCase(t)
	mockProjectsRepo := mockcontract.NewMockProjectsRepository(t)

	mockProjectsRepo.EXPECT().List(mock.Anything).Return([]domain.ProjectExtended{}, nil).Maybe()
	mockTeamsUseCase.EXPECT().GetTeamsByUserID(mock.Anything, mock.Anything).Return([]domain.Team{}, nil).Maybe()
	mockPermissionsSvc.EXPECT().GetAccessibleProjects(mock.Anything, mock.Anything).Return([]domain.ProjectExtended{}, nil).Maybe()

	srv := New(&ServiceParams{
		SecretKey:        []byte("secret"),
		AccessTTL:        time.Second,
		RefreshTTL:       time.Second,
		ResetPasswordTTL: time.Minute,
	}, mockPermissionsSvc, mockTeamsUseCase, mockProjectsRepo)
	user := domain.User{
		ID: 123,
	}

	accessToken, err := srv.AccessToken(&user)
	require.NoError(t, err)

	refreshToken, err := srv.RefreshToken(&user)
	require.NoError(t, err)

	t.Run("valid access token", func(t *testing.T) {
		claims, err := srv.VerifyToken(accessToken, domain.TokenTypeAccess)
		require.NoError(t, err)
		require.Equal(t, user.ID, domain.UserID(claims.UserID))
		require.Equal(t, domain.TokenTypeAccess, claims.TokenType)
	})

	t.Run("valid refresh token", func(t *testing.T) {
		claims, err := srv.VerifyToken(refreshToken, domain.TokenTypeRefresh)
		require.NoError(t, err)
		require.Equal(t, user.ID, domain.UserID(claims.UserID))
		require.Equal(t, domain.TokenTypeRefresh, claims.TokenType)
	})

	t.Run("valid reset password token", func(t *testing.T) {
		resetPasswordToken, _, err := srv.ResetPasswordToken(&user)
		require.NoError(t, err)

		claims, err := srv.VerifyToken(resetPasswordToken, domain.TokenTypeResetPassword)
		require.NoError(t, err)
		require.Equal(t, user.ID, domain.UserID(claims.UserID))
		require.Equal(t, domain.TokenTypeResetPassword, claims.TokenType)
	})

	t.Run("wrong token type", func(t *testing.T) {
		claims, err := srv.VerifyToken(refreshToken, domain.TokenTypeAccess)
		require.Error(t, err)
		require.ErrorIs(t, err, domain.ErrInvalidToken)
		require.Nil(t, claims)
	})

	t.Run("expired access token", func(t *testing.T) {
		time.Sleep(time.Second * 2)
		claims, err := srv.VerifyToken(accessToken, domain.TokenTypeAccess)
		require.Error(t, err)
		require.ErrorIs(t, err, domain.ErrInvalidToken)
		require.Nil(t, claims)
	})

	t.Run("wrong signing method", func(t *testing.T) {
		now := time.Now()
		token, err := jwt.NewWithClaims(jwt.SigningMethodPS256, &domain.TokenClaims{
			StandardClaims: jwt.StandardClaims{
				ExpiresAt: now.Add(time.Minute).Unix(),
				IssuedAt:  now.Unix(),
			},
			TokenType: domain.TokenTypeAccess,
			UserID:    uint(user.ID),
		}).SignedString(test.LoadRSAPrivateKeyFromDisk("test/sample_key"))
		require.NoError(t, err)
		require.NotEmpty(t, token)

		claims, err := srv.VerifyToken(token, domain.TokenTypeAccess)
		require.Error(t, err)
		require.ErrorIs(t, err, domain.ErrInvalidToken)
		require.Nil(t, claims)
	})
}

func TestService_GenerateUserPermissions(t *testing.T) {
	t.Run("superuser permissions", func(t *testing.T) {
		mockPermissionsSvc := mockcontract.NewMockPermissionsService(t)
		mockTeamsUseCase := mockcontract.NewMockTeamsUseCase(t)
		mockProjectsRepo := mockcontract.NewMockProjectsRepository(t)

		// Setup mocks for superuser
		mockProjectsRepo.EXPECT().
			List(mock.Anything).
			Return([]domain.ProjectExtended{}, nil)

		mockTeamsUseCase.EXPECT().
			List(mock.Anything).
			Return([]domain.Team{}, nil)

		srv := New(&ServiceParams{
			SecretKey:        []byte("secret"),
			AccessTTL:        time.Minute,
			RefreshTTL:       time.Minute,
			ResetPasswordTTL: time.Minute,
		}, mockPermissionsSvc, mockTeamsUseCase, mockProjectsRepo)

		user := &domain.User{
			ID:          123,
			Username:    "admin",
			IsSuperuser: true,
		}

		// Generate access token for superuser
		token, err := srv.AccessToken(user)
		require.NoError(t, err)
		require.NotEmpty(t, token)

		// Check claims
		claims, err := srv.VerifyToken(token, domain.TokenTypeAccess)
		require.NoError(t, err)
		require.True(t, claims.IsSuperuser)
		require.True(t, claims.Permissions.CanCreateProjects)
		require.True(t, claims.Permissions.CanCreateTeams)
		require.True(t, claims.Permissions.CanManageUsers)
	})

	t.Run("regular user with team permissions", func(t *testing.T) {
		mockPermissionsSvc := mockcontract.NewMockPermissionsService(t)
		mockTeamsUseCase := mockcontract.NewMockTeamsUseCase(t)
		mockProjectsRepo := mockcontract.NewMockProjectsRepository(t)

		srv := New(&ServiceParams{
			SecretKey:        []byte("secret"),
			AccessTTL:        time.Minute,
			RefreshTTL:       time.Minute,
			ResetPasswordTTL: time.Minute,
		}, mockPermissionsSvc, mockTeamsUseCase, mockProjectsRepo)

		user := &domain.User{
			ID:          456,
			Username:    "user",
			IsSuperuser: false,
		}

		teamID := domain.TeamID(1)
		projectID := domain.ProjectID(1)

		// Setup mocks
		mockTeamsUseCase.EXPECT().
			GetTeamsByUserID(mock.Anything, user.ID).
			Return([]domain.Team{{ID: teamID, Name: "Test Team"}}, nil)

		mockTeamsUseCase.EXPECT().
			GetMembers(mock.Anything, teamID).
			Return([]domain.TeamMember{
				{UserID: user.ID, TeamID: teamID, Role: domain.RoleAdmin},
			}, nil)

		mockProjectsRepo.EXPECT().
			List(mock.Anything).
			Return([]domain.ProjectExtended{
				{
					Project: domain.Project{
						ID:     projectID,
						Name:   "Test Project",
						TeamID: &teamID,
					},
				},
			}, nil)

		mockPermissionsSvc.EXPECT().
			GetAccessibleProjects(mock.Anything, mock.Anything).
			Return([]domain.ProjectExtended{
				{
					Project: domain.Project{
						ID:     projectID,
						Name:   "Test Project",
						TeamID: &teamID,
					},
				},
			}, nil)

		// Generate access token
		token, err := srv.AccessToken(user)
		require.NoError(t, err)
		require.NotEmpty(t, token)

		// Check claims
		claims, err := srv.VerifyToken(token, domain.TokenTypeAccess)
		require.NoError(t, err)
		require.False(t, claims.IsSuperuser)
		require.True(t, claims.Permissions.CanCreateProjects) // Admin can create projects
		require.True(t, claims.Permissions.CanCreateTeams)    // Admin can create teams
		require.False(t, claims.Permissions.CanManageUsers)   // Regular user cannot manage users

		// Check team roles
		require.Equal(t, domain.RoleAdmin, claims.Permissions.TeamRoles[teamID])

		// Check project permissions
		projectPerm, exists := claims.Permissions.ProjectPermissions[projectID]
		require.True(t, exists)
		require.True(t, projectPerm.CanRead)
		require.True(t, projectPerm.CanWrite)
		require.True(t, projectPerm.CanDelete)
		require.True(t, projectPerm.CanManage)
		require.Equal(t, domain.RoleAdmin, projectPerm.TeamRole)
	})
}
