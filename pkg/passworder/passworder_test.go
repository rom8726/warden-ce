package passworder

import (
	"testing"

	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

func TestPasswordHash(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		password, err := PasswordHash("some_password")
		require.NoError(t, err)
		require.NotEmpty(t, password)
	})

	t.Run("error (too long password)", func(t *testing.T) {
		password, err := PasswordHash(string(make([]byte, 73)))
		require.Error(t, err)
		require.ErrorIs(t, err, bcrypt.ErrPasswordTooLong)
		require.Empty(t, password)
	})
}

func TestMustPasswordHash(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		defer func() {
			r := recover()
			require.Nil(t, r)
		}()
		hash := MustPasswordHash("some_password")
		require.NotEmpty(t, hash)
	})

	t.Run("panic", func(t *testing.T) {
		defer func() {
			r := recover()
			require.NotNil(t, r)
		}()
		hash := MustPasswordHash(string(make([]byte, 73)))
		require.NotEmpty(t, hash)
	})
}

func TestValidatePassword(t *testing.T) {
	t.Run("valid", func(t *testing.T) {
		isValid, err := ValidatePassword("some_password", "$2a$10$U7Mv.l4/.wq5ZCrDmGMW2.rYNVdfe75qFgxvay4i8gVjpSYvZ3Q72")
		require.NoError(t, err)
		require.True(t, isValid)
	})

	t.Run("not valid", func(t *testing.T) {
		isValid, err := ValidatePassword("some_password", "$2a$10$Eaq8MIH12RM3Ro.2fRD8F.qKeX9ualdY5ccqAfmi4byEj7vNE5ofW")
		require.NoError(t, err)
		require.False(t, isValid)
	})

	t.Run("error", func(t *testing.T) {
		isValid, err := ValidatePassword("some_password", "-")
		require.Error(t, err)
		require.False(t, isValid)
	})
}
