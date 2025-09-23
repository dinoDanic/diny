package cmd

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

type confirmModel struct {
	message string
	choice  bool
	done    bool
}

func (m confirmModel) Init() tea.Cmd {
	return nil
}

func (m confirmModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "y", "Y":
			m.choice = true
			m.done = true
			return m, tea.Quit
		case "n", "N", "q", "ctrl+c", "esc":
			m.choice = false
			m.done = true
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m confirmModel) View() string {
	if m.done {
		return ""
	}
	return fmt.Sprintf("%s\n\n(y)es / (n)o: ", m.message)
}

func confirmPrompt(message string) bool {
	model := confirmModel{message: message}
	p := tea.NewProgram(model)

	finalModel, err := p.Run()
	if err != nil {
		fmt.Printf("Error running prompt: %v\n", err)
		os.Exit(1)
	}

	return finalModel.(confirmModel).choice
}
