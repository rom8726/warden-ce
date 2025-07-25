package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"github.com/spf13/cobra"

	envelopeconsumer "github.com/rom8726/warden/internal/envelope-consumer"
	"github.com/rom8726/warden/internal/envelope-consumer/config"
)

var ConsumerCmd = &cobra.Command{
	Use:   "consumer",
	Short: "Run consumer",
	RunE: func(cmd *cobra.Command, args []string) error {
		return runConsumerCommand(cmd.Context(), args)
	},
}

var envFile string

func init() {
	ConsumerCmd.PersistentFlags().StringVarP(
		&envFile,
		"env-file",
		"e",
		"",
		"path to env file",
	)
}

func runConsumerCommand(ctx context.Context, _ []string) error {
	cfg, err := config.New(envFile)
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}

	loggerHandler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: &cfg.Logger,
	})
	logger := slog.New(loggerHandler)
	slog.SetDefault(logger)

	app, err := envelopeconsumer.NewApp(ctx, cfg, logger)
	if err != nil {
		return fmt.Errorf("create app: %w", err)
	}
	defer app.Close()

	if err := app.Run(ctx); err != nil {
		return fmt.Errorf("run app: %w", err)
	}

	return nil
}
