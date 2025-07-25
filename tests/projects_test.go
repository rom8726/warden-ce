//go:build integration

package tests

import (
	"testing"

	"github.com/rom8726/warden/tests/runner"
)

func TestProjectsAPI(t *testing.T) {
	cfg := runner.Config{
		CasesDir: "./cases/projects",
	}
	runner.Run(t, &cfg)
}
