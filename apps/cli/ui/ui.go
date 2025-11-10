/*
Package ui provides styled terminal UI components for the diny CLI.

# Quick Start

The easiest way to display styled messages is using the convenience functions:

	ui.SuccessMsg("Operation completed successfully!")
	ui.ErrorMsg("Something went wrong")
	ui.WarningMsg("Please check your configuration")
	ui.InfoMsg("Processing your request...")

For messages with titles:

	ui.SuccessBox("Success", "Your commit has been created")
	ui.ErrorBox("Error", "Failed to connect to server")
	ui.WarningBox("Warning", "Configuration file not found")
	ui.PrimaryBox("Info", "Loading your settings...")

For advanced usage with custom options:

	ui.Box(ui.BoxOptions{
		Title:   "Custom Box",
		Message: "Full control over styling",
		Variant: ui.Success,
	})

Available functions:
  - PrimaryBox(title, message) - Blue/primary styled box with title
  - SuccessBox(title, message) - Green success box with title
  - ErrorBox(title, message)   - Red error box with title
  - WarningBox(title, message) - Orange warning box with title
  - InfoMsg(message)           - Blue/primary message without title
  - SuccessMsg(message)        - Green success message without title
  - ErrorMsg(message)          - Red error message without title
  - WarningMsg(message)        - Orange warning message without title
  - RenderTitle(text)          - Display a primary styled title only
  - WithSpinner(msg, func)     - Run a function with a loading spinner
  - GetHuhPrimaryTheme()       - Get themed huh form styles
*/
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
			Foreground(theme.PrimaryForeground).
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

func PrimaryBox(title, message string) {
	Box(BoxOptions{Title: title, Message: message, Variant: Primary})
}

func SuccessBox(title, message string) {
	Box(BoxOptions{Title: title, Message: message, Variant: Success})
}

func ErrorBox(title, message string) {
	Box(BoxOptions{Title: title, Message: message, Variant: Error})
}

func WarningBox(title, message string) {
	Box(BoxOptions{Title: title, Message: message, Variant: Warning})
}

func InfoMsg(message string) {
	Box(BoxOptions{Message: message, Variant: Primary})
}

func SuccessMsg(message string) {
	Box(BoxOptions{Message: message, Variant: Success})
}

func ErrorMsg(message string) {
	Box(BoxOptions{Message: message, Variant: Error})
}

// WarningMsg displays a warning box with just a message (no title)
func WarningMsg(message string) {
	Box(BoxOptions{Message: message, Variant: Warning})
}

func WithSpinner(message string, fn func() error) error {
	var actionErr error
	theme := GetCurrentTheme()

	titleStyle := lipgloss.NewStyle().
		Foreground(theme.PrimaryForeground).
		Bold(true)

	err := spinner.New().
		Title(titleStyle.Render(message)).
		Type(spinner.Dots).
		Style(lipgloss.NewStyle().Foreground(theme.PrimaryForeground)).
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

func PrintThemeList() {
	themes := []struct {
		name     string
		themeKey string
	}{
		{"Gruvbox Light", "gruvbox-light"},
		{"GitHub Light", "github-light"},
		{"Solarized Light", "solarized-light"},
		{"Flexoki Light", "flexoki-light"},
		{"Everforest Dark", "everforest-dark"},
		{"Solarized Dark", "solarized-dark"},
		{"Monokai", "monokai"},
		{"One Dark", "onedark"},
		{"Gruvbox Dark", "gruvbox-dark"},
		{"Dracula", "dracula"},
		{"Nord", "nord"},
		{"Tokyo Night", "tokyo"},
		{"Catppuccin Mocha", "catppuccin"},
		{"Flexoki Dark", "flexoki-dark"},
	}

	originalTheme := LoadTheme()

	for _, t := range themes {
		SetTheme(t.themeKey)
		theme := GetCurrentTheme()

		themeTitle := t.name
		if t.themeKey == originalTheme {
			themeTitle = t.name + " (current)"
		}

		titleStyle := lipgloss.NewStyle().
			Foreground(theme.PrimaryForeground).
			Bold(true)

		primaryBox := lipgloss.NewStyle().
			Foreground(theme.PrimaryForeground).
			Background(theme.PrimaryBackground).
			BorderLeft(true).
			BorderStyle(lipgloss.ThickBorder()).
			BorderForeground(theme.PrimaryForeground).
			Padding(0, 2)

		successBox := lipgloss.NewStyle().
			Foreground(theme.SuccessForeground).
			Background(theme.SuccessBackground).
			BorderLeft(true).
			BorderStyle(lipgloss.ThickBorder()).
			BorderForeground(theme.SuccessForeground).
			Padding(0, 2)

		errorBox := lipgloss.NewStyle().
			Foreground(theme.ErrorForeground).
			Background(theme.ErrorBackground).
			BorderLeft(true).
			BorderStyle(lipgloss.ThickBorder()).
			BorderForeground(theme.ErrorForeground).
			Padding(0, 2)

		warningBox := lipgloss.NewStyle().
			Foreground(theme.WarningForeground).
			Background(theme.WarningBackground).
			BorderLeft(true).
			BorderStyle(lipgloss.ThickBorder()).
			BorderForeground(theme.WarningForeground).
			Padding(0, 2)

		fmt.Println(titleStyle.Render(themeTitle))

		boxes := lipgloss.JoinHorizontal(
			lipgloss.Top,
			primaryBox.Render("Primary")+"  ",
			successBox.Render("Success")+"  ",
			errorBox.Render("Error")+"  ",
			warningBox.Render("Warning"),
		)

		fmt.Println(boxes)

		separator := lipgloss.NewStyle().
			Foreground(theme.MutedForeground).
			Render(strings.Repeat("─", 60))
		fmt.Println(separator)
		fmt.Println()
	}

	if originalTheme != "" {
		SetTheme(originalTheme)
	}
}
