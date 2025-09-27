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
	fmt.Println(RenderTitle("Diny - AI Commit Message Generator"))
	fmt.Println()

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
	if err == nil && userConfig != nil {
		fmt.Println(RenderConfigBox(formatConfig(*userConfig)))
	}

	fmt.Println()
	fmt.Println(RenderStep("Analyzing staged changes..."))
	fmt.Println()

	commitMessage, note, err := CreateCommitMessage(diff, userConfig)

	if err != nil {
		fmt.Println(RenderError(fmt.Sprintf("Error generating commit message: %v", err)))
		os.Exit(1)
	}

	HandleCommitFlow(commitMessage, note, diff, userConfig)
}

func formatConfig(userConfig config.UserConfig) string {
	result := fmt.Sprintf("Emoji: %t\n", userConfig.UseEmoji)
	result += fmt.Sprintf("Conventional: %t\n", userConfig.UseConventional)
	result += fmt.Sprintf("Tone: %s\n", userConfig.Tone)
	result += fmt.Sprintf("Length: %s", userConfig.Length)
	return result
}
