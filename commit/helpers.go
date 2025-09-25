package commit

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/charmbracelet/huh"
	"github.com/dinoDanic/diny/config"
)

func HandleCommitFlow(commitMessage, fullPrompt string, userConfig config.UserConfig) {
	HandleCommitFlowWithHistory(commitMessage, fullPrompt, userConfig, []string{})
}

func HandleCommitFlowWithHistory(commitMessage, fullPrompt string, userConfig config.UserConfig, previousMessages []string) {
	fmt.Println()
	fmt.Print(commitMessage)
	fmt.Println()
	fmt.Println()

	choice := choicePrompt("🦕 Choose your next task:")
	fmt.Println()

	switch choice {
	case "commit":
		fmt.Println("🦕 Creating commit...")
		commitCmd := exec.Command("git", "commit", "--no-verify", "-m", commitMessage)
		err := commitCmd.Run()
		if err != nil {
			fmt.Printf("❌ Commit failed: %v\n", err)
			os.Exit(1)
		}
		fmt.Println()
		fmt.Println("🦕 Success! Commit added to history.")
	case "regenerate":
		fmt.Println("🦕 Generating new commit message...")

		modifiedPrompt := fullPrompt
		if len(previousMessages) > 0 {
			modifiedPrompt += "\n\nPrevious commit messages that were not satisfactory:\n"
			for i, msg := range previousMessages {
				modifiedPrompt += fmt.Sprintf("%d. %s\n", i+1, msg)
			}
			modifiedPrompt += "\nPlease generate a different commit message that avoids the style and approach of the previous ones."
		} else {
			modifiedPrompt += "\n\nPlease provide an alternative commit message with a different approach or focus."
		}

		newCommitMessage, err := CreateCommitMessage(modifiedPrompt, userConfig)
		if err != nil {
			fmt.Printf("💥 Error generating commit message: %v\n", err)
			os.Exit(1)
		}

		updatedHistory := append(previousMessages, commitMessage)
		HandleCommitFlowWithHistory(newCommitMessage, fullPrompt, userConfig, updatedHistory)
	case "exit":
		fmt.Println()
		fmt.Println("🦕 Goodbye!")
		os.Exit(0)
	}
}

func choicePrompt(message string) string {
	var choice string

	err := huh.NewSelect[string]().
		Title(message).
		Options(
			huh.NewOption("Commit", "commit"),
			huh.NewOption("Generate different message", "regenerate"),
			huh.NewOption("Exit", "exit"),
		).
		Value(&choice).
		Run()

	if err != nil {
		fmt.Printf("Error running prompt: %v\n", err)
		os.Exit(1)
	}

	return choice
}
