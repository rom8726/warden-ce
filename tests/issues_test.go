//go:build integration

package tests

import (
	"testing"

	"github.com/rom8726/warden/tests/runner"
)

func TestIssuesAPI(t *testing.T) {
	cfg := runner.Config{
		CasesDir: "./cases/issues",
	}
	runner.Run(t, &cfg)
}
