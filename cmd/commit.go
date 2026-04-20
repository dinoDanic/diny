package cmd

import (
	"github.com/dinoDanic/diny/prompts"
	"github.com/dinoDanic/diny/tui/app"
	"github.com/dinoDanic/diny/update"
	"github.com/spf13/cobra"
)

var commitCmd = &cobra.Command{
	Use:   "commit",
	Short: "Generate commit messages from staged changes",
	Long: `generate and apply a commit message from staged changes.

Diny reads your staged changes and propose a commit message, and lets
you commit, edit, regenerate, or refine it—all`,
	Run: func(cmd *cobra.Command, args []string) {
		checker := update.NewUpdateChecker(Version)
		updateCh := checker.CheckAsync()

		result := app.Run(AppConfig, Version)

		if result.CommitSucceeded {
			prompts.MaybeShow(AppConfig)
		}

		checker.PromptIfAvailable(updateCh)
	},
}

func init() {
	rootCmd.AddCommand(commitCmd)
}
