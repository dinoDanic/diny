package app

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/dinoDanic/diny/ui"
)

func brandStyle() lipgloss.Style {
	t := ui.GetCurrentTheme()
	return lipgloss.NewStyle().
		Bold(true).
		Foreground(t.PrimaryBackground).
		Background(t.PrimaryForeground).
		Padding(0, 1)
}

func metaStyle() lipgloss.Style {
	t := ui.GetCurrentTheme()
	return lipgloss.NewStyle().
		Foreground(t.MutedForeground)
}

func sectionTitleStyle() lipgloss.Style {
	t := ui.GetCurrentTheme()
	return lipgloss.NewStyle().
		Foreground(t.PrimaryForeground).
		Bold(true)
}

func fileAddedStyle() lipgloss.Style {
	t := ui.GetCurrentTheme()
	return lipgloss.NewStyle().
		Foreground(t.SuccessForeground)
}

func fileModifiedStyle() lipgloss.Style {
	t := ui.GetCurrentTheme()
	return lipgloss.NewStyle().
		Foreground(t.WarningForeground)
}

func fileDeletedStyle() lipgloss.Style {
	t := ui.GetCurrentTheme()
	return lipgloss.NewStyle().
		Foreground(t.ErrorForeground)
}

func fileRenamedStyle() lipgloss.Style {
	t := ui.GetCurrentTheme()
	return lipgloss.NewStyle().
		Foreground(t.PrimaryForeground)
}

func commitMessageStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		PaddingLeft(2)
}

func footerKeyStyle() lipgloss.Style {
	t := ui.GetCurrentTheme()
	return lipgloss.NewStyle().
		Foreground(t.PrimaryForeground).
		Bold(true)
}

func footerDescStyle() lipgloss.Style {
	t := ui.GetCurrentTheme()
	return lipgloss.NewStyle().
		Foreground(t.MutedForeground)
}

func footerGroupStyle() lipgloss.Style {
	t := ui.GetCurrentTheme()
	return lipgloss.NewStyle().
		Foreground(t.MutedForeground).
		Width(10)
}

func statusSuccessStyle() lipgloss.Style {
	t := ui.GetCurrentTheme()
	return lipgloss.NewStyle().
		Foreground(t.SuccessForeground)
}

func statusErrorStyle() lipgloss.Style {
	t := ui.GetCurrentTheme()
	return lipgloss.NewStyle().
		Foreground(t.ErrorForeground)
}

func spinnerStyle() lipgloss.Style {
	t := ui.GetCurrentTheme()
	return lipgloss.NewStyle().
		Foreground(t.PrimaryForeground)
}

func helpBoxStyle() lipgloss.Style {
	t := ui.GetCurrentTheme()
	return lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(t.MutedForeground).
		Padding(1, 2)
}

func errorStyle() lipgloss.Style {
	t := ui.GetCurrentTheme()
	return lipgloss.NewStyle().
		Foreground(t.ErrorForeground)
}

func noStagedStyle() lipgloss.Style {
	t := ui.GetCurrentTheme()
	return lipgloss.NewStyle().
		Foreground(t.WarningForeground)
}

func successBigStyle() lipgloss.Style {
	t := ui.GetCurrentTheme()
	return lipgloss.NewStyle().
		Foreground(t.SuccessForeground).
		Bold(true)
}

func indentStyle() lipgloss.Style {
	return lipgloss.NewStyle().PaddingLeft(3)
}

