package ingestserver

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"time"

	"github.com/IBM/sarama"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"github.com/rom8726/di"
	"golang.org/x/sync/errgroup"

	commonconfig "github.com/rom8726/warden/internal/common/config"
	"github.com/rom8726/warden/internal/common/techserver"
	generatedserver "github.com/rom8726/warden/internal/generated/ingestserver"
	"github.com/rom8726/warden/internal/ingest-server/api/rest"
	"github.com/rom8726/warden/internal/ingest-server/api/rest/middlewares"
	"github.com/rom8726/warden/internal/ingest-server/config"
	"github.com/rom8726/warden/internal/ingest-server/contract"
	"github.com/rom8726/warden/internal/ingest-server/services/envelopequeueproducer"
	envelopeusecase "github.com/rom8726/warden/internal/ingest-server/usecases/envelope"
	projectsusecase "github.com/rom8726/warden/internal/ingest-server/usecases/projects"
	storeeventusecase "github.com/rom8726/warden/internal/ingest-server/usecases/storeevent"
	"github.com/rom8726/warden/internal/repository/projects"
	"github.com/rom8726/warden/internal/services/storeeventqueueproducer"
	"github.com/rom8726/warden/pkg/db"
	"github.com/rom8726/warden/pkg/httpserver"
	pkgmiddlewares "github.com/rom8726/warden/pkg/httpserver/middlewares"
	"github.com/rom8726/warden/pkg/kafka"
)

const (
	ctxTimeout = 30 * time.Second
)

type App struct {
	Config *config.Config
	Logger *slog.Logger

	PostgresPool       *pgxpool.Pool
	KafkaClient        sarama.Client
	KafkaAsyncProducer *kafka.Producer
	RedisClient        *redis.Client

	APIServer *httpserver.Server

	container *di.Container
	diApp     *di.App
}

func NewApp(ctx context.Context, cfg *config.Config, logger *slog.Logger) (*App, error) {
	ctxWithTimeout, cancel := context.WithTimeout(context.Background(), ctxTimeout)
	defer cancel()

	pgPool, err := commonconfig.NewPostgresConnPool(ctxWithTimeout, &cfg.Postgres)
	if err != nil {
		return nil, fmt.Errorf("create postgres pool: %w", err)
	}

	kafkaClient, kafkaAsyncProducer, err := commonconfig.NewKafka(ctxWithTimeout, &cfg.Kafka)
	if err != nil {
		return nil, fmt.Errorf("create kafka producer: %w", err)
	}

	// Create Kafka topics for events and exceptions
	if err := commonconfig.CreateKafkaTopics(cfg.Kafka.Brokers); err != nil {
		return nil, fmt.Errorf("create kafka topics: %w", err)
	}

	// Create a Redis client
	redisClient, err := commonconfig.NewRedisClient(ctxWithTimeout, &cfg.Redis)
	if err != nil {
		return nil, fmt.Errorf("create redis client: %w", err)
	}

	container := di.New()
	diApp := di.NewApp(container)

	app := &App{
		Config:             cfg,
		Logger:             logger,
		container:          container,
		diApp:              diApp,
		PostgresPool:       pgPool,
		KafkaClient:        kafkaClient,
		KafkaAsyncProducer: kafkaAsyncProducer,
		RedisClient:        redisClient,
	}

	app.registerComponents()

	app.APIServer, err = app.newAPIServer(ctx)
	if err != nil {
		return nil, fmt.Errorf("create API server: %w", err)
	}

	return app, nil
}

func (app *App) Run(ctx context.Context) error {
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

func (app *App) Close() {
	if app.PostgresPool != nil {
		app.PostgresPool.Close()
	}

	if app.KafkaClient != nil {
		_ = app.KafkaClient.Close()
	}

	if app.KafkaAsyncProducer != nil {
		_ = app.KafkaAsyncProducer.Close()
	}

	if app.RedisClient != nil {
		_ = app.RedisClient.Close()
	}
}

func (app *App) registerComponents() {
	// Register the transaction manager
	app.registerComponent(db.NewTxManager).Arg(app.PostgresPool)

	app.registerComponent(kafka.NewTopicProducerCreator).Arg(app.KafkaAsyncProducer)

	// Register repositories
	app.registerComponent(projects.New).Arg(app.PostgresPool)

	// Register use cases
	app.registerComponent(projectsusecase.New)
	app.registerComponent(envelopeusecase.New)
	app.registerComponent(storeeventusecase.New)

	// Register services
	app.registerComponent(envelopequeueproducer.New)

	var topicProducerCreator *kafka.TopicProducerCreator
	if err := app.container.Resolve(&topicProducerCreator); err != nil {
		panic(err)
	}

	app.registerComponent(storeeventqueueproducer.New).Arg(topicProducerCreator)

	// Register API components
	app.registerComponent(rest.NewSecurityHandler)
	app.registerComponent(rest.New)
}

func (app *App) registerComponent(constructor any) *di.Provider {
	return app.container.Provide(constructor)
}

func (app *App) newAPIServer(ctx context.Context) (*httpserver.Server, error) {
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

	// Convert config to middleware config
	adaptiveThrottleConfig := &middlewares.AdaptiveThrottleConfig{
		RPSWindow:            app.Config.RateLimit.RPSWindow,
		StatsRefreshInterval: app.Config.RateLimit.StatsRefreshInterval,
		RateLimit:            app.Config.RateLimit.RateLimit,
	}

	// Get the project repository from the DI container
	var projectsRepo contract.ProjectsRepository
	if err := app.container.Resolve(&projectsRepo); err != nil {
		return nil, fmt.Errorf("resolve projects repository: %w", err)
	}

	// Create the adaptive throttle middleware
	adaptiveThrottleResult := middlewares.AdaptiveThrottle(ctx, app.RedisClient, projectsRepo, adaptiveThrottleConfig)

	// Middleware chain:
	// CORS → ProjectID → AdaptiveThrottle → API implementation
	handler := pkgmiddlewares.CORSMdw(
		middlewares.WithProjectID(
			adaptiveThrottleResult.Handler(
				genServer,
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
