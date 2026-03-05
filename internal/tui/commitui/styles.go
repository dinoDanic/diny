package commitui

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/dinoDanic/diny/ui"
)

func getTheme() lipgloss.Style {
	return lipgloss.NewStyle()
}

func getSpinnerStyle() lipgloss.Style {
	t := ui.GetCurrentTheme()
	return lipgloss.NewStyle().Foreground(t.PrimaryForeground)
}

func headerStyle() lipgloss.Style {
	t := ui.GetCurrentTheme()
	return lipgloss.NewStyle().
		Bold(true).
		Foreground(t.PrimaryForeground).
		Padding(0, 1)
}

func headerLabelStyle() lipgloss.Style {
	t := ui.GetCurrentTheme()
	return lipgloss.NewStyle().
		Bold(true).
		Foreground(t.PrimaryBackground).
		Background(t.PrimaryForeground).
		Padding(0, 1)
}

func headerCountStyle() lipgloss.Style {
	t := ui.GetCurrentTheme()
	return lipgloss.NewStyle().
		Foreground(t.MutedForeground).
		Padding(0, 1)
}

func filePaneStyle(width, height int) lipgloss.Style {
	t := ui.GetCurrentTheme()
	return lipgloss.NewStyle().
		Width(width).
		Height(height).
		BorderRight(true).
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(t.MutedForeground).
		Padding(0, 1)
}

func messagePaneStyle(width, height int) lipgloss.Style {
	return lipgloss.NewStyle().
		Width(width).
		Height(height).
		Padding(0, 1)
}

func paneTitleStyle() lipgloss.Style {
	t := ui.GetCurrentTheme()
	return lipgloss.NewStyle().
		Bold(true).
		Foreground(t.PrimaryForeground).
		MarginBottom(1)
}

func footerStyle(width int) lipgloss.Style {
	t := ui.GetCurrentTheme()
	return lipgloss.NewStyle().
		Width(width).
		Foreground(t.MutedForeground).
		Padding(0, 1)
}

func footerKeyStyle() lipgloss.Style {
	t := ui.GetCurrentTheme()
	return lipgloss.NewStyle().
		Bold(true).
		Foreground(t.PrimaryForeground)
}

func footerDescStyle() lipgloss.Style {
	t := ui.GetCurrentTheme()
	return lipgloss.NewStyle().
		Foreground(t.MutedForeground)
}

func footerSepStyle() lipgloss.Style {
	t := ui.GetCurrentTheme()
	return lipgloss.NewStyle().
		Foreground(t.MutedForeground)
}

func statusSuccessStyle() lipgloss.Style {
	t := ui.GetCurrentTheme()
	return lipgloss.NewStyle().
		Foreground(t.SuccessForeground).
		Padding(0, 1)
}

func statusErrorStyle() lipgloss.Style {
	t := ui.GetCurrentTheme()
	return lipgloss.NewStyle().
		Foreground(t.ErrorForeground).
		Padding(0, 1)
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

func filePathStyle() lipgloss.Style {
	return lipgloss.NewStyle()
}

func helpOverlayStyle() lipgloss.Style {
	t := ui.GetCurrentTheme()
	return lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(t.PrimaryForeground).
		Padding(1, 2).
		Foreground(t.PrimaryForeground)
}

func helpKeyStyle() lipgloss.Style {
	t := ui.GetCurrentTheme()
	return lipgloss.NewStyle().
		Bold(true).
		Foreground(t.PrimaryForeground).
		Width(10)
}

func helpDescStyle() lipgloss.Style {
	t := ui.GetCurrentTheme()
	return lipgloss.NewStyle().
		Foreground(t.MutedForeground)
}

func variantActiveStyle() lipgloss.Style {
	t := ui.GetCurrentTheme()
	return lipgloss.NewStyle().
		Bold(true).
		Foreground(t.PrimaryForeground)
}

func variantInactiveStyle() lipgloss.Style {
	t := ui.GetCurrentTheme()
	return lipgloss.NewStyle().
		Foreground(t.MutedForeground)
}

func loadingStyle() lipgloss.Style {
	t := ui.GetCurrentTheme()
	return lipgloss.NewStyle().
		Foreground(t.PrimaryForeground).
		Bold(true).
		Padding(0, 1)
}
