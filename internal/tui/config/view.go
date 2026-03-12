package configtui

import (
	"fmt"
	"strings"

	"github.com/dinoDanic/diny/internal/tui/shared"
)

func (m model) View() string {
	var b strings.Builder

	b.WriteString("\n")
	b.WriteString(shared.RenderHeader(m.version, m.repoName, m.branchName, m.width))
	b.WriteString("\n")

	switch m.state {
	case statePicker:
		b.WriteString(m.renderPicker())
	case stateMenu:
		b.WriteString(m.renderMenu())
	case stateEditing:
		b.WriteString(m.renderEditing())
	}

	return b.String()
}

func (m model) renderPicker() string {
	indent := indentStyle()
	var b strings.Builder

	b.WriteString(indent.Render(sectionTitleStyle().Render("Which config to edit?")))
	b.WriteString("\n\n")

	for i, opt := range m.pickerOptions {
		cursor := "  "
		radio := "○"
		if i == m.pickerCursor {
			cursor = "> "
			radio = "●"
		}
		line := cursor + radio + " " + opt.label
		if i == m.pickerCursor {
			b.WriteString(indent.Render(sectionTitleStyle().Render(line)))
		} else {
			b.WriteString(indent.Render(metaStyle().Render(line)))
		}
		b.WriteString("\n")
	}

	b.WriteString("\n")
	keys := []struct{ key, desc string }{
		{"↑/k", "up"}, {"↓/j", "down"}, {"enter", "select"}, {"q", "quit"},
	}
	var parts []string
	for _, k := range keys {
		parts = append(parts, footerKeyStyle().Render(k.key)+" "+footerDescStyle().Render(k.desc))
	}
	b.WriteString(indent.Render(strings.Join(parts, "  ")))
	b.WriteString("\n")

	return b.String()
}

func (m model) renderMenu() string {
	indent := indentStyle()
	var b strings.Builder

	title := configTypeLabel(m.configType) + " Config"
	b.WriteString(indent.Render(sectionTitleStyle().Render(title)))
	b.WriteString("\n\n")

	for i, f := range m.fields {
		cursor := "  "
		if i == m.cursor {
			cursor = "> "
		}
		val := f.value
		if val == "" {
			val = "—"
		}
		line := fmt.Sprintf("%s%-20s %s", cursor, f.label, val)
		if i == m.cursor {
			b.WriteString(indent.Render(sectionTitleStyle().Render(line)))
		} else {
			b.WriteString(indent.Render(metaStyle().Render(line)))
		}
		b.WriteString("\n")
	}

	b.WriteString("\n")

	if m.statusMessage != "" {
		var s string
		if m.statusIsError {
			s = statusErrorStyle().Render(m.statusMessage)
		} else {
			s = statusSuccessStyle().Render(m.statusMessage)
		}
		b.WriteString(indent.Render(s))
		b.WriteString("\n")
	}

	keys := []struct{ key, desc string }{
		{"↑/k", "up"}, {"↓/j", "down"}, {"enter", "edit"}, {"s", "save"}, {"w", "save+quit"}, {"q", "quit"},
	}
	var parts []string
	for _, k := range keys {
		parts = append(parts, footerKeyStyle().Render(k.key)+" "+footerDescStyle().Render(k.desc))
	}
	b.WriteString(indent.Render(strings.Join(parts, "  ")))
	b.WriteString("\n")

	return b.String()
}

func (m model) renderEditing() string {
	indent := indentStyle()
	var b strings.Builder

	f := m.fields[m.activeField]
	b.WriteString(indent.Render(sectionTitleStyle().Render(f.label)))
	b.WriteString("\n\n")

	switch f.kind {
	case kindSelect:
		for i, opt := range f.options {
			cursor := "  "
			if i == m.optionCursor {
				cursor = "> "
			}
			line := cursor + opt
			if i == m.optionCursor {
				b.WriteString(indent.Render(sectionTitleStyle().Render(line)))
			} else {
				b.WriteString(indent.Render(metaStyle().Render(line)))
			}
			b.WriteString("\n")
		}
		b.WriteString("\n")
		keys := []struct{ key, desc string }{
			{"↑/k", "up"}, {"↓/j", "down"}, {"enter", "confirm"}, {"esc", "cancel"},
		}
		var parts []string
		for _, k := range keys {
			parts = append(parts, footerKeyStyle().Render(k.key)+" "+footerDescStyle().Render(k.desc))
		}
		b.WriteString(indent.Render(strings.Join(parts, "  ")))
		b.WriteString("\n")

	case kindText:
		b.WriteString(indent.Render(m.textinput.View()))
		b.WriteString("\n\n")
		b.WriteString(indent.Render(metaStyle().Render("enter confirm  esc cancel")))
		b.WriteString("\n")
	}

	return b.String()
}

func configTypeLabel(ct string) string {
	switch ct {
	case "global":
		return "Global"
	case "versioned":
		return "Versioned"
	case "local":
		return "Local"
	}
	return ct
}
