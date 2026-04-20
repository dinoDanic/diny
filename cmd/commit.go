package cmd

import (
	"github.com/dinoDanic/diny/prompts"
	"github.com/dinoDanic/diny/tui/app"
	"github.com/dinoDanic/diny/update"
	"github.com/dinoDanic/diny/version"
	"github.com/spf13/cobra"
)

var commitCmd = &cobra.Command{
	Use:   "commit",
	Short: "Generate commit messages from staged changes",
	Long: `generate and apply a commit message from staged changes.

Diny reads your staged changes and propose a commit message, and lets
you commit, edit, regenerate, or refine it—all`,
	Run: func(cmd *cobra.Command, args []string) {
		checker := update.NewUpdateChecker(version.Get())
		updateCh := checker.CheckAsync()

		noVerify, _ := cmd.Flags().GetBool("no-verify")
		push, _ := cmd.Flags().GetBool("push")
		print, _ := cmd.Flags().GetBool("print")

		result := app.Run(AppConfig, version.Get(), app.Options{
			NoVerify: noVerify,
			Push:     push,
			Print:    print,
		})

		if result.CommitSucceeded {
			if !checker.PromptIfAvailable(updateCh) {
				prompts.MaybeShow(AppConfig)
			}
		}
	},
}

func init() {
	commitCmd.Flags().Bool("no-verify", false, "Skip pre-commit and commit-msg hooks on every commit")
	commitCmd.Flags().Bool("push", false, "Push after committing (after the final commit when splitting)")
	commitCmd.Flags().Bool("print", false, "Print the generated message to stdout (incompatible with split)")
	rootCmd.AddCommand(commitCmd)
}
