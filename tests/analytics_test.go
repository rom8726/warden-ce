//go:build integration

package tests

import (
	"testing"

	"github.com/rom8726/warden/tests/runner"
)

func TestAnalyticsAPI(t *testing.T) {
	cfg := runner.Config{
		CasesDir: "./cases/analytics",
	}
	runner.Run(t, &cfg)
}
