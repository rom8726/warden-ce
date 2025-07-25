package users

import (
	"database/sql"
	"time"

	"github.com/rom8726/warden/internal/domain"
)

type userModel struct {
	ID               uint           `db:"id"`
	Username         string         `db:"username"`
	Email            string         `db:"email"`
	PasswordHash     string         `db:"password_hash"`
	IsSuperuser      bool           `db:"is_superuser"`
	IsActive         bool           `db:"is_active"`
	IsTmpPassword    bool           `db:"is_tmp_password"`
	TwoFAEnabled     bool           `db:"two_fa_enabled"`
	TwoFASecret      sql.NullString `db:"two_fa_secret"`
	TwoFAConfirmedAt *time.Time     `db:"two_fa_confirmed_at"`
	CreatedAt        time.Time      `db:"created_at"`
	UpdatedAt        time.Time      `db:"updated_at"`
	LastLogin        *time.Time     `db:"last_login"`
}

func (m *userModel) toDomain() domain.User {
	return domain.User{
		ID:               domain.UserID(m.ID),
		Username:         m.Username,
		Email:            m.Email,
		PasswordHash:     m.PasswordHash,
		IsSuperuser:      m.IsSuperuser,
		IsActive:         m.IsActive,
		IsTmpPassword:    m.IsTmpPassword,
		TwoFAEnabled:     m.TwoFAEnabled,
		TwoFASecret:      m.TwoFASecret.String,
		TwoFAConfirmedAt: m.TwoFAConfirmedAt,
		CreatedAt:        m.CreatedAt,
		UpdatedAt:        m.UpdatedAt,
		LastLogin:        m.LastLogin,
	}
}
