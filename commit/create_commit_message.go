package commit

import (
	"github.com/dinoDanic/diny/config"
	"github.com/dinoDanic/diny/groq"
)

func CreateCommitMessage(prompt string, userConfig config.UserConfig) (string, error) {
	commitMessage, err := groq.CreateCommitMessageWithGroq(prompt, userConfig)

	if err != nil {
		return "", err
	}

	return commitMessage, nil
}
