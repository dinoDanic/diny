package commitui

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/dinoDanic/diny/config"
)

// New returns a commitui model as a tea.Model for embedding in a parent TUI.
func New(cfg *config.Config, noVerify bool) tea.Model {
	return newModel(cfg, noVerify)
}

func Run(cfg *config.Config, noVerify bool) {
	m := newModel(cfg, noVerify)
	p := tea.NewProgram(m, tea.WithAltScreen())

	finalModel, err := p.Run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error running TUI: %v\n", err)
		os.Exit(1)
	}

	if fm, ok := finalModel.(model); ok {
		if fm.statusMessage != "" && !fm.statusIsError {
			fmt.Println(fm.statusMessage)
		}
		if fm.err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", fm.err)
			os.Exit(1)
		}
	}
}
