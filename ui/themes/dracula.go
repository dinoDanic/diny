package themes

import (
	"github.com/charmbracelet/lipgloss"
)

func Dracula() *Theme {
	return &Theme{
		Name:              "Dracula",
		IsDark:            true,
		PrimaryForeground: lipgloss.Color("#BD93F9"),
		PrimaryBackground: lipgloss.Color("#282A36"),
		SuccessForeground: lipgloss.Color("#50FA7B"),
		SuccessBackground: lipgloss.Color("#1A2820"),
		ErrorForeground:   lipgloss.Color("#FF5555"),
		ErrorBackground:   lipgloss.Color("#2E1E1E"),
		WarningForeground: lipgloss.Color("#F1FA8C"),
		WarningBackground: lipgloss.Color("#2E2A1E"),
		MutedForeground:   lipgloss.Color("#6272A4"),
	}
}
