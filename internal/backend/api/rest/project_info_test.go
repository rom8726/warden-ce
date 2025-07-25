package rest

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/rom8726/warden/internal/domain"
	generatedapi "github.com/rom8726/warden/internal/generated/server"
	mockcontract "github.com/rom8726/warden/test_mocks/internal_/backend/contract"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestRestAPI_GetProject(t *testing.T) {
	t.Run("successful project info", func(t *testing.T) {
		mockProjectsUseCase := mockcontract.NewMockProjectsUseCase(t)
		mockPermissionsService := mockcontract.NewMockPermissionsService(t)

		api := &RestAPI{
			projectsUseCase:    mockProjectsUseCase,
			permissionsService: mockPermissionsService,
		}

		params := generatedapi.GetProjectParams{
			ProjectID: 123,
		}

		teamID := domain.TeamID(456)
		teamName := "Test Team"

		expectedProject := &domain.ProjectExtended{
			Project: domain.Project{
				ID:          domain.ProjectID(123),
				Name:        "Test Project",
				PublicKey:   "public_key",
				Description: "Test Description",
				TeamID:      &teamID,
				CreatedAt:   time.Now(),
			},
			TeamName: &teamName,
		}

		mockPermissionsService.EXPECT().
			CanAccessProject(mock.Anything, domain.ProjectID(123)).
			Return(nil)

		mockProjectsUseCase.EXPECT().
			GetProjectExtended(mock.Anything, domain.ProjectID(123)).
			Return(*expectedProject, nil)

		resp, err := api.GetProject(context.Background(), params)

		require.NoError(t, err)
		require.NotNil(t, resp)

		projectResp, ok := resp.(*generatedapi.ProjectResponse)
		require.True(t, ok)
		assert.Equal(t, uint(123), projectResp.Project.ID)
		assert.Equal(t, "Test Project", projectResp.Project.Name)
		assert.Equal(t, "public_key", projectResp.Project.PublicKey)
		assert.Equal(t, "Test Description", projectResp.Project.Description)
		assert.True(t, projectResp.Project.TeamID.Set)
		assert.Equal(t, uint(456), projectResp.Project.TeamID.Value)
		assert.True(t, projectResp.Project.TeamName.Set)
		assert.Equal(t, "Test Team", projectResp.Project.TeamName.Value)
	})

	t.Run("project without team", func(t *testing.T) {
		mockProjectsUseCase := mockcontract.NewMockProjectsUseCase(t)
		mockPermissionsService := mockcontract.NewMockPermissionsService(t)

		api := &RestAPI{
			projectsUseCase:    mockProjectsUseCase,
			permissionsService: mockPermissionsService,
		}

		params := generatedapi.GetProjectParams{
			ProjectID: 123,
		}

		expectedProject := &domain.ProjectExtended{
			Project: domain.Project{
				ID:          domain.ProjectID(123),
				Name:        "Test Project",
				PublicKey:   "public_key",
				Description: "Test Description",
				TeamID:      nil,
				CreatedAt:   time.Now(),
			},
		}

		mockPermissionsService.EXPECT().
			CanAccessProject(mock.Anything, domain.ProjectID(123)).
			Return(nil)

		mockProjectsUseCase.EXPECT().
			GetProjectExtended(mock.Anything, domain.ProjectID(123)).
			Return(*expectedProject, nil)

		resp, err := api.GetProject(context.Background(), params)

		require.NoError(t, err)
		require.NotNil(t, resp)

		projectResp, ok := resp.(*generatedapi.ProjectResponse)
		require.True(t, ok)
		assert.Equal(t, uint(123), projectResp.Project.ID)
		assert.Equal(t, "Test Project", projectResp.Project.Name)
		assert.Equal(t, "public_key", projectResp.Project.PublicKey)
		assert.False(t, projectResp.Project.TeamID.Set)
		assert.False(t, projectResp.Project.TeamName.Set)
	})

	t.Run("permission denied", func(t *testing.T) {
		mockPermissionsService := mockcontract.NewMockPermissionsService(t)

		api := &RestAPI{
			permissionsService: mockPermissionsService,
		}

		params := generatedapi.GetProjectParams{
			ProjectID: 123,
		}

		mockPermissionsService.EXPECT().
			CanAccessProject(mock.Anything, domain.ProjectID(123)).
			Return(domain.ErrPermissionDenied)

		resp, err := api.GetProject(context.Background(), params)

		require.NoError(t, err)
		require.NotNil(t, resp)

		errorResp, ok := resp.(*generatedapi.ErrorPermissionDenied)
		require.True(t, ok)
		assert.Equal(t, "permission denied", errorResp.Error.Message.Value)
	})

	t.Run("unauthorized", func(t *testing.T) {
		mockPermissionsService := mockcontract.NewMockPermissionsService(t)

		api := &RestAPI{
			permissionsService: mockPermissionsService,
		}

		params := generatedapi.GetProjectParams{
			ProjectID: 123,
		}

		mockPermissionsService.EXPECT().
			CanAccessProject(mock.Anything, domain.ProjectID(123)).
			Return(domain.ErrUserNotFound)

		resp, err := api.GetProject(context.Background(), params)

		require.NoError(t, err)
		require.NotNil(t, resp)

		errorResp, ok := resp.(*generatedapi.ErrorUnauthorized)
		require.True(t, ok)
		assert.Equal(t, "unauthorized", errorResp.Error.Message.Value)
	})

	t.Run("project not found", func(t *testing.T) {
		mockProjectsUseCase := mockcontract.NewMockProjectsUseCase(t)
		mockPermissionsService := mockcontract.NewMockPermissionsService(t)

		api := &RestAPI{
			projectsUseCase:    mockProjectsUseCase,
			permissionsService: mockPermissionsService,
		}

		params := generatedapi.GetProjectParams{
			ProjectID: 123,
		}

		mockPermissionsService.EXPECT().
			CanAccessProject(mock.Anything, domain.ProjectID(123)).
			Return(nil)

		mockProjectsUseCase.EXPECT().
			GetProjectExtended(mock.Anything, domain.ProjectID(123)).
			Return(domain.ProjectExtended{}, domain.ErrEntityNotFound)

		resp, err := api.GetProject(context.Background(), params)

		require.NoError(t, err)
		require.NotNil(t, resp)

		errorResp, ok := resp.(*generatedapi.ErrorNotFound)
		require.True(t, ok)
		assert.Equal(t, domain.ErrEntityNotFound.Error(), errorResp.Error.Message.Value)
	})

	t.Run("permission check failed with unexpected error", func(t *testing.T) {
		mockPermissionsService := mockcontract.NewMockPermissionsService(t)

		api := &RestAPI{
			permissionsService: mockPermissionsService,
		}

		params := generatedapi.GetProjectParams{
			ProjectID: 123,
		}

		unexpectedErr := errors.New("database error")
		mockPermissionsService.EXPECT().
			CanAccessProject(mock.Anything, domain.ProjectID(123)).
			Return(unexpectedErr)

		resp, err := api.GetProject(context.Background(), params)

		require.Error(t, err)
		require.Nil(t, resp)
		assert.Equal(t, unexpectedErr, err)
	})

	t.Run("get project failed with unexpected error", func(t *testing.T) {
		mockProjectsUseCase := mockcontract.NewMockProjectsUseCase(t)
		mockPermissionsService := mockcontract.NewMockPermissionsService(t)

		api := &RestAPI{
			projectsUseCase:    mockProjectsUseCase,
			permissionsService: mockPermissionsService,
		}

		params := generatedapi.GetProjectParams{
			ProjectID: 123,
		}

		unexpectedErr := errors.New("database error")
		mockPermissionsService.EXPECT().
			CanAccessProject(mock.Anything, domain.ProjectID(123)).
			Return(nil)

		mockProjectsUseCase.EXPECT().
			GetProjectExtended(mock.Anything, domain.ProjectID(123)).
			Return(domain.ProjectExtended{}, unexpectedErr)

		resp, err := api.GetProject(context.Background(), params)

		require.Error(t, err)
		require.Nil(t, resp)
		assert.Equal(t, unexpectedErr, err)
	})
}
