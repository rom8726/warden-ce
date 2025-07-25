//go:build integration

package tests

import (
	"testing"

	"github.com/rom8726/warden/internal/backend"
	"github.com/rom8726/warden/tests/runner"
)

func TestSDKEnvelopeAPI(t *testing.T) {
	t.SkipNow()
	t.SkipNow()
	cfg := runner.Config{
		CasesDir: "./cases/sentry/envelope",
		AfterReq: func(app *backend.App) error {
			return nil
		},
	}
	runner.Run(t, &cfg)
}

func TestSDKStoreAPI(t *testing.T) {
	t.SkipNow()
	t.SkipNow()
	cfg := runner.Config{
		CasesDir: "./cases/sentry/store",
		AfterReq: func(app *backend.App) error {
			return nil
		},
	}
	runner.Run(t, &cfg)
}
