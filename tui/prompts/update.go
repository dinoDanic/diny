package prompts

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

// Rating

func (m ratingModel) Init() tea.Cmd { return nil }

func (m ratingModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		return m, nil
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(ratingOptions)-1 {
				m.cursor++
			}
		case "1":
			m.value = 1
			return m, tea.Quit
		case "2":
			m.value = 2
			return m, tea.Quit
		case "3":
			m.value = 3
			return m, tea.Quit
		case "0":
			m.value = 0
			return m, tea.Quit
		case "enter":
			m.value = ratingOptions[m.cursor].value
			return m, tea.Quit
		case "q", "esc", "ctrl+c":
			m.cancelled = true
			return m, tea.Quit
		}
	}
	return m, nil
}

// Star

func (m starModel) Init() tea.Cmd { return nil }

func (m starModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		return m, nil
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(starOptions)-1 {
				m.cursor++
			}
		case "1":
			m.outcome = "starred"
			return m, tea.Quit
		case "2":
			m.outcome = "already_given"
			return m, tea.Quit
		case "0":
			m.outcome = "dismissed"
			return m, tea.Quit
		case "enter":
			m.outcome = starOptions[m.cursor].key
			return m, tea.Quit
		case "q", "esc", "ctrl+c":
			m.cancelled = true
			return m, tea.Quit
		}
	}
	return m, nil
}

// Feedback

func (m feedbackModel) Init() tea.Cmd { return m.textarea.Cursor.BlinkCmd() }

func (m feedbackModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		if msg.Width > 10 {
			m.textarea.SetWidth(msg.Width - 6)
		}
		return m, nil
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			text := strings.TrimSpace(m.textarea.Value())
			if text != "" {
				m.text = text
				m.submitted = true
			} else {
				m.cancelled = true
			}
			return m, tea.Quit
		case "esc", "ctrl+c":
			m.cancelled = true
			return m, tea.Quit
		}
	}
	var cmd tea.Cmd
	m.textarea, cmd = m.textarea.Update(msg)
	return m, cmd
}
