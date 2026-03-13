package cmd

import (
	"github.com/dinoDanic/diny/internal/tui/app"
	"github.com/spf13/cobra"
)

var commitCmd = &cobra.Command{
	Use:   "commit",
	Short: "Generate commit messages from staged changes",
	Long: `generate and apply a commit message from staged changes.

Diny reads your staged changes and propose a commit message, and lets
you commit, edit, regenerate, or refine it—all`,
	Run: func(cmd *cobra.Command, args []string) {
		app.Run(AppConfig, Version)
	},
}

func init() {
	rootCmd.AddCommand(commitCmd)
}
