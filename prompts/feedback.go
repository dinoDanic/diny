package prompts

import (
	"time"

	"github.com/dinoDanic/diny/config"
	tuiprompts "github.com/dinoDanic/diny/tui/prompts"
)

// ShowFeedback runs the feedback TUI and mutates state.
// Returns the submitted text when the user confirmed non-empty input,
// or "" when cancelled (state unchanged).
func ShowFeedback(state *State, cfg *config.Config) string {
	res := tuiprompts.RunFeedback(cfg)
	if res.Cancelled {
		return ""
	}

	now := time.Now()
	if res.Submitted && res.Text != "" {
		state.Prompts.Feedback.Status = StatusSubmitted
		state.Prompts.Feedback.AnsweredAt = &now
		return res.Text
	}

	return ""
}
