package usernotifications

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

func (r *Repository) Create(
	ctx context.Context,
	userID domain.UserID,
	notificationType domain.UserNotificationType,
	content json.RawMessage,
) (domain.UserNotification, error) {
	executor := r.getExecutor(ctx)

	const query = `
		INSERT INTO user_notifications (user_id, type, content, created_at, updated_at)
		VALUES ($1, $2, $3, NOW(), NOW())
		RETURNING *`

	rows, err := executor.Query(ctx, query, userID, notificationType, content)
	if err != nil {
		return domain.UserNotification{}, fmt.Errorf("insert user notification: %w", err)
	}
	defer rows.Close()

	model, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[userNotificationModel])
	if err != nil {
		return domain.UserNotification{}, fmt.Errorf("collect user notification: %w", err)
	}

	return model.toDomain(), nil
}

func (r *Repository) GetByID(ctx context.Context, id domain.UserNotificationID) (domain.UserNotification, error) {
	executor := r.getExecutor(ctx)

	const query = `SELECT * FROM user_notifications WHERE id = $1 LIMIT 1`

	rows, err := executor.Query(ctx, query, id)
	if err != nil {
		return domain.UserNotification{}, fmt.Errorf("query user notification by ID: %w", err)
	}
	defer rows.Close()

	model, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[userNotificationModel])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.UserNotification{}, domain.ErrEntityNotFound
		}

		return domain.UserNotification{}, fmt.Errorf("collect user notification: %w", err)
	}

	return model.toDomain(), nil
}

func (r *Repository) GetByUserID(
	ctx context.Context,
	userID domain.UserID,
	limit, offset uint,
) ([]domain.UserNotification, error) {
	executor := r.getExecutor(ctx)

	const query = `
		SELECT * FROM user_notifications
		WHERE user_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3`

	rows, err := executor.Query(ctx, query, userID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("query user notifications: %w", err)
	}
	defer rows.Close()

	models, err := pgx.CollectRows(rows, pgx.RowToStructByName[userNotificationModel])
	if err != nil {
		return nil, fmt.Errorf("collect user notifications: %w", err)
	}

	notifications := make([]domain.UserNotification, 0, len(models))
	for _, model := range models {
		notifications = append(notifications, model.toDomain())
	}

	return notifications, nil
}

func (r *Repository) GetUnreadCount(ctx context.Context, userID domain.UserID) (uint, error) {
	executor := r.getExecutor(ctx)

	const query = `
		SELECT COUNT(*) FROM user_notifications
		WHERE user_id = $1 AND is_read = false`

	var count uint
	err := executor.QueryRow(ctx, query, userID).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("count unread notifications: %w", err)
	}

	return count, nil
}

func (r *Repository) MarkAsRead(ctx context.Context, id domain.UserNotificationID) error {
	executor := r.getExecutor(ctx)

	const query = `
		UPDATE user_notifications
		SET is_read = true, updated_at = NOW()
		WHERE id = $1`

	_, err := executor.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("mark notification as read: %w", err)
	}

	return nil
}

func (r *Repository) MarkAllAsRead(ctx context.Context, userID domain.UserID) error {
	executor := r.getExecutor(ctx)

	const query = `
		UPDATE user_notifications
		SET is_read = true, updated_at = NOW()
		WHERE user_id = $1 AND is_read = false`

	_, err := executor.Exec(ctx, query, userID)
	if err != nil {
		return fmt.Errorf("mark all notifications as read: %w", err)
	}

	return nil
}

func (r *Repository) DeleteOld(
	ctx context.Context,
	maxAge time.Duration,
	limit uint,
) (uint, error) {
	executor := r.getExecutor(ctx)

	const query = `
		DELETE FROM user_notifications
		WHERE created_at < NOW() - INTERVAL '1 day' * $1
		LIMIT $2`

	result, err := executor.Exec(ctx, query, int(maxAge.Hours()/24), limit)
	if err != nil {
		return 0, fmt.Errorf("delete old notifications: %w", err)
	}

	return uint(result.RowsAffected()), nil //nolint:gosec // it's ok
}

func (r *Repository) GetPendingEmailNotifications(ctx context.Context, limit uint) ([]domain.UserNotification, error) {
	executor := r.getExecutor(ctx)

	const query = `
		SELECT * FROM user_notifications
		WHERE email_sent = false
		ORDER BY created_at ASC
		LIMIT $1`

	rows, err := executor.Query(ctx, query, limit)
	if err != nil {
		return nil, fmt.Errorf("query pending email notifications: %w", err)
	}
	defer rows.Close()

	models, err := pgx.CollectRows(rows, pgx.RowToStructByName[userNotificationModel])
	if err != nil {
		return nil, fmt.Errorf("collect pending email notifications: %w", err)
	}

	notifications := make([]domain.UserNotification, 0, len(models))
	for _, model := range models {
		notifications = append(notifications, model.toDomain())
	}

	return notifications, nil
}

func (r *Repository) MarkEmailAsSent(ctx context.Context, id domain.UserNotificationID) error {
	executor := r.getExecutor(ctx)

	const query = `
		UPDATE user_notifications
		SET email_sent = true, updated_at = NOW()
		WHERE id = $1`

	_, err := executor.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("mark email as sent: %w", err)
	}

	return nil
}

func (r *Repository) MarkEmailAsFailed(ctx context.Context, id domain.UserNotificationID, reason string) error {
	executor := r.getExecutor(ctx)

	const query = `
		UPDATE user_notifications
		SET email_sent = true, updated_at = NOW()
		WHERE id = $1`

	_, err := executor.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("mark email as failed: %w", err)
	}

	return nil
}

//nolint:ireturn // it's ok here
func (r *Repository) getExecutor(ctx context.Context) db.Tx {
	if tx := db.TxFromContext(ctx); tx != nil {
		return tx
	}

	return r.db
}
