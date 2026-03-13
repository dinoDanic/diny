package yolo

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/dinoDanic/diny/tui/shared"
)

func (m model) View() string {
	var b strings.Builder

	b.WriteString("\n")
	b.WriteString(shared.RenderHeader(m.version, m.repoName, m.branchName, m.width))
	b.WriteString("\n")

	switch m.state {
	case stateStaging:
		b.WriteString(m.renderStaging())
	case stateGenerating:
		b.WriteString(m.renderGenerating())
	case stateCommitting:
		b.WriteString(m.renderCommitting())
	case stateSuccess:
		b.WriteString(m.renderSuccess())
	case stateNothingToCommit:
		b.WriteString(m.renderNothingToCommit())
	case stateError:
		b.WriteString(m.renderError())
	}

	return b.String()
}

func (m model) renderStaging() string {
	indent := indentStyle()
	return indent.Render(m.loader.View()) + "\n"
}

func (m model) renderGenerating() string {
	indent := indentStyle()
	var b strings.Builder

	if len(m.stagedFiles) > 0 {
		b.WriteString("\n")
		b.WriteString(m.renderStagedFiles())
		b.WriteString("\n")
	}

	b.WriteString(indent.Render(m.loader.View()))
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

	b.WriteString(indent.Render(m.loader.View()))
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

	successLine := "Committed! Pushed!"
	if m.hash != "" {
		successLine = fmt.Sprintf("%s  Committed! Pushed!", m.hash)
	}
	b.WriteString(indent.Render(successBigStyle().Render(successLine)))
	b.WriteString("\n\n")

	return b.String()
}

func (m model) renderNothingToCommit() string {
	indent := indentStyle()
	return "\n" + indent.Render(warningStyle().Render("No changes to commit.")) + "\n"
}

func (m model) renderError() string {
	indent := indentStyle()
	var b strings.Builder

	b.WriteString("\n")
	b.WriteString(indent.Render(errorStyle().Render("Error: "+m.err.Error())))
	b.WriteString("\n\n")
	b.WriteString(indent.Render(footerKeyStyle().Render("q") + " " + footerDescStyle().Render("quit")))
	b.WriteString("\n")

	return b.String()
}

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
