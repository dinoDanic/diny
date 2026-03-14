package changelog

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/dinoDanic/diny/git"
	"github.com/dinoDanic/diny/tui/shared"
)

var modeMenuItems = []struct {
	label string
	value string
}{
	{"Between two tags", "tag"},
	{"Between two commits", "commit"},
}

func (m model) View() string {
	var b strings.Builder

	b.WriteString("\n")
	b.WriteString(shared.RenderHeader(m.version, m.repoName, m.branchName, m.width))
	b.WriteString("\n")

	switch m.state {
	case stateModeSelect:
		b.WriteString(m.renderModeSelect())
	case stateLoadingRefs:
		b.WriteString(m.renderLoading())
	case stateSelectNewerRef:
		title := "Select newer ref (to):"
		b.WriteString(m.renderRefList(title, m.currentListLabels()))
	case stateSelectOlderRef:
		title := "Select older ref (from):"
		b.WriteString(m.renderRefList(title, m.currentListLabels()))
	case stateGenerating, stateRegenerating:
		b.WriteString(m.renderLoading())
	case stateResults:
		b.WriteString(m.renderResults())
	case stateNoCommits:
		b.WriteString(m.renderNoCommits())
	case stateError:
		b.WriteString(m.renderError())
	}

	return b.String()
}

func (m model) renderModeSelect() string {
	indent := indentStyle()
	var b strings.Builder

	b.WriteString(indent.Render(sectionTitleStyle().Render("Generate a changelog for...")))
	b.WriteString("\n\n")

	for i, item := range modeMenuItems {
		var line string
		if i == m.modeCursor {
			line = footerKeyStyle().Render("▶") + "  " +
				metaStyle().Render(fmt.Sprintf("%d", i+1)) + "  " +
				sectionTitleStyle().Render(item.label)
		} else {
			line = "   " +
				metaStyle().Render(fmt.Sprintf("%d", i+1)) + "  " +
				item.label
		}
		b.WriteString(indent.Render(line))
		b.WriteString("\n")
	}

	b.WriteString("\n")
	b.WriteString(indent.Render(
		footerKeyStyle().Render("j/k") + " " + footerDescStyle().Render("move") + "  " +
			footerKeyStyle().Render("enter") + " " + footerDescStyle().Render("confirm") + "  " +
			footerKeyStyle().Render("q") + " " + footerDescStyle().Render("quit"),
	))
	b.WriteString("\n")

	return b.String()
}

func (m model) renderRefList(title string, items []string) string {
	indent := indentStyle()
	var b strings.Builder

	b.WriteString(indent.Render(sectionTitleStyle().Render(title)))
	b.WriteString("\n\n")

	end := m.listOffset + listPageSize
	if end > len(items) {
		end = len(items)
	}
	visible := items[m.listOffset:end]

	for i, item := range visible {
		idx := m.listOffset + i
		var line string
		if idx == m.listCursor {
			line = footerKeyStyle().Render("▶") + "  " + item
		} else {
			line = "   " + metaStyle().Render(item)
		}
		b.WriteString(indent.Render(line))
		b.WriteString("\n")
	}

	if len(items) > m.listOffset+listPageSize {
		remaining := len(items) - m.listOffset - listPageSize
		b.WriteString(indent.Render(metaStyle().Render(fmt.Sprintf("  ... %d more", remaining))))
		b.WriteString("\n")
	}

	b.WriteString("\n")
	b.WriteString(indent.Render(
		footerKeyStyle().Render("j/k") + " " + footerDescStyle().Render("move") + "  " +
			footerKeyStyle().Render("enter") + " " + footerDescStyle().Render("confirm") + "  " +
			footerKeyStyle().Render("esc") + " " + footerDescStyle().Render("back"),
	))
	b.WriteString("\n")

	return b.String()
}

func (m model) renderLoading() string {
	indent := indentStyle()
	return indent.Render(m.loader.View()) + "\n"
}

func (m model) renderResults() string {
	indent := indentStyle()
	var b strings.Builder

	b.WriteString("\n")
	b.WriteString(indent.Render(sectionTitleStyle().Render("Changelog: " + m.rangeLabel)))
	b.WriteString("\n\n")

	w := m.width - 6
	if w < 40 {
		w = 40
	}
	resultStyle := lipgloss.NewStyle().Width(w)
	b.WriteString(indent.Render(resultStyle.Render(m.result)))
	b.WriteString("\n\n")

	if m.statusMessage != "" {
		var statusLine string
		if m.statusIsError {
			statusLine = errorStyle().Render(m.statusMessage)
		} else {
			statusLine = statusSuccessStyle().Render(m.statusMessage)
		}
		b.WriteString(indent.Render(statusLine))
		b.WriteString("\n")
	}

	b.WriteString(indent.Render(
		footerKeyStyle().Render("c") + " " + footerDescStyle().Render("copy") + "  " +
			footerKeyStyle().Render("s") + " " + footerDescStyle().Render("save") + "  " +
			footerKeyStyle().Render("r") + " " + footerDescStyle().Render("regen") + "  " +
			footerKeyStyle().Render("n") + " " + footerDescStyle().Render("new") + "  " +
			footerKeyStyle().Render("q") + " " + footerDescStyle().Render("quit"),
	))
	b.WriteString("\n")

	return b.String()
}

func (m model) renderNoCommits() string {
	indent := indentStyle()
	var b strings.Builder

	b.WriteString("\n")
	b.WriteString(indent.Render(warningStyle().Render(
		fmt.Sprintf("No commits found between %s.", m.rangeLabel),
	)))
	b.WriteString("\n\n")
	b.WriteString(indent.Render(
		footerKeyStyle().Render("n") + " " + footerDescStyle().Render("new") + "  " +
			footerKeyStyle().Render("q") + " " + footerDescStyle().Render("quit"),
	))
	b.WriteString("\n")

	return b.String()
}

func (m model) renderError() string {
	indent := indentStyle()
	var b strings.Builder

	b.WriteString("\n")
	b.WriteString(indent.Render(errorStyle().Render("Error: " + m.err.Error())))
	b.WriteString("\n\n")
	b.WriteString(indent.Render(footerKeyStyle().Render("q") + " " + footerDescStyle().Render("quit")))
	b.WriteString("\n")

	return b.String()
}

// currentListLabels returns display strings for the current list selection state.
func (m model) currentListLabels() []string {
	switch m.state {
	case stateSelectNewerRef:
		if m.mode == "tag" {
			return m.tags
		}
		return commitLabels(m.commits)
	case stateSelectOlderRef:
		if m.mode == "tag" {
			return m.olderTags
		}
		return commitLabels(m.commits)
	}
	return nil
}

func commitLabels(commits []git.CommitInfo) []string {
	labels := make([]string, len(commits))
	for i, c := range commits {
		labels[i] = c.SHA + "  " + c.Message
	}
	return labels
}
