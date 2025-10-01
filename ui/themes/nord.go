package themes

import (
	"github.com/charmbracelet/lipgloss"
)

func Nord() *Theme {
	return &Theme{
		Name:              "Nord",
		IsDark:            true,
		PrimaryForeground: lipgloss.Color("#88C0D0"),
		PrimaryBackground: lipgloss.Color("#2E3440"),
		SuccessForeground: lipgloss.Color("#A3BE8C"),
		SuccessBackground: lipgloss.Color("#1A2820"),
		ErrorForeground:   lipgloss.Color("#BF616A"),
		ErrorBackground:   lipgloss.Color("#2E1E1E"),
		WarningForeground: lipgloss.Color("#EBCB8B"),
		WarningBackground: lipgloss.Color("#2E2A1E"),
		MutedForeground:   lipgloss.Color("#4C566A"),
	}
}
