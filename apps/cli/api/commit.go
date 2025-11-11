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

	// TODO: Implement actual API call to backend
	// For now, return a mock commit message based on config preferences
	var message string

	if cfg.Emoji {
		message += "✨ "
	}

	if cfg.Conventional {
		message += "feat: "
	}

	switch cfg.Length {
	case config.Short:
		message += "implement new feature"
	case config.Normal:
		message += "implement new feature\n\nAdds functionality based on staged changes"
	case config.Long:
		message += "implement new feature\n\nAdds comprehensive functionality based on staged changes.\nThis includes detailed implementation of the feature with proper\nerror handling and documentation."
	default:
		message += "implement new feature"
	}

	return message, nil
}
