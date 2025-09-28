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
	Short: "Visualize and analyze your git commit patterns with AI insights",
	Long: `Advanced timeline analysis of your commit history with AI-powered insights.

Analyze commits from various time periods and get detailed statistics:
- Today's commits with hourly breakdown
- Weekly/monthly productivity patterns  
- Specific date ranges with custom filters
- Repository-wide historical analysis
- Team collaboration patterns (when used with --team flag)

Features:
  • Smart pattern recognition for commit types
  • Conventional commit compliance analysis  
  • Code quality trends based on commit messages
  • Productivity insights and recommendations
  • Visual timeline with emoji indicators
  • Export reports in multiple formats (JSON, CSV, HTML)

Examples:
  diny timeline                         # Analyze today's commits
  diny timeline --week                  # Show this week's patterns
  diny timeline --since "2024-01-01"   # Custom date range
  diny timeline --author "john.doe"    # Filter by specific author
  diny timeline --export json          # Export detailed report`,
	Run: func(cmd *cobra.Command, args []string) {
		timeline.Main()
	},
}

func init() {
	rootCmd.AddCommand(timelineCmd)

	// Time period flags
	timelineCmd.Flags().BoolP("today", "t", false, "Show today's commits only")
	timelineCmd.Flags().BoolP("week", "w", false, "Show this week's commit timeline")
	timelineCmd.Flags().BoolP("month", "m", false, "Show this month's patterns")
	timelineCmd.Flags().StringP("since", "s", "", "Show commits since date (YYYY-MM-DD)")
	timelineCmd.Flags().StringP("until", "u", "", "Show commits until date (YYYY-MM-DD)")

	// Filtering flags
	timelineCmd.Flags().StringP("author", "a", "", "Filter commits by author email or name")
	timelineCmd.Flags().StringSliceP("type", "", []string{}, "Filter by commit types (feat, fix, docs, etc.)")
	timelineCmd.Flags().BoolP("team", "", false, "Include team collaboration analysis")

	// Output and display flags
	timelineCmd.Flags().StringP("export", "e", "", "Export format (json, csv, html)")
	timelineCmd.Flags().BoolP("detailed", "d", false, "Show detailed commit analysis")
	timelineCmd.Flags().BoolP("visual", "v", true, "Show visual timeline with emojis")
	timelineCmd.Flags().IntP("limit", "l", 50, "Maximum number of commits to analyze")

	// Analysis flags
	timelineCmd.Flags().BoolP("stats", "", true, "Include statistical analysis")
	timelineCmd.Flags().BoolP("trends", "", false, "Show productivity trends and insights")
}
