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
	"time"
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

// GetCommitsToday returns commit messages from today
func GetCommitsToday() ([]string, error) {
	today := time.Now().Format("2006-01-02")
	return GetCommitsByDate(today)
}

// GetCommitsByDate returns commit messages from a specific date
func GetCommitsByDate(date string) ([]string, error) {
	startDate := date + " 00:00:00"
	endDate := date + " 23:59:59"
	return GetCommitsByDateRange(startDate, endDate)
}

// GetCommitsByDateRange returns commit messages between two dates
func GetCommitsByDateRange(startDate, endDate string) ([]string, error) {
	cmd := exec.Command("git", "log",
		"--since="+startDate,
		"--until="+endDate,
		"--pretty=format:%s",
		"--no-merges")

	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to get git log: %w", err)
	}

	if len(output) == 0 {
		return []string{}, nil
	}

	commits := strings.Split(strings.TrimSpace(string(output)), "\n")

	// Filter out empty lines
	var filteredCommits []string
	for _, commit := range commits {
		if strings.TrimSpace(commit) != "" {
			filteredCommits = append(filteredCommits, strings.TrimSpace(commit))
		}
	}

	return filteredCommits, nil
}
