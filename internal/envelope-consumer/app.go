package envelopeconsumer

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/IBM/sarama"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"github.com/rom8726/di"
	"golang.org/x/sync/errgroup"

	commonconfig "github.com/rom8726/warden/internal/common/config"
	"github.com/rom8726/warden/internal/common/techserver"
	"github.com/rom8726/warden/internal/domain"
	"github.com/rom8726/warden/internal/envelope-consumer/config"
	"github.com/rom8726/warden/internal/envelope-consumer/contract"
	cacheservice "github.com/rom8726/warden/internal/envelope-consumer/services/cache"
	"github.com/rom8726/warden/internal/envelope-consumer/services/cachemanager"
	"github.com/rom8726/warden/internal/envelope-consumer/services/envelopequeueprocessor"
	"github.com/rom8726/warden/internal/envelope-consumer/services/storeeventqueueprocessor"
	envelopeusecase "github.com/rom8726/warden/internal/envelope-consumer/usecases/envelope"
	eventsusecase "github.com/rom8726/warden/internal/envelope-consumer/usecases/events"
	storeeventusecase "github.com/rom8726/warden/internal/envelope-consumer/usecases/storeevent"
	"github.com/rom8726/warden/internal/infra"
	"github.com/rom8726/warden/internal/repository/events"
	"github.com/rom8726/warden/internal/repository/issuereleases"
	"github.com/rom8726/warden/internal/repository/issues"
	"github.com/rom8726/warden/internal/repository/notificationsqueue"
	"github.com/rom8726/warden/internal/repository/releases"
	"github.com/rom8726/warden/internal/services/storeeventqueueproducer"
	"github.com/rom8726/warden/pkg/db"
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
	ClickHouse         infra.ClickHouseConn
	RedisClient        *redis.Client

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

	kafkaClient, kafkaAsyncProducer, err := commonconfig.NewKafka(ctx, &cfg.Kafka)
	if err != nil {
		return nil, fmt.Errorf("create kafka producer: %w", err)
	}

	// Create Kafka topics for events and exceptions
	if err := commonconfig.CreateKafkaTopics(cfg.Kafka.Brokers); err != nil {
		return nil, fmt.Errorf("create kafka topics: %w", err)
	}

	clickHouseClient, err := commonconfig.NewClickHouseClient(ctx, &cfg.ClickHouse)
	if err != nil {
		return nil, fmt.Errorf("create clickhouse client: %w", err)
	}

	redisClient, err := commonconfig.NewRedisClient(ctx, &cfg.Redis)
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
		ClickHouse:         clickHouseClient,
		RedisClient:        redisClient,
	}

	if err := app.registerComponents(); err != nil {
		return nil, fmt.Errorf("register components: %w", err)
	}

	return app, nil
}

func (app *App) Run(ctx context.Context) error {
	techServer, err := techserver.NewTechServer(&app.Config.TechServer)
	if err != nil {
		return fmt.Errorf("create tech server: %w", err)
	}

	app.Logger.Info("Start consumer")

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

	if app.KafkaClient != nil {
		_ = app.KafkaClient.Close()
	}

	if app.KafkaAsyncProducer != nil {
		_ = app.KafkaAsyncProducer.Close()
	}

	// Close cache service
	var cacheService contract.CacheService
	if err := app.container.Resolve(&cacheService); err == nil && cacheService != nil {
		if err := cacheService.Close(context.Background()); err != nil {
			app.Logger.Error("Failed to close cache service", "error", err)
		}
	}
}

func (app *App) registerComponents() error {
	// Register the transaction manager
	app.registerComponent(db.NewTxManager).Arg(app.PostgresPool)

	// Register ClickHouse connection
	app.registerComponent(func() *infra.ClickHouseConnImpl {
		return infra.NewClickHouseConn(app.ClickHouse)
	})

	app.registerComponent(kafka.NewTopicProducerCreator).Arg(app.KafkaAsyncProducer)

	// Kafka
	eventProducer := kafka.NewTopicProducer(app.KafkaAsyncProducer, domain.EventsKafkaTopic)

	// for /envelope
	envelopeHighConsumer, err := kafka.NewConsumer(
		app.Config.Kafka.Brokers,
		domain.EnvelopeTopicHigh,
		app.Config.Kafka.ConsumerGroupID,
	)
	if err != nil {
		return fmt.Errorf("create envelope high consumer: %w", err)
	}
	envelopeNormalConsumer, err := kafka.NewConsumer(
		app.Config.Kafka.Brokers,
		domain.EnvelopeTopicNormal,
		app.Config.Kafka.ConsumerGroupID,
	)
	if err != nil {
		return fmt.Errorf("create envelope normal consumer: %w", err)
	}
	envelopeLowConsumer, err := kafka.NewConsumer(
		app.Config.Kafka.Brokers,
		domain.EnvelopeTopicLow,
		app.Config.Kafka.ConsumerGroupID,
	)
	if err != nil {
		return fmt.Errorf("create envelope low consumer: %w", err)
	}

	// for /store
	storeEventHighConsumer, err := kafka.NewConsumer(
		app.Config.Kafka.Brokers,
		domain.StoreEventTopicHigh,
		app.Config.Kafka.ConsumerGroupID,
	)
	if err != nil {
		return fmt.Errorf("create store event high consumer: %w", err)
	}
	storeEventNormalConsumer, err := kafka.NewConsumer(
		app.Config.Kafka.Brokers,
		domain.StoreEventTopicNormal,
		app.Config.Kafka.ConsumerGroupID,
	)
	if err != nil {
		return fmt.Errorf("create store event normal consumer: %w", err)
	}
	storeEventLowConsumer, err := kafka.NewConsumer(
		app.Config.Kafka.Brokers,
		domain.StoreEventTopicLow,
		app.Config.Kafka.ConsumerGroupID,
	)
	if err != nil {
		return fmt.Errorf("create store event low consumer: %w", err)
	}

	// Register cache manager
	app.registerComponent(cachemanager.New).Arg(&app.Config.Cache)
	app.registerComponent(cacheservice.New).Arg(&app.Config.Cache)

	// Register repositories
	app.registerComponent(issues.New).Arg(app.PostgresPool)
	app.registerComponent(events.New).Arg(eventProducer)
	app.registerComponent(releases.New).Arg(app.PostgresPool)
	app.registerComponent(notificationsqueue.New).Arg(app.PostgresPool)
	app.registerComponent(issuereleases.New).Arg(app.PostgresPool)

	// Register use cases
	app.registerComponent(envelopeusecase.New)
	app.registerComponent(eventsusecase.New)
	app.registerComponent(storeeventusecase.New)

	// Register services
	var topicProducerCreator *kafka.TopicProducerCreator
	if err := app.container.Resolve(&topicProducerCreator); err != nil {
		panic(err)
	}

	app.registerComponent(storeeventqueueproducer.New).Arg(topicProducerCreator)
	app.registerComponent(envelopequeueprocessor.New).Arg([]contract.DataConsumer{
		envelopeHighConsumer,
		envelopeNormalConsumer,
		envelopeLowConsumer,
	})
	app.registerComponent(storeeventqueueprocessor.New).Arg([]contract.DataConsumer{
		storeEventHighConsumer,
		storeEventNormalConsumer,
		storeEventLowConsumer,
	})

	var envelopeConsumers *envelopequeueprocessor.Service
	if err := app.container.Resolve(&envelopeConsumers); err != nil {
		panic(err)
	}

	// Resolve consumers to start
	var storeEventConsumers *storeeventqueueprocessor.Service
	if err := app.container.Resolve(&storeEventConsumers); err != nil {
		panic(err)
	}

	return nil
}

func (app *App) registerComponent(constructor any) *di.Provider {
	return app.container.Provide(constructor)
}
