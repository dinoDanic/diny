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
	logoStyle := lipgloss.NewStyle().Foreground(t.PrimaryForeground).PaddingRight(2)

	logo := logoStyle.Render("▗▄▄▄▄▖\n▐████▌\n▝▀▀▀▀▘")
	logoW := lipgloss.Width(logoStyle.Render("▗▄▄▄▄▖"))

	pwd, _ := os.Getwd()
	metaAvail := width - 2 - logoW

	var parts []string
	if repoName != "" {
		parts = append(parts, repoName)
	}
	if branchName != "" {
		parts = append(parts, "⎇ "+branchName)
	}

	versionText := metaStyle.Render("v" + version)
	row1 := versionText
	if len(parts) > 0 {
		rightText := metaStyle.Render(strings.Join(parts, "  "))
		gap := metaAvail - lipgloss.Width(versionText) - lipgloss.Width(rightText)
		if gap > 0 {
			row1 = versionText + strings.Repeat(" ", gap) + rightText
		}
	}

	metaBlock := strings.Join([]string{row1, metaStyle.Render(pwd), ""}, "\n")
	sep := metaStyle.Render(strings.Repeat("─", width))

	return strings.Join([]string{
		indent.Render(lipgloss.JoinHorizontal(lipgloss.Top, logo, metaBlock)),
		sep,
	}, "\n") + "\n"
}
