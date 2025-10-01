package commit

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/huh"
	"github.com/dinoDanic/diny/config"
	"github.com/dinoDanic/diny/git"
	"github.com/dinoDanic/diny/ui"
)

func HandleCommitFlow(commitMessage, fullPrompt string, userConfig *config.UserConfig) {
	HandleCommitFlowWithHistory(commitMessage, fullPrompt, userConfig, []string{})
}

func HandleCommitFlowWithHistory(commitMessage, fullPrompt string, userConfig *config.UserConfig, previousMessages []string) {

	ui.Box(ui.BoxOptions{Title: "Commit message", Message: commitMessage})

	choice := choicePrompt()

	switch choice {
	case "commit":
		commitCmd := exec.Command("git", "commit", "--no-verify", "-m", commitMessage)
		err := commitCmd.Run()
		if err != nil {
			ui.Box(ui.BoxOptions{Message: fmt.Sprintf("Commit failed: %v", err), Variant: ui.Error})
			os.Exit(1)
		}
		ui.Box(ui.BoxOptions{Message: "Commited!", Variant: ui.Success})
		fmt.Println()
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
		err := ui.WithSpinner("Generating alternative commit message...", func() error {
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
		err := ui.WithSpinner("Refining commit message with your feedback...", func() error {
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
		os.Exit(0)
	}
}

func choicePrompt() string {
	var choice string

	err := huh.NewSelect[string]().
		Title("What would you like to do next?").
		Description("Select an option using arrow keys or j,k and press Enter").
		Options(
			huh.NewOption("Commit this message", "commit"),
			huh.NewOption("Edit in $EDITOR", "edit"),
			huh.NewOption("Save as draft", "save"),
			huh.NewOption("Generate different message", "regenerate"),
			huh.NewOption("Refine with feedback", "custom"),
			huh.NewOption("Exit", "exit"),
		).
		Value(&choice).
		Height(8).
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

func openInEditor(message string) (string, error) {
	editor := git.GetGitEditor()

	tmpFile, err := os.CreateTemp("", "diny-commit-*.txt")
	if err != nil {
		return "", fmt.Errorf("failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())
	defer tmpFile.Close()

	if _, err := tmpFile.WriteString(message); err != nil {
		return "", fmt.Errorf("failed to write to temp file: %v", err)
	}
	tmpFile.Close()

	editorArgs := strings.Fields(editor)
	editorCmd := editorArgs[0]
	args := append(editorArgs[1:], tmpFile.Name())

	cmd := exec.Command(editorCmd, args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("editor exited with error: %v", err)
	}

	editedContent, err := os.ReadFile(tmpFile.Name())
	if err != nil {
		return "", fmt.Errorf("failed to read edited file: %v", err)
	}

	return strings.TrimSpace(string(editedContent)), nil
}

func saveDraft(message string) error {
	repoRoot, err := git.FindGitRoot()
	if err != nil {
		return fmt.Errorf("failed to find git repository: %v", err)
	}

	draftFiles := []string{
		"COMMIT_EDITMSG",         // Standard git, tig, magit
		"PREPARE_COMMIT_MSG",     // Git hooks & some GUIs
		"LAZYGIT_PENDING_COMMIT", // lazygit
	}

	var errors []string
	successCount := 0

	for _, file := range draftFiles {
		filePath := filepath.Join(repoRoot, ".git", file)
		if err := os.WriteFile(filePath, []byte(message), 0644); err != nil {
			errors = append(errors, fmt.Sprintf("%s: %v", file, err))
		} else {
			successCount++
		}
	}

	if successCount == 0 {
		return fmt.Errorf("failed to write to any draft files: %s", strings.Join(errors, ", "))
	}

	if len(errors) > 0 {
		return fmt.Errorf("partial success - some files failed: %s", strings.Join(errors, ", "))
	}

	return nil
}
