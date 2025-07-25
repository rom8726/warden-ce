package users

import (
	"context"
	"errors"

	"github.com/rom8726/warden/internal/domain"
)

// AuthProvider defines the interface for authentication providers.
type AuthProvider interface {
	// Authenticate authenticates a user with the given credentials
	// Returns the user if authentication is successful, or an error otherwise
	Authenticate(ctx context.Context, username, password string) (*domain.User, error)

	// CanHandle returns true if this provider can handle the given username
	CanHandle(username string) bool
}

// AuthProviderChain is a chain of authentication providers that will be tried in order.
type AuthProviderChain struct {
	providers []AuthProvider
}

// NewAuthProviderChain creates a new authentication provider chain.
func NewAuthProviderChain(providers ...AuthProvider) *AuthProviderChain {
	return &AuthProviderChain{
		providers: providers,
	}
}

// Authenticate until one succeeds or all fail.
func (c *AuthProviderChain) Authenticate(ctx context.Context, username, password string) (*domain.User, error) {
	var lastErr error

	for _, provider := range c.providers {
		if !provider.CanHandle(username) {
			continue
		}

		user, err := provider.Authenticate(ctx, username, password)
		if err == nil {
			return user, nil
		}

		// Only save the error if it's not a "not found" or "invalid credentials" error
		if !errors.Is(err, domain.ErrEntityNotFound) && !errors.Is(err, domain.ErrInvalidPassword) {
			lastErr = err
		}
	}

	if lastErr != nil {
		return nil, lastErr
	}

	return nil, domain.ErrInvalidCredentials
}

// CanHandle returns true if any provider in the chain can handle the username.
func (c *AuthProviderChain) CanHandle(username string) bool {
	for _, provider := range c.providers {
		if provider.CanHandle(username) {
			return true
		}
	}

	return false
}
