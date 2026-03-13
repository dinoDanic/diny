package cmd

import (
	"os"

	"github.com/dinoDanic/diny/commit"
	"github.com/dinoDanic/diny/git"
	"github.com/dinoDanic/diny/ui"
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
	err := ui.WithSpinner("Staging all changes...", func() error {
		return git.AddAll()
	})
	if err != nil {
		ui.Error("Failed to stage changes: %v", err)
		os.Exit(1)
	}

	diff, err := git.GetGitDiff()
	if err != nil {
		ui.Error("Failed to get git diff: %v", err)
		os.Exit(1)
	}
	if len(diff) == 0 {
		ui.Warning("No changes to commit.")
		os.Exit(0)
	}

	var commitMessage string
	err = ui.WithSpinner("Generating commit message...", func() error {
		var genErr error
		commitMessage, genErr = commit.CreateCommitMessage(diff, AppConfig)
		return genErr
	})
	if err != nil {
		ui.Error("Failed to generate commit message: %v", err)
		os.Exit(1)
	}

	ui.Box("Commit message", commitMessage)

	hash, err := commit.TryCommit(commitMessage, true, true, AppConfig)
	if err != nil {
		ui.Error("%v", err)
		os.Exit(1)
	}

	ui.Success("Committed!")
	if hash != "" {
		ui.Success("Hash: %s (copied)", hash)
	}
	ui.Success("Pushed!")
}
