package commit

import (
	"github.com/dinoDanic/diny/config"
	"github.com/dinoDanic/diny/groq"
)

func CreateCommitMessage(gitDiff string, cfg *config.Config) (string, error) {
	commitMessage, err := groq.CreateCommitMessageWithGroq(gitDiff, cfg)

	if err != nil {
		return "", err
	}

	return commitMessage, nil
}
