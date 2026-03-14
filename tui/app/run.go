package app

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/dinoDanic/diny/config"
)

func Run(cfg *config.Config, version string) {
	m := newModel(cfg, version)
	p := tea.NewProgram(m) // NO alt-screen — inline rendering
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
