package commit

import (
	"fmt"
	"os"

	"github.com/dinoDanic/diny/config"
	"github.com/dinoDanic/diny/git"
	"github.com/dinoDanic/diny/ui"
	"github.com/spf13/cobra"
)

func Main(cmd *cobra.Command, args []string) {
	printMode, _ := cmd.Flags().GetBool("print")

	diff, userConfig := getCommitData(printMode)

	if printMode {
		commitMessage, err := CreateCommitMessage(diff, userConfig)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error generating commit message: %v\n", err)
			os.Exit(1)
		}
		fmt.Print(commitMessage)
		return
	}

	fmt.Println()

	var commitMessage string
	err := ui.WithSpinner("Generating your commit message...", func() error {
		var genErr error
		commitMessage, genErr = CreateCommitMessage(diff, userConfig)
		return genErr
	})

	if err != nil {
		ui.RenderError(fmt.Sprintf("%v", err))
		os.Exit(1)
	}

	HandleCommitFlow(commitMessage, diff, userConfig)
}

func getCommitData(isQuietMode bool) (string, *config.UserConfig) {
	gitDiff, err := git.GetGitDiff()

	fmt.Println()

	if err != nil {
		if isQuietMode {
			fmt.Fprintf(os.Stderr, "Failed to get git diff: %v\n", err)
		} else {
			ui.RenderError(fmt.Sprintf("Failed to get git diff: %v", err))
		}
		os.Exit(1)
	}

	if len(gitDiff) == 0 {
		if isQuietMode {
			fmt.Fprintf(os.Stderr, "No staged changes found. Stage files first with `git add`.\n")
		} else {
			ui.RenderWarning("No staged changes found. Stage files first with `git add`.")
		}
		os.Exit(0)
	}

	// Load config
	userConfig, err := config.Load()
	if err != nil {
		if isQuietMode {
			fmt.Fprintf(os.Stderr, "Failed to load config: %v\n", err)
		} else {
			ui.RenderError(fmt.Sprintf("Failed to load config: %v", err))
		}
		os.Exit(1)
	}

	return gitDiff, userConfig
}
