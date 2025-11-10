/*
Copyright © 2025 NAME HERE dino.danic@gmail.com
*/
package git

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

func FindGitRoot() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	for {
		gitPath := filepath.Join(dir, ".git")
		if _, err := os.Stat(gitPath); err == nil {
			return dir, nil
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			return "", fmt.Errorf("not in a git repository")
		}
		dir = parent
	}
}

func GetGitDiff() (string, error) {
	gitDiffCmd := exec.Command("git", "diff", "--cached",
		"-U3", "--no-color", "--ignore-all-space", "--ignore-blank-lines",
		":(exclude)*.lock", ":(exclude)*package-lock.json", ":(exclude)*yarn.lock",
		":(exclude)node_modules/", ":(exclude)dist/", ":(exclude)build/")

	gitDiff, err := gitDiffCmd.Output()

	if err != nil {
		return "", err
	}

	return string(gitDiff), nil

}
