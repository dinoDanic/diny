package timeline

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
	case stateDateSelect:
		b.WriteString(m.renderDateSelect())
	case stateEnterDate:
		b.WriteString(m.renderTextInput("Enter date (YYYY-MM-DD):"))
	case stateEnterStartDate:
		b.WriteString(m.renderTextInput("Enter start date (YYYY-MM-DD):"))
	case stateEnterEndDate:
		b.WriteString(m.renderTextInput("Enter end date (YYYY-MM-DD):"))
	case stateFetching, stateRegenerating:
		b.WriteString(m.renderLoading())
	case stateResults:
		b.WriteString(m.renderResults())
	case stateFeedbackInput:
		b.WriteString(m.renderFeedbackInput())
	case stateNoCommits:
		b.WriteString(m.renderNoCommits())
	case stateError:
		b.WriteString(m.renderError())
	}

	return b.String()
}

var dateMenuItems = []string{"Today", "Specific date", "Date range"}

func (m model) renderDateSelect() string {
	indent := indentStyle()
	var b strings.Builder

	b.WriteString(indent.Render(sectionTitleStyle().Render("Select timeline period")))
	b.WriteString("\n\n")

	for i, item := range dateMenuItems {
		var line string
		if i == m.dateCursor {
			line = footerKeyStyle().Render("▶") + "  " +
				metaStyle().Render(fmt.Sprintf("%d", i+1)) + "  " +
				sectionTitleStyle().Render(item)
		} else {
			line = "   " +
				metaStyle().Render(fmt.Sprintf("%d", i+1)) + "  " +
				item
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

func (m model) renderTextInput(title string) string {
	indent := indentStyle()
	var b strings.Builder

	b.WriteString(indent.Render(sectionTitleStyle().Render(title)))
	b.WriteString("\n")
	b.WriteString(indent.Render(m.textinput.View()))
	b.WriteString("\n\n")
	b.WriteString(indent.Render(
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

	b.WriteString(m.renderCommitList())
	b.WriteString("\n")
	b.WriteString(m.renderAnalysis())
	b.WriteString("\n")

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
			footerKeyStyle().Render("f") + " " + footerDescStyle().Render("feedback") + "  " +
			footerKeyStyle().Render("n") + " " + footerDescStyle().Render("new") + "  " +
			footerKeyStyle().Render("q") + " " + footerDescStyle().Render("quit"),
	))
	b.WriteString("\n")

	return b.String()
}

func (m model) renderFeedbackInput() string {
	indent := indentStyle()
	var b strings.Builder

	b.WriteString(m.renderCommitList())
	b.WriteString("\n")
	b.WriteString(m.renderAnalysis())
	b.WriteString("\n")

	b.WriteString(indent.Render(sectionTitleStyle().Render("Feedback")))
	b.WriteString("\n")
	b.WriteString(indent.Render(m.textinput.View()))
	b.WriteString("\n\n")
	b.WriteString(indent.Render(
		footerKeyStyle().Render("enter") + " " + footerDescStyle().Render("submit") + "  " +
			footerKeyStyle().Render("esc") + " " + footerDescStyle().Render("cancel"),
	))
	b.WriteString("\n")

	return b.String()
}

func (m model) renderNoCommits() string {
	indent := indentStyle()
	var b strings.Builder

	b.WriteString("\n")
	b.WriteString(indent.Render(warningStyle().Render(
		fmt.Sprintf("No commits found for %s.", m.dateRange),
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

func (m model) renderCommitList() string {
	indent := indentStyle()
	var b strings.Builder

	b.WriteString("\n")
	b.WriteString(indent.Render(sectionTitleStyle().Render(
		fmt.Sprintf("Commits (%d) — %s", len(m.commits), m.dateRange),
	)))
	b.WriteString("\n")

	for i, c := range m.commits {
		line := metaStyle().Render(fmt.Sprintf("%d.", i+1)) + "  " + c
		b.WriteString(indent.Render(commitMessageStyle().Render(line)))
		b.WriteString("\n")
	}

	return b.String()
}

func (m model) renderAnalysis() string {
	indent := indentStyle()
	var b strings.Builder

	b.WriteString(indent.Render(sectionTitleStyle().Render("Analysis")))
	b.WriteString("\n")

	w := m.width - 6
	if w < 40 {
		w = 40
	}
	analysisStyle := lipgloss.NewStyle().Width(w)
	b.WriteString(indent.Render(analysisStyle.Render(m.analysis)))
	b.WriteString("\n")

	return b.String()
}

// successBigStyle is defined in styles.go but referenced here to suppress unused warning
var _ = successBigStyle
