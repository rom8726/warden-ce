package projects

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

func (r *Repository) GetByID(ctx context.Context, id domain.ProjectID) (domain.Project, error) {
	executor := r.getExecutor(ctx)

	const query = `SELECT * FROM projects WHERE id = $1 LIMIT 1`

	rows, err := executor.Query(ctx, query, id)
	if err != nil {
		return domain.Project{}, fmt.Errorf("query project by ID: %w", err)
	}
	defer rows.Close()

	project, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[projectModel])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Project{}, domain.ErrEntityNotFound
		}

		return domain.Project{}, fmt.Errorf("collect project: %w", err)
	}

	return project.toDomain(), nil
}

func (r *Repository) GetByPublicKey(ctx context.Context, publicKey string) (domain.Project, error) {
	executor := r.getExecutor(ctx)

	const query = `SELECT * FROM projects WHERE public_key = $1 LIMIT 1`

	rows, err := executor.Query(ctx, query, publicKey)
	if err != nil {
		return domain.Project{}, fmt.Errorf("query project by public key: %w", err)
	}
	defer rows.Close()

	project, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[projectModel])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Project{}, domain.ErrEntityNotFound
		}

		return domain.Project{}, fmt.Errorf("collect project: %w", err)
	}

	return project.toDomain(), nil
}

func (r *Repository) Create(ctx context.Context, project *domain.ProjectDTO) (domain.ProjectID, error) {
	executor := r.getExecutor(ctx)

	const query = `
INSERT INTO projects (name, description, public_key, team_id, created_at)
VALUES ($1, $2, $3, $4, $5)
RETURNING id`

	var id uint
	err := executor.QueryRow(ctx, query,
		project.Name,
		project.Description,
		project.PublicKey,
		project.TeamID,
		time.Now(),
	).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("insert project: %w", err)
	}

	return domain.ProjectID(id), nil
}

func (r *Repository) ValidateProjectKey(ctx context.Context, projectID domain.ProjectID, key string) (bool, error) {
	executor := r.getExecutor(ctx)

	const query = `SELECT 1 FROM projects WHERE id = $1 AND public_key = $2 AND archived_at IS NULL LIMIT 1`

	rows, err := executor.Query(ctx, query, projectID, key)
	if err != nil {
		return false, fmt.Errorf("validate project key: %w", err)
	}
	defer rows.Close()

	// Check if any rows were returned
	if !rows.Next() {
		return false, nil
	}

	return true, nil
}

func (r *Repository) List(ctx context.Context) ([]domain.ProjectExtended, error) {
	executor := r.getExecutor(ctx)

	const query = `
SELECT
    p.id,
    p.public_key,
    p.name,
    p.description,
    p.created_at,
    p.team_id,
    p.archived_at,
    tm.name AS team_name
FROM projects p
LEFT JOIN teams tm ON p.team_id = tm.id
WHERE p.archived_at IS NULL
ORDER BY p.id, tm.id
`

	rows, err := executor.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("query projects: %w", err)
	}
	defer rows.Close()

	listModels, err := pgx.CollectRows(rows, pgx.RowToStructByName[projectModelExtended])
	if err != nil {
		return nil, fmt.Errorf("collect projects: %w", err)
	}

	projects := make([]domain.ProjectExtended, 0, len(listModels))
	for i := range listModels {
		model := listModels[i]
		projects = append(projects, model.toDomain())
	}

	return projects, nil
}

func (r *Repository) RecentProjects(
	ctx context.Context,
	userID domain.UserID,
	limit uint,
) ([]domain.ProjectExtended, error) {
	if limit == 0 {
		limit = 5
	}

	executor := r.getExecutor(ctx)
	const query = `
WITH stats AS (
    SELECT project_id,
          SUM(total_events) AS occurrences,
          MAX(last_seen) AS last_seen,
          SUM(
            CASE
                WHEN level = 'debug' THEN 1
                WHEN level = 'info' THEN 10
                WHEN level = 'warning' THEN 30
                WHEN level = 'error' THEN 100
                WHEN level = 'exception' THEN 150
                WHEN level = 'fatal' THEN 160
                ELSE 0
            END
          ) AS rating
    FROM issues
    WHERE last_seen >= (NOW() - INTERVAL '3 hours')
    GROUP BY project_id
    ORDER BY rating DESC, occurrences DESC, last_seen DESC
    LIMIT $2
)
SELECT
    p.id,
    p.name,
    p.description,
    p.public_key,
    p.created_at,
    p.archived_at,
    t.id AS team_id,
    t.name AS team_name
FROM projects p
LEFT JOIN teams t ON p.team_id = t.id
LEFT JOIN team_members tm ON tm.team_id = p.team_id AND tm.user_id = $1
LEFT JOIN users u ON u.id = $1
LEFT JOIN stats s ON s.project_id = p.id
WHERE
    p.archived_at IS NULL AND (
        p.team_id IS NULL OR (
            CASE
                WHEN u.is_superuser THEN TRUE
                ELSE tm.user_id IS NOT NULL
            END
        )
    )
ORDER BY s.rating DESC NULLS LAST,
         s.occurrences DESC NULLS LAST,
         s.last_seen DESC NULLS LAST,
         p.name
LIMIT $2`

	rows, err := executor.Query(ctx, query, userID, limit)
	if err != nil {
		return nil, fmt.Errorf("query recent projects: %w", err)
	}
	defer rows.Close()

	listModels, err := pgx.CollectRows(rows, pgx.RowToStructByName[projectModelExtended])
	if err != nil {
		return nil, fmt.Errorf("collect recent projects: %w", err)
	}

	projects := make([]domain.ProjectExtended, 0, len(listModels))
	for i := range listModels {
		model := listModels[i]
		projects = append(projects, model.toDomain())
	}

	return projects, nil
}

func (r *Repository) Update(ctx context.Context, id domain.ProjectID, name, description string) error {
	executor := r.getExecutor(ctx)

	const query = `
UPDATE projects
	SET name = $1, description = $2
WHERE id = $3`

	_, err := executor.Exec(ctx, query, name, description, id)
	if err != nil {
		return fmt.Errorf("failed to update project: %w", err)
	}

	return nil
}

func (r *Repository) Archive(ctx context.Context, id domain.ProjectID) error {
	executor := r.getExecutor(ctx)

	const query = `
UPDATE projects
	SET archived_at = NOW(), team_id = NULL
WHERE id = $1 AND archived_at IS NULL`

	result, err := executor.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to archive project: %w", err)
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		// Check if the project exists
		exists, err := r.projectExists(ctx, id)
		if err != nil {
			return fmt.Errorf("check if project exists: %w", err)
		}
		if !exists {
			return domain.ErrEntityNotFound
		}
		// Project exists but was already archived
		return nil
	}

	return nil
}

func (r *Repository) projectExists(ctx context.Context, id domain.ProjectID) (bool, error) {
	executor := r.getExecutor(ctx)

	const query = `SELECT 1 FROM projects WHERE id = $1 LIMIT 1`

	rows, err := executor.Query(ctx, query, id)
	if err != nil {
		return false, fmt.Errorf("query project existence: %w", err)
	}
	defer rows.Close()

	return rows.Next(), nil
}

func (r *Repository) GetProjectIDs(ctx context.Context) ([]domain.ProjectID, error) {
	executor := r.getExecutor(ctx)

	const query = `SELECT id FROM projects WHERE archived_at IS NULL`

	rows, err := executor.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("query project IDs: %w", err)
	}
	defer rows.Close()

	var projectIDs []domain.ProjectID
	for rows.Next() {
		var id uint
		if err := rows.Scan(&id); err != nil {
			return nil, fmt.Errorf("scan project ID: %w", err)
		}
		projectIDs = append(projectIDs, domain.ProjectID(id))
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate project IDs: %w", err)
	}

	return projectIDs, nil
}

func (r *Repository) Count(ctx context.Context) (uint, error) {
	executor := r.getExecutor(ctx)

	const query = "SELECT COUNT(*) FROM projects"
	var count uint
	err := executor.QueryRow(ctx, query).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("query projects count: %w", err)
	}

	return count, nil
}

//nolint:ireturn // it's ok here
func (r *Repository) getExecutor(ctx context.Context) db.Tx {
	if tx := db.TxFromContext(ctx); tx != nil {
		return tx
	}

	return r.db
}
