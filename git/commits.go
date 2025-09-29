/*
Copyright © 2025 NAME HERE dino.danic@gmail.com
*/
package git

import (
	"fmt"
	"os/exec"
	"strings"
	"time"
)

func GetCommitsToday() ([]string, error) {
	today := time.Now().Format("2006-01-02")
	return GetCommitsByDate(today)
}

func GetCommitsByDate(date string) ([]string, error) {
	startDate := date + " 00:00:00"
	endDate := date + " 23:59:59"
	return GetCommitsByDateRange(startDate, endDate)
}

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

	var filteredCommits []string
	for _, commit := range commits {
		if strings.TrimSpace(commit) != "" {
			filteredCommits = append(filteredCommits, strings.TrimSpace(commit))
		}
	}

	return filteredCommits, nil
}
