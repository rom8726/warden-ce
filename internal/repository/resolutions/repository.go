package resolutions

import (
	"context"
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

func (r *Repository) Create(ctx context.Context, resolutionDTO domain.ResolutionDTO) (domain.Resolution, error) {
	executor := r.getExecutor(ctx)

	model := fromDTO(resolutionDTO)

	// If status is resolved, set the resolved_at timestamp
	if resolutionDTO.Status == domain.IssueStatusResolved {
		now := time.Now()
		model.ResolvedAt = &now
	}

	const query = `
INSERT INTO resolutions (project_id, issue_id, status, resolved_by, resolved_at, comment, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING *`

	rows, err := executor.Query(ctx, query,
		model.ProjectID,
		model.IssueID,
		model.Status,
		model.ResolvedBy,
		model.ResolvedAt,
		model.Comment,
		model.CreatedAt,
		model.UpdatedAt,
	)
	if err != nil {
		return domain.Resolution{}, fmt.Errorf("insert resolution: %w", err)
	}
	defer rows.Close()

	resolution, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[resolutionModel])
	if err != nil {
		return domain.Resolution{}, fmt.Errorf("collect resolution: %w", err)
	}

	return resolution.toDomain(), nil
}

func (r *Repository) GetByID(ctx context.Context, id domain.ResolutionID) (domain.Resolution, error) {
	executor := r.getExecutor(ctx)

	const query = `SELECT * FROM resolutions WHERE id = $1 LIMIT 1`

	rows, err := executor.Query(ctx, query, id)
	if err != nil {
		return domain.Resolution{}, fmt.Errorf("query resolution by ID: %w", err)
	}
	defer rows.Close()

	resolution, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[resolutionModel])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Resolution{}, domain.ErrEntityNotFound
		}

		return domain.Resolution{}, fmt.Errorf("collect resolution: %w", err)
	}

	return resolution.toDomain(), nil
}

func (r *Repository) Update(ctx context.Context, resolution *domain.Resolution) error {
	executor := r.getExecutor(ctx)

	model := fromDomain(resolution)

	const query = `
UPDATE resolutions
SET project_id = $1, issue_id = $2, status = $3, resolved_by = $4, resolved_at = $5, comment = $6, updated_at = $7
WHERE id = $8`

	_, err := executor.Exec(ctx, query,
		model.ProjectID,
		model.IssueID,
		model.Status,
		model.ResolvedBy,
		model.ResolvedAt,
		model.Comment,
		time.Now(), // Always update the updated_at timestamp
		model.ID,
	)
	if err != nil {
		return fmt.Errorf("update resolution: %w", err)
	}

	return nil
}

func (r *Repository) Delete(ctx context.Context, id domain.ResolutionID) error {
	executor := r.getExecutor(ctx)

	const query = `
DELETE FROM resolutions
WHERE id = $1`

	_, err := executor.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("delete resolution: %w", err)
	}

	return nil
}

func (r *Repository) List(ctx context.Context, projectID domain.ProjectID) ([]domain.Resolution, error) {
	executor := r.getExecutor(ctx)

	const query = `
SELECT *
FROM resolutions
WHERE project_id = $1
ORDER BY updated_at DESC`

	rows, err := executor.Query(ctx, query, projectID)
	if err != nil {
		return nil, fmt.Errorf("query resolutions: %w", err)
	}
	defer rows.Close()

	listModels, err := pgx.CollectRows(rows, pgx.RowToStructByName[resolutionModel])
	if err != nil {
		return nil, fmt.Errorf("collect resolutions: %w", err)
	}

	resolutions := make([]domain.Resolution, 0, len(listModels))
	for i := range listModels {
		model := listModels[i]
		resolutions = append(resolutions, model.toDomain())
	}

	return resolutions, nil
}

func (r *Repository) GetByIssueID(ctx context.Context, issueID domain.IssueID) ([]domain.Resolution, error) {
	executor := r.getExecutor(ctx)

	const query = `
SELECT *
FROM resolutions
WHERE issue_id = $1
ORDER BY resolved_at DESC`

	rows, err := executor.Query(ctx, query, issueID)
	if err != nil {
		return nil, fmt.Errorf("query resolutions by issue ID: %w", err)
	}
	defer rows.Close()

	listModels, err := pgx.CollectRows(rows, pgx.RowToStructByName[resolutionModel])
	if err != nil {
		return nil, fmt.Errorf("collect resolutions: %w", err)
	}

	resolutions := make([]domain.Resolution, 0, len(listModels))
	for i := range listModels {
		model := listModels[i]
		resolutions = append(resolutions, model.toDomain())
	}

	return resolutions, nil
}

//nolint:ireturn // it's ok here
func (r *Repository) getExecutor(ctx context.Context) db.Tx {
	if tx := db.TxFromContext(ctx); tx != nil {
		return tx
	}

	return r.db
}
