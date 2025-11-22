package commit

import (
	"fmt"
	"os"

	"github.com/dinoDanic/diny/config"
	"github.com/dinoDanic/diny/git"
	"github.com/dinoDanic/diny/ui"
	"github.com/spf13/cobra"
)

func Main(cmd *cobra.Command, args []string, cfg *config.Config) {
	printMode, _ := cmd.Flags().GetBool("print")
	noVerify, _ := cmd.Flags().GetBool("no-verify")

	diff := getGitDiff(printMode)

	if printMode {
		commitMessage, err := CreateCommitMessage(diff, cfg)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error generating commit message: %v\n", err)
			os.Exit(1)
		}
		fmt.Print(commitMessage)
		return
	}

	var commitMessage string
	err := ui.WithSpinner("Generating your commit message...", func() error {
		var genErr error
		commitMessage, genErr = CreateCommitMessage(diff, cfg)
		return genErr
	})

	if err != nil {
		ui.Box(ui.BoxOptions{Message: fmt.Sprintf("%v", err), Variant: ui.Error})
		os.Exit(1)
	}

	if noVerify {
		ui.Box(ui.BoxOptions{Title: "Commit message", Message: commitMessage})
		ExecuteCommit(commitMessage, false, true)
		return
	}

	HandleCommitFlow(commitMessage, diff, cfg)
}

func getGitDiff(isQuietMode bool) string {
	gitDiff, err := git.GetGitDiff()

	if err != nil {
		if isQuietMode {
			fmt.Fprintf(os.Stderr, "Failed to get git diff: %v\n", err)
		} else {
			ui.Box(ui.BoxOptions{Message: fmt.Sprintf("Failed to get git diff: %v", err), Variant: ui.Error})
		}
		os.Exit(1)
	}

	if len(gitDiff) == 0 {
		if isQuietMode {
			fmt.Fprintf(os.Stderr, "No staged changes found. Stage files first with `git add`.\n")
		} else {
			ui.Box(ui.BoxOptions{Message: "No staged changes found. Stage files first with `git add`.", Variant: ui.Warning})
		}
		os.Exit(0)
	}

	return gitDiff
}
