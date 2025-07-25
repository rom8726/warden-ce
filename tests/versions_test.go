package tests

import (
	"testing"

	"github.com/rom8726/warden/tests/runner"
)

func TestVersions(t *testing.T) {
	t.SkipNow()

	cfg := runner.Config{
		CasesDir: "./cases/versions",
		UsesOTP:  true,
	}
	runner.Run(t, &cfg)
}
