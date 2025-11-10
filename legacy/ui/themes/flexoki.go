package themes

import (
	"github.com/charmbracelet/lipgloss"
)

func FlexokiDark() *Theme {
	return &Theme{
		Name:              "Flexoki Dark",
		IsDark:            true,
		PrimaryForeground: lipgloss.Color("#8B7EC8"),
		PrimaryBackground: lipgloss.Color("#282726"),
		SuccessForeground: lipgloss.Color("#879A39"),
		SuccessBackground: lipgloss.Color("#282726"),
		ErrorForeground:   lipgloss.Color("#D14D41"),
		ErrorBackground:   lipgloss.Color("#282726"),
		WarningForeground: lipgloss.Color("#D0A215"),
		WarningBackground: lipgloss.Color("#282726"),
		MutedForeground:   lipgloss.Color("#6F6E69"),
	}
}

func FlexokiLight() *Theme {
	return &Theme{
		Name:              "Flexoki Light",
		IsDark:            false,
		PrimaryForeground: lipgloss.Color("#5E409D"),
		PrimaryBackground: lipgloss.Color("#FFFCF0"),
		SuccessForeground: lipgloss.Color("#66800B"),
		SuccessBackground: lipgloss.Color("#F2F0E5"),
		ErrorForeground:   lipgloss.Color("#AF3029"),
		ErrorBackground:   lipgloss.Color("#F2F0E5"),
		WarningForeground: lipgloss.Color("#AD8301"),
		WarningBackground: lipgloss.Color("#F2F0E5"),
		MutedForeground:   lipgloss.Color("#9F9D96"),
	}
}
