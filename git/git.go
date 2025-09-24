/*
Copyright © 2025 NAME HERE dino.danic@gmail.com
*/
package git

import (
	"fmt"
	"os"
	"path/filepath"
)

// FindGitRoot finds the git repository root directory
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
