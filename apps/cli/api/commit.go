/* Copyright © 2025 dinoDanic dino.danic@gmail.com */
package api

import (
	"fmt"

	"github.com/dinoDanic/diny/cli/config"
)

type CommitRequest struct {
	GitDiff string               `json:"gitDiff"`
	Version string               `json:"version"`
	Config  *config.CommitConfig `json:"config"`
}

type CommitResponse struct {
	CommitMessage string `json:"commitMessage"`
	Error         string `json:"error,omitempty"`
}

// TODO: Replace this fake implementation with actual API call to backend
func CreateCommitMessage(gitDiff string, cfg *config.CommitConfig, version string) (string, error) {
	if gitDiff == "" {
		return "", fmt.Errorf("git diff is empty")
	}

	var message = "response"

	return message, nil
}
