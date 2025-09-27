package commit

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/charmbracelet/huh"
	"github.com/dinoDanic/diny/config"
)

func HandleCommitFlow(commitMessage, note, fullPrompt string, userConfig *config.UserConfig) {
	HandleCommitFlowWithHistory(commitMessage, note, fullPrompt, userConfig, []string{})
}

func HandleCommitFlowWithHistory(commitMessage, note, fullPrompt string, userConfig *config.UserConfig, previousMessages []string) {
	fmt.Println()

	if note != "" {
		fmt.Println(RenderNote(note))
	}

	fmt.Println(RenderCommitMessage(commitMessage))
	fmt.Println()

	choice := choicePrompt("Choose your next action:")

	switch choice {
	case "commit":
		fmt.Println(RenderStep("Creating commit..."))
		commitCmd := exec.Command("git", "commit", "--no-verify", "-m", commitMessage)
		err := commitCmd.Run()
		if err != nil {
			fmt.Println(RenderError(fmt.Sprintf("Commit failed: %v", err)))
			os.Exit(1)
		}
		fmt.Println()
		fmt.Println(RenderSuccess("Commit successfully added to history!"))
	case "regenerate":
		fmt.Println()
		fmt.Println(RenderStep("Generating alternative commit message..."))

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

		newCommitMessage, newNote, err := CreateCommitMessage(modifiedPrompt, userConfig)
		if err != nil {
			fmt.Println(RenderError(fmt.Sprintf("Error: %v", err)))
			os.Exit(1)
		}

		updatedHistory := append(previousMessages, commitMessage)
		HandleCommitFlowWithHistory(newCommitMessage, newNote, fullPrompt, userConfig, updatedHistory)
	case "custom":
		fmt.Println()
		customInput := customInputPrompt("What changes would you like to see in the commit message?")
		fmt.Println()
		fmt.Println(RenderStep("Generating commit message with your feedback..."))

		modifiedPrompt := fullPrompt + fmt.Sprintf("\n\nCurrent commit message:\n%s\n\nUser feedback: %s\n\nPlease generate a new commit message that addresses the user's feedback.", commitMessage, customInput)

		newCommitMessage, newNote, err := CreateCommitMessage(modifiedPrompt, userConfig)
		if err != nil {
			fmt.Println(RenderError(fmt.Sprintf("Error: %v", err)))
			os.Exit(1)
		}

		updatedHistory := append(previousMessages, commitMessage)
		HandleCommitFlowWithHistory(newCommitMessage, newNote, fullPrompt, userConfig, updatedHistory)
	case "exit":
		fmt.Println()
		fmt.Println(RenderInfo("Thanks for using Diny! 👋"))
		os.Exit(0)
	}
}

func choicePrompt(message string) string {
	var choice string

	err := huh.NewSelect[string]().
		Title("🦕 "+message).
		Description("Select an option using arrow keys and press Enter").
		Options(
			huh.NewOption("✅ Commit this message", "commit"),
			huh.NewOption("🔄 Generate different message", "regenerate"),
			huh.NewOption("✏️  Refine with feedback", "custom"),
			huh.NewOption("❌ Exit", "exit"),
		).
		Value(&choice).
		Height(6).
		Run()

	if err != nil {
		fmt.Println(RenderError(fmt.Sprintf("Error running prompt: %v", err)))
		os.Exit(1)
	}

	return choice
}

func customInputPrompt(message string) string {
	var input string

	err := huh.NewInput().
		Title("✏️  " + message).
		Description("Provide specific feedback to improve the commit message").
		Placeholder("e.g., make it shorter, use conventional format, focus on the bug fix...").
		CharLimit(200).
		Value(&input).
		Run()

	if err != nil {
		fmt.Println(RenderError(fmt.Sprintf("Error running prompt: %v", err)))
		os.Exit(1)
	}

	return input
}
