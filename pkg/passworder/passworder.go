// Package passworder provides functionality for creating password hashes and for verifying passwords
package passworder

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

func PasswordHash(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hash), nil
}

func MustPasswordHash(password string) string {
	hash, err := PasswordHash(password)
	if err != nil {
		panic(err)
	}

	return hash
}

func ValidatePassword(password, hashedPassword string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return false, nil
		}

		return false, err
	}

	return true, nil
}
