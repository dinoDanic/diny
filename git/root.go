/*
Copyright © 2025 NAME HERE dino.danic@gmail.com
*/
package git

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
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

func FindGitDir() (string, error) {
	repoRoot, err := FindGitRoot()
	if err != nil {
		return "", err
	}

	gitPath := filepath.Join(repoRoot, ".git")
	info, err := os.Stat(gitPath)
	if err != nil {
		return "", err
	}

	if info.IsDir() {
		return gitPath, nil
	}

	data, err := os.ReadFile(gitPath)
	if err != nil {
		return "", err
	}

	content := strings.TrimSpace(string(data))
	if content == "" {
		return "", fmt.Errorf(".git file is empty")
	}

	if idx := strings.IndexByte(content, '\n'); idx >= 0 {
		content = content[:idx]
	}

	const prefix = "gitdir:"
	if !strings.HasPrefix(content, prefix) {
		return "", fmt.Errorf("invalid .git file: missing gitdir")
	}

	path := strings.TrimSpace(strings.TrimPrefix(content, prefix))
	if path == "" {
		return "", fmt.Errorf("invalid .git file: empty gitdir")
	}

	if !filepath.IsAbs(path) {
		path = filepath.Join(repoRoot, path)
	}

	return filepath.Clean(path), nil
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
