package runner

import (
	"bytes"
	"context"
	"encoding/base32"
	"encoding/base64"
	"fmt"
	"log/slog"
	"net"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/pquerna/otp"
	"github.com/pquerna/otp/totp"
	"github.com/rom8726/pgfixtures"
	"github.com/rom8726/testy"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/clickhouse"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
	"gopkg.in/yaml.v3"

	"github.com/rom8726/warden/internal/backend"
	"github.com/rom8726/warden/internal/backend/config"
	"github.com/rom8726/warden/pkg/crypt"
)

type Config struct {
	CasesDir string
	UsesOTP  bool

	BeforeReq func(app *backend.App) error
	AfterReq  func(app *backend.App) error
}

func Run(t *testing.T, testCfg *Config) {
	t.Helper()

	env := NewEnv()

	var pgCconnStr string
	var err error
	// Postgres ---------------------------------------------------------------
	dbType := pgfixtures.PostgreSQL
	pgContainer, pgDown := startPostgres(t)
	defer pgDown()

	pgCconnStr, err = pgContainer.ConnectionString(t.Context(), "sslmode=disable")
	require.NoError(t, err)
	env.Set("WARDEN_POSTGRES_PORT", extractPort(pgCconnStr))
	// Redis ------------------------------------------------------------------
	redisC, redisDown := startRedis(t)
	defer redisDown()

	redisPort, err := redisC.MappedPort(t.Context(), "6379")
	require.NoError(t, err)
	env.Set("WARDEN_REDIS_PORT", redisPort.Port())

	// Kafka ------------------------------------------------------------------
	zkC, zkDown := startZookeeper(t)
	defer zkDown()

	kafkaC, kafkaDown := startKafka(t, zkC)
	defer kafkaDown()

	kafkaPort, err := kafkaC.MappedPort(t.Context(), "9092")
	require.NoError(t, err)
	env.Set("WARDEN_KAFKA_BROKERS", "localhost:"+kafkaPort.Port())

	// ClickHouse -------------------------------------------------------------
	clickC, clickDown := startClickHouse(t)
	defer clickDown()

	clickPort, err := clickC.MappedPort(t.Context(), "9000")
	require.NoError(t, err)
	env.Set("WARDEN_CLICKHOUSE_HOST", "localhost")
	env.Set("WARDEN_CLICKHOUSE_PORT", clickPort.Port())
	env.Set("WARDEN_CLICKHOUSE_USER", "default")
	env.Set("WARDEN_CLICKHOUSE_PASSWORD", "password")
	env.Set("WARDEN_CLICKHOUSE_DATABASE", "warden")

	chConnStr, err := clickC.ConnectionString(t.Context())
	require.NoError(t, err)

	// MailHog ----------------------------------------------------------------
	mailC, mailDown := startMailHog(t)
	defer mailDown()
	mailPort, err := mailC.MappedPort(t.Context(), "1025")
	require.NoError(t, err)
	env.Set("WARDEN_MAILER_ADDR", "localhost:"+mailPort.Port())

	mailPort, _ = mailC.MappedPort(t.Context(), "8025")
	fmt.Println("MailHog UI: http://localhost:" + mailPort.Port())

	// Config and App initialization ------------------------------------------
	env.SetUp()
	defer env.CleanUp()

	cfg, err := config.New("")
	if err != nil {
		t.Fatal(err)
	}

	loggerHandler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: &cfg.Logger,
	})
	logger := slog.New(loggerHandler)
	slog.SetDefault(logger)

	time.Sleep(time.Second * 10)

	if err := upPGMigrations(pgCconnStr, cfg.Postgres.MigrationsDir); err != nil {
		t.Fatal(err)
	}
	if err := upClickHouseMigrations(chConnStr, cfg.ClickHouse.MigrationsDir); err != nil {
		t.Fatal(err)
	}

	app, err := backend.NewApp(t.Context(), cfg, logger)
	if err != nil {
		t.Fatal(err)
	}
	defer app.Close()

	if testCfg.UsesOTP {
		modifiedFixtures := setValidFASecretsInFixtures(t, "./fixtures")

		defer func() {
			// --- reset modified fixtures ---
			if len(modifiedFixtures) > 0 {
				for _, filePath := range modifiedFixtures {
					resetFixtureFile(t, filePath)
				}
			}
		}()
	}

	if err := app.TruncateCHTables(t.Context()); err != nil {
		t.Fatal(err)
	}

	var (
		beforeReq func() error
		afterReq  func() error
	)
	if testCfg.BeforeReq != nil {
		beforeReq = func() error {
			return testCfg.BeforeReq(app)
		}
	}
	if testCfg.AfterReq != nil {
		afterReq = func() error {
			return testCfg.AfterReq(app)
		}
	}

	testyCfg := testy.Config{
		Handler:     app.APIServer.Handler,
		DBType:      dbType,
		CasesDir:    testCfg.CasesDir,
		FixturesDir: "./fixtures",
		ConnStr:     pgCconnStr,
		BeforeReq:   beforeReq,
		AfterReq:    afterReq,
	}
	testy.Run(t, &testyCfg)
}

// Postgres -----------------------------------------------------------------.
func startPostgres(t *testing.T) (*postgres.PostgresContainer, func()) {
	t.Helper()

	container, err := postgres.Run(t.Context(),
		"postgres:16",
		postgres.WithDatabase("test_db"),
		postgres.WithUsername("user"),
		postgres.WithPassword("password"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(10*time.Second),
		),
	)
	require.NoError(t, err)

	return container, func() {
		if err := container.Terminate(context.Background()); err != nil {
			t.Fatalf("terminate postgres: %v", err)
		}
	}
}

// Redis --------------------------------------------------------------------.
func startRedis(t *testing.T) (testcontainers.Container, func()) {
	t.Helper()

	container, err := testcontainers.GenericContainer(t.Context(), testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Name:         "warden-redis-test",
			Image:        "redis:7-alpine",
			ExposedPorts: []string{"6379/tcp"},
			Env: map[string]string{
				"REDIS_PASSWORD": "password",
			},
			WaitingFor: wait.ForLog("Ready to accept connections"),
		},
		Started: true,
	})
	require.NoError(t, err)

	return container, func() {
		if err := container.Terminate(context.Background()); err != nil {
			t.Fatalf("terminate redis: %v", err)
		}
	}
}

// ZooKeeper -----------------------------------------------------------------
func startZookeeper(t *testing.T) (testcontainers.Container, func()) {
	t.Helper()

	c, err := testcontainers.GenericContainer(t.Context(), testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Name:         "warden-zookeeper-test",
			Image:        "bitnami/zookeeper:3.9.3",
			ExposedPorts: []string{"2181/tcp"},
			Env: map[string]string{
				"ALLOW_ANONYMOUS_LOGIN": "yes",
			},
			WaitingFor: wait.ForLog("Started AdminServer on address").WithStartupTimeout(30 * time.Second),
		},
		Started: true,
	})
	require.NoError(t, err)

	return c, func() {
		if err := c.Terminate(context.Background()); err != nil {
			t.Fatalf("terminate zookeeper: %v", err)
		}
	}
}

// Kafka --------------------------------------------------------------------
func startKafka(t *testing.T, zk testcontainers.Container) (testcontainers.Container, func()) {
	t.Helper()

	zkIP, err := zk.ContainerIP(t.Context())
	require.NoError(t, err)

	zookeeperConnect := fmt.Sprintf("%s:2181", zkIP)

	c, err := testcontainers.GenericContainer(t.Context(), testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Name:         "warden-kafka-test",
			Image:        "bitnami/kafka:3.6.0",
			ExposedPorts: []string{"9092/tcp"},
			Env: map[string]string{
				"KAFKA_BROKER_ID":                     "1",
				"KAFKA_CFG_ZOOKEEPER_CONNECT":         zookeeperConnect,
				"KAFKA_CFG_ADVERTISED_LISTENERS":      "PLAINTEXT://localhost:9092",
				"KAFKA_CFG_LISTENERS":                 "PLAINTEXT://:9092",
				"KAFKA_CFG_AUTO_CREATE_TOPICS_ENABLE": "true",
				"ALLOW_PLAINTEXT_LISTENER":            "yes",
			},
			WaitingFor: wait.ForLog("[KafkaServer id=1] started (kafka.server.KafkaServer)").
				WithStartupTimeout(3 * time.Minute),
		},
		Started: true,
	})
	require.NoError(t, err)

	return c, func() {
		if err := c.Terminate(context.Background()); err != nil {
			t.Fatalf("terminate kafka: %v", err)
		}
	}
}

// ClickHouse ----------------------------------------------------------------
func startClickHouse(t *testing.T) (*clickhouse.ClickHouseContainer, func()) {
	t.Helper()

	c, err := clickhouse.Run(t.Context(),
		"clickhouse/clickhouse-server:latest",
		clickhouse.WithDatabase("warden"),
		clickhouse.WithPassword("password"),
		clickhouse.WithUsername("default"),
		testcontainers.WithEnv(map[string]string{"CLICKHOUSE_DEFAULT_ACCESS_MANAGEMENT": "1"}),
		testcontainers.WithWaitStrategy(
			wait.ForLog("/entrypoint.sh: create database").
				WithStartupTimeout(time.Minute),
		),
	)
	require.NoError(t, err)

	return c, func() {
		if err := c.Terminate(context.Background()); err != nil {
			t.Fatalf("terminate clickhouse: %v", err)
		}
	}
}

// MailHog ----------------------------------------------------------------
func startMailHog(t *testing.T) (testcontainers.Container, func()) {
	t.Helper()

	c, err := testcontainers.GenericContainer(t.Context(), testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Name:         "warden-mailhog-test",
			Image:        "mailhog/mailhog:latest",
			ExposedPorts: []string{"1025/tcp", "8025/tcp"},
			WaitingFor:   wait.ForListeningPort("1025"),
		},
		Started: true,
	})
	require.NoError(t, err)

	return c, func() {
		if err := c.Terminate(context.Background()); err != nil {
			t.Fatalf("terminate mailhog: %v", err)
		}
	}
}

func extractPort(connStr string) string {
	// Check if it's a MySQL connection string (user:password@tcp(host:port)/dbname)
	if strings.Contains(connStr, "@tcp(") {
		start := strings.Index(connStr, "@tcp(")
		if start == -1 {
			return ""
		}
		start += 5 // Skip "@tcp("

		end := strings.Index(connStr[start:], ")")
		if end == -1 {
			return ""
		}

		hostPort := connStr[start : start+end]
		_, port, _ := net.SplitHostPort(hostPort)

		return port
	}

	// Otherwise, assume it's a PostgreSQL connection string (postgres://user:password@host:port/dbname)
	u, err := url.Parse(connStr)
	if err != nil {
		return ""
	}

	host, port, _ := net.SplitHostPort(u.Host)
	if host == "" {
		return ""
	}

	return port
}

func setValidFASecretsInFixtures(t *testing.T, fixturesDir string) []string {
	t.Helper()

	var modifiedFiles []string

	secretKey := os.Getenv("WARDEN_JWT_SECRET_KEY")
	if secretKey == "" {
		t.Fatal("WARDEN_JWT_SECRET_KEY is not set")
	}

	files, err := filepath.Glob(filepath.Join(fixturesDir, "*.yml"))
	if err != nil {
		t.Fatalf("failed to list fixture files: %v", err)
	}

	for _, filePath := range files {
		content, err := os.ReadFile(filePath)
		if err != nil {
			t.Fatalf("failed to read file %s: %v", filePath, err)
		}

		var data map[string]any
		if err := yaml.Unmarshal(content, &data); err != nil {
			t.Fatalf("failed to parse YAML in %s: %v", filePath, err)
		}

		modified := false

		usersRaw, ok := data["public.users"]
		if !ok {
			continue
		}

		users, ok := usersRaw.([]any)
		if !ok {
			t.Fatalf("invalid format for public.users in %s", filePath)
		}

		for _, u := range users {
			user, ok := u.(map[string]any)
			if !ok {
				continue
			}

			twoFA, enabled := user["two_fa_enabled"].(bool)
			email, hasEmail := user["email"].(string)
			if enabled && twoFA && hasEmail && email != "" {
				enc := generateValid2FASecret(t, secretKey)
				user["two_fa_secret"] = enc
				modified = true
			}
		}

		if modified {
			var buf bytes.Buffer
			encoder := yaml.NewEncoder(&buf)
			encoder.SetIndent(2)
			if err := encoder.Encode(data); err != nil {
				t.Fatalf("failed to marshal YAML for %s: %v", filePath, err)
			}
			_ = encoder.Close()

			if err := os.WriteFile(filePath, buf.Bytes(), 0644); err != nil {
				t.Fatalf("failed to write updated file %s: %v", filePath, err)
			}

			modifiedFiles = append(modifiedFiles, filePath)
		}
	}

	return modifiedFiles
}

func generateValid2FASecret(t *testing.T, secret string) string {
	t.Helper()

	generatedSecret := findOTPSecret(t, "123456")
	encSecret, err := crypt.EncryptAESGCM([]byte(generatedSecret), []byte(secret))
	if err != nil {
		t.Fatal(err)
	}

	return base64.StdEncoding.EncodeToString(encSecret)
}

func findOTPSecret(t *testing.T, targetCode string) string {
	t.Helper()

	now := time.Now()

	for i := 0; i < 10_000_000; i++ {
		seed := fmt.Sprintf("SECRET%d", i)
		secret := base32.StdEncoding.EncodeToString([]byte(seed))
		secret = strings.TrimRight(secret, "=")

		code, err := totp.GenerateCodeCustom(secret, now, totp.ValidateOpts{
			Period:    30,
			Skew:      0,
			Digits:    otp.DigitsSix,
			Algorithm: otp.AlgorithmSHA1,
		})
		if err != nil {
			panic(err)
		}

		if code == targetCode && totp.Validate(targetCode, secret) {
			fmt.Printf(">>> Found valid 2FA OTP secret [iteration: %d]\n", i+1)

			return secret
		}
	}

	t.Fatal("Failed to find a valid 2FA OTP secret")

	return ""
}

func resetFixtureFile(t *testing.T, filePath string) {
	t.Helper()

	cmd := exec.Command("git", "checkout", "--", filePath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	fmt.Printf("⏪ Resetting fixture: %s\n", filePath)

	_ = cmd.Run()
}
