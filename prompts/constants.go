package prompts

const (
	// MinCommitsBeforePrompt is the number of successful commits before any prompt can appear.
	MinCommitsBeforePrompt = 5

	// TriggerProbability is the chance a prompt appears on an eligible commit.
	TriggerProbability = 0.15

	// GitHubRepoURL is opened in the browser when the user chooses to star.
	GitHubRepoURL = "https://github.com/dinoDanic/diny"
)
