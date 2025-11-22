/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/dinoDanic/diny/timeline"
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
		timeline.Main(AppConfig)
	},
}

func init() {
	rootCmd.AddCommand(timelineCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// timelineCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// timelineCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
