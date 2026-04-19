package timeline

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

func commitMessageStyle() lipgloss.Style {
	return lipgloss.NewStyle().PaddingLeft(2)
}

func successBigStyle() lipgloss.Style {
	t := ui.GetCurrentTheme()
	return lipgloss.NewStyle().
		Foreground(t.SuccessForeground).
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

func pickerFocusedStyle() lipgloss.Style {
	t := ui.GetCurrentTheme()
	return lipgloss.NewStyle().
		Foreground(t.PrimaryForeground).
		Bold(true).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(t.PrimaryForeground).
		Align(lipgloss.Center).
		Width(8)
}

func pickerNormalStyle() lipgloss.Style {
	t := ui.GetCurrentTheme()
	return lipgloss.NewStyle().
		Foreground(t.MutedForeground).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(t.MutedForeground).
		Align(lipgloss.Center).
		Width(8)
}
