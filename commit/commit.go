package commit

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/dinoDanic/diny/config"
	"github.com/spf13/cobra"
)

func Main(cmd *cobra.Command, args []string) {
	fmt.Println()
	fmt.Println(RenderTitle("diny commiting"))

	gitDiffCmd := exec.Command("git", "diff", "--cached",
		"-U0", "--no-color", "--ignore-all-space", "--ignore-blank-lines",
		":(exclude)*.lock", ":(exclude)*package-lock.json", ":(exclude)*yarn.lock",
		":(exclude)node_modules/", ":(exclude)dist/", ":(exclude)build/")

	gitDiff, err := gitDiffCmd.Output()

	if err != nil {
		fmt.Println(RenderError(fmt.Sprintf("Failed to get git diff: %v", err)))
		os.Exit(1)
	}

	if len(gitDiff) == 0 {
		fmt.Println(RenderWarning("No staged changes found. Stage files first with `git add`."))
		os.Exit(0)
	}

	diff := string(gitDiff)

	userConfig, err := config.Load()

	var commitMessage, note string
	err = WithSpinner("Generating your commit message...", func() error {
		var genErr error
		commitMessage, note, genErr = CreateCommitMessage(diff, userConfig)
		return genErr
	})

	if err != nil {
		fmt.Println(RenderError(fmt.Sprintf("Error generating commit message: %v", err)))
		os.Exit(1)
	}

	HandleCommitFlow(commitMessage, note, diff, userConfig)
}
