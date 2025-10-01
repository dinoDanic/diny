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

func getTerminalWidth() int {
	width, _, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil || width == 0 {
		return 80
	}
	return width
}

func getBaseBoxStyle() lipgloss.Style {
	width := getTerminalWidth() - 1
	return lipgloss.NewStyle().
		BorderLeft(true).
		BorderStyle(lipgloss.ThickBorder()).
		Padding(1, 2).
		Width(width)
}

func getBoxStyleByVariant(variant BoxVariant) lipgloss.Style {
	base := getBaseBoxStyle()
	theme := GetCurrentTheme()

	switch variant {
	case Success:
		return base.
			Background(theme.SuccessBackground).
			Foreground(theme.SuccessForeground).
			BorderForeground(theme.SuccessForeground)
	case Error:
		return base.
			Background(theme.ErrorBackground).
			Foreground(theme.ErrorForeground).
			BorderForeground(theme.ErrorForeground)
	case Warning:
		return base.
			Background(theme.WarningBackground).
			Foreground(theme.WarningForeground).
			BorderForeground(theme.WarningForeground)
	case Primary:
		fallthrough
	default:
		return base.
			Background(theme.PrimaryBackground).
			BorderForeground(theme.PrimaryForeground)
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
	theme := GetCurrentTheme()

	switch variant {
	case Success:
		return base.Foreground(theme.SuccessForeground)
	case Error:
		return base.Foreground(theme.ErrorForeground)
	case Warning:
		return base.Foreground(theme.WarningForeground)
	case Primary:
		fallthrough
	default:
		return base.Foreground(theme.PrimaryForeground)
	}
}

func RenderTitle(text string) {
	Box(BoxOptions{Title: text})
}

func WithSpinner(message string, fn func() error) error {
	var actionErr error
	theme := GetCurrentTheme()

	spinnerStyle := lipgloss.NewStyle().
		Foreground(theme.PrimaryForeground)

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
	t := GetCurrentTheme()
	theme := huh.ThemeBase()
	theme.Focused.Base = theme.Focused.Base.
		BorderForeground(t.PrimaryForeground).PaddingTop(1).PaddingBottom(1).Bold(true)
	theme.Focused.Title = theme.Focused.Title.Foreground(t.PrimaryForeground)
	theme.Focused.Description = theme.Focused.Description.Foreground(t.MutedForeground)
	theme.Focused.SelectedOption = theme.Focused.SelectedOption.
		Foreground(t.PrimaryForeground).Bold(true)
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
