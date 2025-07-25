//go:build integration

package tests

import (
	"testing"

	"github.com/rom8726/warden/tests/runner"
)

func TestAuthAPI(t *testing.T) {
	cfg := runner.Config{
		CasesDir: "./cases/auth",
		UsesOTP:  true,
	}
	runner.Run(t, &cfg)
}
