package main

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/rom8726/warden/internal/installer"
)

var InstallCmd = &cobra.Command{
	Use:   "install",
	Short: "Run installer",
	RunE: func(cmd *cobra.Command, args []string) error {
		return runInstallCommand(cmd.Context(), args)
	},
}

func runInstallCommand(ctx context.Context, _ []string) error {
	app := installer.New()

	return app.Run(ctx)
}
