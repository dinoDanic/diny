package prompts

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/dinoDanic/diny/config"
)

type RatingResult struct {
	Value     int  // 1..3 if rated, 0 if dismissed
	Cancelled bool // true if user closed without answering (keep pending)
}

type StarResult struct {
	Outcome   string // "starred" | "already_given" | "dismissed"
	Cancelled bool   // true if user closed without answering (keep pending)
}

type FeedbackResult struct {
	Text      string
	Submitted bool
	Cancelled bool // true if user closed without submitting (keep pending)
}

func RunRating(cfg *config.Config) RatingResult {
	p := tea.NewProgram(newRatingModel(cfg))
	final, err := p.Run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		return RatingResult{Cancelled: true}
	}
	if m, ok := final.(ratingModel); ok {
		if m.cancelled {
			return RatingResult{Cancelled: true}
		}
		return RatingResult{Value: m.value}
	}
	return RatingResult{Cancelled: true}
}

func RunStar(cfg *config.Config) StarResult {
	p := tea.NewProgram(newStarModel(cfg))
	final, err := p.Run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		return StarResult{Cancelled: true}
	}
	if m, ok := final.(starModel); ok {
		if m.cancelled {
			return StarResult{Cancelled: true}
		}
		return StarResult{Outcome: m.outcome}
	}
	return StarResult{Cancelled: true}
}

// PrintThanks prints a themed one-line acknowledgement after a prompt is answered.
func PrintThanks(message string) {
	fmt.Println(indentStyle().Render(thanksStyle().Render(message)))
}

func RunFeedback(cfg *config.Config) FeedbackResult {
	p := tea.NewProgram(newFeedbackModel(cfg))
	final, err := p.Run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		return FeedbackResult{Cancelled: true}
	}
	if m, ok := final.(feedbackModel); ok {
		if m.cancelled {
			return FeedbackResult{Cancelled: true}
		}
		return FeedbackResult{Text: m.text, Submitted: m.submitted}
	}
	return FeedbackResult{Cancelled: true}
}
