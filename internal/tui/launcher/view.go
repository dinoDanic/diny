package launcher

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/dinoDanic/diny/ui"
)

func (m model) View() string {
	if m.width == 0 {
		return ""
	}

	left := m.renderLeftPanel()
	right := m.renderRightPanel()

	return lipgloss.JoinHorizontal(lipgloss.Top, left, right)
}

func (m model) renderLeftPanel() string {
	t := ui.GetCurrentTheme()

	headerLabel := lipgloss.NewStyle().
		Bold(true).
		Foreground(t.PrimaryBackground).
		Background(t.PrimaryForeground).
		Padding(0, 1).
		Render("diny")

	var itemLines []string
	for i, it := range m.items {
		prefix := "  "
		style := lipgloss.NewStyle().Foreground(t.MutedForeground)
		if m.activeIndex == i {
			prefix = lipgloss.NewStyle().Foreground(t.PrimaryForeground).Render("→ ")
			style = lipgloss.NewStyle().Bold(true).Foreground(t.PrimaryForeground)
		} else if m.cursor == i {
			prefix = lipgloss.NewStyle().Foreground(t.PrimaryForeground).Render("▶ ")
			style = lipgloss.NewStyle().Bold(true).Foreground(t.PrimaryForeground)
		}
		itemLines = append(itemLines, prefix+style.Render(it.title))
	}

	fk := lipgloss.NewStyle().Bold(true).Foreground(t.PrimaryForeground)
	fd := lipgloss.NewStyle().Foreground(t.MutedForeground)
	footer := fk.Render("↑↓") + fd.Render(" nav  ") + fk.Render("q") + fd.Render(" quit")

	content := lipgloss.JoinVertical(lipgloss.Left,
		headerLabel,
		"",
		strings.Join(itemLines, "\n"),
		"",
		footer,
	)

	return lipgloss.NewStyle().
		Width(leftPanelWidth).
		Height(m.height).
		BorderRight(true).
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(t.MutedForeground).
		Padding(1, 1).
		Render(content)
}

func (m model) renderRightPanel() string {
	rw := m.rightWidth()

	if m.subModel != nil {
		return m.subModel.View()
	}

	// Welcome screen: show description of the item under the cursor.
	t := ui.GetCurrentTheme()
	if m.cursor < 0 || m.cursor >= len(m.items) {
		return ""
	}
	it := m.items[m.cursor]
	title := lipgloss.NewStyle().Bold(true).Foreground(t.PrimaryForeground).Render(it.title)
	desc := lipgloss.NewStyle().Foreground(t.MutedForeground).Render(it.description)
	hint := lipgloss.NewStyle().Foreground(t.MutedForeground).Render("Press enter to open")

	content := lipgloss.NewStyle().Padding(0, 2).Render(
		lipgloss.JoinVertical(lipgloss.Left, title, desc, "", hint),
	)
	return lipgloss.Place(rw, m.height, lipgloss.Left, lipgloss.Center, content)
}
