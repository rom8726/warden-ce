package config

import (
	"log/slog"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLogger_Level(t *testing.T) {
	tests := []struct {
		name     string
		level    string
		expected slog.Level
		panics   bool
	}{
		{
			name:     "debug level",
			level:    "debug",
			expected: slog.LevelDebug,
		},
		{
			name:     "info level",
			level:    "info",
			expected: slog.LevelInfo,
		},
		{
			name:     "warn level",
			level:    "warn",
			expected: slog.LevelWarn,
		},
		{
			name:     "error level",
			level:    "error",
			expected: slog.LevelError,
		},
		{
			name:   "invalid level",
			level:  "invalid",
			panics: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger := &Logger{Lvl: tt.level}

			if tt.panics {
				assert.Panics(t, func() {
					logger.Level()
				})
			} else {
				result := logger.Level()
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}

func TestPostgres_ConnString(t *testing.T) {
	tests := []struct {
		name     string
		postgres Postgres
		expected string
	}{
		{
			name: "with user and password",
			postgres: Postgres{
				User:     "testuser",
				Password: "testpass",
				Host:     "localhost",
				Port:     "5432",
				Database: "testdb",
			},
			expected: "postgres://testuser:testpass@localhost:5432/testdb?connect_timeout=10&sslmode=disable",
		},
		{
			name: "with user only",
			postgres: Postgres{
				User:     "testuser",
				Password: "",
				Host:     "localhost",
				Port:     "5432",
				Database: "testdb",
			},
			expected: "postgres://testuser:@localhost:5432/testdb?connect_timeout=10&sslmode=disable",
		},
		{
			name: "without user",
			postgres: Postgres{
				User:     "",
				Password: "testpass",
				Host:     "localhost",
				Port:     "5432",
				Database: "testdb",
			},
			expected: "postgres://localhost:5432/testdb?connect_timeout=10&sslmode=disable",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.postgres.ConnString()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestPostgres_ConnStringWithPoolSize(t *testing.T) {
	postgres := Postgres{
		User:     "testuser",
		Password: "testpass",
		Host:     "localhost",
		Port:     "5432",
		Database: "testdb",
		MaxConns: 25,
	}

	expected := "postgres://testuser:testpass@localhost:5432/testdb?connect_timeout=10&sslmode=disable&pool_max_conns=25"
	result := postgres.ConnStringWithPoolSize()

	assert.Equal(t, expected, result)
}

func TestPostgres_MigrationConnString(t *testing.T) {
	tests := []struct {
		name     string
		postgres Postgres
		expected string
	}{
		{
			name: "with migration host and port",
			postgres: Postgres{
				User:          "testuser",
				Password:      "testpass",
				Host:          "pgbouncer",
				Port:          "6432",
				Database:      "testdb",
				MigrationHost: "postgresql",
				MigrationPort: "5432",
			},
			expected: "postgres://testuser:testpass@postgresql:5432/testdb?connect_timeout=10&sslmode=disable",
		},
		{
			name: "without migration host and port (fallback to regular)",
			postgres: Postgres{
				User:     "testuser",
				Password: "testpass",
				Host:     "localhost",
				Port:     "5432",
				Database: "testdb",
			},
			expected: "postgres://testuser:testpass@localhost:5432/testdb?connect_timeout=10&sslmode=disable",
		},
		{
			name: "with only migration host",
			postgres: Postgres{
				User:          "testuser",
				Password:      "testpass",
				Host:          "pgbouncer",
				Port:          "6432",
				Database:      "testdb",
				MigrationHost: "postgresql",
			},
			expected: "postgres://testuser:testpass@postgresql:6432/testdb?connect_timeout=10&sslmode=disable",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.postgres.MigrationConnString()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestClickHouse_ConnString(t *testing.T) {
	clickhouse := ClickHouse{
		User:     "testuser",
		Password: "testpass",
		Host:     "localhost",
		Port:     9000,
		Database: "testdb",
	}

	expected := "clickhouse://testuser:testpass@localhost:9000/testdb"
	result := clickhouse.ConnString()

	assert.Equal(t, expected, result)
}
