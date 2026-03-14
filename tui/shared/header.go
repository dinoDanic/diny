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
	logoStyle := lipgloss.NewStyle().
		Foreground(t.PrimaryForeground).
		Background(t.PrimaryBackground).
		Padding(0, 1).
		MarginRight(2)

	logo := logoStyle.Render("diny")

	pwd, _ := os.Getwd()

	parts := []string{"v" + version}
	if repoName != "" {
		parts = append(parts, repoName)
	}
	if branchName != "" {
		parts = append(parts, "⎇ "+branchName)
	}
	row1 := metaStyle.Render(strings.Join(parts, "  "))

	metaBlock := strings.Join([]string{row1, metaStyle.Render(pwd), ""}, "\n")
	sep := metaStyle.Render(strings.Repeat("─", width))

	return strings.Join([]string{
		indent.Render(lipgloss.JoinHorizontal(lipgloss.Top, logo, metaBlock)),
		sep,
	}, "\n") + "\n"
}
