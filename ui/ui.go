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

	for _, t := range themes {
		SetTheme(t.themeKey)
		theme := GetCurrentTheme()

		themeTitle := t.name

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

	SetTheme("catppuccin")
	Box(BoxOptions{
		Message: "Set theme in config file. Open with: diny config",
		Variant: Primary,
	})
}
