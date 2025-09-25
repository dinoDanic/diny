package commit

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/charmbracelet/huh"
	"github.com/dinoDanic/diny/config"
	"github.com/dinoDanic/diny/helpers"
	"github.com/spf13/cobra"
)

func Main(cmd *cobra.Command, args []string) {
	gitDiffCmd := exec.Command("git", "diff", "--cached",
		"-U0", "--no-color", "--ignore-all-space", "--ignore-blank-lines",
		":(exclude)*.lock", ":(exclude)*package-lock.json", ":(exclude)*yarn.lock",
		":(exclude)node_modules/", ":(exclude)dist/", ":(exclude)build/")

	gitDiff, err := gitDiffCmd.Output()

	if err != nil {
		fmt.Printf("❌ Failed to get git diff: %v\n", err)
		os.Exit(1)
	}

	if len(gitDiff) == 0 {
		fmt.Println("🦕 No staged changes found. Stage files first with `git add`.")
		os.Exit(0)
	}

	diff := string(gitDiff)

	userConfig := config.Load()

	systemPrompt := helpers.BuildSystemPrompt(userConfig)
	fullPrompt := diff + systemPrompt

	fmt.Println()
	config.PrintConfiguration(userConfig)
	fmt.Println()
	fmt.Println("🦕 Generating commit message...")

	commitMessage, err := CreateCommitMessage(fullPrompt, userConfig)

	if err != nil {
		fmt.Printf("💥 Error generating commit message: %v\n", err)
		os.Exit(1)
	}

	fmt.Println()
	fmt.Print(commitMessage)
	fmt.Println()
	fmt.Println()

	confirmed := confirmPrompt("🦕 Do you want to commit with this message?")
	fmt.Println()

	if confirmed {
		fmt.Println("🦕 Creating commit...")
		commitCmd := exec.Command("git", "commit", "--no-verify", "-m", commitMessage)
		err := commitCmd.Run()
		if err != nil {
			fmt.Printf("❌ Commit failed: %v\n", err)
			os.Exit(1)
		}
		fmt.Println()
		fmt.Println("🦕 Success! Commit added to history.")
	} else {
		fmt.Println()
		fmt.Println("🦕 Commit cancelled.")
	}
}

func confirmPrompt(message string) bool {
	var confirmed bool

	err := huh.NewConfirm().
		Title(message).
		Affirmative("Yes").
		Negative("No").
		Value(&confirmed).
		Run()

	if err != nil {
		fmt.Printf("Error running prompt: %v\n", err)
		os.Exit(1)
	}

	return confirmed
}
