package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"github.com/spf13/cobra"

	issuenotificator "github.com/rom8726/warden/internal/issue-notificator"
	"github.com/rom8726/warden/internal/issue-notificator/config"
)

var NotificatorCmd = &cobra.Command{
	Use:   "notificator",
	Short: "Run notificator",
	RunE: func(cmd *cobra.Command, args []string) error {
		return runNotificatorCommand(cmd.Context(), args)
	},
}

var envFile string

func init() {
	NotificatorCmd.PersistentFlags().StringVarP(
		&envFile,
		"env-file",
		"e",
		"",
		"path to env file",
	)
}

func runNotificatorCommand(ctx context.Context, _ []string) error {
	cfg, err := config.New(envFile)
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}

	loggerHandler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: &cfg.Logger,
	})
	logger := slog.New(loggerHandler)
	slog.SetDefault(logger)

	app, err := issuenotificator.NewApp(ctx, cfg, logger)
	if err != nil {
		return fmt.Errorf("create app: %w", err)
	}
	defer app.Close()

	if err := app.Run(ctx); err != nil {
		return fmt.Errorf("run app: %w", err)
	}

	return nil
}
