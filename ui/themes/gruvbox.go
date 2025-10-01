package themes

import (
	"github.com/charmbracelet/lipgloss"
)

func GruvboxDark() *Theme {
	return &Theme{
		Name:              "Gruvbox Dark",
		IsDark:            true,
		PrimaryForeground: lipgloss.Color("#D3869B"),
		PrimaryBackground: lipgloss.Color("#282828"),
		SuccessForeground: lipgloss.Color("#B8BB26"),
		SuccessBackground: lipgloss.Color("#1A2820"),
		ErrorForeground:   lipgloss.Color("#FB4934"),
		ErrorBackground:   lipgloss.Color("#2E1E1E"),
		WarningForeground: lipgloss.Color("#FABD2F"),
		WarningBackground: lipgloss.Color("#2E2A1E"),
		MutedForeground:   lipgloss.Color("#928374"),
	}
}

func GruvboxLight() *Theme {
	return &Theme{
		Name:              "Gruvbox Light",
		IsDark:            false,
		PrimaryForeground: lipgloss.Color("#B16286"),
		PrimaryBackground: lipgloss.Color("#FBF1C7"),
		SuccessForeground: lipgloss.Color("#79740E"),
		SuccessBackground: lipgloss.Color("#D5C4A1"),
		ErrorForeground:   lipgloss.Color("#CC241D"),
		ErrorBackground:   lipgloss.Color("#EBDBB2"),
		WarningForeground: lipgloss.Color("#D79921"),
		WarningBackground: lipgloss.Color("#EBDBB2"),
		MutedForeground:   lipgloss.Color("#7C6F64"),
	}
}
