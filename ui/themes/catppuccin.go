package themes

import (
	"github.com/charmbracelet/lipgloss"
)

type Theme struct {
	Name              string
	IsDark            bool
	PrimaryForeground lipgloss.Color
	PrimaryBackground lipgloss.Color
	SuccessForeground lipgloss.Color
	SuccessBackground lipgloss.Color
	ErrorForeground   lipgloss.Color
	ErrorBackground   lipgloss.Color
	WarningForeground lipgloss.Color
	WarningBackground lipgloss.Color
	MutedForeground   lipgloss.Color
}

func Catppuccin() *Theme {
	return &Theme{
		Name:              "Catppuccin Mocha",
		IsDark:            true,
		PrimaryForeground: lipgloss.Color("#A78BFA"),
		PrimaryBackground: lipgloss.Color("#1E1B2E"),
		SuccessForeground: lipgloss.Color("#5FD787"),
		SuccessBackground: lipgloss.Color("#1A2820"),
		ErrorForeground:   lipgloss.Color("#F87171"),
		ErrorBackground:   lipgloss.Color("#2E1E1E"),
		WarningForeground: lipgloss.Color("#FACC15"),
		WarningBackground: lipgloss.Color("#2E2A1E"),
		MutedForeground:   lipgloss.Color("#6C7086"),
	}
}
