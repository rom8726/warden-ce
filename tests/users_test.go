//go:build integration

package tests

import (
	"testing"

	"github.com/rom8726/warden/tests/runner"
)

func TestUsersAPI(t *testing.T) {
	cfg := runner.Config{
		CasesDir: "./cases/users",
		UsesOTP:  true,
	}
	runner.Run(t, &cfg)
}
