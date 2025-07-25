package domain

import (
	"errors"
)

var (
	ErrEntityNotFound        = errors.New("entity not found")
	ErrInvalidToken          = errors.New("invalid token")
	ErrUsernameAlreadyInUse  = errors.New("username already in use")
	ErrEmailAlreadyInUse     = errors.New("email already in use")
	ErrTeamNameAlreadyInUse  = errors.New("team name already in use")
	ErrInvalidPassword       = errors.New("invalid password")
	ErrInvalidCredentials    = errors.New("invalid credentials")
	ErrInactiveUser          = errors.New("inactive user")
	ErrForbidden             = errors.New("forbidden")
	ErrPermissionDenied      = errors.New("permission denied")
	ErrUserNotFound          = errors.New("user not found")
	ErrNoEnvelope            = errors.New("empty envelope")
	ErrInvalidEnvelopeHeader = errors.New("invalid envelope header")
	ErrInvalid2FACode        = errors.New("invalid 2FA code")
	ErrInvalidEmailCode      = errors.New("invalid email code")
	ErrTwoFARequired         = errors.New("2FA required")
	ErrTooMany2FAAttempts    = errors.New("too many 2FA attempts, try later")
	ErrLastOwner             = errors.New("cannot leave team as the last owner")
	ErrTeamHasProjects       = errors.New("team is attached to one or more projects")
)
