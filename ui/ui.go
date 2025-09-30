package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/huh/spinner"
	"github.com/charmbracelet/lipgloss"
)

var (
	PrimaryForeground = lipgloss.Color("#726FF2")
	PrimaryBackground = lipgloss.Color("#25253A")
	SuccessForeground = lipgloss.Color("#00FF87")
	SuccessBackground = lipgloss.Color("#1A3A20")
	ErrorForeground   = lipgloss.Color("#FF5F87")
	ErrorBackground   = lipgloss.Color("#3A2025")
	WarningForeground = lipgloss.Color("#FFAF00")
	WarningBackground = lipgloss.Color("#3A3320")
	MutedForeground   = lipgloss.Color("#6C7086")

	titleStyle = lipgloss.NewStyle().
			Foreground(PrimaryForeground).
			Bold(true).
			MarginTop(1).
			MarginBottom(1).
	   	MarginLeft(2)

	baseBoxStyle = lipgloss.NewStyle().
			BorderLeft(true).
			BorderStyle(lipgloss.ThickBorder()).
			Padding(1, 2).
			MarginTop(1).
			MarginBottom(1).
	   	MarginLeft(2)

	primaryBoxStyle = baseBoxStyle.
			Background(PrimaryBackground).
			BorderForeground(PrimaryForeground)

	errorBoxStyle = baseBoxStyle.
			Background(ErrorBackground).
			Foreground(ErrorForeground).
			BorderForeground(ErrorForeground)

	warningBoxStyle = baseBoxStyle.
			Background(WarningBackground).
			Foreground(WarningForeground).
			BorderForeground(WarningForeground)

	successBoxStyle = baseBoxStyle.
			Background(SuccessBackground).
			Foreground(SuccessForeground).
			BorderForeground(SuccessForeground)
)

func RenderTitle(text string) {
	fmt.Println(titleStyle.Render("🦕 " + text))
	fmt.Println()
}

func RenderError(text string) {
	fmt.Println(errorBoxStyle.Render(strings.TrimSpace(text)))
}

func RenderWarning(text string) {
	fmt.Println(warningBoxStyle.Render(strings.TrimSpace(text)))
}

func RenderSuccess(text string) {
	fmt.Println(successBoxStyle.Render(strings.TrimSpace(text)))
}

func RenderBox(title, content string) {
	innerTitleStyle := titleStyle.UnsetMarginTop().UnsetMarginBottom()
	fmt.Println(primaryBoxStyle.Render(
		innerTitleStyle.Render(title) + "\n\n" + content,
	))
}

func WithSpinner(message string, fn func() error) error {
	var actionErr error

	err := spinner.New().
		Title("🦕 " + message).
		Style(lipgloss.NewStyle().Foreground(PrimaryForeground)).
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

func DebugUI() {
	fmt.Println("=== DINY UI DEBUG ===")
	RenderTitle("Sample Title")
	RenderBox("Primary Box", "This is a primary box with some content to demonstrate the styling and border.")
	RenderError("This is an error message to show how errors are displayed with red styling and border.")
	RenderWarning("This is a warning message to show how warnings are displayed with orange styling and border.")
	RenderSuccess("This is a success message to show how success messages are displayed with green styling and border.")
	fmt.Println("=== END DEBUG ===")
}
