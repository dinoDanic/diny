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

	indent := lipgloss.NewStyle().PaddingLeft(3)
	metaStyle := lipgloss.NewStyle().Foreground(t.MutedForeground)
	sectionTitleStyle := lipgloss.NewStyle().Foreground(t.PrimaryForeground).Bold(true)
	primaryStyle := lipgloss.NewStyle().Foreground(t.PrimaryForeground)

	line1 := sectionTitleStyle.Render("diny") + "  " + metaStyle.Render("v"+version)

	var infoParts []string
	if repoName != "" {
		infoParts = append(infoParts, repoName)
	}
	if branchName != "" {
		infoParts = append(infoParts, branchName)
	}

	pwd, _ := os.Getwd()

	var rows []string
	rows = append(rows, indent.Render(line1))
	if len(infoParts) > 0 {
		rows = append(rows, indent.Render(primaryStyle.Render(strings.Join(infoParts, " • "))))
	}
	rows = append(rows, indent.Render(metaStyle.Render(pwd)))

	sep := metaStyle.Render(strings.Repeat("─", width))
	rows = append(rows, sep, )

	return strings.Join(rows, "\n") + "\n" 
}
