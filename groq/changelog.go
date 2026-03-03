package groq

import "github.com/dinoDanic/diny/config"

func CreateChangelogWithGroq(prompt string, cfg *config.Config) (string, error) {
	return CreateTimelineWithGroq(prompt, cfg)
}
