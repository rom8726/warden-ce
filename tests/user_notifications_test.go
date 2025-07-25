//go:build integration

package tests

import (
	"testing"

	"github.com/rom8726/warden/tests/runner"
)

func TestUserNotifications(t *testing.T) {
	cfg := runner.Config{
		CasesDir: "./cases/user_notifications",
	}
	runner.Run(t, &cfg)
}
