package runner

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"log/slog"

	_ "github.com/ClickHouse/clickhouse-go/v2" //nolint:revive // it's ok
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/clickhouse" //nolint:revive // it's ok
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

func upPGMigrations(connStr, migrationsDir string) error {
	slog.Info("up migrations...")

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return fmt.Errorf("open postgres connection: %w", err)
	}
	defer func() { _ = db.Close() }()

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("create postgres driver: %w", err)
	}

	pgMigrate, err := migrate.NewWithDatabaseInstance("file://"+migrationsDir, "postgres", driver)
	if err != nil {
		return fmt.Errorf("create migrations: %w", err)
	}

	if err := pgMigrate.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			slog.Info("up migrations: no changes")

			return nil
		}

		return fmt.Errorf("up: %w", err)
	}

	slog.Info("up migrations: done")

	return nil
}

func upClickHouseMigrations(connStr, migrationsDir string) error {
	slog.Info("up clickhouse migrations...")

	migrator, err := migrate.New(
		"file://"+migrationsDir,
		connStr,
	)
	if err != nil {
		log.Fatalf("failed to create migrate instance: %v", err)
	}

	if err := migrator.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("up: %w", err)
	}

	slog.Info("up migrations: done")

	return nil
}
