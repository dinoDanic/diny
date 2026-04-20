package app

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/dinoDanic/diny/config"
)

// Options carries flags from the cobra layer into the TUI.
type Options struct {
	NoVerify bool
	Push     bool
	Print    bool
}

// RunResult holds the outcome of the TUI session.
type RunResult struct {
	CommitSucceeded bool
}

func Run(cfg *config.Config, version string, opts Options) RunResult {
	m := newModel(cfg, version, opts)
	p := tea.NewProgram(m) // NO alt-screen — inline rendering
	finalModel, err := p.Run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	if fm, ok := finalModel.(model); ok && (fm.state == stateSuccess || fm.state == stateSplitSuccess) {
		return RunResult{CommitSucceeded: true}
	}

	return RunResult{}
}
