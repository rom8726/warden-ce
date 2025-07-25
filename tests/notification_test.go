//go:build integration

package tests

import (
	"testing"

	"github.com/rom8726/warden/tests/runner"
)

func TestNotificationSettingsAPI(t *testing.T) {
	cfg := runner.Config{
		CasesDir: "./cases/notification/settings",
	}
	runner.Run(t, &cfg)
}

func TestNotificationRulesAPI(t *testing.T) {
	cfg := runner.Config{
		CasesDir: "./cases/notification/rules",
	}
	runner.Run(t, &cfg)
}
