package scheduler

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
	"github.com/rom8726/warden/internal/domain"
	"github.com/rom8726/warden/internal/infra"
	"github.com/rom8726/warden/internal/repository/events"
	"github.com/rom8726/warden/internal/repository/issues"
	"github.com/rom8726/warden/internal/repository/notificationsqueue"
	"github.com/rom8726/warden/internal/repository/projects"
	"github.com/rom8726/warden/internal/repository/releases"
	"github.com/rom8726/warden/internal/repository/releasestats"
	"github.com/rom8726/warden/internal/repository/teams"
	"github.com/rom8726/warden/internal/repository/usernotifications"
	"github.com/rom8726/warden/internal/repository/users"
	"github.com/rom8726/warden/internal/scheduler/config"
	"github.com/rom8726/warden/internal/scheduler/scheduler"
	"github.com/rom8726/warden/internal/scheduler/scheduler/jobs"
	"github.com/rom8726/warden/internal/scheduler/usecases/analytics"
	usernotificationsusecase "github.com/rom8726/warden/internal/scheduler/usecases/usernotifications"
	"github.com/rom8726/warden/internal/services/notification-channels/email"
	"github.com/rom8726/warden/pkg/db"
	"github.com/rom8726/warden/pkg/kafka"
)

const (
	ctxTimeout = 30 * time.Second
)

type App struct {
	Config *config.Config
	Logger *slog.Logger

	PostgresPool *pgxpool.Pool
	ClickHouse   infra.ClickHouseConn

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

	return app, nil
}

func (app *App) Run(ctx context.Context) error {
	techServer, err := techserver.NewTechServer(&app.Config.TechServer)
	if err != nil {
		return fmt.Errorf("create tech server: %w", err)
	}

	app.Logger.Info("Start scheduler...")

	group, groupCtx := errgroup.WithContext(ctx)
	group.Go(func() error { return techServer.ListenAndServe(groupCtx) })
	group.Go(func() error { return app.diApp.Run(groupCtx) })

	return group.Wait()
}

func (app *App) Close() {
	if app.PostgresPool != nil {
		app.PostgresPool.Close()
	}

	if app.ClickHouse != nil {
		_ = app.ClickHouse.Close()
	}
}

func (app *App) registerComponents() {
	// Register the transaction manager
	app.registerComponent(db.NewTxManager).Arg(app.PostgresPool)

	// Register ClickHouse connection
	app.registerComponent(func() *infra.ClickHouseConnImpl {
		return infra.NewClickHouseConn(app.ClickHouse)
	})

	eventsProducer := kafka.NewTopicProducer(kafka.NewNopKafkaProducer(), domain.EventsKafkaTopic)

	// Register repositories
	app.registerComponent(issues.New).Arg(app.PostgresPool)
	app.registerComponent(notificationsqueue.New).Arg(app.PostgresPool)
	app.registerComponent(usernotifications.New).Arg(app.PostgresPool)
	app.registerComponent(releases.New).Arg(app.PostgresPool)
	app.registerComponent(releasestats.New).Arg(app.PostgresPool)
	app.registerComponent(events.New).Arg(app.PostgresPool).Arg(eventsProducer)
	app.registerComponent(projects.New).Arg(app.PostgresPool)
	app.registerComponent(teams.New).Arg(app.PostgresPool)
	app.registerComponent(users.New).Arg(app.PostgresPool)

	// Register use cases
	app.registerComponent(analytics.New)
	app.registerComponent(usernotificationsusecase.New)

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

	// Scheduler
	app.registerComponent(scheduler.New)
	app.registerComponent(jobs.NewAlertsDaily)
	app.registerComponent(jobs.NewAlertsSummary)
	app.registerComponent(jobs.NewAnalyticsStats)
	app.registerComponent(jobs.NewNotificationsQueueCleaner)
	app.registerComponent(jobs.NewUserNotificationsCleaner)

	// Resolve background scheduler
	var schedulerSrv *scheduler.Scheduler
	if err := app.container.Resolve(&schedulerSrv); err != nil {
		panic(err)
	}

	var alertsDailyJob *jobs.AlertsDailyJob
	if err := app.container.Resolve(&alertsDailyJob); err != nil {
		panic(err)
	}
	if err := schedulerSrv.Register(alertsDailyJob, &scheduler.CronConfigDailyAlerts{}); err != nil {
		panic(err)
	}

	var alertsSummaryJob *jobs.AlertsSummary
	if err := app.container.Resolve(&alertsSummaryJob); err != nil {
		panic(err)
	}
	if err := schedulerSrv.Register(alertsSummaryJob, &scheduler.CronConfigSummaryAlerts{}); err != nil {
		panic(err)
	}

	var analyticsStatsJob *jobs.AnalyticsStatsJob
	if err := app.container.Resolve(&analyticsStatsJob); err != nil {
		panic(err)
	}
	if err := schedulerSrv.Register(analyticsStatsJob, &scheduler.CronAnalyticsStats{}); err != nil {
		panic(err)
	}

	var notifQueueCleanerJob *jobs.NotificationsQueueCleanerJob
	if err := app.container.Resolve(&notifQueueCleanerJob); err != nil {
		panic(err)
	}
	if err := schedulerSrv.Register(notifQueueCleanerJob, &scheduler.CronNotificationsCleaner{}); err != nil {
		panic(err)
	}

	var userNotifCleanerJob *jobs.UserNotificationsCleanerJob
	if err := app.container.Resolve(&userNotifCleanerJob); err != nil {
		panic(err)
	}
	if err := schedulerSrv.Register(userNotifCleanerJob, &scheduler.CronUserNotificationsCleaner{}); err != nil {
		panic(err)
	}
}

func (app *App) registerComponent(constructor any) *di.Provider {
	return app.container.Provide(constructor)
}
