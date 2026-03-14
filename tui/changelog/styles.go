package changelog

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/dinoDanic/diny/ui"
)

func indentStyle() lipgloss.Style {
	return lipgloss.NewStyle().PaddingLeft(3)
}

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

func errorStyle() lipgloss.Style {
	t := ui.GetCurrentTheme()
	return lipgloss.NewStyle().Foreground(t.ErrorForeground)
}

func statusSuccessStyle() lipgloss.Style {
	t := ui.GetCurrentTheme()
	return lipgloss.NewStyle().Foreground(t.SuccessForeground)
}

func warningStyle() lipgloss.Style {
	t := ui.GetCurrentTheme()
	return lipgloss.NewStyle().Foreground(t.WarningForeground)
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
