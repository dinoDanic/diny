package commit

import (
	"fmt"
	"os"

	"github.com/dinoDanic/diny/config"
	"github.com/dinoDanic/diny/groq"
	"github.com/dinoDanic/diny/ollama"
	"github.com/dinoDanic/diny/ui"
)

func CreateCommitMessage(gitDiff string, userConfig *config.UserConfig) (string, error) {
	configService := config.GetService()

	if configService.IsUsingLocalAPI() {
		return ollama.CreateCommitMessage(gitDiff, userConfig, buildCommitPrompt)
	}

	return groq.CreateCommitMessageWithGroq(gitDiff, userConfig)
}

func HandleCommitFlow(commitMessage, fullPrompt string, userConfig *config.UserConfig) {
	HandleCommitFlowWithHistory(commitMessage, fullPrompt, userConfig, []string{})
}

func HandleCommitFlowWithHistory(commitMessage, fullPrompt string, userConfig *config.UserConfig, previousMessages []string) {

	ui.Box(ui.BoxOptions{Title: "Commit message", Message: commitMessage})

	choice := choicePrompt()

	switch choice {
	case "commit":
		executeCommit(commitMessage, false)
	case "commit-push":
		executeCommit(commitMessage, true)
	case "edit":
		editedMessage, err := openInEditor(commitMessage)
		if err != nil {
			ui.Box(ui.BoxOptions{Message: fmt.Sprintf("Failed to open editor: %v", err), Variant: ui.Error})
			HandleCommitFlowWithHistory(commitMessage, fullPrompt, userConfig, previousMessages)
			return
		}
		if editedMessage != commitMessage && editedMessage != "" {
			HandleCommitFlowWithHistory(editedMessage, fullPrompt, userConfig, previousMessages)
		} else {
			HandleCommitFlowWithHistory(commitMessage, fullPrompt, userConfig, previousMessages)
		}
	case "save":
		if err := saveDraft(commitMessage); err != nil {
			ui.Box(ui.BoxOptions{Message: fmt.Sprintf("Failed to save draft: %v", err), Variant: ui.Error})
			HandleCommitFlowWithHistory(commitMessage, fullPrompt, userConfig, previousMessages)
			return
		}

		ui.Box(ui.BoxOptions{Message: "Draft saved!", Variant: ui.Success})
	case "regenerate":
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
		configService := config.GetService()
		apiConfig := configService.GetAPIConfig()

		var spinnerMessage string
		if apiConfig.Provider == config.LocalOllama {
			spinnerMessage = "Regenerating locally..."
		} else {
			spinnerMessage = "Regenerating via Diny cloud..."
		}

		err := ui.WithSpinner(spinnerMessage, func() error {
			var genErr error
			newCommitMessage, genErr = CreateCommitMessage(modifiedPrompt, userConfig)
			return genErr
		})
		if err != nil {
			ui.Box(ui.BoxOptions{Message: fmt.Sprintf("Error: %v", err), Variant: ui.Error})
			os.Exit(1)
		}

		updatedHistory := append(previousMessages, commitMessage)
		HandleCommitFlowWithHistory(newCommitMessage, fullPrompt, userConfig, updatedHistory)
	case "custom":
		customInput := customInputPrompt("What changes would you like to see in the commit message?")

		modifiedPrompt := fullPrompt + fmt.Sprintf("\n\nCurrent commit message:\n%s\n\nUser feedback: %s\n\nPlease generate a new commit message that addresses the user's feedback.", commitMessage, customInput)

		var newCommitMessage string
		configService := config.GetService()
		apiConfig := configService.GetAPIConfig()

		var spinnerMessage string
		if apiConfig.Provider == config.LocalOllama {
			spinnerMessage = "Refining locally with your feedback..."
		} else {
			spinnerMessage = "Refining via Diny cloud with your feedback..."
		}

		err := ui.WithSpinner(spinnerMessage, func() error {
			var genErr error
			newCommitMessage, genErr = CreateCommitMessage(modifiedPrompt, userConfig)
			return genErr
		})
		if err != nil {
			ui.Box(ui.BoxOptions{Message: fmt.Sprintf("Error: %v", err), Variant: ui.Error})
			os.Exit(1)
		}

		updatedHistory := append(previousMessages, commitMessage)
		HandleCommitFlowWithHistory(newCommitMessage, fullPrompt, userConfig, updatedHistory)
	case "exit":
		ui.RenderTitle("Bye!")
		fmt.Println()
		os.Exit(0)
	}
}
