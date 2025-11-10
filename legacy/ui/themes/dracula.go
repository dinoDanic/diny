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
		SuccessBackground: lipgloss.Color("#1E2F26"),
		ErrorForeground:   lipgloss.Color("#FF5555"),
		ErrorBackground:   lipgloss.Color("#332A2C"),
		WarningForeground: lipgloss.Color("#F1FA8C"),
		WarningBackground: lipgloss.Color("#333125"),
		MutedForeground:   lipgloss.Color("#6272A4"),
	}
}
