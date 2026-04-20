package prompts

import (
	"runtime"
	"strconv"

	"github.com/dinoDanic/diny/config"
	"github.com/dinoDanic/diny/feedback"
	"github.com/dinoDanic/diny/git"
	"github.com/dinoDanic/diny/version"
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

	// Phase 2: always show when eligible (TriggerProbability = 1.0).
	// Phase 3 will add random roll here.

	// When both are pending, show rating first (deterministic for Phase 2).
	// Phase 3 will add 50/50 random selection.
	if ratingPending {
		showRatingPrompt(state)
	} else if starPending {
		showStarPrompt(state)
	}
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
