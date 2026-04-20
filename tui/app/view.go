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
	case stateFilePicker:
		b.WriteString(m.renderFilePicker())
	case stateError:
		b.WriteString(m.renderError())
	case stateSplitGenerating:
		b.WriteString(m.renderSplitGenerating())
	case stateSplitPlan:
		b.WriteString(m.renderSplitPlan())
	case stateSplitCommitting:
		b.WriteString(m.renderSplitCommitting())
	case stateSplitSuccess:
		b.WriteString(m.renderSplitSuccess())
	case stateSplitFailure:
		b.WriteString(m.renderSplitFailure())
	case stateSplitFeedback:
		b.WriteString(m.renderSplitFeedback())
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

	if m.currentTip != "" {
		b.WriteString(indent.Render(metaStyle().Render("• tip: '" + m.currentTip + "'")))
		b.WriteString("\n")
	}

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

	if m.currentTip != "" {
		b.WriteString(indent.Render(metaStyle().Render("• tip: '" + m.currentTip + "'")))
		b.WriteString("\n")
		b.WriteString("\n")
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
		{"x", "Manage staged/unstaged files"},
		{"S", "Split staged changes into multiple commits"},
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

	if m.commitProgress != "" {
		grey := lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
		b.WriteString(indent.Render(grey.Render(m.loader.SpinnerFrame() + " " + m.commitProgress)))
		b.WriteString("\n")
	}

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

func (m model) renderFilePicker() string {
	indent := indentStyle()
	var b strings.Builder

	b.WriteString("\n")
	b.WriteString(indent.Render(sectionTitleStyle().Render("Manage Staged Files")))
	b.WriteString("\n\n")

	for i, e := range m.fileEntries {
		cursor := "  "
		if i == m.filePickerCursor {
			cursor = "> "
		}

		var indicator string
		var style lipgloss.Style
		switch {
		case e.currentStaged && e.wantStaged:
			indicator = "[✓]"
			style = sectionTitleStyle()
		case e.currentStaged && !e.wantStaged:
			indicator = "[-]"
			style = fileDeletedStyle()
		case !e.currentStaged && e.wantStaged:
			indicator = "[+]"
			style = fileAddedStyle()
		default:
			indicator = "[ ]"
			style = metaStyle()
		}

		var statusIcon string
		switch e.status {
		case "A":
			statusIcon = "+"
		case "M":
			statusIcon = "~"
		case "D":
			statusIcon = "-"
		case "R":
			statusIcon = ">"
		default:
			statusIcon = "?"
		}

		line := cursor + indicator + " " + style.Render(statusIcon+" "+e.path)
		b.WriteString(indent.Render(line))
		b.WriteString("\n")
	}

	b.WriteString("\n")
	keys := []struct{ key, desc string }{
		{"↑/k", "up"}, {"↓/j", "down"}, {"space", "toggle"}, {"a", "all"}, {"enter", "apply"}, {"esc", "cancel"},
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

func (m model) renderSplitGenerating() string {
	indent := indentStyle()
	var b strings.Builder

	b.WriteString("\n")
	b.WriteString(m.renderStagedFiles())
	b.WriteString("\n")
	b.WriteString(indent.Render(sectionTitleStyle().Render("Planning split")))
	b.WriteString("\n")
	b.WriteString(indent.Render(m.loader.View()))
	b.WriteString("\n")
	b.WriteString(indent.Render(metaStyle().Render("grouping staged files into logical commits...")))
	b.WriteString("\n")
	return b.String()
}

func (m model) renderSplitPlan() string {
	indent := indentStyle()
	var b strings.Builder

	total := len(m.splitPlan)

	title := fmt.Sprintf("Split plan — %d commit(s)", total)
	if m.splitMoveMode {
		if m.splitMovePickDest {
			title += " — pick destination group"
		} else {
			title += " — move mode"
		}
	}
	b.WriteString("\n")
	b.WriteString(indent.Render(sectionTitleStyle().Render(title)))
	b.WriteString("\n\n")

	if m.splitRegenerating {
		b.WriteString(indent.Render(m.loader.View()))
		b.WriteString("\n")
		b.WriteString(indent.Render(metaStyle().Render("regenerating plan...")))
		b.WriteString("\n\n")
	}

	for i, g := range m.splitPlan {
		groupSelected := i == m.splitCursor
		destHighlight := m.splitMoveMode && m.splitMovePickDest && i == m.splitMoveDestIdx && i != m.splitCursor
		expanded := m.splitExpanded[i]

		cursor := "  "
		if groupSelected {
			cursor = "> "
		} else if destHighlight {
			cursor = "→ "
		}

		caret := "▸"
		if expanded {
			caret = "▾"
		}

		firstLine := g.Message
		if idx := strings.Index(g.Message, "\n"); idx >= 0 {
			firstLine = g.Message[:idx]
		}

		header := fmt.Sprintf("[%d/%d] %s — %s", g.Order, total, g.Type, firstLine)
		var headerStyle lipgloss.Style
		switch {
		case groupSelected:
			headerStyle = sectionTitleStyle()
		case destHighlight:
			headerStyle = splitDestStyle()
		default:
			headerStyle = lipgloss.NewStyle()
		}
		line := cursor + caret + " " + headerStyle.Render(header)
		b.WriteString(indent.Render(line))
		b.WriteString("\n")

		if expanded {
			for fi, f := range g.Files {
				status := m.stagedStatus(f)
				var statusStyle lipgloss.Style
				switch status {
				case "A":
					statusStyle = fileAddedStyle()
				case "M":
					statusStyle = fileModifiedStyle()
				case "D":
					statusStyle = fileDeletedStyle()
				case "R":
					statusStyle = fileRenamedStyle()
				default:
					statusStyle = metaStyle()
				}

				marker := "  "
				if m.splitMoveMode && !m.splitMovePickDest && groupSelected && fi == m.splitMoveFileIdx {
					marker = splitMoveCursorStyle().Render("▶ ")
				}

				fileLine := "   " + marker + statusStyle.Render(status+" "+f)
				b.WriteString(indent.Render(fileLine))
				b.WriteString("\n")
			}
		}
	}

	b.WriteString("\n")

	if m.statusMessage != "" {
		b.WriteString(m.renderStatus())
	}

	var keys []struct{ key, desc string }
	if m.splitMoveMode {
		if m.splitMovePickDest {
			keys = []struct{ key, desc string }{
				{"↑/k", "prev group"},
				{"↓/j", "next group"},
				{"enter", "move here"},
				{"esc", "back"},
			}
		} else {
			if len(m.splitPlan) <= 9 {
				keys = []struct{ key, desc string }{
					{"↑/k", "prev file"},
					{"↓/j", "next file"},
					{"1-9", "reassign to group"},
					{"enter", "pick dest"},
					{"esc", "cancel"},
				}
			} else {
				keys = []struct{ key, desc string }{
					{"↑/k", "prev file"},
					{"↓/j", "next file"},
					{"enter", "pick dest"},
					{"esc", "cancel"},
				}
			}
		}
	} else {
		keys = []struct{ key, desc string }{
			{"↑/k", "up"},
			{"↓/j", "down"},
			{"enter/space", "expand"},
			{"e", "edit msg"},
			{"m", "move file"},
			{"r", "regen"},
			{"f", "regen w/ feedback"},
			{"c", "confirm all"},
			{"esc/q", "cancel"},
		}
	}
	var parts []string
	for _, k := range keys {
		parts = append(parts, footerKeyStyle().Render(k.key)+" "+footerDescStyle().Render(k.desc))
	}
	b.WriteString(indent.Render(strings.Join(parts, "  ")))
	b.WriteString("\n")

	return b.String()
}

func (m model) renderSplitCommitting() string {
	indent := indentStyle()
	var b strings.Builder

	b.WriteString("\n")
	b.WriteString(indent.Render(sectionTitleStyle().Render(fmt.Sprintf("Splitting into %d commit(s)", len(m.splitPlan)))))
	b.WriteString("\n")
	b.WriteString(indent.Render(m.loader.View()))
	b.WriteString("\n")
	return b.String()
}

func (m model) renderSplitSuccess() string {
	indent := indentStyle()
	var b strings.Builder

	b.WriteString("\n")
	headline := fmt.Sprintf("Created %d commit(s)!", len(m.splitHashes))
	if m.splitPushed {
		headline += " Pushed!"
	}
	b.WriteString(indent.Render(successBigStyle().Render(headline)))
	b.WriteString("\n\n")

	for i, g := range m.splitPlan {
		hash := ""
		if i < len(m.splitHashes) {
			hash = m.splitHashes[i]
		}
		firstLine := g.Message
		if idx := strings.Index(g.Message, "\n"); idx >= 0 {
			firstLine = g.Message[:idx]
		}
		line := fmt.Sprintf("  %s %s", hash, firstLine)
		b.WriteString(indent.Render(line))
		b.WriteString("\n")
	}

	b.WriteString("\n")
	return b.String()
}

func (m model) renderSplitFeedback() string {
	indent := indentStyle()
	var b strings.Builder

	b.WriteString("\n")
	b.WriteString(indent.Render(sectionTitleStyle().Render("Plan feedback")))
	b.WriteString("\n")
	b.WriteString(indent.Render(metaStyle().Render("Describe what's off so the model regroups differently on the next plan.")))
	b.WriteString("\n\n")
	b.WriteString(indent.Render(m.textinput.View()))
	b.WriteString("\n\n")
	b.WriteString(indent.Render(footerKeyStyle().Render("enter") + " " + footerDescStyle().Render("submit") + "  " + footerKeyStyle().Render("esc") + " " + footerDescStyle().Render("cancel")))
	b.WriteString("\n")
	return b.String()
}

func (m model) renderSplitFailure() string {
	indent := indentStyle()
	var b strings.Builder

	if m.splitFailure == nil {
		return ""
	}
	f := m.splitFailure

	b.WriteString("\n")
	b.WriteString(indent.Render(errorStyle().Render("Split stopped — one commit failed")))
	b.WriteString("\n\n")

	committed := len(f.committedHashes)
	if committed > 0 {
		b.WriteString(indent.Render(sectionTitleStyle().Render(fmt.Sprintf("Landed (%d):", committed))))
		b.WriteString("\n")
		for i := 0; i < committed && i < len(m.splitPlan); i++ {
			hash := f.committedHashes[i]
			firstLine := m.splitPlan[i].Message
			if idx := strings.Index(firstLine, "\n"); idx >= 0 {
				firstLine = firstLine[:idx]
			}
			b.WriteString(indent.Render(fmt.Sprintf("  %s %s", hash, firstLine)))
			b.WriteString("\n")
		}
		b.WriteString("\n")
	} else {
		b.WriteString(indent.Render(metaStyle().Render("No commits landed before the failure.")))
		b.WriteString("\n\n")
	}

	if f.failedIndex >= 0 && f.failedIndex < len(m.splitPlan) {
		g := m.splitPlan[f.failedIndex]
		firstLine := g.Message
		if idx := strings.Index(firstLine, "\n"); idx >= 0 {
			firstLine = firstLine[:idx]
		}
		b.WriteString(indent.Render(sectionTitleStyle().Render(fmt.Sprintf("Failed group [%d/%d] %s — %s", g.Order, len(m.splitPlan), g.Type, firstLine))))
		b.WriteString("\n")
	}

	if f.failedStderr != "" {
		b.WriteString(indent.Render(metaStyle().Render("git stderr:")))
		b.WriteString("\n")
		for _, line := range strings.Split(f.failedStderr, "\n") {
			b.WriteString(indent.Render("  " + statusErrorStyle().Render(line)))
			b.WriteString("\n")
		}
		b.WriteString("\n")
	}

	if len(f.failedFiles) > 0 {
		b.WriteString(indent.Render(sectionTitleStyle().Render("Still staged (failed group's files):")))
		b.WriteString("\n")
		for _, path := range f.failedFiles {
			b.WriteString(indent.Render("  " + metaStyle().Render(path)))
			b.WriteString("\n")
		}
		b.WriteString("\n")
	}

	if len(f.remainingFiles) > 0 {
		b.WriteString(indent.Render(sectionTitleStyle().Render(fmt.Sprintf("Unstaged — groups that never ran (%d files):", len(f.remainingFiles)))))
		b.WriteString("\n")
		for _, path := range f.remainingFiles {
			b.WriteString(indent.Render("  " + metaStyle().Render(path)))
			b.WriteString("\n")
		}
		b.WriteString("\n")
	}

	b.WriteString(indent.Render(footerKeyStyle().Render("enter/q") + " " + footerDescStyle().Render("quit")))
	b.WriteString("\n")

	return b.String()
}

func (m model) stagedStatus(path string) string {
	for _, f := range m.stagedFiles {
		if f.Path == path {
			return f.Status
		}
	}
	return "?"
}

func (m model) renderFooter() string {
	indent := indentStyle()

	type kb struct{ key, desc string }

	row1 := []kb{
		{"enter", "commit"},
		{"n", "no-verify"},
		{"p", "push"},
		{"r", "regen"},
		{"v", "variants"},
		{"f", "feedback"},
	}
	row2 := []kb{
		{"E", "$EDITOR"},
		{"x", "files"},
		{"S", "split"},
		{"q", "quit"},
		{"?", "more"},
	}

	render := func(keys []kb) string {
		var parts []string
		for _, k := range keys {
			parts = append(parts, footerKeyStyle().Render(k.key)+" "+footerDescStyle().Render(k.desc))
		}
		return indent.Render(strings.Join(parts, "  "))
	}

	return render(row1) + "\n" + render(row2) + "\n"
}

