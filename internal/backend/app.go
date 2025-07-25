package backend

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rom8726/di"
	"golang.org/x/sync/errgroup"

	"github.com/rom8726/warden/internal/backend/api/rest"
	"github.com/rom8726/warden/internal/backend/api/rest/middlewares"
	"github.com/rom8726/warden/internal/backend/config"
	"github.com/rom8726/warden/internal/backend/contract"
	ratelimiter2fa "github.com/rom8726/warden/internal/backend/services/2fa/ratelimiter"
	"github.com/rom8726/warden/internal/backend/services/permissions"
	"github.com/rom8726/warden/internal/backend/services/tokenizer"
	"github.com/rom8726/warden/internal/backend/usecases/analytics"
	eventsusecases "github.com/rom8726/warden/internal/backend/usecases/events"
	issuesusecases "github.com/rom8726/warden/internal/backend/usecases/issues"
	notificationsusecases "github.com/rom8726/warden/internal/backend/usecases/notifications"
	projectsusecase "github.com/rom8726/warden/internal/backend/usecases/projects"
	settingsusecase "github.com/rom8726/warden/internal/backend/usecases/settings"
	teamsusecases "github.com/rom8726/warden/internal/backend/usecases/teams"
	usernotificationsusecase "github.com/rom8726/warden/internal/backend/usecases/usernotifications"
	usersusecase "github.com/rom8726/warden/internal/backend/usecases/users"
	versionsusecase "github.com/rom8726/warden/internal/backend/usecases/versions"
	commonconfig "github.com/rom8726/warden/internal/common/config"
	"github.com/rom8726/warden/internal/common/techserver"
	"github.com/rom8726/warden/internal/domain"
	generatedserver "github.com/rom8726/warden/internal/generated/server"
	"github.com/rom8726/warden/internal/infra"
	"github.com/rom8726/warden/internal/repository/events"
	"github.com/rom8726/warden/internal/repository/issuereleases"
	"github.com/rom8726/warden/internal/repository/issues"
	"github.com/rom8726/warden/internal/repository/notifications"
	"github.com/rom8726/warden/internal/repository/notificationsqueue"
	"github.com/rom8726/warden/internal/repository/projects"
	"github.com/rom8726/warden/internal/repository/releases"
	"github.com/rom8726/warden/internal/repository/releasestats"
	"github.com/rom8726/warden/internal/repository/resolutions"
	"github.com/rom8726/warden/internal/repository/settings"
	"github.com/rom8726/warden/internal/repository/teams"
	"github.com/rom8726/warden/internal/repository/usernotifications"
	"github.com/rom8726/warden/internal/repository/users"
	"github.com/rom8726/warden/internal/services/notification-channels/email"
	"github.com/rom8726/warden/internal/services/notification-channels/mattermost"
	"github.com/rom8726/warden/internal/services/notification-channels/pachca"
	"github.com/rom8726/warden/internal/services/notification-channels/slack"
	"github.com/rom8726/warden/internal/services/notification-channels/telegram"
	"github.com/rom8726/warden/internal/services/notification-channels/webhook"
	"github.com/rom8726/warden/pkg/db"
	"github.com/rom8726/warden/pkg/httpserver"
	pkgmiddlewares "github.com/rom8726/warden/pkg/httpserver/middlewares"
	"github.com/rom8726/warden/pkg/kafka"
	"github.com/rom8726/warden/pkg/passworder"
)

const (
	ctxTimeout = 30 * time.Second
)

type App struct {
	Config *config.Config
	Logger *slog.Logger

	PostgresPool *pgxpool.Pool
	ClickHouse   infra.ClickHouseConn

	APIServer *httpserver.Server

	container *di.Container
	diApp     *di.App
}

func NewApp(ctx context.Context, cfg *config.Config, logger *slog.Logger) (*App, error) {
	ctx, cancel := context.WithTimeout(ctx, ctxTimeout)
	defer cancel()

	pgPool, err := commonconfig.NewPostgresConnPool(ctx, &cfg.Postgres)
	if err != nil {
		return nil, fmt.Errorf("create postgres pool: %w", err)
	}

	clickHouseClient, err := commonconfig.NewClickHouseClient(ctx, &cfg.ClickHouse)
	if err != nil {
		return nil, fmt.Errorf("create clickhouse client: %w", err)
	}

	container := di.New()
	diApp := di.NewApp(container)

	app := &App{
		Config:       cfg,
		Logger:       logger,
		container:    container,
		diApp:        diApp,
		PostgresPool: pgPool,
		ClickHouse:   clickHouseClient,
	}

	app.registerComponents()

	app.APIServer, err = app.newAPIServer()
	if err != nil {
		return nil, fmt.Errorf("create API server: %w", err)
	}

	return app, nil
}

func (app *App) RegisterComponent(constructor any) *di.Provider {
	return app.container.Provide(constructor)
}

func (app *App) ResolveComponent(target any) error {
	return app.container.Resolve(target)
}

func (app *App) ResolveComponentsToStruct(target any) error {
	return app.container.ResolveToStruct(target)
}

func (app *App) Run(ctx context.Context) error {
	// Check and create superuser if needed
	if app.Config.AdminEmail != "" {
		if err := app.ensureSuperuser(ctx); err != nil {
			app.Logger.Error("Failed to ensure superuser exists", "error", err)
		}
	}

	techServer, err := techserver.NewTechServer(&app.Config.TechServer)
	if err != nil {
		return fmt.Errorf("create tech server: %w", err)
	}

	app.Logger.Info("Start API server")

	group, groupCtx := errgroup.WithContext(ctx)
	group.Go(func() error { return app.APIServer.ListenAndServe(groupCtx) })
	group.Go(func() error { return techServer.ListenAndServe(groupCtx) })
	group.Go(func() error { return app.diApp.Run(groupCtx) })

	return group.Wait()
}

func (app *App) TruncateCHTables(ctx context.Context) error {
	return app.ClickHouse.Exec(ctx, "TRUNCATE TABLE events")
}

func (app *App) Close() {
	if app.PostgresPool != nil {
		app.PostgresPool.Close()
	}

	if app.ClickHouse != nil {
		_ = app.ClickHouse.Close()
	}
}

func (app *App) registerComponent(constructor any) *di.Provider {
	return app.container.Provide(constructor)
}

//nolint:gocyclo // it's ok
func (app *App) registerComponents() {
	// Register the transaction manager
	app.registerComponent(db.NewTxManager).Arg(app.PostgresPool)

	// Kafka
	eventsProducer := kafka.NewTopicProducer(kafka.NewNopKafkaProducer(), domain.EventsKafkaTopic)

	// Register ClickHouse connection
	app.registerComponent(func() *infra.ClickHouseConnImpl {
		return infra.NewClickHouseConn(app.ClickHouse)
	})

	// Register repositories
	app.registerComponent(projects.New).Arg(app.PostgresPool)
	app.registerComponent(events.New).Arg(eventsProducer)
	app.registerComponent(users.New).Arg(app.PostgresPool)
	app.registerComponent(issues.New).Arg(app.PostgresPool)
	app.registerComponent(teams.New).Arg(app.PostgresPool)
	app.registerComponent(resolutions.New).Arg(app.PostgresPool)
	app.registerComponent(notifications.New).Arg(app.PostgresPool)
	app.registerComponent(notificationsqueue.New).Arg(app.PostgresPool)
	app.registerComponent(releases.New).Arg(app.PostgresPool)
	app.registerComponent(releasestats.New).Arg(app.PostgresPool)
	app.registerComponent(issuereleases.New).Arg(app.PostgresPool)
	app.registerComponent(settings.New).Arg(app.PostgresPool)
	app.registerComponent(usernotifications.New).Arg(app.PostgresPool)

	// Register permissions service
	app.registerComponent(permissions.New)

	// Register channels
	app.registerComponent(mattermost.New).Arg(&mattermost.ServiceParams{
		BaseURL: app.Config.FrontendURL,
	})
	app.registerComponent(webhook.New).Arg(app.Config.FrontendURL)
	app.registerComponent(telegram.New).Arg(&telegram.ServiceParams{
		BaseURL: app.Config.FrontendURL,
	})
	app.registerComponent(slack.New).Arg(&slack.ServiceParams{
		BaseURL: app.Config.FrontendURL,
	})
	app.registerComponent(pachca.New).Arg(&pachca.ServiceParams{
		BaseURL: app.Config.FrontendURL,
	})
	app.registerComponent(email.New).Arg(&email.Config{
		SMTPHost:      app.Config.Mailer.Addr,
		Username:      app.Config.Mailer.User,
		Password:      app.Config.Mailer.Password,
		CertFile:      app.Config.Mailer.CertFile,
		KeyFile:       app.Config.Mailer.KeyFile,
		AllowInsecure: app.Config.Mailer.AllowInsecure,
		UseTLS:        app.Config.Mailer.UseTLS,
		BaseURL:       app.Config.FrontendURL,
		From:          app.Config.Mailer.From,
	})

	// Resolve channels
	var emailChannel *email.Service
	if err := app.container.Resolve(&emailChannel); err != nil {
		panic(err)
	}

	var mattermostChannel *mattermost.Service
	if err := app.container.Resolve(&mattermostChannel); err != nil {
		panic(err)
	}

	var webhookChannel *webhook.Service
	if err := app.container.Resolve(&webhookChannel); err != nil {
		panic(err)
	}

	var telegramChannel *telegram.Service
	if err := app.container.Resolve(&telegramChannel); err != nil {
		panic(err)
	}

	var slackChannel *slack.Service
	if err := app.container.Resolve(&slackChannel); err != nil {
		panic(err)
	}

	var pachcaChannel *pachca.Service
	if err := app.container.Resolve(&pachcaChannel); err != nil {
		panic(err)
	}

	// Register use cases
	app.registerComponent(eventsusecases.New)
	app.registerComponent(issuesusecases.New)
	app.registerComponent(teamsusecases.New)
	app.registerComponent(projectsusecase.New)
	app.registerComponent(notificationsusecases.New).Arg([]contract.NotificationChannel{
		emailChannel,
		mattermostChannel,
		webhookChannel,
		telegramChannel,
		slackChannel,
		pachcaChannel,
	})
	app.registerComponent(analytics.New)
	app.registerComponent(settingsusecase.New).Arg(app.Config.SecretKey)
	app.registerComponent(usernotificationsusecase.New)

	// Register versions service
	app.registerComponent(versionsusecase.New)

	app.registerComponent(usersusecase.New).Arg([]usersusecase.AuthProvider{})

	// Register services
	app.registerComponent(tokenizer.New).Arg(&tokenizer.ServiceParams{
		SecretKey:        []byte(app.Config.JWTSecretKey),
		AccessTTL:        app.Config.AccessTokenTTL,
		RefreshTTL:       app.Config.RefreshTokenTTL,
		ResetPasswordTTL: app.Config.ResetPasswordTTL,
	})
	app.registerComponent(ratelimiter2fa.New)

	// Register API components
	app.registerComponent(rest.NewSecurityHandler)
	app.registerComponent(rest.New).Arg(app.Config)
}

func (app *App) newAPIServer() (*httpserver.Server, error) {
	cfg := app.Config.APIServer

	var restAPI generatedserver.Handler
	if err := app.container.Resolve(&restAPI); err != nil {
		return nil, fmt.Errorf("resolve REST API service component: %w", err)
	}

	var securityHandler generatedserver.SecurityHandler
	if err := app.container.Resolve(&securityHandler); err != nil {
		return nil, fmt.Errorf("resolve API security handler component: %w", err)
	}

	genServer, err := generatedserver.NewServer(restAPI, securityHandler)
	if err != nil {
		return nil, fmt.Errorf("create API server: %w", err)
	}

	var tokenizerSrv contract.Tokenizer
	if err := app.container.Resolve(&tokenizerSrv); err != nil {
		return nil, fmt.Errorf("resolve tokenizer service component: %w", err)
	}

	var usersSrv contract.UsersUseCase
	if err := app.container.Resolve(&usersSrv); err != nil {
		return nil, fmt.Errorf("resolve users service component: %w", err)
	}

	var permService contract.PermissionsService
	if err := app.container.Resolve(&permService); err != nil {
		return nil, fmt.Errorf("resolve permissions service component: %w", err)
	}

	// Middleware chain:
	// CORS → RAW -> Auth → ProjectAccess → ProjectManagement → IssueAccess → IssueManagement → API implementation
	handler := pkgmiddlewares.CORSMdw(
		middlewares.WithRawRequest(
			middlewares.AuthMiddleware(tokenizerSrv, usersSrv)(
				middlewares.ProjectAccess(permService)(
					middlewares.ProjectManagement(permService)(
						middlewares.IssueAccess(permService)(
							middlewares.IssueManagement(permService)(
								genServer,
							),
						),
					),
				),
			),
		),
	)

	lis, err := net.Listen("tcp", cfg.Addr)
	if err != nil {
		return nil, fmt.Errorf("listen %q: %w", cfg.Addr, err)
	}

	return &httpserver.Server{
		Listener:     lis,
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
		IdleTimeout:  cfg.IdleTimeout,
		Handler:      handler,
	}, nil
}

// ensureSuperuser checks if a user with the admin email exists and creates one if not.
func (app *App) ensureSuperuser(ctx context.Context) error {
	app.Logger.Info("Checking if superuser exists")

	var usersRepo contract.UsersRepository
	if err := app.container.Resolve(&usersRepo); err != nil {
		return fmt.Errorf("resolve users repository: %w", err)
	}

	// Check if user with admin email already exists
	_, err := usersRepo.GetByEmail(ctx, app.Config.AdminEmail)
	if err == nil {
		// User already exists
		app.Logger.Info("Superuser already exists")

		return nil
	}

	// Extract username from email (part before @)
	username := app.Config.AdminEmail
	for i, c := range username {
		if c == '@' {
			username = username[:i]

			break
		}
	}

	// Hash the temporary password
	passwordHash, err := passworder.PasswordHash(app.Config.AdminTmpPassword)
	if err != nil {
		return fmt.Errorf("hash password: %w", err)
	}

	// Create the superuser
	userDTO := domain.UserDTO{
		Username:      username,
		Email:         app.Config.AdminEmail,
		PasswordHash:  passwordHash,
		IsSuperuser:   true,
		IsTmpPassword: true,
	}

	user, err := usersRepo.Create(ctx, userDTO)
	if err != nil {
		return fmt.Errorf("create superuser: %w", err)
	}

	app.Logger.Info("Created superuser", "id", user.ID, "username", user.Username)

	return nil
}
