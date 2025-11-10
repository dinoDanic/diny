package themes

import (
	"github.com/charmbracelet/lipgloss"
)

func SolarizedDark() *Theme {
	return &Theme{
		Name:              "Solarized Dark",
		IsDark:            true,
		PrimaryForeground: lipgloss.Color("#268BD2"),
		PrimaryBackground: lipgloss.Color("#002B36"),
		SuccessForeground: lipgloss.Color("#859900"),
		SuccessBackground: lipgloss.Color("#073642"),
		ErrorForeground:   lipgloss.Color("#DC322F"),
		ErrorBackground:   lipgloss.Color("#073642"),
		WarningForeground: lipgloss.Color("#B58900"),
		WarningBackground: lipgloss.Color("#073642"),
		MutedForeground:   lipgloss.Color("#586E75"),
	}
}

func SolarizedLight() *Theme {
	return &Theme{
		Name:              "Solarized Light",
		IsDark:            false,
		PrimaryForeground: lipgloss.Color("#268BD2"),
		PrimaryBackground: lipgloss.Color("#FDF6E3"),
		SuccessForeground: lipgloss.Color("#859900"),
		SuccessBackground: lipgloss.Color("#EEE8D5"),
		ErrorForeground:   lipgloss.Color("#DC322F"),
		ErrorBackground:   lipgloss.Color("#EEE8D5"),
		WarningForeground: lipgloss.Color("#B58900"),
		WarningBackground: lipgloss.Color("#EEE8D5"),
		MutedForeground:   lipgloss.Color("#93A1A1"),
	}
}
