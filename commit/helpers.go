package commit

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/charmbracelet/huh"
	"github.com/dinoDanic/diny/config"
	"github.com/dinoDanic/diny/ui"
)

func HandleCommitFlow(commitMessage, fullPrompt string, userConfig *config.UserConfig) {
	HandleCommitFlowWithHistory(commitMessage, fullPrompt, userConfig, []string{})
}

func HandleCommitFlowWithHistory(commitMessage, fullPrompt string, userConfig *config.UserConfig, previousMessages []string) {

	ui.RenderBox("Commit message", commitMessage)

	choice := choicePrompt("What would you like to do next?")

	switch choice {
	case "commit":
		ui.RenderTitle("Creating commit...")
		commitCmd := exec.Command("git", "commit", "--no-verify", "-m", commitMessage)
		err := commitCmd.Run()
		if err != nil {
			ui.RenderError(fmt.Sprintf("Commit failed: %v", err))
			os.Exit(1)
		}
		fmt.Println()
		ui.RenderTitle("Commit successfully added to history!")
	case "regenerate":
		fmt.Println()

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

		var newCommitMessage string
		err := ui.WithSpinner("Generating alternative commit message...", func() error {
			var genErr error
			newCommitMessage, genErr = CreateCommitMessage(modifiedPrompt, userConfig)
			return genErr
		})
		if err != nil {
			ui.RenderError(fmt.Sprintf("Error: %v", err))
			os.Exit(1)
		}

		updatedHistory := append(previousMessages, commitMessage)
		HandleCommitFlowWithHistory(newCommitMessage, fullPrompt, userConfig, updatedHistory)
	case "custom":
		fmt.Println()
		customInput := customInputPrompt("What changes would you like to see in the commit message?")
		fmt.Println()

		modifiedPrompt := fullPrompt + fmt.Sprintf("\n\nCurrent commit message:\n%s\n\nUser feedback: %s\n\nPlease generate a new commit message that addresses the user's feedback.", commitMessage, customInput)

		var newCommitMessage string
		err := ui.WithSpinner("Refining commit message with your feedback...", func() error {
			var genErr error
			newCommitMessage, genErr = CreateCommitMessage(modifiedPrompt, userConfig)
			return genErr
		})
		if err != nil {
			ui.RenderError(fmt.Sprintf("Error: %v", err))
			os.Exit(1)
		}

		updatedHistory := append(previousMessages, commitMessage)
		HandleCommitFlowWithHistory(newCommitMessage, fullPrompt, userConfig, updatedHistory)
	case "exit":
		ui.RenderTitle("Thanks for using Diny!")
		fmt.Println()
		os.Exit(0)
	}
}

func choicePrompt(message string) string {
	var choice string

	err := huh.NewSelect[string]().
		Title("🦕 "+message).
		Description("Select an option using arrow keys and press Enter").
		Options(
			huh.NewOption("Commit this message", "commit"),
			huh.NewOption("Generate different message", "regenerate"),
			huh.NewOption("Refine with feedback", "custom"),
			huh.NewOption("Exit", "exit"),
		).
		Value(&choice).
		Height(6).
		Run()

	if err != nil {
		ui.RenderError(fmt.Sprintf("Error running prompt: %v", err))
		os.Exit(1)
	}

	return choice
}

func customInputPrompt(message string) string {
	var input string

	err := huh.NewInput().
		Title("🦕 " + message).
		Description("Provide specific feedback to improve the commit message").
		Placeholder("e.g., make it shorter, use conventional format, focus on the bug fix...").
		CharLimit(200).
		Value(&input).
		Run()

	if err != nil {
		ui.RenderError(fmt.Sprintf("Error running prompt: %v", err))
		os.Exit(1)
	}

	return input
}
