package issuenotificator

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rom8726/di"
	"golang.org/x/sync/errgroup"

	commonconfig "github.com/rom8726/warden/internal/common/config"
	"github.com/rom8726/warden/internal/common/techserver"
	"github.com/rom8726/warden/internal/issue-notificator/config"
	"github.com/rom8726/warden/internal/issue-notificator/notificator"
	notificationsusecase "github.com/rom8726/warden/internal/issue-notificator/usecases/notifications"
	"github.com/rom8726/warden/internal/repository/issues"
	"github.com/rom8726/warden/internal/repository/notifications"
	"github.com/rom8726/warden/internal/repository/notificationsqueue"
	"github.com/rom8726/warden/internal/repository/projects"
	"github.com/rom8726/warden/internal/repository/teams"
	"github.com/rom8726/warden/internal/repository/users"
	"github.com/rom8726/warden/internal/services/notification-channels/email"
	"github.com/rom8726/warden/internal/services/notification-channels/mattermost"
	"github.com/rom8726/warden/internal/services/notification-channels/pachca"
	"github.com/rom8726/warden/internal/services/notification-channels/slack"
	"github.com/rom8726/warden/internal/services/notification-channels/telegram"
	"github.com/rom8726/warden/internal/services/notification-channels/webhook"
	"github.com/rom8726/warden/pkg/db"
)

const (
	ctxTimeout = 30 * time.Second
)

type App struct {
	Config *config.Config
	Logger *slog.Logger

	PostgresPool *pgxpool.Pool

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

	container := di.New()
	diApp := di.NewApp(container)

	app := &App{
		Config:       cfg,
		Logger:       logger,
		container:    container,
		diApp:        diApp,
		PostgresPool: pgPool,
	}

	app.registerComponents()

	return app, nil
}

func (app *App) Run(ctx context.Context) error {
	techServer, err := techserver.NewTechServer(&app.Config.TechServer)
	if err != nil {
		return fmt.Errorf("create tech server: %w", err)
	}

	app.Logger.Info("Start worker...")

	group, groupCtx := errgroup.WithContext(ctx)
	group.Go(func() error { return techServer.ListenAndServe(groupCtx) })
	group.Go(func() error { return app.diApp.Run(groupCtx) })

	return group.Wait()
}

func (app *App) Close() {
	if app.PostgresPool != nil {
		app.PostgresPool.Close()
	}
}

func (app *App) registerComponents() {
	// Register the transaction manager
	app.registerComponent(db.NewTxManager).Arg(app.PostgresPool)

	// Register repositories
	app.registerComponent(projects.New).Arg(app.PostgresPool)
	app.registerComponent(teams.New).Arg(app.PostgresPool)
	app.registerComponent(users.New).Arg(app.PostgresPool)
	app.registerComponent(issues.New).Arg(app.PostgresPool)
	app.registerComponent(notifications.New).Arg(app.PostgresPool)
	app.registerComponent(notificationsqueue.New).Arg(app.PostgresPool)

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

	// Register use cases
	app.registerComponent(notificationsusecase.New)

	// Resolve channels
	var emailService *email.Service
	if err := app.container.Resolve(&emailService); err != nil {
		panic(err)
	}

	var mattermostService *mattermost.Service
	if err := app.container.Resolve(&mattermostService); err != nil {
		panic(err)
	}

	var webhookService *webhook.Service
	if err := app.container.Resolve(&webhookService); err != nil {
		panic(err)
	}

	var telegramService *telegram.Service
	if err := app.container.Resolve(&telegramService); err != nil {
		panic(err)
	}

	var slackService *slack.Service
	if err := app.container.Resolve(&slackService); err != nil {
		panic(err)
	}

	var pachcaService *pachca.Service
	if err := app.container.Resolve(&pachcaService); err != nil {
		panic(err)
	}

	app.registerComponent(notificator.New).Arg([]notificator.Channel{
		emailService,
		mattermostService,
		webhookService,
		telegramService,
		slackService,
		pachcaService,
	}).Arg(app.Config.Notificator.WorkerCount)
	var notificatorSrv *notificator.Service
	if err := app.container.Resolve(&notificatorSrv); err != nil {
		panic(err)
	}
}

func (app *App) registerComponent(constructor any) *di.Provider {
	return app.container.Provide(constructor)
}
