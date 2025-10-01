package themes

import (
	"github.com/charmbracelet/lipgloss"
)

func OneDark() *Theme {
	return &Theme{
		Name:              "One Dark",
		IsDark:            true,
		PrimaryForeground: lipgloss.Color("#C678DD"),
		PrimaryBackground: lipgloss.Color("#282C34"),
		SuccessForeground: lipgloss.Color("#98C379"),
		SuccessBackground: lipgloss.Color("#1A2820"),
		ErrorForeground:   lipgloss.Color("#E06C75"),
		ErrorBackground:   lipgloss.Color("#2E1E1E"),
		WarningForeground: lipgloss.Color("#E5C07B"),
		WarningBackground: lipgloss.Color("#2E2A1E"),
		MutedForeground:   lipgloss.Color("#5C6370"),
	}
}
