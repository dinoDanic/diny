package commit

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/dinoDanic/diny/git"
	"github.com/dinoDanic/diny/ui"
)

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

func executeCommit(commitMessage string, push bool) {
	var output []byte
	var err error

	spinnerErr := ui.WithSpinner("Committing...", func() error {
		commitCmd := exec.Command("git", "commit", "-m", commitMessage)
		output, err = commitCmd.CombinedOutput()
		return err
	})

	if spinnerErr != nil {
		if len(output) > 0 {
			fmt.Fprint(os.Stderr, string(output))
		}
		ui.Box(ui.BoxOptions{Message: fmt.Sprintf("Commit failed: %v", spinnerErr), Variant: ui.Error})
		os.Exit(1)
	}
	ui.Box(ui.BoxOptions{Message: "Commited!", Variant: ui.Success})

	if push {
		var pushOutput []byte
		var pushErr error

		pushSpinnerErr := ui.WithSpinner("Pushing...", func() error {
			pushCmd := exec.Command("git", "push")
			pushOutput, pushErr = pushCmd.CombinedOutput()
			return pushErr
		})

		if pushSpinnerErr != nil {
			if len(pushOutput) > 0 {
				fmt.Fprint(os.Stderr, string(pushOutput))
			}
			ui.Box(ui.BoxOptions{Message: fmt.Sprintf("Push failed: %v", pushSpinnerErr), Variant: ui.Error})
			os.Exit(1)
		}
		ui.Box(ui.BoxOptions{Message: "Pushed!", Variant: ui.Success})
	}

	fmt.Println()
}
