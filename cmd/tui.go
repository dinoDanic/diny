package cmd

import (
	"github.com/dinoDanic/diny/internal/tui/launcher"
	"github.com/spf13/cobra"
)

var tuiCmd = &cobra.Command{
	Use:   "tui",
	Short: "Open the interactive TUI launcher",
	Run: func(cmd *cobra.Command, args []string) {
		launcher.Run(AppConfig)
	},
}

func init() {
	rootCmd.AddCommand(tuiCmd)
}
