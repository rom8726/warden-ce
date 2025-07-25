package users

import (
	"context"
	"errors"
	"fmt"

	"github.com/rom8726/warden/internal/backend/contract"
	"github.com/rom8726/warden/internal/domain"
	"github.com/rom8726/warden/pkg/passworder"
)

// LocalAuthProvider This is the default authentication provider that checks against the local database.
type LocalAuthProvider struct {
	repo contract.UsersRepository
}

// NewLocalAuthProvider creates a new local authentication provider.
func NewLocalAuthProvider(repo contract.UsersRepository) *LocalAuthProvider {
	return &LocalAuthProvider{
		repo: repo,
	}
}

// Authenticate authenticates a user with the given credentials.
//
//nolint:nestif // need refactoring
func (p *LocalAuthProvider) Authenticate(ctx context.Context, username, password string) (*domain.User, error) {
	// First try to find by username
	user, err := p.repo.GetByUsername(ctx, username)
	if err != nil {
		// If not found by username, try by email
		if errors.Is(err, domain.ErrEntityNotFound) {
			user, err = p.repo.GetByEmail(ctx, username)
			if err != nil {
				return nil, fmt.Errorf("get user by email: %w", err)
			}
		} else {
			return nil, fmt.Errorf("get user by username: %w", err)
		}
	}

	// Check if the user is active
	if !user.IsActive {
		return nil, domain.ErrInactiveUser
	}

	// Check if the password is valid
	isValid, err := passworder.ValidatePassword(password, user.PasswordHash)
	if err != nil {
		return nil, fmt.Errorf("validate password: %w", err)
	}
	if !isValid {
		return nil, domain.ErrInvalidPassword
	}

	return &user, nil
}

// CanHandle always returns true as this is the default provider.
func (p *LocalAuthProvider) CanHandle(_ string) bool {
	return true
}
