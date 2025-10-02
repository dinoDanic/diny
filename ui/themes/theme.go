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
