package themes

import (
	"github.com/charmbracelet/lipgloss"
)

func GithubLight() *Theme {
	return &Theme{
		Name:              "GitHub Light",
		IsDark:            false,
		PrimaryForeground: lipgloss.Color("#0969DA"),
		PrimaryBackground: lipgloss.Color("#FFFFFF"),
		SuccessForeground: lipgloss.Color("#1A7F37"),
		SuccessBackground: lipgloss.Color("#DDF4FF"),
		ErrorForeground:   lipgloss.Color("#CF222E"),
		ErrorBackground:   lipgloss.Color("#FFEBE9"),
		WarningForeground: lipgloss.Color("#9A6700"),
		WarningBackground: lipgloss.Color("#FFF8C5"),
		MutedForeground:   lipgloss.Color("#57606A"),
	}
}
