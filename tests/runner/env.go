package runner

import (
	"fmt"
	"os"
)

func NewEnv() Env {
	return Env{
		collidedVars: make(map[string]bool),
		redefinedVars: map[string]string{
			"WARDEN_ENVIRONMENT":               "test",
			"WARDEN_SECRET_KEY":                "secret_123456789",
			"WARDEN_API_SERVER_ADDR":           ":8080",
			"WARDEN_TECH_SERVER_ADDR":          ":8081",
			"WARDEN_POSTGRES_HOST":             "localhost",
			"WARDEN_POSTGRES_PORT":             "5432",
			"WARDEN_POSTGRES_USER":             "user",
			"WARDEN_POSTGRES_PASSWORD":         "password",
			"WARDEN_POSTGRES_DATABASE":         "test_db",
			"WARDEN_POSTGRES_MIGRATIONS_DIR":   "../migrations/postgresql",
			"WARDEN_REDIS_HOST":                "localhost",
			"WARDEN_REDIS_PORT":                "6379",
			"WARDEN_REDIS_PASSWORD":            "password",
			"WARDEN_REDIS_DB":                  "0",
			"WARDEN_KAFKA_BROKERS":             "localhost:9092",
			"WARDEN_KAFKA_CLIENT_ID":           "app",
			"WARDEN_CLICKHOUSE_HOST":           "localhost",
			"WARDEN_CLICKHOUSE_PORT":           "9000",
			"WARDEN_CLICKHOUSE_DATABASE":       "warden",
			"WARDEN_CLICKHOUSE_USER":           "default",
			"WARDEN_CLICKHOUSE_PASSWORD":       "password",
			"WARDEN_CLICKHOUSE_TIMEOUT":        "10s",
			"WARDEN_CLICKHOUSE_MIGRATIONS_DIR": "../migrations/clickhouse",
			"WARDEN_EMAIL_SMTP_HOST":           "warden-mailhog",
			"WARDEN_EMAIL_SMTP_PORT":           "1025",
			"WARDEN_EMAIL_FROM":                "noreply@warden.local",
			"WARDEN_EMAIL_FROM_NAME":           "Warden",
			"WARDEN_MAILER_ADDR":               "mailhog:1025",
			"WARDEN_MAILER_USER":               "warden",
			"WARDEN_MAILER_PASSWORD":           "WardenQwe321!",
			"WARDEN_MAILER_FROM":               "noreply@warden.local",
			"WARDEN_MAILER_ALLOW_INSECURE":     "true",
			"WARDEN_MAILER_USE_TLS":            "false",
			"WARDEN_JWT_SECRET_KEY":            "secret_key123456",
			"WARDEN_ACCESS_TOKEN_TTL":          "3h",
			"WARDEN_REFRESH_TOKEN_TTL":         "168h",
			"WARDEN_RESET_PASSWORD_TTL":        "8h",
		},
	}
}

type Env struct {
	redefinedVars map[string]string
	collidedVars  map[string]bool
}

func (e *Env) SetUp() {
	var err error
	for key, value := range e.redefinedVars {
		if envVar := os.Getenv(key); envVar != "" {
			e.redefinedVars[key] = envVar
			e.collidedVars[key] = true

			continue
		}
		if err = os.Setenv(key, value); err != nil {
			err = fmt.Errorf("can't clear ENV %s: %w", key, err)
			panic(err)
		}
	}
}

func (e *Env) CleanUp() {
	var err error
	for key := range e.redefinedVars {
		if _, ok := e.collidedVars[key]; ok {
			continue
		}
		err = os.Unsetenv(key)
		if err != nil {
			err = fmt.Errorf("can't clear ENV %s: %w", key, err)
			panic(err)
		}
	}
}

func (e *Env) Set(key, val string) {
	e.redefinedVars[key] = val
}

func (e *Env) Get(key string) string {
	return e.redefinedVars[key]
}

func (*Env) SetMock(key, server string) {
	if err := os.Setenv(key, server); err != nil {
		err = fmt.Errorf("can't set mock ENV %s: %w", key, err)
		panic(err)
	}
}
