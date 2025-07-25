package rest

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/rom8726/warden/internal/domain"
	generatedapi "github.com/rom8726/warden/internal/generated/server"
	mockcontract "github.com/rom8726/warden/test_mocks/internal_/backend/contract"
)

func TestRestAPI_CreateTeam(t *testing.T) {
	t.Run("successful team creation", func(t *testing.T) {
		mockTeamsUseCase := mockcontract.NewMockTeamsUseCase(t)

		api := &RestAPI{
			teamsUseCase: mockTeamsUseCase,
		}

		req := &generatedapi.CreateTeamRequest{
			Name: "New Team",
		}

		expectedTeam := &domain.Team{
			ID:        domain.TeamID(123),
			Name:      "New Team",
			CreatedAt: time.Now(),
			Members: []domain.TeamMember{
				{
					TeamID: domain.TeamID(123),
					UserID: domain.UserID(1),
					Role:   domain.RoleOwner,
				},
			},
		}

		mockTeamsUseCase.EXPECT().
			Create(mock.Anything, domain.TeamDTO{
				Name: "New Team",
			}).
			Return(*expectedTeam, nil)

		resp, err := api.CreateTeam(context.Background(), req)

		require.NoError(t, err)
		require.NotNil(t, resp)

		teamResp, ok := resp.(*generatedapi.CreateTeamResponse)
		require.True(t, ok)
		assert.Equal(t, uint(123), teamResp.Team.ID)
		assert.Equal(t, "New Team", teamResp.Team.Name)
		assert.Len(t, teamResp.Team.Members, 1)
		assert.Equal(t, uint(1), teamResp.Team.Members[0].UserID)
		assert.Equal(t, generatedapi.TeamMemberRoleOwner, teamResp.Team.Members[0].Role)
	})

	t.Run("successful team creation without description", func(t *testing.T) {
		mockTeamsUseCase := mockcontract.NewMockTeamsUseCase(t)

		api := &RestAPI{
			teamsUseCase: mockTeamsUseCase,
		}

		req := &generatedapi.CreateTeamRequest{
			Name: "New Team",
		}

		expectedTeam := &domain.Team{
			ID:        domain.TeamID(123),
			Name:      "New Team",
			CreatedAt: time.Now(),
			Members:   []domain.TeamMember{},
		}

		mockTeamsUseCase.EXPECT().
			Create(mock.Anything, domain.TeamDTO{
				Name: "New Team",
			}).
			Return(*expectedTeam, nil)

		resp, err := api.CreateTeam(context.Background(), req)

		require.NoError(t, err)
		require.NotNil(t, resp)

		teamResp, ok := resp.(*generatedapi.CreateTeamResponse)
		require.True(t, ok)
		assert.Equal(t, uint(123), teamResp.Team.ID)
		assert.Equal(t, "New Team", teamResp.Team.Name)
		assert.Len(t, teamResp.Team.Members, 0)
	})

	t.Run("team name already exists", func(t *testing.T) {
		mockTeamsUseCase := mockcontract.NewMockTeamsUseCase(t)

		api := &RestAPI{
			teamsUseCase: mockTeamsUseCase,
		}

		req := &generatedapi.CreateTeamRequest{
			Name: "Existing Team",
		}

		mockTeamsUseCase.EXPECT().
			Create(mock.Anything, domain.TeamDTO{
				Name: "Existing Team",
			}).
			Return(domain.Team{}, domain.ErrTeamNameAlreadyInUse)

		resp, err := api.CreateTeam(context.Background(), req)

		require.NoError(t, err)
		require.NotNil(t, resp)

		errorResp, ok := resp.(*generatedapi.ErrorBadRequest)
		require.True(t, ok)
		assert.Equal(t, domain.ErrTeamNameAlreadyInUse.Error(), errorResp.Error.Message.Value)
	})

	t.Run("create team failed with unexpected error", func(t *testing.T) {
		mockTeamsUseCase := mockcontract.NewMockTeamsUseCase(t)

		api := &RestAPI{
			teamsUseCase: mockTeamsUseCase,
		}

		req := &generatedapi.CreateTeamRequest{
			Name: "New Team",
		}

		unexpectedErr := errors.New("database error")
		mockTeamsUseCase.EXPECT().
			Create(mock.Anything, domain.TeamDTO{
				Name: "New Team",
			}).
			Return(domain.Team{}, unexpectedErr)

		resp, err := api.CreateTeam(context.Background(), req)

		require.Error(t, err)
		require.NotNil(t, resp)
		assert.Equal(t, unexpectedErr, err)

		errorResp, ok := resp.(*generatedapi.ErrorInternalServerError)
		require.True(t, ok)
		assert.Equal(t, "Failed to create team", errorResp.Error.Message.Value)
	})
}
