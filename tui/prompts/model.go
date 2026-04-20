package prompts

import (
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/dinoDanic/diny/config"
)

// Rating prompt.

type ratingOption struct {
	value int // 1..3 for rate, 0 for dismiss
	label string
}

var ratingOptions = []ratingOption{
	{3, "Good"},
	{2, "Meh"},
	{1, "Bad"},
	{0, "Don't ask again"},
}

type ratingModel struct {
	cfg    *config.Config
	width  int
	cursor int

	// Set on exit.
	cancelled bool
	value     int // 1..3 if rated, 0 if dismissed
}

func newRatingModel(cfg *config.Config) ratingModel {
	return ratingModel{
		cfg:    cfg,
		cursor: 0,
	}
}

// Star prompt.

type starOption struct {
	key     string // "starred" | "already_given" | "dismissed"
	quickKey rune
	label   string
}

var starOptions = []starOption{
	{"starred", '1', "Star it"},
	{"already_given", '2', "Already starred"},
	{"dismissed", '0', "Don't ask again"},
}

type starModel struct {
	cfg    *config.Config
	width  int
	cursor int

	cancelled bool
	outcome   string // "starred" | "already_given" | "dismissed"
}

func newStarModel(cfg *config.Config) starModel {
	return starModel{
		cfg:    cfg,
		cursor: 0,
	}
}

// Feedback prompt.

type feedbackModel struct {
	cfg      *config.Config
	width    int
	textarea textarea.Model

	cancelled bool
	submitted bool
	text      string
}

func newFeedbackModel(cfg *config.Config) feedbackModel {
	ta := textarea.New()
	ta.Placeholder = "Type your feedback or feature request…"
	ta.SetHeight(4)
	ta.SetWidth(60)
	ta.Focus()

	return feedbackModel{
		cfg:      cfg,
		textarea: ta,
	}
}
