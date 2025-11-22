/*
Copyright © 2025 dinoDanic dino.danic@gmail.com
*/
package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/dinoDanic/diny/config"
	"github.com/dinoDanic/diny/ui"
	"github.com/spf13/cobra"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Open config file in your default editor",
	Long: `Open the Diny configuration file in your default editor.

If the config file doesn't exist, it will be created with default values.

The configuration file is located at ~/.config/diny/config.yaml`,
	Run: func(cmd *cobra.Command, args []string) {
		openConfig()
	},
}

func openConfig() {
	configPath := config.GetConfigPath()
	if configPath == "" {
		ui.Box(ui.BoxOptions{Message: "Could not determine config path", Variant: ui.Error})
		os.Exit(1)
	}

	result, err := config.LoadOrRecover("")
	if err != nil {
		ui.Box(ui.BoxOptions{
			Message: fmt.Sprintf("Failed to load config: %v", err),
			Variant: ui.Error,
		})
		os.Exit(1)
	}

	if result.ValidationErr != "" {
		ui.Box(ui.BoxOptions{
			Title:   "Config Validation Error",
			Message: result.ValidationErr,
			Variant: ui.Error,
		})
	}
	if result.RecoveryMsg != "" {
		ui.Box(ui.BoxOptions{
			Message: result.RecoveryMsg,
			Variant: ui.Warning,
		})
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
