package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"github.com/spf13/cobra"

	"github.com/rom8726/warden/internal/ingest-server"
	"github.com/rom8726/warden/internal/ingest-server/config"
)

var ServerCmd = &cobra.Command{
	Use:   "server",
	Short: "Run server",
	RunE: func(cmd *cobra.Command, args []string) error {
		return runServerCommand(cmd.Context(), args)
	},
}

var envFile string

func init() {
	ServerCmd.PersistentFlags().StringVarP(
		&envFile,
		"env-file",
		"e",
		"",
		"path to env file",
	)
}

func runServerCommand(ctx context.Context, _ []string) error {
	cfg, err := config.New(envFile)
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}

	loggerHandler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: &cfg.Logger,
	})
	logger := slog.New(loggerHandler)
	slog.SetDefault(logger)

	app, err := ingestserver.NewApp(ctx, cfg, logger)
	if err != nil {
		return fmt.Errorf("create app: %w", err)
	}
	defer app.Close()

	if err := app.Run(ctx); err != nil {
		return fmt.Errorf("run app: %w", err)
	}

	return nil
}
