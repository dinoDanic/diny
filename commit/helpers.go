package commit

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/atotto/clipboard"
	"github.com/dinoDanic/diny/config"
	"github.com/dinoDanic/diny/git"
)

// TryCommit runs git commit and returns the short hash on success.
func TryCommit(message string, push bool, noVerify bool, cfg *config.Config) (string, error) {
	var commitCmd *exec.Cmd
	if noVerify {
		commitCmd = exec.Command("git", "commit", "--no-verify", "-m", message)
	} else {
		commitCmd = exec.Command("git", "commit", "-m", message)
	}
	output, err := commitCmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("commit failed: %s", strings.TrimSpace(string(output)))
	}

	var hash string
	if cfg != nil && cfg.Commit.HashAfterCommit {
		hashCmd := exec.Command("git", "rev-parse", "--short", "HEAD")
		hashOutput, hashErr := hashCmd.Output()
		if hashErr == nil {
			hash = strings.TrimSpace(string(hashOutput))
			_ = clipboard.WriteAll(hash)
		}
	}

	if push {
		pushCmd := exec.Command("git", "push")
		pushOut, pushErr := pushCmd.CombinedOutput()
		if pushErr != nil {
			return hash, fmt.Errorf("committed but push failed: %s", strings.TrimSpace(string(pushOut)))
		}
	}

	return hash, nil
}

// SaveDraft writes the commit message to git draft files for use by other tools.
func SaveDraft(message string) error {
	gitDir, err := git.FindGitDir()
	if err != nil {
		return fmt.Errorf("failed to find git repository: %v", err)
	}

	draftFiles := []string{
		"COMMIT_EDITMSG",
		"PREPARE_COMMIT_MSG",
		"LAZYGIT_PENDING_COMMIT",
	}

	var errors []string
	successCount := 0

	for _, file := range draftFiles {
		filePath := filepath.Join(gitDir, file)
		if err := os.WriteFile(filePath, []byte(message), 0644); err != nil {
			errors = append(errors, fmt.Sprintf("%s: %v", file, err))
		} else {
			successCount++
		}
	}

	if successCount == 0 {
		return fmt.Errorf("failed to write to any draft files: %s", strings.Join(errors, ", "))
	}

	return nil
}
