package commit

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

var (
	// Color scheme
	primaryColor    = lipgloss.Color("#00D7FF")
	successColor    = lipgloss.Color("#00FF87")
	errorColor      = lipgloss.Color("#FF5F87")
	warningColor    = lipgloss.Color("#FFAF00")
	mutedColor      = lipgloss.Color("#6C7086")
	backgroundColor = lipgloss.Color("#1E1E2E")

	// Base styles
	baseStyle = lipgloss.NewStyle().
			Padding(0, 1)

	// Header styles
	titleStyle = lipgloss.NewStyle().
			Foreground(primaryColor).
			Bold(true).
			Padding(0, 1).
			MarginBottom(1)

	// Message styles
	successStyle = lipgloss.NewStyle().
			Foreground(successColor).
			Bold(true).
			Padding(0, 1)

	errorStyle = lipgloss.NewStyle().
			Foreground(errorColor).
			Bold(true).
			Padding(0, 1)

	warningStyle = lipgloss.NewStyle().
			Foreground(warningColor).
			Bold(true).
			Padding(0, 1)

	infoStyle = lipgloss.NewStyle().
			Foreground(primaryColor).
			Padding(0, 1)

	// Box styles
	commitBoxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(primaryColor).
			Padding(1, 2).
			MarginTop(1).
			MarginBottom(1)

	configBoxStyle = lipgloss.NewStyle().
			Border(lipgloss.NormalBorder()).
			BorderForeground(mutedColor).
			Padding(1, 2).
			MarginBottom(1)

	noteBoxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(warningColor).
			Padding(1, 2).
			MarginBottom(1)

	// Text styles
	mutedStyle = lipgloss.NewStyle().
			Foreground(mutedColor).
			Italic(true)

	boldStyle = lipgloss.NewStyle().
			Bold(true)
)

// UI component functions
func RenderTitle(text string) string {
	return titleStyle.Render("🦕 " + text)
}

func RenderSuccess(text string) string {
	return successStyle.Render("✅ " + text)
}

func RenderError(text string) string {
	return errorStyle.Render("❌ " + text)
}

func RenderWarning(text string) string {
	return warningStyle.Render("⚠️  " + text)
}

func RenderInfo(text string) string {
	return infoStyle.Render("💡 " + text)
}

func RenderCommitMessage(message string) string {
	return commitBoxStyle.Render(
		boldStyle.Render("Generated Commit Message:") + "\n\n" +
			strings.TrimSpace(message),
	)
}

func RenderNote(note string) string {
	if note == "" {
		return ""
	}
	return noteBoxStyle.Render(
		warningStyle.Render("💡 Note: ") + note,
	)
}

func RenderConfigBox(content string) string {
	return configBoxStyle.Render(
		mutedStyle.Render("Configuration:") + "\n" + content,
	)
}

func RenderSeparator() string {
	return mutedStyle.Render(strings.Repeat("─", 50))
}

func RenderStep(step string) string {
	return infoStyle.Render("🔄 " + step)
}

func RenderProgress() {
	dots := []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}
	spinner := lipgloss.NewStyle().Foreground(primaryColor)
	for i := 0; i < 10; i++ {
		fmt.Printf("\r%s Generating commit message... %s",
			infoStyle.Render("🦕"),
			spinner.Render(dots[i%len(dots)]))
	}
	fmt.Print("\r" + strings.Repeat(" ", 50) + "\r")
}
