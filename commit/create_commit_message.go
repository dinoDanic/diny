package commit

import (
	"github.com/dinoDanic/diny/config"
	"github.com/dinoDanic/diny/groq"
)

func CreateCommitMessage(gitDiff string, userConfig *config.Config) (string, error) {
	commitMessage, err := groq.CreateCommitMessageWithGroq(gitDiff, userConfig)

	if err != nil {
		return "", err
	}

	return commitMessage, nil
}
