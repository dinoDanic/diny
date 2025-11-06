package commit

import (
	"fmt"

	"github.com/dinoDanic/diny/config"
)

func buildCommitPrompt(gitDiff string, userConfig *config.UserConfig) string {
	prompt := "You are a commit message generator. Generate a clear, concise commit message based on the following git diff.\n\n"

	prompt += "IMPORTANT: Output ONLY the commit message text. Do not include any explanations, descriptions, or meta-commentary about the commit message.\n\n"

	if userConfig.UseConventional {
		prompt += "Format: Use Conventional Commits format (type(scope): description)\n"
	}

	if userConfig.UseEmoji {
		prompt += "Style: Include appropriate emoji prefixes\n"
	}

	prompt += fmt.Sprintf("Tone: %s\n", userConfig.Tone)
	prompt += fmt.Sprintf("Length: %s\n", userConfig.Length)

	prompt += "\nGit diff:\n" + gitDiff + "\n\n"
	prompt += "Output the commit message now (no other text):"

	return prompt
}
