package themes

import "github.com/charmbracelet/lipgloss"

var EverforestDark = Theme{
	Name:              "Everforest Dark",
	IsDark:            true,
	PrimaryForeground: lipgloss.Color("#D3C6AA"),
	PrimaryBackground: lipgloss.Color("#2D353B"),
	SuccessForeground: lipgloss.Color("#A7C080"),
	SuccessBackground: lipgloss.Color("#2D353B"),
	ErrorForeground:   lipgloss.Color("#E67E80"),
	ErrorBackground:   lipgloss.Color("#2D353B"),
	WarningForeground: lipgloss.Color("#DBBC7F"),
	WarningBackground: lipgloss.Color("#2D353B"),
	MutedForeground:   lipgloss.Color("#859289"),
}
