package ui

import (
	"fmt"
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
			Bold(true)

	primaryBoxStyle = lipgloss.NewStyle().
			Background(lipgloss.Color("#25253A")).
			Padding(1, 3).
			MarginBottom(1)

	errorBoxStyle = lipgloss.NewStyle().
			Background(lipgloss.Color("#3A2025")).
			Foreground(ErrorColor).
			Bold(true).
			Padding(1, 3).
			MarginBottom(1)

	warningBoxStyle = lipgloss.NewStyle().
			Background(lipgloss.Color("#3A3320")).
			Foreground(WarningColor).
			Bold(true).
			Padding(1, 3).
			MarginBottom(1)
)

func RenderTitle(text string) {
	fmt.Println(titleStyle.Render("🦕 " + text))
}

func RenderError(text string) {
	fmt.Println(errorBoxStyle.Render(strings.TrimSpace(text)))
}

func RenderWarning(text string) {
	fmt.Println(warningBoxStyle.Render(strings.TrimSpace(text)))
}

func RenderBox(title, content string) {
	fmt.Println(primaryBoxStyle.Render(
		titleStyle.Render(title) + "\n\n" + content,
	))
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
