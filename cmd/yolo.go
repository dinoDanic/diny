package cmd

import (
	"github.com/dinoDanic/diny/tui/yolo"
	"github.com/spf13/cobra"
)

var yoloCmd = &cobra.Command{
	Use:   "yolo",
	Short: "Stage everything, generate a commit message, commit, and push",
	Long: `The yolo command runs the full lazy-commit flow in one shot:
stage all changes, generate an AI commit message, commit without hooks, and push.
Fully non-interactive.`,
	Run: func(cmd *cobra.Command, args []string) {
		runYolo()
	},
}

func init() {
	rootCmd.AddCommand(yoloCmd)
}

func runYolo() {
	yolo.Run(AppConfig, Version)
}
