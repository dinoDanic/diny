/*
Copyright © 2025 dinoDanic dino.danic@gmail.com
*/
package cmd

import (
	"github.com/dinoDanic/diny/config"
	"github.com/dinoDanic/diny/git"
	configtui "github.com/dinoDanic/diny/tui/config"
	"github.com/spf13/cobra"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Open interactive TUI config editor",
	Long: `Open the Diny configuration file in an interactive TUI editor.

If in a git repository, you can choose between:
  - Global config: ~/.config/diny/config.yaml (applies to all projects)
  - Project config - versioned: .diny.yaml (can be committed, shared with team)
  - Project config - local: <gitdir>/diny/config.yaml (never committed, personal overrides)

Config priority: local > versioned > global (higher priority overrides lower)

If not in a git repository, only global config is available.`,
	Run: func(cmd *cobra.Command, args []string) {
		openConfigTUI()
	},
}

func openConfigTUI() {
	versionedPath := config.GetVersionedProjectConfigPath()
	localPath := config.GetLocalProjectConfigPath()
	inGitRepo := versionedPath != "" && localPath != ""

	repoName := git.GetRepoName()
	branchName, _ := git.GetCurrentBranch()

	var configPath, configType string
	var cfg *config.Config

	if !inGitRepo {
		configPath = config.GetConfigPath()
		configType = "global"
		var err error
		cfg, err = config.Load(configPath)
		if err != nil {
			cfg = &config.Config{}
		}
	}

	configtui.Run(Version, repoName, branchName, configPath, configType, cfg)
}

func init() {
	rootCmd.AddCommand(configCmd)
}
