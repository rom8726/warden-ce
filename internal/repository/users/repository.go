package users

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

func (r *Repository) Create(ctx context.Context, userDTO domain.UserDTO) (domain.User, error) {
	executor := r.getExecutor(ctx)

	const query = `
INSERT INTO users (username, email, password_hash, is_superuser, is_active, created_at, is_tmp_password)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING id, username, email, password_hash, is_superuser,
    is_active, created_at, last_login, is_tmp_password`

	var user userModel
	err := executor.QueryRow(ctx, query,
		userDTO.Username,
		userDTO.Email,
		userDTO.PasswordHash,
		userDTO.IsSuperuser,
		true, // IsActive defaults to true
		time.Now(),
		userDTO.IsTmpPassword,
	).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.PasswordHash,
		&user.IsSuperuser,
		&user.IsActive,
		&user.CreatedAt,
		&user.LastLogin,
		&user.IsTmpPassword,
	)
	if err != nil {
		return domain.User{}, fmt.Errorf("insert user: %w", err)
	}

	return user.toDomain(), nil
}

func (r *Repository) GetByID(ctx context.Context, id domain.UserID) (domain.User, error) {
	executor := r.getExecutor(ctx)

	const query = `SELECT * FROM users WHERE id = $1 LIMIT 1`

	rows, err := executor.Query(ctx, query, id)
	if err != nil {
		return domain.User{}, fmt.Errorf("query user by ID: %w", err)
	}
	defer rows.Close()

	user, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[userModel])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.User{}, domain.ErrEntityNotFound
		}

		return domain.User{}, fmt.Errorf("collect user: %w", err)
	}

	return user.toDomain(), nil
}

func (r *Repository) ExistsByID(ctx context.Context, id domain.UserID) (bool, error) {
	_, err := r.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, domain.ErrEntityNotFound) {
			return false, nil
		}

		return false, err
	}

	return true, nil
}

func (r *Repository) GetByUsername(ctx context.Context, username string) (domain.User, error) {
	executor := r.getExecutor(ctx)

	const query = `SELECT * FROM users WHERE username = $1 LIMIT 1`

	rows, err := executor.Query(ctx, query, username)
	if err != nil {
		return domain.User{}, fmt.Errorf("query user by username: %w", err)
	}
	defer rows.Close()

	user, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[userModel])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.User{}, domain.ErrEntityNotFound
		}

		return domain.User{}, fmt.Errorf("collect user: %w", err)
	}

	return user.toDomain(), nil
}

func (r *Repository) GetByEmail(ctx context.Context, email string) (domain.User, error) {
	executor := r.getExecutor(ctx)

	const query = `SELECT * FROM users WHERE email = $1 LIMIT 1`

	rows, err := executor.Query(ctx, query, email)
	if err != nil {
		return domain.User{}, fmt.Errorf("query user by email: %w", err)
	}
	defer rows.Close()

	user, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[userModel])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.User{}, domain.ErrEntityNotFound
		}

		return domain.User{}, fmt.Errorf("collect user: %w", err)
	}

	return user.toDomain(), nil
}

func (r *Repository) FetchByIDs(ctx context.Context, ids []domain.UserID) ([]domain.User, error) {
	executor := r.getExecutor(ctx)

	const query = `SELECT * FROM users WHERE id = ANY($1)`
	rows, err := executor.Query(ctx, query, ids)
	if err != nil {
		return nil, fmt.Errorf("query users by IDs: %w", err)
	}
	defer rows.Close()

	listModels, err := pgx.CollectRows(rows, pgx.RowToStructByName[userModel])
	if err != nil {
		return nil, fmt.Errorf("collect users: %w", err)
	}

	users := make([]domain.User, 0, len(listModels))
	for i := range listModels {
		model := listModels[i]
		users = append(users, model.toDomain())
	}

	return users, nil
}

func (r *Repository) Update(ctx context.Context, user *domain.User) error {
	executor := r.getExecutor(ctx)

	const query = `
UPDATE users
SET username = $1, email = $2, password_hash = $3, is_superuser = $4, is_active = $5, last_login = $6,
    is_tmp_password = $7, updated_at = NOW()
WHERE id = $8`

	tag, err := executor.Exec(ctx, query,
		user.Username,
		user.Email,
		user.PasswordHash,
		user.IsSuperuser,
		user.IsActive,
		user.LastLogin,
		user.IsTmpPassword,
		user.ID,
	)
	if err != nil {
		return fmt.Errorf("update user: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return domain.ErrEntityNotFound
	}

	return nil
}

func (r *Repository) Delete(ctx context.Context, id domain.UserID) error {
	executor := r.getExecutor(ctx)

	const query = `
DELETE FROM users
WHERE id = $1`

	_, err := executor.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("delete user: %w", err)
	}

	return nil
}

func (r *Repository) List(ctx context.Context) ([]domain.User, error) {
	executor := r.getExecutor(ctx)

	const query = `SELECT * FROM users ORDER BY id`

	rows, err := executor.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("query users: %w", err)
	}
	defer rows.Close()

	listModels, err := pgx.CollectRows(rows, pgx.RowToStructByName[userModel])
	if err != nil {
		return nil, fmt.Errorf("collect users: %w", err)
	}

	users := make([]domain.User, 0, len(listModels))
	for i := range listModels {
		model := listModels[i]
		users = append(users, model.toDomain())
	}

	return users, nil
}

func (r *Repository) UpdateLastLogin(ctx context.Context, id domain.UserID) error {
	executor := r.getExecutor(ctx)
	const query = `UPDATE users SET last_login = NOW(), updated_at = NOW() WHERE id = $1;`

	_, err := executor.Exec(ctx, query, id)

	return err
}

func (r *Repository) UpdatePassword(ctx context.Context, id domain.UserID, passwordHash string) error {
	executor := r.getExecutor(ctx)
	const query = `UPDATE users SET password_hash = $1, is_tmp_password = false, updated_at = NOW() WHERE id = $2;`

	_, err := executor.Exec(ctx, query, passwordHash, id)

	return err
}

// Update2FA updates only 2FA-related fields for a user.
func (r *Repository) Update2FA(
	ctx context.Context,
	id domain.UserID,
	enabled bool,
	secret string,
	confirmedAt *time.Time,
) error {
	executor := r.getExecutor(ctx)

	const query = `
UPDATE users
SET two_fa_enabled = $1, two_fa_secret = $2, two_fa_confirmed_at = $3, updated_at = NOW()
WHERE id = $4`

	tag, err := executor.Exec(ctx, query,
		enabled,
		secret,
		confirmedAt,
		id,
	)
	if err != nil {
		return fmt.Errorf("update 2fa: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return domain.ErrEntityNotFound
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
