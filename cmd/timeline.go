/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	tuitimeline "github.com/dinoDanic/diny/tui/timeline"
	"github.com/dinoDanic/diny/version"
	"github.com/spf13/cobra"
)

// timelineCmd represents the timeline command
var timelineCmd = &cobra.Command{
	Use:   "timeline",
	Short: "Analyze your commit message patterns and history",
	Long: `Analyze and display patterns from your commit history.

You can analyze commits from:
- Today's commits
- A specific date
- A date range

This will show you statistics about your commit message style,
including conventional commit usage, average length, and common patterns.`,
	Run: func(cmd *cobra.Command, args []string) {
		tuitimeline.Run(AppConfig, version.Get())
	},
}

func init() {
	rootCmd.AddCommand(timelineCmd)
}
