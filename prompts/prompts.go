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

	// Always increment commit count on successful commit.
	state.Prompts.CommitCount++
	_ = SaveState(state)

	// Gate: prompts.enabled config flag.
	if !cfg.Prompts.Enabled {
		return
	}

	// Gate: non-interactive environments (CI, piped stdout).
	if !isInteractive() {
		return
	}

	// Gate: minimum commits.
	if state.Prompts.CommitCount < MinCommitsBeforePrompt {
		return
	}

	ratingPending := state.Prompts.Rating.Status == StatusPending
	starPending := state.Prompts.Star.Status == StatusPending

	// Gate: at least one prompt must be pending.
	if !ratingPending && !starPending {
		return
	}

	// Gate: random roll — only show on ~15% of eligible commits.
	if rand.Float64() >= TriggerProbability {
		return
	}

	// Pick which prompt to show (at most one per invocation).
	switch {
	case ratingPending && starPending:
		if rand.Intn(2) == 0 {
			showRatingPrompt(state)
		} else {
			showStarPrompt(state)
		}
	case ratingPending:
		showRatingPrompt(state)
	case starPending:
		showStarPrompt(state)
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

func showStarPrompt(state *State) {
	outcome := ShowStar(state)
	_ = SaveState(state)

	// All star outcomes are POSTed (starred, already_given, dismissed).
	// Empty string means read error — don't POST.
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

func showRatingPrompt(state *State) {
	value := ShowRating(state)
	_ = SaveState(state)

	// Only POST if the user rated (1-5), not on dismiss (0) or error (-1).
	if value >= 1 && value <= 5 {
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
