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
		SuccessBackground: lipgloss.Color("#2E3A35"),
		ErrorForeground:   lipgloss.Color("#BF616A"),
		ErrorBackground:   lipgloss.Color("#3B3537"),
		WarningForeground: lipgloss.Color("#EBCB8B"),
		WarningBackground: lipgloss.Color("#3B3A35"),
		MutedForeground:   lipgloss.Color("#4C566A"),
	}
}
