package prompts

import (
	"bufio"
	"fmt"
	"os"
	"time"

	"github.com/dinoDanic/diny/ui"
)

// ShowRating displays the rating prompt and handles the user's keypress.
// Returns the rating value (1-5) if rated, 0 if dismissed, -1 if invalid/error.
func ShowRating(state *State) int {
	ui.Box("How's diny working for you?",
		"[1] 😐  [2] 🙂  [3] 😀  [4] 🤩  [5] ❤️   ·   [0] Don't ask again")

	for {
		r := readSingleRune()
		if r == 0 {
			return -1 // read error
		}

		switch r {
		case '1', '2', '3', '4', '5':
			value := int(r - '0')
			now := time.Now()
			state.Prompts.Rating.Status = StatusRated
			state.Prompts.Rating.Value = value
			state.Prompts.Rating.AnsweredAt = &now
			fmt.Println()
			ui.Success("Thanks!")
			return value

		case '0':
			now := time.Now()
			state.Prompts.Rating.Status = StatusDismissed
			state.Prompts.Rating.AnsweredAt = &now
			return 0

		default:
			// Ignore invalid keys, re-read.
			continue
		}
	}
}

// readSingleRune reads a single rune from stdin in raw-ish mode.
func readSingleRune() rune {
	reader := bufio.NewReader(os.Stdin)
	r, _, err := reader.ReadRune()
	if err != nil {
		return 0
	}
	return r
}
