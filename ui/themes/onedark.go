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
		SuccessBackground: lipgloss.Color("#232D28"),
		ErrorForeground:   lipgloss.Color("#E06C75"),
		ErrorBackground:   lipgloss.Color("#2D2528"),
		WarningForeground: lipgloss.Color("#E5C07B"),
		WarningBackground: lipgloss.Color("#2F2C26"),
		MutedForeground:   lipgloss.Color("#5C6370"),
	}
}
