package prompts

const (
	// MinCommitsBeforePrompt is the number of successful commits before any prompt can appear.
	MinCommitsBeforePrompt = 5

	// TriggerProbability is the chance a prompt appears on an eligible commit.
	// Set to 1.0 for Phase 1 (always show when eligible) — will become 0.15 in Phase 3.
	TriggerProbability = 1.0

	// GitHubRepoURL is opened in the browser when the user chooses to star.
	GitHubRepoURL = "https://github.com/dinoDanic/diny"
)
