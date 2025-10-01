package themes

import (
	"github.com/charmbracelet/lipgloss"
)

func Tokyo() *Theme {
	return &Theme{
		Name:              "Tokyo Night",
		IsDark:            true,
		PrimaryForeground: lipgloss.Color("#7AA2F7"),
		PrimaryBackground: lipgloss.Color("#1A1B26"),
		SuccessForeground: lipgloss.Color("#9ECE6A"),
		SuccessBackground: lipgloss.Color("#1A2820"),
		ErrorForeground:   lipgloss.Color("#F7768E"),
		ErrorBackground:   lipgloss.Color("#2E1E1E"),
		WarningForeground: lipgloss.Color("#E0AF68"),
		WarningBackground: lipgloss.Color("#2E2A1E"),
		MutedForeground:   lipgloss.Color("#565F89"),
	}
}
