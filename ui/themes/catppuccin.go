package themes

import (
	"github.com/charmbracelet/lipgloss"
)

func Catppuccin() *Theme {
	return &Theme{
		Name:              "Catppuccin Mocha",
		IsDark:            true,
		PrimaryForeground: lipgloss.Color("#CBA6F7"),
		PrimaryBackground: lipgloss.Color("#1E1E2E"),
		SuccessForeground: lipgloss.Color("#A6E3A1"),
		SuccessBackground: lipgloss.Color("#1E2D26"),
		ErrorForeground:   lipgloss.Color("#F38BA8"),
		ErrorBackground:   lipgloss.Color("#2D1E24"),
		WarningForeground: lipgloss.Color("#F9E2AF"),
		WarningBackground: lipgloss.Color("#2D2A1E"),
		MutedForeground:   lipgloss.Color("#6C7086"),
	}
}
