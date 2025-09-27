package commit

import (
	"github.com/dinoDanic/diny/config"
	"github.com/dinoDanic/diny/groq"
)

func CreateCommitMessage(gitDiff string, userConfig *config.UserConfig) (string, string, error) {
	commitMessage, note, err := groq.CreateCommitMessageWithGroq(gitDiff, userConfig)

	if err != nil {
		return "", "", err
	}

	return commitMessage, note, nil
}
