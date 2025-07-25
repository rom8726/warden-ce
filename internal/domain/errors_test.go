package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDomainErrors(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		expected string
	}{
		{
			name:     "ErrEntityNotFound",
			err:      ErrEntityNotFound,
			expected: "entity not found",
		},
		{
			name:     "ErrInvalidToken",
			err:      ErrInvalidToken,
			expected: "invalid token",
		},
		{
			name:     "ErrUsernameAlreadyInUse",
			err:      ErrUsernameAlreadyInUse,
			expected: "username already in use",
		},
		{
			name:     "ErrEmailAlreadyInUse",
			err:      ErrEmailAlreadyInUse,
			expected: "email already in use",
		},
		{
			name:     "ErrTeamNameAlreadyInUse",
			err:      ErrTeamNameAlreadyInUse,
			expected: "team name already in use",
		},
		{
			name:     "ErrInvalidPassword",
			err:      ErrInvalidPassword,
			expected: "invalid password",
		},
		{
			name:     "ErrInvalidCredentials",
			err:      ErrInvalidCredentials,
			expected: "invalid credentials",
		},
		{
			name:     "ErrInactiveUser",
			err:      ErrInactiveUser,
			expected: "inactive user",
		},
		{
			name:     "ErrForbidden",
			err:      ErrForbidden,
			expected: "forbidden",
		},
		{
			name:     "ErrPermissionDenied",
			err:      ErrPermissionDenied,
			expected: "permission denied",
		},
		{
			name:     "ErrUserNotFound",
			err:      ErrUserNotFound,
			expected: "user not found",
		},
		{
			name:     "ErrNoEnvelope",
			err:      ErrNoEnvelope,
			expected: "empty envelope",
		},
		{
			name:     "ErrInvalidEnvelopeHeader",
			err:      ErrInvalidEnvelopeHeader,
			expected: "invalid envelope header",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.err.Error())
		})
	}
}
