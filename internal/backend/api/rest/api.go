//nolint:interfacebloat // it's ok here
package rest

import (
	"context"
	"errors"
	"net/http"

	"github.com/ogen-go/ogen/ogenerrors"

	"github.com/rom8726/warden/internal/backend/config"
	"github.com/rom8726/warden/internal/backend/contract"
	generatedapi "github.com/rom8726/warden/internal/generated/server"
)

var _ generatedapi.Handler = (*RestAPI)(nil)

type RestAPI struct {
	config                   *config.Config
	tokenizer                contract.Tokenizer
	usersUseCase             contract.UsersUseCase
	projectsUseCase          contract.ProjectsUseCase
	eventUseCase             contract.EventUseCase
	issueUseCase             contract.IssueUseCase
	teamsUseCase             contract.TeamsUseCase
	notificationsUseCase     contract.NotificationsUseCase
	analyticsUseCase         contract.AnalyticsUseCase
	permissionsService       contract.PermissionsService
	settingsUseCase          contract.SettingsUseCase
	userNotificationsUseCase contract.UserNotificationsUseCase
	versionsUseCase          contract.VersionsUseCase
}

func New(
	config *config.Config,
	usersService contract.UsersUseCase,
	tokenizer contract.Tokenizer,
	projectsUseCase contract.ProjectsUseCase,
	eventUseCase contract.EventUseCase,
	issueUseCase contract.IssueUseCase,
	teamsUseCase contract.TeamsUseCase,
	analyticsUseCase contract.AnalyticsUseCase,
	notificationsUseCase contract.NotificationsUseCase,
	permissionsService contract.PermissionsService,
	settingsUseCase contract.SettingsUseCase,
	userNotificationsUseCase contract.UserNotificationsUseCase,
	versionsUseCase contract.VersionsUseCase,
) *RestAPI {
	return &RestAPI{
		config:                   config,
		usersUseCase:             usersService,
		tokenizer:                tokenizer,
		projectsUseCase:          projectsUseCase,
		eventUseCase:             eventUseCase,
		issueUseCase:             issueUseCase,
		teamsUseCase:             teamsUseCase,
		notificationsUseCase:     notificationsUseCase,
		analyticsUseCase:         analyticsUseCase,
		permissionsService:       permissionsService,
		settingsUseCase:          settingsUseCase,
		userNotificationsUseCase: userNotificationsUseCase,
		versionsUseCase:          versionsUseCase,
	}
}

func (r *RestAPI) NewError(_ context.Context, err error) *generatedapi.ErrorStatusCode {
	code := http.StatusInternalServerError
	errMessage := err.Error()

	var secError *ogenerrors.SecurityError
	if errors.As(err, &secError) {
		code = http.StatusUnauthorized
		errMessage = "unauthorized"
	}

	return &generatedapi.ErrorStatusCode{
		StatusCode: code,
		Response: generatedapi.Error{
			Error: generatedapi.ErrorError{
				Message: generatedapi.NewOptString(errMessage),
			},
		},
	}
}
