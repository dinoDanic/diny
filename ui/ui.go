package ui

import (
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/huh/spinner"
	"github.com/charmbracelet/lipgloss"
	"golang.org/x/term"
)

type BoxVariant string

const (
	Primary BoxVariant = "primary"
	Success BoxVariant = "success"
	Error   BoxVariant = "error"
	Warning BoxVariant = "warning"
)

type BoxOptions struct {
	Title   string
	Message string
	Variant BoxVariant
}

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
	width := getTerminalWidth() - 1
	return lipgloss.NewStyle().
		BorderLeft(true).
		BorderStyle(lipgloss.ThickBorder()).
		Padding(1, 2).
		// MarginTop(1).
		// MarginBottom(1).
		Width(width)
}

func getBoxStyleByVariant(variant BoxVariant) lipgloss.Style {
	base := getBaseBoxStyle()

	switch variant {
	case Success:
		return base.
			Background(SuccessBackground).
			Foreground(SuccessForeground).
			BorderForeground(SuccessForeground)
	case Error:
		return base.
			Background(ErrorBackground).
			Foreground(ErrorForeground).
			BorderForeground(ErrorForeground)
	case Warning:
		return base.
			Background(WarningBackground).
			Foreground(WarningForeground).
			BorderForeground(WarningForeground)
	case Primary:
		fallthrough
	default:
		return base.
			Background(PrimaryBackground).
			BorderForeground(PrimaryForeground)
	}
}

func Box(opts BoxOptions) {
	if opts.Variant == "" {
		opts.Variant = Primary
	}

	style := getBoxStyleByVariant(opts.Variant)

	var content string
	if opts.Title != "" && opts.Message != "" {
		titleStyle := getTitleStyleByVariant(opts.Variant)
		content = titleStyle.Render(opts.Title) + "\n\n" + strings.TrimSpace(opts.Message)
	} else if opts.Title != "" {
		titleStyle := getTitleStyleByVariant(opts.Variant)
		content = titleStyle.Render(opts.Title)
	} else if opts.Message != "" {
		content = strings.TrimSpace(opts.Message)
	}

	if content != "" {
		fmt.Println(style.Render(content))
	}
}

func getTitleStyleByVariant(variant BoxVariant) lipgloss.Style {
	base := lipgloss.NewStyle().Bold(true)

	switch variant {
	case Success:
		return base.Foreground(SuccessForeground)
	case Error:
		return base.Foreground(ErrorForeground)
	case Warning:
		return base.Foreground(WarningForeground)
	case Primary:
		fallthrough
	default:
		return base.Foreground(PrimaryForeground)
	}
}

func RenderTitle(text string) {
	Box(BoxOptions{Title: "🦕 " + text})
}

func WithSpinner(message string, fn func() error) error {
	var actionErr error

	spinnerStyle := lipgloss.NewStyle().
		Foreground(PrimaryForeground)

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

func GetHuhPrimaryTheme() *huh.Theme {
	theme := huh.ThemeBase()
	theme.Focused.Base = theme.Focused.Base.
		BorderForeground(PrimaryForeground).
		Background(PrimaryBackground)
	theme.Focused.Title = theme.Focused.Title.Foreground(PrimaryForeground).Bold(true)
	theme.Focused.Description = theme.Focused.Description.Foreground(MutedForeground)
	theme.Focused.SelectedOption = theme.Focused.SelectedOption.
		Foreground(PrimaryForeground).
		Background(PrimaryBackground)
	theme.Focused.TextInput.Cursor = theme.Focused.TextInput.Cursor.Foreground(PrimaryForeground)
	theme.Focused.TextInput.Prompt = theme.Focused.TextInput.Prompt.
		Foreground(PrimaryForeground).
		Background(PrimaryBackground)
	theme.Focused.TextInput.Placeholder = theme.Focused.TextInput.Placeholder.Foreground(MutedForeground)
	theme.Focused.TextInput.Text = theme.Focused.TextInput.Text.Background(PrimaryBackground)
	return theme
}

func DebugUI() {
	fmt.Println("=== DINY UI DEBUG ===")
	RenderTitle("Sample Title")
	Box(BoxOptions{Title: "Primary Box", Message: "This is a primary box with some content to demonstrate the styling and border.", Variant: Primary})
	Box(BoxOptions{Title: "Error Box", Message: "This is an error message to show how errors are displayed with red styling and border.", Variant: Error})
	Box(BoxOptions{Title: "Warning Box", Message: "This is a warning message to show how warnings are displayed with orange styling and border.", Variant: Warning})
	Box(BoxOptions{Title: "Success Box", Message: "This is a success message to show how success messages are displayed with green styling and border.", Variant: Success})
	fmt.Println("=== END DEBUG ===")
}
