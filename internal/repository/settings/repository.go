package settings

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/rom8726/warden/internal/domain"
	"github.com/rom8726/warden/pkg/db"
)

// Repository implements domain.SettingRepository.
type Repository struct {
	db db.Tx
}

// New creates a new settings repository.
func New(pool *pgxpool.Pool) *Repository {
	return &Repository{
		db: pool,
	}
}

// GetByName retrieves a setting by name.
func (r *Repository) GetByName(ctx context.Context, name string) (*domain.Setting, error) {
	executor := r.getExecutor(ctx)

	const query = `
		SELECT id, name, value, description, created_at, updated_at
		FROM settings
		WHERE name = $1
		LIMIT 1
	`

	rows, err := executor.Query(ctx, query, name)
	if err != nil {
		return nil, fmt.Errorf("failed to get setting by name %s: %w", name, err)
	}
	defer rows.Close()

	model, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[settingModel])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrEntityNotFound
		}

		return nil, fmt.Errorf("failed to collect setting: %w", err)
	}

	return model.toDomain(), nil
}

// SetByName creates or updates a setting by name.
func (r *Repository) SetByName(ctx context.Context, name string, value interface{}, description string) error {
	executor := r.getExecutor(ctx)

	valueJSON, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("failed to marshal setting value: %w", err)
	}

	const query = `
		INSERT INTO settings (name, value, description)
		VALUES ($1, $2, $3)
		ON CONFLICT (name) 
		DO UPDATE SET 
			value = EXCLUDED.value,
			description = EXCLUDED.description,
			updated_at = NOW()
	`

	_, err = executor.Exec(ctx, query, name, valueJSON, description)
	if err != nil {
		return fmt.Errorf("failed to set setting %s: %w", name, err)
	}

	return nil
}

// DeleteByName deletes a setting by name.
func (r *Repository) DeleteByName(ctx context.Context, name string) error {
	executor := r.getExecutor(ctx)

	const query = `DELETE FROM settings WHERE name = $1`

	tag, err := executor.Exec(ctx, query, name)
	if err != nil {
		return fmt.Errorf("failed to delete setting %s: %w", name, err)
	}

	if tag.RowsAffected() == 0 {
		return domain.ErrEntityNotFound
	}

	return nil
}

// List retrieves all settings.
func (r *Repository) List(ctx context.Context) ([]*domain.Setting, error) {
	executor := r.getExecutor(ctx)

	const query = `
		SELECT id, name, value, description, created_at, updated_at
		FROM settings
		ORDER BY name
	`

	rows, err := executor.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to list settings: %w", err)
	}
	defer rows.Close()

	models, err := pgx.CollectRows(rows, pgx.RowToStructByName[settingModel])
	if err != nil {
		return nil, fmt.Errorf("failed to collect settings: %w", err)
	}

	settings := make([]*domain.Setting, len(models))
	for i, model := range models {
		settings[i] = model.toDomain()
	}

	return settings, nil
}

//nolint:ireturn // it's ok here
func (r *Repository) getExecutor(ctx context.Context) db.Tx {
	if tx := db.TxFromContext(ctx); tx != nil {
		return tx
	}

	return r.db
}
