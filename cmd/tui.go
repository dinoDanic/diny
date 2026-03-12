package cmd

import (
	"github.com/dinoDanic/diny/internal/tui/app"
	"github.com/spf13/cobra"
)

var tuiCmd = &cobra.Command{
	Use:   "tui",
	Short: "Interactive TUI for generating and committing messages",
	Long: `Launch an interactive terminal UI that auto-detects staged changes,
generates a commit message, and lets you commit, edit, regenerate,
or refine — all with single-key shortcuts.`,
	Run: func(cmd *cobra.Command, args []string) {
		app.Run(AppConfig, Version)
	},
}

func init() {
	rootCmd.AddCommand(tuiCmd)
}
