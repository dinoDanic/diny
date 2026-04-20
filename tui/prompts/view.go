package prompts

import (
	"fmt"
	"strings"
)

func (m ratingModel) View() string {
	indent := indentStyle()
	var b strings.Builder

	b.WriteString(indent.Render(sectionTitleStyle().Render("How's diny working for you?")))
	b.WriteString("\n\n")

	for i, opt := range ratingOptions {
		var line string
		numLabel := fmt.Sprintf("%d", opt.value)
		if i == m.cursor {
			line = cursorStyle().Render("▶") + "  " +
				cursorStyle().Render(numLabel) + "  " +
				sectionTitleStyle().Render(opt.label)
		} else {
			line = "   " + metaStyle().Render(numLabel) + "  " + opt.label
		}
		b.WriteString(indent.Render(line))
		b.WriteString("\n")
	}

	b.WriteString("\n")
	b.WriteString(indent.Render(
		footerKeyStyle().Render("j/k") + " " + footerDescStyle().Render("move") + "  " +
			footerKeyStyle().Render("enter") + " " + footerDescStyle().Render("confirm") + "  " +
			footerKeyStyle().Render("1-3") + " " + footerDescStyle().Render("quick pick") + "  " +
			footerKeyStyle().Render("0") + " " + footerDescStyle().Render("never again") + "  " +
			footerKeyStyle().Render("esc") + " " + footerDescStyle().Render("close"),
	))
	b.WriteString("\n")

	return b.String()
}

func (m starModel) View() string {
	indent := indentStyle()
	var b strings.Builder

	b.WriteString(indent.Render(sectionTitleStyle().Render("Diny is free & open source. A GitHub star helps a lot!")))
	b.WriteString("\n\n")

	for i, opt := range starOptions {
		var line string
		numLabel := string(opt.quickKey)
		if i == m.cursor {
			line = cursorStyle().Render("▶") + "  " +
				cursorStyle().Render(numLabel) + "  " +
				sectionTitleStyle().Render(opt.label)
		} else {
			line = "   " + metaStyle().Render(numLabel) + "  " + opt.label
		}
		b.WriteString(indent.Render(line))
		b.WriteString("\n")
	}

	b.WriteString("\n")
	b.WriteString(indent.Render(
		footerKeyStyle().Render("j/k") + " " + footerDescStyle().Render("move") + "  " +
			footerKeyStyle().Render("enter") + " " + footerDescStyle().Render("confirm") + "  " +
			footerKeyStyle().Render("1/2") + " " + footerDescStyle().Render("quick pick") + "  " +
			footerKeyStyle().Render("0") + " " + footerDescStyle().Render("never again") + "  " +
			footerKeyStyle().Render("esc") + " " + footerDescStyle().Render("close"),
	))
	b.WriteString("\n")

	return b.String()
}

func (m feedbackModel) View() string {
	indent := indentStyle()
	var b strings.Builder

	b.WriteString(indent.Render(sectionTitleStyle().Render("Got any feedback or feature requests?")))
	b.WriteString(" ")
	b.WriteString(metaStyle().Render("(optional)"))
	b.WriteString("\n\n")

	b.WriteString(indent.Render(m.textarea.View()))
	b.WriteString("\n\n")

	b.WriteString(indent.Render(
		footerKeyStyle().Render("enter") + " " + footerDescStyle().Render("submit") + "  " +
			footerKeyStyle().Render("esc") + " " + footerDescStyle().Render("skip") + "  " +
			footerKeyStyle().Render("ctrl+c") + " " + footerDescStyle().Render("quit"),
	))
	b.WriteString("\n")

	return b.String()
}
