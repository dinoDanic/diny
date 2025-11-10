package main

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
	"github.com/dinoDanic/diny/ui"
)

func main() {
	themes := []struct {
		name     string
		themeKey string
	}{
		{"Catppuccin Mocha", "catppuccin"},
		{"Tokyo Night", "tokyo"},
		{"Nord", "nord"},
		{"Dracula", "dracula"},
		{"Gruvbox Dark", "gruvbox-dark"},
		{"One Dark", "onedark"},
		{"Monokai", "monokai"},
		{"Solarized Dark", "solarized-dark"},
		{"Solarized Light", "solarized-light"},
		{"GitHub Light", "github-light"},
		{"Gruvbox Light", "gruvbox-light"},
	}

	ui.RenderTitle("Available Themes")

	for _, t := range themes {
		ui.SetTheme(t.themeKey)
		theme := ui.GetCurrentTheme()

		nameStyle := lipgloss.NewStyle().
			Foreground(theme.PrimaryForeground).
			Bold(true)

		keyStyle := lipgloss.NewStyle().
			Foreground(theme.MutedForeground)

		primaryBox := lipgloss.NewStyle().
			Background(theme.PrimaryBackground).
			Foreground(theme.PrimaryForeground).
			Padding(0, 2)

		successBox := lipgloss.NewStyle().
			Background(theme.SuccessBackground).
			Foreground(theme.SuccessForeground).
			Padding(0, 2)

		errorBox := lipgloss.NewStyle().
			Background(theme.ErrorBackground).
			Foreground(theme.ErrorForeground).
			Padding(0, 2)

		warningBox := lipgloss.NewStyle().
			Background(theme.WarningBackground).
			Foreground(theme.WarningForeground).
			Padding(0, 2)

		fmt.Printf("\n%s %s\n", nameStyle.Render(theme.Name), keyStyle.Render("("+t.themeKey+")"))
		fmt.Printf("  %s %s %s %s\n",
			primaryBox.Render("Primary"),
			successBox.Render("Success"),
			errorBox.Render("Error"),
			warningBox.Render("Warning"),
		)
	}

	fmt.Println()
}
