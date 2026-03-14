package configtui

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/dinoDanic/diny/ui"
)

func metaStyle() lipgloss.Style {
	t := ui.GetCurrentTheme()
	return lipgloss.NewStyle().Foreground(t.MutedForeground)
}

func sectionTitleStyle() lipgloss.Style {
	t := ui.GetCurrentTheme()
	return lipgloss.NewStyle().
		Foreground(t.PrimaryForeground).
		Bold(true)
}

func footerKeyStyle() lipgloss.Style {
	t := ui.GetCurrentTheme()
	return lipgloss.NewStyle().
		Foreground(t.PrimaryForeground).
		Bold(true)
}

func footerDescStyle() lipgloss.Style {
	t := ui.GetCurrentTheme()
	return lipgloss.NewStyle().Foreground(t.MutedForeground)
}

func statusSuccessStyle() lipgloss.Style {
	t := ui.GetCurrentTheme()
	return lipgloss.NewStyle().Foreground(t.SuccessForeground)
}

func statusErrorStyle() lipgloss.Style {
	t := ui.GetCurrentTheme()
	return lipgloss.NewStyle().Foreground(t.ErrorForeground)
}

func indentStyle() lipgloss.Style {
	return lipgloss.NewStyle().PaddingLeft(3)
}
