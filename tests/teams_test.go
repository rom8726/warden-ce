//go:build integration

package tests

import (
	"testing"

	"github.com/rom8726/warden/tests/runner"
)

func TestTeamsAPI(t *testing.T) {
	cfg := runner.Config{
		CasesDir: "./cases/teams",
	}
	runner.Run(t, &cfg)
}
