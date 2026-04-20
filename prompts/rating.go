package prompts

import (
	"time"

	"github.com/dinoDanic/diny/config"
	tuiprompts "github.com/dinoDanic/diny/tui/prompts"
)

// ShowRating runs the rating TUI and mutates state.
// Returns 1..3 if rated, 0 if explicitly dismissed, -1 if cancelled (state unchanged).
func ShowRating(state *State, cfg *config.Config) int {
	res := tuiprompts.RunRating(cfg)
	if res.Cancelled {
		return -1
	}

	now := time.Now()
	switch {
	case res.Value >= 1 && res.Value <= 3:
		state.Prompts.Rating.Status = StatusRated
		state.Prompts.Rating.Value = res.Value
		state.Prompts.Rating.AnsweredAt = &now
		return res.Value
	case res.Value == 0:
		state.Prompts.Rating.Status = StatusDismissed
		state.Prompts.Rating.AnsweredAt = &now
		return 0
	}
	return -1
}
