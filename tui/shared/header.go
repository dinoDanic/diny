package shared

import (
	"os"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/dinoDanic/diny/ui"
)

// RenderHeader renders the standard diny TUI header.
// Matches the layout in app/view.go renderHeader().
func RenderHeader(version, repoName, branchName string, width int) string {
	t := ui.GetCurrentTheme()

	indent := lipgloss.NewStyle().PaddingLeft(2)
	metaStyle := lipgloss.NewStyle().Foreground(t.MutedForeground)
	badgeStyle := lipgloss.NewStyle().
		Background(t.PrimaryForeground).
		Foreground(t.PrimaryBackground).
		Padding(0, 1).Bold(true)

	// Row 1: left = badge + version, right = repo ⎇ branch
	left := badgeStyle.Render("diny") + "  " + metaStyle.Render("v"+version)

	var right string
	var parts []string
	if repoName != "" {
		parts = append(parts, repoName)
	}
	if branchName != "" {
		parts = append(parts, "⎇ "+branchName)
	}
	if len(parts) > 0 {
		right = metaStyle.Render(strings.Join(parts, "  "))
	}

	// Fill gap between left and right (accounting for 2-char indent)
	innerWidth := width - 2
	row1 := left
	if right != "" {
		gap := innerWidth - lipgloss.Width(left) - lipgloss.Width(right)
		if gap > 0 {
			row1 = left + strings.Repeat(" ", gap) + right
		}
	}

	pwd, _ := os.Getwd()
	sep := metaStyle.Render(strings.Repeat("─", width))

	return strings.Join([]string{
		indent.Render(row1),
		indent.Render(metaStyle.Render(pwd)),
		sep,
	}, "\n") + "\n"
}
