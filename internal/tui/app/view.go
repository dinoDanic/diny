package app

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/dinoDanic/diny/ui"
)

func (m model) View() string {
	var b strings.Builder

	b.WriteString("\n")
	b.WriteString(m.renderHeader())
	b.WriteString("\n")

	switch m.state {
	case stateWelcome:
		b.WriteString(m.renderWelcome())
	case stateGenerating:
		b.WriteString(m.renderGenerating())
	case stateReady:
		b.WriteString(m.renderReady())
	case stateFeedback:
		b.WriteString(m.renderFeedback())
	case stateEditing:
		b.WriteString(m.renderEditing())
	case stateHelp:
		b.WriteString(m.renderHelp())
	case stateCommitting:
		b.WriteString(m.renderCommitting())
	case stateSuccess:
		b.WriteString(m.renderSuccess())
	case stateNoStaged:
		b.WriteString(m.renderNoStaged())
	case stateError:
		b.WriteString(m.renderError())
	}

	return b.String()
}

func (m model) renderHeader() string {
	indent := indentStyle()
	brand := brandStyle().Render(" diny ")
	version := metaStyle().Render("v" + m.version)

	line1 := indent.Render(brand + "  " + version)

	var metaParts []string
	if m.repoName != "" {
		metaParts = append(metaParts, m.repoName)
	}
	if m.branchName != "" {
		metaParts = append(metaParts, m.branchName)
	}
	if len(m.stagedFiles) > 0 {
		metaParts = append(metaParts, fmt.Sprintf("%d files staged", len(m.stagedFiles)))
	}

	line2 := ""
	if len(metaParts) > 0 {
		sep := metaStyle().Render(" \u00b7 ")
		rendered := make([]string, len(metaParts))
		for i, p := range metaParts {
			rendered[i] = metaStyle().Render(p)
		}
		line2 = indent.Render(strings.Join(rendered, sep))
	}

	if line2 != "" {
		return line1 + "\n" + line2 + "\n"
	}
	return line1 + "\n"
}

func (m model) renderWelcome() string {
	indent := indentStyle()
	return "\n" + indent.Render(m.spinner.View()+" Loading...") + "\n"
}

func (m model) renderGenerating() string {
	indent := indentStyle()
	t := ui.GetCurrentTheme()
	genStyle := lipgloss.NewStyle().Foreground(t.PrimaryForeground)
	var b strings.Builder

	if len(m.stagedFiles) > 0 {
		b.WriteString("\n")
		b.WriteString(m.renderStagedFiles())
		b.WriteString("\n")
	}

	b.WriteString(indent.Render(genStyle.Render(m.spinner.View() + " Generating commit message...")))
	b.WriteString("\n")
	return b.String()
}

func (m model) renderReady() string {
	var b strings.Builder

	b.WriteString("\n")
	b.WriteString(m.renderStagedFiles())
	b.WriteString("\n")
	b.WriteString(m.renderCommitMessage())
	b.WriteString("\n")

	if m.statusMessage != "" {
		b.WriteString(m.renderStatus())
	}

	b.WriteString(m.renderFooter())
	b.WriteString("\n")

	return b.String()
}

func (m model) renderFeedback() string {
	var b strings.Builder
	indent := indentStyle()

	b.WriteString("\n")
	b.WriteString(m.renderStagedFiles())
	b.WriteString("\n")
	b.WriteString(m.renderCommitMessage())
	b.WriteString("\n")

	b.WriteString(indent.Render(sectionTitleStyle().Render("Feedback")))
	b.WriteString("\n")
	b.WriteString(indent.Render(m.textinput.View()))
	b.WriteString("\n\n")
	b.WriteString(indent.Render(metaStyle().Render("enter submit  esc cancel")))
	b.WriteString("\n")

	return b.String()
}

func (m model) renderEditing() string {
	var b strings.Builder
	indent := indentStyle()

	b.WriteString("\n")
	b.WriteString(m.renderStagedFiles())
	b.WriteString("\n")

	b.WriteString(indent.Render(sectionTitleStyle().Render("Edit Commit Message")))
	b.WriteString("\n")
	b.WriteString(indent.Render(m.textarea.View()))
	b.WriteString("\n\n")
	b.WriteString(indent.Render(metaStyle().Render("esc accept")))
	b.WriteString("\n")

	return b.String()
}

func (m model) renderHelp() string {
	var b strings.Builder

	b.WriteString("\n")

	helpContent := []struct {
		key  string
		desc string
	}{
		{"enter", "Commit"},
		{"n", "Commit (skip hooks / no-verify)"},
		{"p", "Commit and push"},
		{"r", "Regenerate commit message"},
		{"f", "Refine with feedback"},
		{"e", "Edit inline"},
		{"E", "Edit in $EDITOR"},
		{"s", "Save as draft"},
		{"y", "Copy to clipboard"},
		{"?", "Toggle help"},
		{"q", "Quit"},
	}

	var lines []string
	for _, h := range helpContent {
		key := footerKeyStyle().Render(fmt.Sprintf("%-8s", h.key))
		desc := metaStyle().Render(h.desc)
		lines = append(lines, "  "+key+" "+desc)
	}

	content := sectionTitleStyle().Render("Keyboard Shortcuts") + "\n\n" + strings.Join(lines, "\n")
	box := helpBoxStyle().Render(content)
	b.WriteString(indentStyle().Render(box))
	b.WriteString("\n\n")
	b.WriteString(indentStyle().Render(metaStyle().Render("Press any key to close")))
	b.WriteString("\n")

	return b.String()
}

func (m model) renderCommitting() string {
	indent := indentStyle()
	var b strings.Builder

	b.WriteString("\n")
	b.WriteString(m.renderStagedFiles())
	b.WriteString("\n")
	b.WriteString(m.renderCommitMessage())
	b.WriteString("\n")

	b.WriteString(indent.Render(m.spinner.View() + " Committing..."))
	b.WriteString("\n")
	return b.String()
}

func (m model) renderSuccess() string {
	indent := indentStyle()
	var b strings.Builder

	b.WriteString("\n")
	b.WriteString(m.renderStagedFiles())
	b.WriteString("\n")
	b.WriteString(m.renderCommitMessage())
	b.WriteString("\n")

	b.WriteString(indent.Render(successBigStyle().Render(m.statusMessage)))
	b.WriteString("\n\n")
	return b.String()
}

func (m model) renderNoStaged() string {
	indent := indentStyle()
	var b strings.Builder

	b.WriteString("\n")
	b.WriteString(indent.Render(noStagedStyle().Render("No staged changes found.")))
	b.WriteString("\n")
	b.WriteString(indent.Render(metaStyle().Render("Stage files with: git add <files>")))
	b.WriteString("\n\n")
	b.WriteString(indent.Render(footerKeyStyle().Render("q")+" "+footerDescStyle().Render("quit")))
	b.WriteString("\n")
	return b.String()
}

func (m model) renderError() string {
	indent := indentStyle()
	var b strings.Builder

	b.WriteString("\n")
	b.WriteString(indent.Render(errorStyle().Render("Error: "+m.err.Error())))
	b.WriteString("\n\n")
	b.WriteString(indent.Render(footerKeyStyle().Render("q")+" "+footerDescStyle().Render("quit")))
	b.WriteString("\n")
	return b.String()
}

// Shared render helpers

func (m model) renderStagedFiles() string {
	indent := indentStyle()
	var b strings.Builder

	b.WriteString(indent.Render(sectionTitleStyle().Render("Staged Files")))
	b.WriteString("\n")

	for _, f := range m.stagedFiles {
		var icon string
		var style lipgloss.Style
		switch f.Status {
		case "A":
			icon = "+"
			style = fileAddedStyle()
		case "M":
			icon = "~"
			style = fileModifiedStyle()
		case "D":
			icon = "-"
			style = fileDeletedStyle()
		case "R":
			icon = ">"
			style = fileRenamedStyle()
		default:
			icon = "?"
			style = metaStyle()
		}
		line := style.Render(fmt.Sprintf("  %s %s", icon, f.Path))
		b.WriteString(indent.Render(line))
		b.WriteString("\n")
	}

	return b.String()
}

func (m model) renderCommitMessage() string {
	indent := indentStyle()
	var b strings.Builder

	b.WriteString(indent.Render(sectionTitleStyle().Render("Commit Message")))
	b.WriteString("\n")
	b.WriteString(indent.Render(commitMessageStyle().Render(m.commitMessage)))
	b.WriteString("\n")

	return b.String()
}

func (m model) renderStatus() string {
	indent := indentStyle()
	var style lipgloss.Style
	if m.statusIsError {
		style = statusErrorStyle()
	} else {
		style = statusSuccessStyle()
	}
	return indent.Render(style.Render(m.statusMessage)) + "\n"
}

func (m model) renderFooter() string {
	indent := indentStyle()
	keys := []struct {
		key  string
		desc string
	}{
		{"enter", "commit"},
		{"n", "no-verify"},
		{"p", "push"},
		{"r", "regen"},
		{"f", "feedback"},
		{"e", "edit"},
		{"E", "$EDITOR"},
		{"s", "draft"},
		{"y", "copy"},
		{"?", "help"},
		{"q", "quit"},
	}

	var parts []string
	for _, k := range keys {
		parts = append(parts, footerKeyStyle().Render(k.key)+" "+footerDescStyle().Render(k.desc))
	}

	return indent.Render(strings.Join(parts, "  ")) + "\n"
}

