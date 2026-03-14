package cmd

import (
	tuichangelog "github.com/dinoDanic/diny/tui/changelog"
	"github.com/spf13/cobra"
)

var changelogCmd = &cobra.Command{
	Use:   "changelog",
	Short: "Generate an AI-powered changelog for your repository",
	Run: func(cmd *cobra.Command, args []string) {
		tuichangelog.Run(AppConfig, Version)
	},
}

func init() {
	rootCmd.AddCommand(changelogCmd)
}
