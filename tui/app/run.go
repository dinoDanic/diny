package app

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/dinoDanic/diny/config"
)

// RunResult holds the outcome of the TUI session.
type RunResult struct {
	CommitSucceeded bool
}

func Run(cfg *config.Config, version string) RunResult {
	m := newModel(cfg, version)
	p := tea.NewProgram(m) // NO alt-screen — inline rendering
	finalModel, err := p.Run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	if fm, ok := finalModel.(model); ok && fm.state == stateSuccess {
		return RunResult{CommitSucceeded: true}
	}

	return RunResult{}
}
