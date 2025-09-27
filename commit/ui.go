package commit

import (
	"strings"

	"github.com/charmbracelet/huh/spinner"
	"github.com/charmbracelet/lipgloss"
)

var (
	// Color scheme
	primaryColor    = lipgloss.Color("#726FF2")
	successColor    = lipgloss.Color("#00FF87")
	errorColor      = lipgloss.Color("#FF5F87")
	warningColor    = lipgloss.Color("#FFAF00")
	mutedColor      = lipgloss.Color("#6C7086")
	backgroundColor = lipgloss.Color("#1E1E2E")

	// Base styles
	baseStyle = lipgloss.NewStyle().
			Padding(0, 1)

	// Header styles
	titleStyle = lipgloss.NewStyle().
			Foreground(primaryColor).
			Bold(true).
			Padding(0, 1).
			MarginBottom(1)

	// Message styles
	successStyle = lipgloss.NewStyle().
			Foreground(successColor).
			Bold(true).
			Padding(0, 1)

	errorStyle = lipgloss.NewStyle().
			Foreground(errorColor).
			Bold(true).
			Padding(0, 1)

	warningStyle = lipgloss.NewStyle().
			Foreground(warningColor).
			Bold(true).
			Padding(0, 1)

	infoStyle = lipgloss.NewStyle().
			Foreground(primaryColor).
			Padding(0, 1)

	// Box styles
	commitBoxStyle = lipgloss.NewStyle().
			Background(lipgloss.Color("#25253A")).
			Padding(1, 2)

	configBoxStyle = lipgloss.NewStyle().
			Border(lipgloss.NormalBorder()).
			BorderForeground(mutedColor).
			Padding(1, 2).
			MarginBottom(1)

	noteBoxStyle = lipgloss.NewStyle().
			Background(lipgloss.Color("#3A3320")).
			Padding(1, 2)

	errorBoxStyle = lipgloss.NewStyle().
			Background(lipgloss.Color("#3A2025")).
			Padding(1, 2)

	successBoxStyle = lipgloss.NewStyle().
			Background(lipgloss.Color("#1A3A20")).
			Padding(1, 2)

	warningBoxStyle = lipgloss.NewStyle().
			Background(lipgloss.Color("#3A3320")).
			Padding(1, 2)

	// Text styles
	mutedStyle = lipgloss.NewStyle().
			Foreground(mutedColor).
			Italic(true)

	boldStyle = lipgloss.NewStyle().
			Bold(true)
)

// UI component functions
func RenderTitle(text string) string {
	return titleStyle.Render("🦕 " + text)
}

func RenderSuccess(text string) string {
	return successBoxStyle.Render(successStyle.Render(text))
}

func RenderError(text string) string {
	return errorBoxStyle.Render(errorStyle.Render(text))
}

func RenderWarning(text string) string {
	return warningBoxStyle.Render(warningStyle.Render(text))
}

func RenderInfo(text string) string {
	return infoStyle.Render(text)
}

func RenderCommitMessage(message string) string {
	return commitBoxStyle.Render(
		boldStyle.Render("Commit Message:") + "\n\n" +
			strings.TrimSpace(message),
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

func RenderConfigBox(content string) string {
	return configBoxStyle.Render(
		mutedStyle.Render("Configuration:") + "\n" + content,
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
		Style(lipgloss.NewStyle().Foreground(primaryColor)).
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
