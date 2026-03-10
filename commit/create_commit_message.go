package commit

import (
	"github.com/dinoDanic/diny/ai"
	"github.com/dinoDanic/diny/config"
)

func CreateCommitMessage(gitDiff string, cfg *config.Config) (string, error) {
	return ai.GenerateCommitMessage(gitDiff, cfg)
}
