package themes

import (
	"github.com/charmbracelet/lipgloss"
)

func Monokai() *Theme {
	return &Theme{
		Name:              "Monokai",
		IsDark:            true,
		PrimaryForeground: lipgloss.Color("#AE81FF"),
		PrimaryBackground: lipgloss.Color("#272822"),
		SuccessForeground: lipgloss.Color("#A6E22E"),
		SuccessBackground: lipgloss.Color("#232D25"),
		ErrorForeground:   lipgloss.Color("#F92672"),
		ErrorBackground:   lipgloss.Color("#2D2225"),
		WarningForeground: lipgloss.Color("#E6DB74"),
		WarningBackground: lipgloss.Color("#2E2D23"),
		MutedForeground:   lipgloss.Color("#75715E"),
	}
}
