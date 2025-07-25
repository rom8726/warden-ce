package releasestats

import (
	"context"
	"encoding/json"
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

func (r *Repository) GetByProjectAndRelease(
	ctx context.Context,
	projectID domain.ProjectID,
	release string,
) (domain.ReleaseStats, error) {
	executor := r.getExecutor(ctx)

	const query = `SELECT * FROM release_stats WHERE project_id = $1 AND release = $2 LIMIT 1`

	rows, err := executor.Query(ctx, query, projectID, release)
	if err != nil {
		return domain.ReleaseStats{}, fmt.Errorf("query release_stats by project and release: %w", err)
	}
	defer rows.Close()

	stats, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[releaseStatsModel])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.ReleaseStats{}, domain.ErrEntityNotFound
		}

		return domain.ReleaseStats{}, fmt.Errorf("collect release_stats: %w", err)
	}

	return stats.toDomain()
}

func (r *Repository) ListByProject(ctx context.Context, projectID domain.ProjectID) ([]domain.ReleaseStats, error) {
	executor := r.getExecutor(ctx)

	const query = `SELECT * FROM release_stats WHERE project_id = $1 ORDER BY generated_at DESC, release_id DESC`

	rows, err := executor.Query(ctx, query, projectID)
	if err != nil {
		return nil, fmt.Errorf("query release_stats by project: %w", err)
	}
	defer rows.Close()

	entities, err := pgx.CollectRows(rows, pgx.RowToStructByName[releaseStatsModel])
	if err != nil {
		return nil, fmt.Errorf("collect release_stats: %w", err)
	}

	stats := make([]domain.ReleaseStats, 0, len(entities))
	for i := range entities {
		stat, err := entities[i].toDomain()
		if err != nil {
			return nil, fmt.Errorf("convert release_stats: %w", err)
		}

		stats = append(stats, stat)
	}

	return stats, nil
}

func (r *Repository) Create(ctx context.Context, stats domain.ReleaseStats) error {
	executor := r.getExecutor(ctx)

	const query = `
INSERT INTO release_stats (
    project_id, release_id, release, generated_at,
    known_issues_total, new_issues_total, regressions_total,
    resolved_in_version_total, fixed_new_in_version_total, fixed_old_in_version_total,
    avg_fix_time, median_fix_time, p95_fix_time,
    severity_distribution, users_affected
) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15)`

	severityJSON, err := json.Marshal(stats.SeverityDistribution)
	if err != nil {
		return fmt.Errorf("marshal severity distribution: %w", err)
	}

	var avgFixTime, medianFixTime, p95FixTime time.Duration
	if stats.AvgFixTime != nil {
		avgFixTime = *stats.AvgFixTime
	}
	if stats.MedianFixTime != nil {
		medianFixTime = *stats.MedianFixTime
	}
	if stats.P95FixTime != nil {
		p95FixTime = *stats.P95FixTime
	}

	_, err = executor.Exec(ctx, query,
		stats.ProjectID,
		stats.ReleaseID,
		stats.Release,
		stats.GeneratedAt,
		stats.KnownIssuesTotal,
		stats.NewIssuesTotal,
		stats.RegressionsTotal,
		stats.ResolvedInVersionTotal,
		stats.FixedNewInVersionTotal,
		stats.FixedOldInVersionTotal,
		avgFixTime.String(),
		medianFixTime.String(),
		p95FixTime.String(),
		severityJSON,
		stats.UsersAffected,
	)
	if err != nil {
		return fmt.Errorf("insert release_stats: %w", err)
	}

	return nil
}

func (r *Repository) getExecutor(ctx context.Context) db.Tx {
	tx := db.TxFromContext(ctx)
	if tx != nil {
		return tx
	}

	return r.db
}
