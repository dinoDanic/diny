/* Copyright © 2025 dinoDanic dino.danic@gmail.com */
package commit

import (
	"fmt"
	"os"

	"github.com/dinoDanic/diny/cli/config"
	"github.com/dinoDanic/diny/cli/git"
	"github.com/dinoDanic/diny/cli/ui"
	"github.com/spf13/cobra"
)

func Main(cmd *cobra.Command, args []string) {
	printMode, _ := cmd.Flags().GetBool("print")

	diff, commitConfig := getCommitData(printMode)

	if printMode {
		commitMessage, err := CreateCommitMessage(diff, commitConfig)
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
		commitMessage, genErr = CreateCommitMessage(diff, commitConfig)
		return genErr
	})

	if err != nil {
		ui.ErrorBox("Error", fmt.Sprintf("%v", err))
		os.Exit(1)
	}

	HandleCommitFlow(commitMessage, diff, commitConfig)
}

func getCommitData(isQuietMode bool) (string, *config.CommitConfig) {
	gitDiff, err := git.GetGitDiff()

	if err != nil {
		if isQuietMode {
			fmt.Fprintf(os.Stderr, "Failed to get git diff: %v\n", err)
		} else {
			ui.ErrorBox("Failed to get git diff", fmt.Sprintf("%v", err))
		}
		os.Exit(1)
	}

	if len(gitDiff) == 0 {
		if isQuietMode {
			fmt.Fprintf(os.Stderr, "No staged changes found. Stage files first with `git add`.\n")
		} else {
			ui.WarningMsg("No staged changes found. Stage files first with `git add`.")
		}
		os.Exit(0)
	}

	cfg := config.Get()
	commitConfig := &cfg.Commit

	return gitDiff, commitConfig
}
