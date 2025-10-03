package themes

import "github.com/charmbracelet/lipgloss"

var EverforestDark = Theme{
	Name:              "Everforest Dark",
	IsDark:            true,
	PrimaryForeground: lipgloss.Color("#D3C6AA"),
	PrimaryBackground: lipgloss.Color("#2D353B"),
	SuccessForeground: lipgloss.Color("#2D353B"),
	SuccessBackground: lipgloss.Color("#A7C080"),
	ErrorForeground:   lipgloss.Color("#FFFFFF"),
	ErrorBackground:   lipgloss.Color("#E67E80"),
	WarningForeground: lipgloss.Color("#2D353B"),
	WarningBackground: lipgloss.Color("#DBBC7F"),
	MutedForeground:   lipgloss.Color("#859289"),
}
