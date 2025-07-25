package notificationsqueue

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

func (r *Repository) AddNotification(
	ctx context.Context,
	projectID domain.ProjectID,
	issueID domain.IssueID,
	level domain.IssueLevel,
	isNew, wasReactivated bool,
) error {
	executor := r.getExecutor(ctx)
	const query = `
INSERT INTO notifications_queue (project_id, issue_id, is_new, was_reactivated, status, level, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, $6, NOW(), NOW())`
	_, err := executor.Exec(
		ctx,
		query,
		projectID,
		issueID,
		isNew,
		wasReactivated,
		domain.NotificationStatusPending,
		level,
	)
	if err != nil {
		return fmt.Errorf("insert notification: %w", err)
	}

	return err
}

func (r *Repository) GetByID(ctx context.Context, id domain.NotificationID) (domain.Notification, error) {
	executor := r.getExecutor(ctx)
	const query = `
SELECT * FROM notifications_queue
WHERE id = $1
LIMIT 1`
	rows, err := executor.Query(ctx, query, id)
	if err != nil {
		return domain.Notification{}, fmt.Errorf("query: %w", err)
	}
	defer rows.Close()

	model, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[notificationModel])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Notification{}, domain.ErrEntityNotFound
		}

		return domain.Notification{}, fmt.Errorf("collect: %w", err)
	}

	return model.toDomain(), nil
}

func (r *Repository) TakePending(ctx context.Context, limit uint) ([]domain.Notification, error) {
	executor := r.getExecutor(ctx)
	const query = `
SELECT * FROM notifications_queue
WHERE status = $1
ORDER BY created_at ASC
LIMIT $2`
	rows, err := executor.Query(ctx, query, domain.NotificationStatusPending, limit)
	if err != nil {
		return nil, fmt.Errorf("query: %w", err)
	}
	defer rows.Close()

	list, err := pgx.CollectRows(rows, pgx.RowToStructByName[notificationModel])
	if err != nil {
		return nil, fmt.Errorf("collect: %w", err)
	}

	result := make([]domain.Notification, 0, len(list))
	for _, model := range list {
		result = append(result, model.toDomain())
	}

	return result, nil
}

func (r *Repository) MarkAsSent(ctx context.Context, id domain.NotificationID) error {
	executor := r.getExecutor(ctx)
	const query = "UPDATE notifications_queue SET status = 'sent', updated_at = NOW() WHERE id = $1"

	_, err := executor.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("exec update: %w", err)
	}

	return nil
}

func (r *Repository) MarkAsFailed(ctx context.Context, id domain.NotificationID, reason string) error {
	executor := r.getExecutor(ctx)
	const query = "UPDATE notifications_queue SET status = 'failed', fail_reason = $1, updated_at = NOW() WHERE id = $2"

	_, err := executor.Exec(ctx, query, reason, id)
	if err != nil {
		return fmt.Errorf("exec update: %w", err)
	}

	return nil
}

func (r *Repository) MarkAsSkipped(ctx context.Context, id domain.NotificationID, reason string) error {
	executor := r.getExecutor(ctx)
	const query = "UPDATE notifications_queue SET status = 'skipped', fail_reason = $1, updated_at = NOW() WHERE id = $2"

	_, err := executor.Exec(ctx, query, reason, id)
	if err != nil {
		return fmt.Errorf("exec update: %w", err)
	}

	return nil
}

func (r *Repository) DeleteOld(ctx context.Context, maxAge time.Duration, limit uint) (uint, error) {
	executor := r.getExecutor(ctx)
	const query = `
DELETE FROM notifications_queue
WHERE id IN (
    SELECT id
    FROM notifications_queue
    WHERE status != 'pending' AND updated_at < (NOW() - $1::interval)
    LIMIT $2
)`

	tag, err := executor.Exec(ctx, query, maxAge, limit)
	if err != nil {
		return 0, fmt.Errorf("exec update: %w", err)
	}

	return uint(tag.RowsAffected()), nil //nolint:gosec // it's ok
}

//nolint:ireturn // it's ok here
func (r *Repository) getExecutor(ctx context.Context) db.Tx {
	if tx := db.TxFromContext(ctx); tx != nil {
		return tx
	}

	return r.db
}
