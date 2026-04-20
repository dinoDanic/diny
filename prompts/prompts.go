package prompts

import (
	"math/rand"
	"os"
	"runtime"
	"strconv"

	"github.com/dinoDanic/diny/config"
	"github.com/dinoDanic/diny/feedback"
	"github.com/dinoDanic/diny/git"
	"github.com/dinoDanic/diny/version"
	"github.com/mattn/go-isatty"
)

// MaybeShow is called once after a successful commit. It handles all side effects:
// loading/saving state, eligibility checks, prompt display, and HTTP feedback.
func MaybeShow(cfg *config.Config) {
	state := LoadState()

	// Gate: prompts.enabled config flag.
	if !cfg.Prompts.Enabled {
		return
	}

	// Gate: non-interactive environments (CI, piped stdout).
	if !isInteractive() {
		return
	}

	pending := []string{}
	if state.Prompts.Rating.Status == StatusPending {
		pending = append(pending, "rating")
	}
	if state.Prompts.Star.Status == StatusPending {
		pending = append(pending, "star")
	}
	if state.Prompts.Feedback.Status == StatusPending {
		pending = append(pending, "feedback")
	}

	// Gate: at least one prompt must be pending.
	if len(pending) == 0 {
		return
	}

	// Gate: random roll — only show on ~15% of eligible commits.
	if rand.Float64() >= TriggerProbability {
		return
	}

	// Pick which prompt to show (at most one per invocation).
	switch pending[rand.Intn(len(pending))] {
	case "rating":
		showRatingPrompt(state, cfg)
	case "star":
		showStarPrompt(state, cfg)
	case "feedback":
		showFeedbackPrompt(state, cfg)
	}
}

// isInteractive returns false when prompts should be suppressed:
// stdout is not a TTY, or common CI env vars are set.
func isInteractive() bool {
	if !isatty.IsTerminal(os.Stdout.Fd()) && !isatty.IsCygwinTerminal(os.Stdout.Fd()) {
		return false
	}
	for _, env := range []string{"CI", "GITHUB_ACTIONS", "NONINTERACTIVE"} {
		if v := os.Getenv(env); v != "" {
			return false
		}
	}
	return true
}

func showStarPrompt(state *State, cfg *config.Config) {
	outcome := ShowStar(state, cfg)
	_ = SaveState(state)

	if outcome != "" {
		feedback.Send(feedback.Payload{
			Type:     "star",
			Value:    outcome,
			Email:    git.GetGitEmail(),
			Name:     git.GetGitName(),
			Version:  version.Get(),
			System:   runtime.GOOS,
			RepoName: git.GetRepoName(),
		})
	}
}

func showRatingPrompt(state *State, cfg *config.Config) {
	value := ShowRating(state, cfg)
	_ = SaveState(state)

	// Only POST if the user rated (1-3), not on dismiss (0) or error (-1).
	if value >= 1 && value <= 3 {
		feedback.Send(feedback.Payload{
			Type:     "rating",
			Value:    strconv.Itoa(value),
			Email:    git.GetGitEmail(),
			Name:     git.GetGitName(),
			Version:  version.Get(),
			System:   runtime.GOOS,
			RepoName: git.GetRepoName(),
		})
	}
}

func showFeedbackPrompt(state *State, cfg *config.Config) {
	text := ShowFeedback(state, cfg)
	_ = SaveState(state)

	if text != "" {
		feedback.Send(feedback.Payload{
			Type:     "feedback",
			Value:    text,
			Email:    git.GetGitEmail(),
			Name:     git.GetGitName(),
			Version:  version.Get(),
			System:   runtime.GOOS,
			RepoName: git.GetRepoName(),
		})
	}
}
