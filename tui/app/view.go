package app

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
	case stateVariantPicking:
		b.WriteString(m.renderVariantPicking())
	case stateDiffView:
		b.WriteString(m.renderDiffView())
	case stateTypePicker:
		b.WriteString(m.renderTypePicker())
	case stateUnstaging:
		b.WriteString(m.renderUnstaging())
	case stateError:
		b.WriteString(m.renderError())
	}

	return b.String()
}


func (m model) renderWelcome() string {
	indent := indentStyle()
	return "\n" + indent.Render(m.loader.View()) + "\n"
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

	if m.statusMessage != "" {
		b.WriteString(m.renderStatus())
	}

	return b.String()
}

func (m model) renderReady() string {
	var b strings.Builder
	indent := indentStyle()

	b.WriteString("\n")
	b.WriteString(m.renderStagedFiles())
	b.WriteString("\n")
	b.WriteString(m.renderCommitMessage())
	b.WriteString("\n")

	if m.messageHistoryIdx != -1 {
		pos := fmt.Sprintf("message %d of %d", m.messageHistoryIdx+1, len(m.previousMessages))
		b.WriteString(indent.Render(metaStyle().Render(pos)))
		b.WriteString("\n")
	}

	if m.pendingAmend {
		b.WriteString(indent.Render(metaStyle().Render("amend mode")))
		b.WriteString("\n")
	}

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
		{"A", "Regenerate from HEAD diff and amend on commit"},
		{"r", "Regenerate commit message"},
		{"v", "Pick from 3 variants"},
		{"f", "Refine with feedback"},
		{"t", "Force conventional commit type (requires conventional: true)"},
		{"L", "Cycle length: short → normal → long (session only)"},
		{"M", "Toggle emoji on/off (session only)"},
		{"e", "Edit inline"},
		{"E", "Edit in $EDITOR"},
		{"d", "View staged diff"},
		{"[", "Browse previous generated messages"},
		{"]", "Browse forward through message history"},
		{"x", "Unstage files"},
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

	b.WriteString(indent.Render(successBigStyle().Render(m.statusMessage)))
	b.WriteString("\n\n")
	return b.String()
}

func (m model) renderNoStaged() string {
	indent := indentStyle()
	var b strings.Builder

	b.WriteString("\n")

	if len(m.unstagedFiles) == 0 {
		b.WriteString(indent.Render(noStagedStyle().Render("Working tree is clean. Nothing to stage.")))
		b.WriteString("\n\n")
		b.WriteString(indent.Render(footerKeyStyle().Render("q") + " " + footerDescStyle().Render("quit")))
		b.WriteString("\n")
		return b.String()
	}

	b.WriteString(indent.Render(noStagedStyle().Render("No staged files detected. Select files to stage:")))
	b.WriteString("\n\n")
	b.WriteString(indent.Render(sectionTitleStyle().Render("Unstaged Files")))
	b.WriteString("\n")

	for i, f := range m.unstagedFiles {
		cursor := "  "
		if i == m.fileCursor {
			cursor = "> "
		}
		checkbox := "[ ]"
		if i < len(m.fileSelected) && m.fileSelected[i] {
			checkbox = "[x]"
		}

		var style lipgloss.Style
		switch f.Status {
		case "M":
			style = fileModifiedStyle()
		case "D":
			style = fileDeletedStyle()
		default:
			style = metaStyle()
		}

		line := cursor + checkbox + " " + style.Render(f.Path)
		b.WriteString(indent.Render(line))
		b.WriteString("\n")
	}

	b.WriteString("\n")
	keys := []struct{ key, desc string }{
		{"↑/k", "up"}, {"↓/j", "down"}, {"space", "toggle"}, {"a", "all"}, {"enter", "stage"}, {"q", "quit"},
	}
	var parts []string
	for _, k := range keys {
		parts = append(parts, footerKeyStyle().Render(k.key)+" "+footerDescStyle().Render(k.desc))
	}
	b.WriteString(indent.Render(strings.Join(parts, "  ")))
	b.WriteString("\n")

	return b.String()
}

func (m model) renderVariantPicking() string {
	indent := indentStyle()
	var b strings.Builder

	b.WriteString("\n")
	b.WriteString(m.renderStagedFiles())
	b.WriteString("\n")
	b.WriteString(indent.Render(sectionTitleStyle().Render("Pick a Variant")))
	b.WriteString("\n\n")

	for i, v := range m.variants {
		cursor := "  "
		if i == m.variantCursor {
			cursor = "> "
		}

		var style lipgloss.Style
		if i == m.variantCursor {
			style = sectionTitleStyle()
		} else {
			style = metaStyle()
		}

		line := cursor + style.Render(fmt.Sprintf("%d. %s", i+1, v))
		b.WriteString(indent.Render(line))
		b.WriteString("\n\n")
	}

	keys := []struct{ key, desc string }{
		{"↑/k", "up"}, {"↓/j", "down"}, {"1/2/3", "pick"}, {"enter", "select"}, {"esc", "cancel"},
	}
	var parts []string
	for _, k := range keys {
		parts = append(parts, footerKeyStyle().Render(k.key)+" "+footerDescStyle().Render(k.desc))
	}
	b.WriteString(indent.Render(strings.Join(parts, "  ")))
	b.WriteString("\n")

	return b.String()
}

func (m model) renderDiffView() string {
	indent := indentStyle()
	var b strings.Builder

	b.WriteString("\n")
	b.WriteString(indent.Render(sectionTitleStyle().Render("Staged Diff")))
	b.WriteString("\n\n")
	b.WriteString(indent.Render(m.viewport.View()))
	b.WriteString("\n\n")
	b.WriteString(indent.Render(metaStyle().Render("↑/k up  ↓/j down  pgup/pgdn scroll  esc close")))
	b.WriteString("\n")

	return b.String()
}

func (m model) renderTypePicker() string {
	indent := indentStyle()
	var b strings.Builder

	b.WriteString("\n")
	b.WriteString(m.renderStagedFiles())
	b.WriteString("\n")
	b.WriteString(indent.Render(sectionTitleStyle().Render("Force Commit Type")))
	b.WriteString("\n\n")

	for i, t := range conventionalTypes {
		cursor := "  "
		if i == m.typeCursor {
			cursor = "> "
		}
		var style lipgloss.Style
		if i == m.typeCursor {
			style = sectionTitleStyle()
		} else {
			style = metaStyle()
		}
		line := cursor + style.Render(fmt.Sprintf("%d. %s", i+1, t))
		b.WriteString(indent.Render(line))
		b.WriteString("\n")
	}

	b.WriteString("\n")
	keys := []struct{ key, desc string }{
		{"↑/k", "up"}, {"↓/j", "down"}, {"1-8", "pick"}, {"enter", "select"}, {"esc", "cancel"},
	}
	var parts []string
	for _, k := range keys {
		parts = append(parts, footerKeyStyle().Render(k.key)+" "+footerDescStyle().Render(k.desc))
	}
	b.WriteString(indent.Render(strings.Join(parts, "  ")))
	b.WriteString("\n")

	return b.String()
}

func (m model) renderUnstaging() string {
	indent := indentStyle()
	var b strings.Builder

	b.WriteString("\n")
	b.WriteString(indent.Render(sectionTitleStyle().Render("Unstage Files")))
	b.WriteString("\n\n")

	for i, f := range m.stagedFiles {
		cursor := "  "
		if i == m.unstageCursor {
			cursor = "> "
		}
		checkbox := "[ ]"
		if i < len(m.unstageSelected) && m.unstageSelected[i] {
			checkbox = "[x]"
		}
		var style lipgloss.Style
		switch f.Status {
		case "A":
			style = fileAddedStyle()
		case "M":
			style = fileModifiedStyle()
		case "D":
			style = fileDeletedStyle()
		default:
			style = metaStyle()
		}
		line := cursor + checkbox + " " + style.Render(f.Path)
		b.WriteString(indent.Render(line))
		b.WriteString("\n")
	}

	b.WriteString("\n")
	keys := []struct{ key, desc string }{
		{"↑/k", "up"}, {"↓/j", "down"}, {"space", "toggle"}, {"enter", "unstage"}, {"esc", "cancel"}, {"q", "quit"},
	}
	var parts []string
	for _, k := range keys {
		parts = append(parts, footerKeyStyle().Render(k.key)+" "+footerDescStyle().Render(k.desc))
	}
	b.WriteString(indent.Render(strings.Join(parts, "  ")))
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

	firstLine := m.commitMessage
	if idx := strings.Index(m.commitMessage, "\n"); idx >= 0 {
		firstLine = m.commitMessage[:idx]
	}
	if len(firstLine) > 72 {
		warning := fmt.Sprintf("Subject line: %d chars (aim for ≤72)", len(firstLine))
		b.WriteString(indent.Render(statusErrorStyle().Render(warning)))
		b.WriteString("\n")
	}

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

	type kb struct{ key, desc string }
	type group struct {
		label string
		keys  []kb
	}

	groups := []group{
		{"commit", []kb{
			{"enter", "commit"},
			{"n", "no-verify"},
			{"p", "push"},
			{"A", "amend"},
		}},
		{"generate", []kb{
			{"r", "regen"},
			{"v", "variants"},
			{"f", "feedback"},
			{"t", "type"},
			{"L", "length"},
			{"M", "emoji"},
		}},
		{"view", []kb{
			{"d", "diff"},
			{"[/]", "history"},
			{"e", "edit"},
			{"E", "$EDITOR"},
			{"x", "unstage"},
		}},
		{"misc", []kb{
			{"s", "draft"},
			{"y", "copy"},
			{"?", "help"},
			{"q", "quit"},
		}},
	}

	var lines []string
	for _, g := range groups {
		var parts []string
		for _, k := range g.keys {
			parts = append(parts, footerKeyStyle().Render(k.key)+" "+footerDescStyle().Render(k.desc))
		}
		line := footerGroupStyle().Render(g.label) + strings.Join(parts, "  ")
		lines = append(lines, indent.Render(line))
	}

	return strings.Join(lines, "\n") + "\n"
}

