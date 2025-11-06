package commit

import (
	"github.com/dinoDanic/diny/config"
	"github.com/dinoDanic/diny/groq"
	"github.com/dinoDanic/diny/ollama"
)

func CreateCommitMessage(gitDiff string, userConfig *config.UserConfig) (string, error) {
	configService := config.GetService()

	if configService.IsUsingLocalAPI() {
		return ollama.CreateCommitMessage(gitDiff, userConfig)
	}

	return groq.CreateCommitMessageWithGroq(gitDiff, userConfig)
}
