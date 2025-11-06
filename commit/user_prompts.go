package commit

import (
	"fmt"
	"os"

	"github.com/charmbracelet/huh"
	"github.com/dinoDanic/diny/ui"
)

func choicePrompt() string {
	var choice string

	err := huh.NewSelect[string]().
		Title("What would you like to do next?").
		Description("Select an option using arrow keys or j,k and press Enter").
		Options(
			huh.NewOption("Commit this message", "commit"),
			huh.NewOption("Commit and push", "commit-push"),
			huh.NewOption("Edit in $EDITOR", "edit"),
			huh.NewOption("Save as draft", "save"),
			huh.NewOption("Generate different message", "regenerate"),
			huh.NewOption("Refine with feedback", "custom"),
			huh.NewOption("Exit", "exit"),
		).
		Value(&choice).
		Height(9).
		WithTheme(ui.GetHuhPrimaryTheme()).
		Run()

	if err != nil {
		ui.Box(ui.BoxOptions{Message: fmt.Sprintf("Error running prompt: %v", err), Variant: ui.Error})
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
		WithTheme(ui.GetHuhPrimaryTheme()).
		Run()

	if err != nil {
		ui.Box(ui.BoxOptions{Message: fmt.Sprintf("Error running prompt: %v", err), Variant: ui.Error})
		os.Exit(1)
	}

	return input
}
