package releases

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/rom8726/warden/internal/domain"
	"github.com/rom8726/warden/pkg/db"
)

type Repository struct {
	db db.Tx
}

func New(pool *pgxpool.Pool) *Repository {
	return &Repository{
		db: pool,
	}
}

func (r *Repository) GetByID(ctx context.Context, id domain.ReleaseID) (domain.Release, error) {
	executor := r.getExecutor(ctx)

	const query = `SELECT * FROM releases WHERE id = $1 LIMIT 1`

	rows, err := executor.Query(ctx, query, id)
	if err != nil {
		return domain.Release{}, fmt.Errorf("query release by ID: %w", err)
	}
	defer rows.Close()

	release, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[releaseModel])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Release{}, domain.ErrEntityNotFound
		}

		return domain.Release{}, fmt.Errorf("collect release: %w", err)
	}

	return release.toDomain(), nil
}

func (r *Repository) GetByProjectAndVersion(
	ctx context.Context,
	projectID domain.ProjectID,
	version string,
) (domain.Release, error) {
	executor := r.getExecutor(ctx)

	const query = `SELECT * FROM releases WHERE project_id = $1 AND version = $2 LIMIT 1`

	rows, err := executor.Query(ctx, query, projectID, version)
	if err != nil {
		return domain.Release{}, fmt.Errorf("query release by project and version: %w", err)
	}
	defer rows.Close()

	release, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[releaseModel])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Release{}, domain.ErrEntityNotFound
		}

		return domain.Release{}, fmt.Errorf("collect release: %w", err)
	}

	return release.toDomain(), nil
}

func (r *Repository) ListByProject(ctx context.Context, projectID domain.ProjectID) ([]domain.Release, error) {
	executor := r.getExecutor(ctx)

	const query = `SELECT * FROM releases WHERE project_id = $1 ORDER BY released_at DESC, id DESC`

	rows, err := executor.Query(ctx, query, projectID)
	if err != nil {
		return nil, fmt.Errorf("query releases by project: %w", err)
	}
	defer rows.Close()

	entities, err := pgx.CollectRows(rows, pgx.RowToStructByName[releaseModel])
	if err != nil {
		return nil, fmt.Errorf("collect releases: %w", err)
	}

	releases := make([]domain.Release, 0, len(entities))
	for i := range entities {
		releases = append(releases, entities[i].toDomain())
	}

	return releases, nil
}

func (r *Repository) Create(ctx context.Context, release domain.ReleaseDTO) (domain.ReleaseID, error) {
	executor := r.getExecutor(ctx)

	const query = `
INSERT INTO releases (project_id, version, description, released_at, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, $6)
ON CONFLICT (project_id, version) DO UPDATE SET version = EXCLUDED.version
RETURNING id`

	now := time.Now()

	desc := sql.NullString{String: release.Description, Valid: release.Description != ""}
	rows, err := executor.Query(ctx, query,
		release.ProjectID,
		release.Version,
		desc,
		now,
		now,
		now,
	)
	if err != nil {
		return 0, fmt.Errorf("query: %w", err)
	}
	defer rows.Close()

	var id domain.ReleaseID
	if rows.Next() {
		err = rows.Scan(&id)
		if err != nil {
			return 0, fmt.Errorf("scan: %w", err)
		}
	}
	if err = rows.Err(); err != nil {
		return 0, fmt.Errorf("rows: %w", err)
	}

	return id, nil
}

func (r *Repository) getExecutor(ctx context.Context) db.Tx {
	tx := db.TxFromContext(ctx)
	if tx != nil {
		return tx
	}

	return r.db
}
