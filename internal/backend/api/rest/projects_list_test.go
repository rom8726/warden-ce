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

func TestRestAPI_ListProjects(t *testing.T) {
	t.Run("successful projects list", func(t *testing.T) {
		mockProjectsUseCase := mockcontract.NewMockProjectsUseCase(t)
		mockPermissionsService := mockcontract.NewMockPermissionsService(t)

		api := &RestAPI{
			projectsUseCase:    mockProjectsUseCase,
			permissionsService: mockPermissionsService,
		}

		teamID := domain.TeamID(1)
		allProjects := []domain.ProjectExtended{
			{
				Project: domain.Project{
					ID:          domain.ProjectID(1),
					Name:        "Project 1",
					PublicKey:   "public_key1",
					Description: "Description 1",
					TeamID:      nil,
					CreatedAt:   time.Now(),
				},
			},
			{
				Project: domain.Project{
					ID:          domain.ProjectID(2),
					Name:        "Project 2",
					PublicKey:   "public_key2",
					Description: "Description 2",
					TeamID:      &teamID,
					CreatedAt:   time.Now(),
				},
			},
		}

		accessibleProjects := []domain.ProjectExtended{
			allProjects[0],
			allProjects[1],
		}

		mockProjectsUseCase.EXPECT().
			List(mock.Anything).
			Return(allProjects, nil)

		mockPermissionsService.EXPECT().
			GetAccessibleProjects(mock.Anything, allProjects).
			Return(accessibleProjects, nil)

		resp, err := api.ListProjects(context.Background())

		require.NoError(t, err)
		require.NotNil(t, resp)

		listResp, ok := resp.(*generatedapi.ListProjectsResponse)
		require.True(t, ok)
		assert.Len(t, *listResp, 2)

		// Check first project (no team)
		assert.Equal(t, uint(1), (*listResp)[0].ID)
		assert.Equal(t, "Project 1", (*listResp)[0].Name)
		assert.Equal(t, "public_key1", (*listResp)[0].PublicKey)
		assert.Equal(t, "Description 1", (*listResp)[0].Description)
		assert.False(t, (*listResp)[0].TeamID.Set)
		assert.False(t, (*listResp)[0].TeamName.Set)

		// Check second project (with team)
		assert.Equal(t, uint(2), (*listResp)[1].ID)
		assert.Equal(t, "Project 2", (*listResp)[1].Name)
		assert.Equal(t, "public_key2", (*listResp)[1].PublicKey)
		assert.Equal(t, "Description 2", (*listResp)[1].Description)
		assert.True(t, (*listResp)[1].TeamID.Set)
		assert.Equal(t, uint(1), (*listResp)[1].TeamID.Value)
		assert.False(t, (*listResp)[1].TeamName.Set)
	})

	t.Run("empty projects list", func(t *testing.T) {
		mockProjectsUseCase := mockcontract.NewMockProjectsUseCase(t)
		mockPermissionsService := mockcontract.NewMockPermissionsService(t)

		api := &RestAPI{
			projectsUseCase:    mockProjectsUseCase,
			permissionsService: mockPermissionsService,
		}

		allProjects := []domain.ProjectExtended{}
		accessibleProjects := []domain.ProjectExtended{}

		mockProjectsUseCase.EXPECT().
			List(mock.Anything).
			Return(allProjects, nil)

		mockPermissionsService.EXPECT().
			GetAccessibleProjects(mock.Anything, allProjects).
			Return(accessibleProjects, nil)

		resp, err := api.ListProjects(context.Background())

		require.NoError(t, err)
		require.NotNil(t, resp)

		listResp, ok := resp.(*generatedapi.ListProjectsResponse)
		require.True(t, ok)
		assert.Len(t, *listResp, 0)
	})

	t.Run("get all projects failed", func(t *testing.T) {
		mockProjectsUseCase := mockcontract.NewMockProjectsUseCase(t)

		api := &RestAPI{
			projectsUseCase: mockProjectsUseCase,
		}

		unexpectedErr := errors.New("database error")
		mockProjectsUseCase.EXPECT().
			List(mock.Anything).
			Return(nil, unexpectedErr)

		resp, err := api.ListProjects(context.Background())

		require.Error(t, err)
		require.Nil(t, resp)
		assert.Equal(t, unexpectedErr, err)
	})

	t.Run("filter projects failed", func(t *testing.T) {
		mockProjectsUseCase := mockcontract.NewMockProjectsUseCase(t)
		mockPermissionsService := mockcontract.NewMockPermissionsService(t)

		api := &RestAPI{
			projectsUseCase:    mockProjectsUseCase,
			permissionsService: mockPermissionsService,
		}

		allProjects := []domain.ProjectExtended{
			{
				Project: domain.Project{
					ID:        domain.ProjectID(1),
					Name:      "Project 1",
					CreatedAt: time.Now(),
				},
			},
		}

		unexpectedErr := errors.New("permission error")
		mockProjectsUseCase.EXPECT().
			List(mock.Anything).
			Return(allProjects, nil)

		mockPermissionsService.EXPECT().
			GetAccessibleProjects(mock.Anything, allProjects).
			Return(nil, unexpectedErr)

		resp, err := api.ListProjects(context.Background())

		require.Error(t, err)
		require.Nil(t, resp)
		assert.Equal(t, unexpectedErr, err)
	})
}

func TestConvertTeamID(t *testing.T) {
	t.Run("nil team ID", func(t *testing.T) {
		result := convertTeamID(nil)
		assert.False(t, result.Set)
	})

	t.Run("valid team ID", func(t *testing.T) {
		teamID := domain.TeamID(123)
		result := convertTeamID(&teamID)
		assert.True(t, result.Set)
		assert.Equal(t, uint(123), result.Value)
	})
}

func TestConvertTeamName(t *testing.T) {
	t.Run("nil team name", func(t *testing.T) {
		result := convertTeamName(nil)
		assert.False(t, result.Set)
	})

	t.Run("valid team name", func(t *testing.T) {
		teamName := "Test Team"
		result := convertTeamName(&teamName)
		assert.True(t, result.Set)
		assert.Equal(t, "Test Team", result.Value)
	})
}
