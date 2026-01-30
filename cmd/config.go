/*
Copyright © 2025 dinoDanic dino.danic@gmail.com
*/
package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/charmbracelet/huh"
	"github.com/dinoDanic/diny/config"
	"github.com/dinoDanic/diny/ui"
	"github.com/spf13/cobra"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Open config file in your default editor",
	Long: `Open the Diny configuration file in your default editor.

If in a git repository, you can choose between:
  - Global config: ~/.config/diny/config.yaml (applies to all projects)
  - Project config - versioned: .diny.yaml (can be committed, shared with team)
  - Project config - local: <gitdir>/diny/config.yaml (never committed, personal overrides)

Config priority: local > versioned > global (higher priority overrides lower)

Project configs overlay on top of global config, allowing per-project customization.

If not in a git repository, only global config is available.

If the config file doesn't exist, it will be created with default values.`,
	Run: func(cmd *cobra.Command, args []string) {
		openConfig()
	},
}

func openConfig() {
	versionedPath := config.GetVersionedProjectConfigPath()
	localPath := config.GetLocalProjectConfigPath()
	inGitRepo := versionedPath != "" && localPath != ""

	var configPath string
	var configType string

	if inGitRepo {
		var choice string
		err := huh.NewSelect[string]().
			Title("Which config would you like to edit?").
			Description("Select using arrow keys or j,k and press Enter").
			Options(
				huh.NewOption("Global config (~/.config/diny/config.yaml)", "global"),
				huh.NewOption("Project config - versioned (.diny.yaml)", "versioned"),
				huh.NewOption("Project config - local only (<gitdir>/diny/config.yaml)", "local"),
			).
			Value(&choice).
			WithTheme(ui.GetHuhPrimaryTheme()).
			Run()

		if err != nil {
			ui.Box(ui.BoxOptions{
				Message: fmt.Sprintf("Error running prompt: %v", err),
				Variant: ui.Error,
			})
			os.Exit(1)
		}

		switch choice {
		case "versioned":
			if err := config.CreateVersionedProjectConfigIfNeeded(); err != nil {
				ui.Box(ui.BoxOptions{
					Message: fmt.Sprintf("Failed to create versioned project config: %v", err),
					Variant: ui.Error,
				})
				os.Exit(1)
			}
			configPath = versionedPath
			configType = "versioned project"
		case "local":
			if err := config.CreateLocalProjectConfigIfNeeded(); err != nil {
				ui.Box(ui.BoxOptions{
					Message: fmt.Sprintf("Failed to create local project config: %v", err),
					Variant: ui.Error,
				})
				os.Exit(1)
			}
			configPath = localPath
			configType = "local project"
		default: 
			configPath = config.GetConfigPath()
			configType = "global"
		}
	} else {
		configPath = config.GetConfigPath()
		configType = "global"
	}

	if configPath == "" {
		ui.Box(ui.BoxOptions{
			Message: "Could not determine config path",
			Variant: ui.Error,
		})
		os.Exit(1)
	}

	editor := getEditor()
	if editor == "" {
		ui.Box(ui.BoxOptions{
			Message: "No editor found. Set $EDITOR or $VISUAL environment variable",
			Variant: ui.Error,
		})
		os.Exit(1)
	}

	editorArgs := strings.Fields(editor)
	editorCmd := editorArgs[0]
	args := append(editorArgs[1:], configPath)

	execCmd := exec.Command(editorCmd, args...)
	execCmd.Stdin = os.Stdin
	execCmd.Stdout = os.Stdout
	execCmd.Stderr = os.Stderr

	if err := execCmd.Run(); err != nil {
		ui.Box(ui.BoxOptions{
			Message: fmt.Sprintf("Error opening editor: %v", err),
			Variant: ui.Error,
		})
		os.Exit(1)
	}

	ui.Box(ui.BoxOptions{
		Message: fmt.Sprintf("Saved %s config", configType),
		Variant: ui.Success,
	})
}

func getEditor() string {
	if editor := os.Getenv("VISUAL"); editor != "" {
		return editor
	}
	if editor := os.Getenv("EDITOR"); editor != "" {
		return editor
	}
	for _, editor := range []string{"vim", "vi", "nano"} {
		if _, err := exec.LookPath(editor); err == nil {
			return editor
		}
	}
	return ""
}

func init() {
	rootCmd.AddCommand(configCmd)
}
