package configtui

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/dinoDanic/diny/config"
)

// Run launches the interactive config TUI.
// Pass empty configPath/configType when in a git repo to show the picker first.
// Pass a resolved configPath/configType (with non-nil cfg) to go straight to the menu.
func Run(version, repoName, branchName, configPath, configType string, cfg *config.Config) {
	m := newModel(version, repoName, branchName, configPath, configType, cfg)
	p := tea.NewProgram(m) // inline rendering, no alt-screen
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
