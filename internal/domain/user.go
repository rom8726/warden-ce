package domain

import (
	"time"
)

type UserID uint

// User represents a user in the system.
type User struct {
	ID               UserID
	Username         string
	Email            string
	PasswordHash     string
	IsSuperuser      bool
	IsActive         bool
	IsTmpPassword    bool
	TwoFAEnabled     bool
	TwoFASecret      string
	TwoFAConfirmedAt *time.Time
	CreatedAt        time.Time
	UpdatedAt        time.Time
	LastLogin        *time.Time
}

type UserDTO struct {
	Username      string
	Email         string
	PasswordHash  string
	IsSuperuser   bool
	IsTmpPassword bool
}
