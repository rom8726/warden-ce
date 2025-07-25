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

func TestRestAPI_ListTeams(t *testing.T) {
	t.Run("successful teams list", func(t *testing.T) {
		mockTeamsUseCase := mockcontract.NewMockTeamsUseCase(t)

		api := &RestAPI{
			teamsUseCase: mockTeamsUseCase,
		}

		expectedTeams := []domain.Team{
			{
				ID:   domain.TeamID(1),
				Name: "Team 1",
				Members: []domain.TeamMember{
					{
						UserID: domain.UserID(1),
						Role:   domain.RoleOwner,
					},
					{
						UserID: domain.UserID(2),
						Role:   domain.RoleMember,
					},
				},
				CreatedAt: time.Now(),
			},
			{
				ID:   domain.TeamID(2),
				Name: "Team 2",
				Members: []domain.TeamMember{
					{
						UserID: domain.UserID(3),
						Role:   domain.RoleAdmin,
					},
				},
				CreatedAt: time.Now(),
			},
		}

		mockTeamsUseCase.EXPECT().
			List(mock.Anything).
			Return(expectedTeams, nil)

		resp, err := api.ListTeams(context.Background())

		require.NoError(t, err)
		require.NotNil(t, resp)

		listResp, ok := resp.(*generatedapi.ListTeamsResponse)
		require.True(t, ok)
		assert.Len(t, *listResp, 2)

		// Check first team
		assert.Equal(t, uint(1), (*listResp)[0].ID)
		assert.Equal(t, "Team 1", (*listResp)[0].Name)
		assert.Len(t, (*listResp)[0].Members, 2)
		assert.Equal(t, uint(1), (*listResp)[0].Members[0].UserID)
		assert.Equal(t, generatedapi.TeamMemberRoleOwner, (*listResp)[0].Members[0].Role)
		assert.Equal(t, uint(2), (*listResp)[0].Members[1].UserID)
		assert.Equal(t, generatedapi.TeamMemberRoleMember, (*listResp)[0].Members[1].Role)

		// Check second team
		assert.Equal(t, uint(2), (*listResp)[1].ID)
		assert.Equal(t, "Team 2", (*listResp)[1].Name)
		assert.Len(t, (*listResp)[1].Members, 1)
		assert.Equal(t, uint(3), (*listResp)[1].Members[0].UserID)
		assert.Equal(t, generatedapi.TeamMemberRoleAdmin, (*listResp)[1].Members[0].Role)
	})

	t.Run("empty teams list", func(t *testing.T) {
		mockTeamsUseCase := mockcontract.NewMockTeamsUseCase(t)

		api := &RestAPI{
			teamsUseCase: mockTeamsUseCase,
		}

		expectedTeams := []domain.Team{}

		mockTeamsUseCase.EXPECT().
			List(mock.Anything).
			Return(expectedTeams, nil)

		resp, err := api.ListTeams(context.Background())

		require.NoError(t, err)
		require.NotNil(t, resp)

		listResp, ok := resp.(*generatedapi.ListTeamsResponse)
		require.True(t, ok)
		assert.Len(t, *listResp, 0)
	})

	t.Run("team without members", func(t *testing.T) {
		mockTeamsUseCase := mockcontract.NewMockTeamsUseCase(t)

		api := &RestAPI{
			teamsUseCase: mockTeamsUseCase,
		}

		expectedTeams := []domain.Team{
			{
				ID:        domain.TeamID(1),
				Name:      "Empty Team",
				Members:   []domain.TeamMember{},
				CreatedAt: time.Now(),
			},
		}

		mockTeamsUseCase.EXPECT().
			List(mock.Anything).
			Return(expectedTeams, nil)

		resp, err := api.ListTeams(context.Background())

		require.NoError(t, err)
		require.NotNil(t, resp)

		listResp, ok := resp.(*generatedapi.ListTeamsResponse)
		require.True(t, ok)
		assert.Len(t, *listResp, 1)
		assert.Equal(t, uint(1), (*listResp)[0].ID)
		assert.Equal(t, "Empty Team", (*listResp)[0].Name)
		assert.Len(t, (*listResp)[0].Members, 0)
	})

	t.Run("list teams failed", func(t *testing.T) {
		mockTeamsUseCase := mockcontract.NewMockTeamsUseCase(t)

		api := &RestAPI{
			teamsUseCase: mockTeamsUseCase,
		}

		unexpectedErr := errors.New("database error")
		mockTeamsUseCase.EXPECT().
			List(mock.Anything).
			Return(nil, unexpectedErr)

		resp, err := api.ListTeams(context.Background())

		require.Error(t, err)
		require.Nil(t, resp)
		assert.Equal(t, unexpectedErr, err)
	})
}

func TestMapDomainRoleToAPI(t *testing.T) {
	t.Run("owner role", func(t *testing.T) {
		result := mapDomainRoleToAPI(domain.RoleOwner)
		assert.Equal(t, generatedapi.TeamMemberRoleOwner, result)
	})

	t.Run("admin role", func(t *testing.T) {
		result := mapDomainRoleToAPI(domain.RoleAdmin)
		assert.Equal(t, generatedapi.TeamMemberRoleAdmin, result)
	})

	t.Run("member role", func(t *testing.T) {
		result := mapDomainRoleToAPI(domain.RoleMember)
		assert.Equal(t, generatedapi.TeamMemberRoleMember, result)
	})

	t.Run("unknown role defaults to member", func(t *testing.T) {
		result := mapDomainRoleToAPI("unknown")
		assert.Equal(t, generatedapi.TeamMemberRoleMember, result)
	})
}
