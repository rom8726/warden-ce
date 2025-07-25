package config

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"net/url"
	"strconv"
	"time"

	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/IBM/sarama"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"

	"github.com/rom8726/warden/internal/domain"
	"github.com/rom8726/warden/internal/infra"
	"github.com/rom8726/warden/pkg/kafka"
)

type Logger struct {
	Lvl string `default:"info" envconfig:"LEVEL"`
}

func (l *Logger) Level() slog.Level {
	switch l.Lvl {
	case "debug":
		return slog.LevelDebug
	case "info":
		return slog.LevelInfo
	case "warn":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		panic("invalid logger level " + l.Lvl)
	}
}

type Server struct {
	Addr         string        `envconfig:"ADDR" required:"true"`
	ReadTimeout  time.Duration `default:"15s"    envconfig:"READ_TIMEOUT"`
	WriteTimeout time.Duration `default:"30s"    envconfig:"WRITE_TIMEOUT"`
	IdleTimeout  time.Duration `default:"60s"    envconfig:"IDLE_TIMEOUT"`
}

type Postgres struct {
	User            string        `envconfig:"USER"                  required:"true"`
	Password        string        `envconfig:"PASSWORD"              required:"true"`
	Host            string        `envconfig:"HOST"                  required:"true"`
	Port            string        `default:"5432"                    envconfig:"PORT"`
	Database        string        `envconfig:"DATABASE"              required:"true"`
	MaxIdleConnTime time.Duration `default:"5m"                      envconfig:"MAX_IDLE_CONN_TIME"`
	MaxConns        int           `default:"20"                      envconfig:"MAX_CONNS"`
	ConnMaxLifetime time.Duration `default:"10m"                     envconfig:"CONN_MAX_LIFETIME"`
	MigrationsDir   string        `default:"./migrations/postgresql" envconfig:"MIGRATIONS_DIR"`
	MigrationHost   string        `envconfig:"MIGRATION_HOST"`
	MigrationPort   string        `envconfig:"MIGRATION_PORT"`
}

func (db *Postgres) ConnString() string {
	var user *url.Userinfo

	if db.User != "" {
		var pass string

		if db.Password != "" {
			pass = db.Password
		}

		user = url.UserPassword(db.User, pass)
	}

	params := url.Values{}
	params.Set("sslmode", "disable")
	params.Set("connect_timeout", "10")

	uri := url.URL{
		Scheme:   "postgres",
		User:     user,
		Host:     net.JoinHostPort(db.Host, db.Port),
		Path:     db.Database,
		RawQuery: params.Encode(),
	}

	return uri.String()
}

func (db *Postgres) ConnStringWithPoolSize() string {
	connString := db.ConnString()

	return connString + fmt.Sprintf("&pool_max_conns=%d", db.MaxConns)
}

func (db *Postgres) MigrationConnString() string {
	var user *url.Userinfo

	if db.User != "" {
		var pass string

		if db.Password != "" {
			pass = db.Password
		}

		user = url.UserPassword(db.User, pass)
	}

	params := url.Values{}
	params.Set("sslmode", "disable")
	params.Set("connect_timeout", "10")

	// Use migration host and port if specified, otherwise fall back to regular host and port
	host := db.Host
	port := db.Port
	if db.MigrationHost != "" {
		host = db.MigrationHost
	}
	if db.MigrationPort != "" {
		port = db.MigrationPort
	}

	uri := url.URL{
		Scheme:   "postgres",
		User:     user,
		Host:     net.JoinHostPort(host, port),
		Path:     db.Database,
		RawQuery: params.Encode(),
	}

	return uri.String()
}

type ClickHouse struct {
	Host          string        `envconfig:"HOST"                  required:"true"`
	Port          int           `envconfig:"PORT"                  required:"true"`
	Database      string        `envconfig:"DATABASE"              required:"true"`
	User          string        `envconfig:"USER"                  required:"true"`
	Password      string        `envconfig:"PASSWORD"              required:"true"`
	Timeout       time.Duration `default:"10s"                     envconfig:"TIMEOUT"`
	MigrationsDir string        `default:"./migrations/clickhouse" envconfig:"MIGRATIONS_DIR"`
}

func (ch *ClickHouse) ConnString() string {
	return fmt.Sprintf("clickhouse://%s:%s@%s:%d/%s", ch.User, ch.Password, ch.Host, ch.Port, ch.Database)
}

type Redis struct {
	Host     string `envconfig:"HOST"     required:"true"`
	Port     int    `envconfig:"PORT"     required:"true"`
	DB       int    `envconfig:"DB"       required:"true"`
	Password string `envconfig:"PASSWORD" required:"false"`
}

type Kafka struct {
	Brokers         []string      `envconfig:"BROKERS"       required:"true"               split_words:"true"`
	ClientID        string        `envconfig:"CLIENT_ID"     required:"true"`
	ConsumerGroupID string        `default:"warden-consumer" envconfig:"CONSUMER_GROUP_ID"`
	Version         string        `default:"3.6.0"           envconfig:"VERSION"`
	Timeout         time.Duration `default:"10s"             envconfig:"TIMEOUT"`
}

type Mailer struct {
	Addr          string `envconfig:"ADDR"     required:"true"`
	User          string `envconfig:"USER"     required:"true"`
	Password      string `envconfig:"PASSWORD" required:"true"`
	From          string `envconfig:"FROM"     required:"true"`
	AllowInsecure bool   `default:"false"      envconfig:"ALLOW_INSECURE"`
	CertFile      string `default:""           envconfig:"CERT_FILE"`
	KeyFile       string `default:""           envconfig:"KEY_FILE"`
	UseTLS        bool   `default:"false"      envconfig:"USE_TLS"`
}

// CacheConfig holds cache configuration.
type CacheConfig struct {
	Enabled               bool `default:"true"  envconfig:"ENABLED"`
	ReleaseCacheSize      int  `default:"10000" envconfig:"RELEASE_CACHE_SIZE"`
	IssueCacheSize        int  `default:"10000" envconfig:"ISSUE_CACHE_SIZE"`
	IssueReleaseCacheSize int  `default:"10000" envconfig:"ISSUE_RELEASE_CACHE_SIZE"`
}

// Notificator holds notificator configuration.
type Notificator struct {
	WorkerCount int `default:"4" envconfig:"WORKER_COUNT"`
}

func NewPostgresConnPool(ctx context.Context, cfg *Postgres) (*pgxpool.Pool, error) {
	pgCfg, err := pgxpool.ParseConfig(cfg.ConnStringWithPoolSize())
	if err != nil {
		return nil, fmt.Errorf("parse config: %w", err)
	}

	pgCfg.MaxConnLifetime = cfg.ConnMaxLifetime
	pgCfg.MaxConnLifetimeJitter = time.Second * 5
	pgCfg.MaxConnIdleTime = cfg.MaxIdleConnTime
	pgCfg.HealthCheckPeriod = time.Second * 5

	pool, err := pgxpool.NewWithConfig(ctx, pgCfg)
	if err != nil {
		return nil, err
	}

	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("ping: %w", err)
	}

	return pool, nil
}

func NewRedisClient(ctx context.Context, cfg *Redis) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     net.JoinHostPort(cfg.Host, strconv.Itoa(cfg.Port)),
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	err := client.Ping(ctx).Err()
	if err != nil {
		return nil, fmt.Errorf("ping: %w", err)
	}

	return client, nil
}

func NewClickHouseClient(_ context.Context, cfg *ClickHouse) (infra.ClickHouseConn, error) {
	conn, err := clickhouse.Open(&clickhouse.Options{
		Addr: []string{net.JoinHostPort(cfg.Host, strconv.Itoa(cfg.Port))},
		Auth: clickhouse.Auth{
			Database: cfg.Database,
			Username: cfg.User,
			Password: cfg.Password,
		},
		DialContext: func(ctx context.Context, addr string) (net.Conn, error) {
			var d net.Dialer

			return d.DialContext(ctx, "tcp", addr)
		},
		Debug: false,
		Debugf: func(format string, v ...any) {
			fmt.Printf(format+"\n", v...)
		},
		Settings: clickhouse.Settings{
			"max_execution_time": 60,
		},
		Compression: &clickhouse.Compression{
			Method: clickhouse.CompressionLZ4,
		},
		DialTimeout:          cfg.Timeout,
		MaxOpenConns:         5,
		MaxIdleConns:         5,
		ConnMaxLifetime:      time.Duration(10) * time.Minute,
		ConnOpenStrategy:     clickhouse.ConnOpenInOrder,
		BlockBufferSize:      10,
		MaxCompressionBuffer: 10240,
	})
	if err != nil {
		return nil, fmt.Errorf("open connection: %w", err)
	}

	//if err := conn.Ping(ctx); err != nil {
	//	return nil, fmt.Errorf("ping: %w", err)
	//}

	return infra.NewClickHouseConn(conn), nil
}

func NewKafka(_ context.Context, cfg *Kafka) (sarama.Client, *kafka.Producer, error) {
	version, err := sarama.ParseKafkaVersion(cfg.Version)
	if err != nil {
		return nil, nil, fmt.Errorf("parse kafka version: %w", err)
	}

	saramaConfig := sarama.NewConfig()
	saramaConfig.ClientID = cfg.ClientID
	saramaConfig.Version = version
	saramaConfig.Net.DialTimeout = cfg.Timeout
	saramaConfig.Net.ReadTimeout = cfg.Timeout
	saramaConfig.Net.WriteTimeout = cfg.Timeout
	saramaConfig.Producer.RequiredAcks = sarama.WaitForAll
	saramaConfig.Producer.Return.Successes = true

	client, err := sarama.NewClient(cfg.Brokers, saramaConfig)
	if err != nil {
		return nil, nil, fmt.Errorf("new client: %w", err)
	}

	asyncProducer, err := kafka.NewProducer(cfg.Brokers)
	if err != nil {
		return nil, nil, fmt.Errorf("new async producer: %w", err)
	}

	return client, asyncProducer, nil
}

// CreateKafkaTopics creates Kafka topics for events and exceptions.
func CreateKafkaTopics(brokers []string) error {
	// Define topics to create
	topics := []struct {
		name              string
		partitions        int32
		replicationFactor int16
	}{
		{
			name:              domain.EventsKafkaTopic,
			partitions:        8,
			replicationFactor: 1,
		},
		// Envelope topics
		{
			name:              domain.EnvelopeTopicHigh,
			partitions:        8,
			replicationFactor: 1,
		},
		{
			name:              domain.EnvelopeTopicNormal,
			partitions:        8,
			replicationFactor: 1,
		},
		{
			name:              domain.EnvelopeTopicLow,
			partitions:        8,
			replicationFactor: 1,
		},
		{
			name:              domain.StoreEventTopicLow,
			partitions:        8,
			replicationFactor: 1,
		},
		{
			name:              domain.StoreEventTopicNormal,
			partitions:        8,
			replicationFactor: 1,
		},
		{
			name:              domain.StoreEventTopicHigh,
			partitions:        8,
			replicationFactor: 1,
		},
	}

	// Create topic requests
	requests := make([]kafka.CreateTopicRequest, 0, len(topics))
	for _, t := range topics {
		requests = append(requests, kafka.CreateTopicRequest{
			Topic:             t.name,
			Partitions:        t.partitions,
			ReplicationFactor: t.replicationFactor,
		})
	}

	// Create topics
	if err := kafka.Migrate(brokers, requests...); err != nil {
		return fmt.Errorf("migrate kafka topics: %w", err)
	}

	return nil
}
