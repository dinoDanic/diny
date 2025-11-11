/* Copyright © 2025 dinoDanic dino.danic@gmail.com */
package commit

import (
	"github.com/dinoDanic/diny/cli/api"
	"github.com/dinoDanic/diny/cli/config"
	"github.com/dinoDanic/diny/cli/version"
)

func CreateCommitMessage(gitDiff string, cfg *config.CommitConfig) (string, error) {
	commitMessage, err := api.CreateCommitMessage(gitDiff, cfg, version.Get())

	if err != nil {
		return "", err
	}

	return commitMessage, nil
}
