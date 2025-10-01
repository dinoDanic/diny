package ui

import (
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/huh/spinner"
	"github.com/charmbracelet/lipgloss"
	"golang.org/x/term"
)

var (
	PrimaryForeground = lipgloss.Color("#A78BFA")
	PrimaryBackground = lipgloss.Color("#1E1B2E")
	SuccessForeground = lipgloss.Color("#5FD787")
	SuccessBackground = lipgloss.Color("#1A2820")
	ErrorForeground   = lipgloss.Color("#F87171")
	ErrorBackground   = lipgloss.Color("#2E1E1E")
	WarningForeground = lipgloss.Color("#FACC15")
	WarningBackground = lipgloss.Color("#2E2A1E")
	MutedForeground   = lipgloss.Color("#6C7086")
)

func getTerminalWidth() int {
	width, _, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil || width == 0 {
		return 80
	}
	return width
}

func getTitleStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(PrimaryForeground).
		Bold(true).
		MarginTop(1).
		MarginBottom(1)
}

func getBaseBoxStyle() lipgloss.Style {
	width := getTerminalWidth() - 2
	return lipgloss.NewStyle().
		BorderLeft(true).
		BorderStyle(lipgloss.ThickBorder()).
		Padding(1, 2).
		MarginTop(1).
		MarginBottom(1).
		Width(width)
}

func getPrimaryBoxStyle() lipgloss.Style {
	return getBaseBoxStyle().
		Background(PrimaryBackground).
		BorderForeground(PrimaryForeground)
}

func getErrorBoxStyle() lipgloss.Style {
	return getBaseBoxStyle().
		Background(ErrorBackground).
		Foreground(ErrorForeground).
		BorderForeground(ErrorForeground)
}

func getWarningBoxStyle() lipgloss.Style {
	return getBaseBoxStyle().
		Background(WarningBackground).
		Foreground(WarningForeground).
		BorderForeground(WarningForeground)
}

func getSuccessBoxStyle() lipgloss.Style {
	return getBaseBoxStyle().
		Background(SuccessBackground).
		Foreground(SuccessForeground).
		BorderForeground(SuccessForeground)
}

func RenderTitle(text string) {
	fmt.Println(getTitleStyle().Render("🦕 " + text))
}

func RenderError(text string) {
	fmt.Println(getErrorBoxStyle().Render(strings.TrimSpace(text)))
}

func RenderWarning(text string) {
	fmt.Println(getWarningBoxStyle().Render(strings.TrimSpace(text)))
}

func RenderSuccess(text string) {
	fmt.Println(getSuccessBoxStyle().Render(strings.TrimSpace(text)))
}

func RenderBox(title, content string) {
	innerTitleStyle := getTitleStyle()
	fmt.Println(getPrimaryBoxStyle().Render(
		innerTitleStyle.Render(title) + "\n\n" + content,
	))
}

func WithSpinner(message string, fn func() error) error {
	var actionErr error

	spinnerStyle := lipgloss.NewStyle().
		Foreground(PrimaryForeground).
		MarginTop(1)

	err := spinner.New().
		Title("🦕 " + message).
		Style(spinnerStyle).
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
