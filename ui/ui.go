package ui

import (
	"strings"

	"github.com/charmbracelet/huh/spinner"
	"github.com/charmbracelet/lipgloss"
)

var (
	// Color scheme
	PrimaryColor = lipgloss.Color("#726FF2")
	SuccessColor = lipgloss.Color("#00FF87")
	ErrorColor   = lipgloss.Color("#FF5F87")
	WarningColor = lipgloss.Color("#FFAF00")
	MutedColor   = lipgloss.Color("#6C7086")

	titleStyle = lipgloss.NewStyle().
			Foreground(PrimaryColor).
			Bold(true).
			Padding(0, 1).
			MarginBottom(1)

	warningStyle = lipgloss.NewStyle().
			Foreground(WarningColor).
			Bold(true).
			Padding(0, 1)

	infoStyle = lipgloss.NewStyle().
			Foreground(PrimaryColor).
			Padding(0, 1)

	// Box styles
	primaryBoxStyle = lipgloss.NewStyle().
			Background(lipgloss.Color("#25253A")).
			Padding(1, 2)

	noteBoxStyle = lipgloss.NewStyle().
			Background(lipgloss.Color("#3A3320")).
			Padding(1, 2)

	errorBoxStyle = lipgloss.NewStyle().
			Background(lipgloss.Color("#3A2025")).
			Foreground(ErrorColor).
			Bold(true).
			Padding(1, 2)

	successBoxStyle = lipgloss.NewStyle().
			Background(lipgloss.Color("#1A3A20")).
			Foreground(SuccessColor).
			Padding(1, 2)

	warningBoxStyle = lipgloss.NewStyle().
			Background(lipgloss.Color("#3A3320")).
			Foreground(WarningColor).
			Bold(true).
			Padding(1, 2)

	// Text styles
	mutedStyle = lipgloss.NewStyle().
			Foreground(MutedColor).
			Italic(true)

	boldStyle = lipgloss.NewStyle().
			Bold(true)
)

// UI component functions
func RenderTitle(text string) string {
	return titleStyle.Render("🦕 " + text)
}

func RenderSuccess(text string) string {
	return titleStyle.Render("🦕 " + text)
}

func RenderError(text string) string {
	return errorBoxStyle.Render(strings.TrimSpace(text))
}

func RenderWarning(text string) string {
	return warningBoxStyle.Render(strings.TrimSpace(text))
}

func RenderInfo(text string) string {
	return infoStyle.Render(text)
}

func RenderBox(title, content string) string {
	return primaryBoxStyle.Render(
		boldStyle.Render(title+":") + "\n\n" +
			strings.TrimSpace(content),
	)
}

func RenderNote(note string) string {
	if note == "" {
		return ""
	}
	return noteBoxStyle.Render(
		warningStyle.Render("Note: ") + note,
	)
}

func RenderSeparator() string {
	return mutedStyle.Render(strings.Repeat("─", 50))
}

func RenderStep(step string) string {
	return infoStyle.Render(step)
}

func WithSpinner(message string, fn func() error) error {
	var actionErr error

	err := spinner.New().
		Title("🦕 " + message).
		Style(lipgloss.NewStyle().Foreground(PrimaryColor)).
		Type(spinner.Dots).
		Action(func() {
			actionErr = fn()
		}).
		Run()

	if err != nil {
		return err
	}

	return actionErr
}
