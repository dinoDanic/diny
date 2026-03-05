package commitui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

func (m model) View() string {
	if m.width == 0 {
		return ""
	}

	switch m.state {
	case stateLoading:
		return m.renderLoading()
	case stateError:
		return m.renderError()
	case stateHelp:
		return m.renderWithOverlay(m.renderHelp())
	default:
		return m.renderMain()
	}
}

func (m model) renderLoading() string {
	msg := "Generating commit message..."
	if m.statusMessage != "" {
		msg = m.statusMessage
	}
	content := m.spinner.View() + " " + loadingStyle().Render(msg)
	return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, content)
}

func (m model) renderError() string {
	errMsg := "Unknown error"
	if m.err != nil {
		errMsg = m.err.Error()
	}
	content := statusErrorStyle().Render("Error: " + errMsg + "\n\nPress any key to exit.")
	return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, content)
}

func (m model) renderMain() string {
	var sections []string

	sections = append(sections, m.renderHeader())
	sections = append(sections, m.renderBody())
	sections = append(sections, m.renderStatus())
	sections = append(sections, m.renderFooter())

	return lipgloss.JoinVertical(lipgloss.Left, sections...)
}

func (m model) renderHeader() string {
	label := headerLabelStyle().Render(" diny ")
	count := headerCountStyle().Render(fmt.Sprintf("%d file(s) staged", len(m.stagedFiles)))
	return lipgloss.JoinHorizontal(lipgloss.Center, label, " ", count) + "\n"
}

func (m model) renderBody() string {
	_, bodyHeight := m.getLayoutDimensions()

	if m.stackedLayout {
		return m.renderStackedBody(bodyHeight)
	}

	leftWidth := m.getLeftPaneWidth()
	rightWidth := m.getRightPaneWidth()

	leftContent := m.renderFileList(bodyHeight)
	left := filePaneStyle(leftWidth, bodyHeight).Render(leftContent)

	var rightContent string
	switch m.state {
	case stateEditing:
		rightContent = paneTitleStyle().Render("Edit Message") + "\n" + m.textarea.View()
	case stateFeedback:
		rightContent = paneTitleStyle().Render("Commit Message") + "\n" + m.viewport.View() + "\n\n" + m.textinput.View()
	case stateVariants:
		rightContent = m.renderVariants()
	default:
		rightContent = paneTitleStyle().Render("Commit Message") + "\n" + m.viewport.View()
	}
	right := messagePaneStyle(rightWidth, bodyHeight).Render(rightContent)

	return lipgloss.JoinHorizontal(lipgloss.Top, left, right)
}

func (m model) renderStackedBody(height int) string {
	filesSummary := m.renderFileSummary()
	header := paneTitleStyle().Render("Files: " + filesSummary)

	msgHeight := height - 3
	if msgHeight < 1 {
		msgHeight = 1
	}

	var content string
	switch m.state {
	case stateEditing:
		content = paneTitleStyle().Render("Edit Message") + "\n" + m.textarea.View()
	case stateFeedback:
		content = paneTitleStyle().Render("Commit Message") + "\n" + m.viewport.View() + "\n\n" + m.textinput.View()
	case stateVariants:
		content = m.renderVariants()
	default:
		content = paneTitleStyle().Render("Commit Message") + "\n" + m.viewport.View()
	}

	return header + "\n" + content
}

func (m model) renderFileSummary() string {
	added, modified, deleted := 0, 0, 0
	for _, f := range m.stagedFiles {
		switch {
		case strings.HasPrefix(f.Status, "A"):
			added++
		case strings.HasPrefix(f.Status, "D"):
			deleted++
		default:
			modified++
		}
	}

	var parts []string
	if added > 0 {
		parts = append(parts, fileAddedStyle().Render(fmt.Sprintf("+%d", added)))
	}
	if modified > 0 {
		parts = append(parts, fileModifiedStyle().Render(fmt.Sprintf("~%d", modified)))
	}
	if deleted > 0 {
		parts = append(parts, fileDeletedStyle().Render(fmt.Sprintf("-%d", deleted)))
	}
	return strings.Join(parts, " ")
}

func (m model) renderFileList(height int) string {
	title := paneTitleStyle().Render("Staged Files")
	var lines []string
	lines = append(lines, title)

	maxFiles := height - 2
	if maxFiles < 1 {
		maxFiles = 1
	}

	leftWidth := m.getLeftPaneWidth()
	maxPathLen := leftWidth - 6
	if maxPathLen < 10 {
		maxPathLen = 10
	}

	for i, f := range m.stagedFiles {
		if i >= maxFiles {
			remaining := len(m.stagedFiles) - maxFiles
			lines = append(lines, footerDescStyle().Render(fmt.Sprintf("  ... +%d more", remaining)))
			break
		}

		var icon string
		var iconStyle lipgloss.Style
		switch {
		case strings.HasPrefix(f.Status, "A"):
			icon = "+"
			iconStyle = fileAddedStyle()
		case strings.HasPrefix(f.Status, "D"):
			icon = "-"
			iconStyle = fileDeletedStyle()
		default:
			icon = "~"
			iconStyle = fileModifiedStyle()
		}

		path := f.Path
		if len(path) > maxPathLen {
			path = "..." + path[len(path)-maxPathLen+3:]
		}

		line := iconStyle.Render(icon) + " " + filePathStyle().Render(path)
		lines = append(lines, line)
	}

	return strings.Join(lines, "\n")
}

func (m model) renderVariants() string {
	title := paneTitleStyle().Render("Select a variant (j/k to navigate, enter to select, esc to cancel)")
	var lines []string
	lines = append(lines, title)

	for i, v := range m.variants {
		prefix := "  "
		style := variantInactiveStyle()
		if i == m.variantCursor {
			prefix = "> "
			style = variantActiveStyle()
		}
		lines = append(lines, style.Render(fmt.Sprintf("%s%d. %s", prefix, i+1, v)))
		lines = append(lines, "")
	}

	return strings.Join(lines, "\n")
}

func (m model) renderStatus() string {
	if m.statusMessage == "" {
		return ""
	}
	if m.statusIsError {
		return statusErrorStyle().Render(m.statusMessage)
	}
	return statusSuccessStyle().Render(m.statusMessage)
}

func (m model) renderFooter() string {
	if m.state == stateFeedback {
		return footerStyle(m.width).Render(
			footerKeyStyle().Render("enter") + " " + footerDescStyle().Render("submit") +
				footerSepStyle().Render(" | ") +
				footerKeyStyle().Render("esc") + " " + footerDescStyle().Render("cancel"),
		)
	}

	if m.state == stateEditing {
		return footerStyle(m.width).Render(
			footerKeyStyle().Render("esc") + " " + footerDescStyle().Render("accept") +
				footerSepStyle().Render(" | ") +
				footerKeyStyle().Render("ctrl+c") + " " + footerDescStyle().Render("cancel"),
		)
	}

	if m.state == stateVariants {
		return footerStyle(m.width).Render(
			footerKeyStyle().Render("j/k") + " " + footerDescStyle().Render("navigate") +
				footerSepStyle().Render(" | ") +
				footerKeyStyle().Render("enter") + " " + footerDescStyle().Render("select") +
				footerSepStyle().Render(" | ") +
				footerKeyStyle().Render("esc") + " " + footerDescStyle().Render("cancel"),
		)
	}

	sep := footerSepStyle().Render(" | ")
	dsep := footerSepStyle().Render(" || ")

	items := []string{
		footerKeyStyle().Render("enter") + " " + footerDescStyle().Render("commit"),
		sep,
		footerKeyStyle().Render("n") + " " + footerDescStyle().Render("no-verify"),
		sep,
		footerKeyStyle().Render("p") + " " + footerDescStyle().Render("push"),
		dsep,
		footerKeyStyle().Render("r") + " " + footerDescStyle().Render("regen"),
		sep,
		footerKeyStyle().Render("f") + " " + footerDescStyle().Render("feedback"),
		sep,
		footerKeyStyle().Render("e") + " " + footerDescStyle().Render("edit"),
		sep,
		footerKeyStyle().Render("E") + " " + footerDescStyle().Render("$EDITOR"),
		dsep,
		footerKeyStyle().Render("v") + " " + footerDescStyle().Render("variants"),
		sep,
		footerKeyStyle().Render("s") + " " + footerDescStyle().Render("draft"),
		sep,
		footerKeyStyle().Render("y") + " " + footerDescStyle().Render("copy"),
		sep,
		footerKeyStyle().Render("?") + " " + footerDescStyle().Render("help"),
		sep,
		footerKeyStyle().Render("q") + " " + footerDescStyle().Render("quit"),
	}

	return footerStyle(m.width).Render(strings.Join(items, ""))
}

func (m model) renderHelp() string {
	bindings := []struct {
		key  string
		desc string
	}{
		{"enter", "Commit with current message"},
		{"n", "Commit skipping hooks (--no-verify)"},
		{"p", "Commit and push"},
		{"r", "Regenerate commit message"},
		{"f", "Refine with custom feedback"},
		{"e", "Edit message inline"},
		{"E", "Edit in $EDITOR"},
		{"v", "Generate 3 variant messages"},
		{"s", "Save as draft"},
		{"y", "Copy to clipboard"},
		{"?", "Toggle this help"},
		{"q", "Quit"},
	}

	var lines []string
	lines = append(lines, paneTitleStyle().Render("Keyboard Shortcuts"))
	lines = append(lines, "")
	for _, b := range bindings {
		lines = append(lines, helpKeyStyle().Render(b.key)+helpDescStyle().Render(b.desc))
	}

	return strings.Join(lines, "\n")
}

func (m model) renderWithOverlay(overlay string) string {
	return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, helpOverlayStyle().Render(overlay))
}
