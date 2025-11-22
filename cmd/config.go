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
	// Load config to ensure it exists (creates default if not)
	_, err := config.Load("")
	if err != nil {
		fmt.Printf("Error loading config: %v\n", err)
		os.Exit(1)
	}

	configPath := config.GetConfigPath()
	if configPath == "" {
		fmt.Println("Error: could not determine config path")
		os.Exit(1)
	}

	editor := getEditor()
	if editor == "" {
		fmt.Println("Error: no editor found. Set $EDITOR or $VISUAL environment variable")
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
		fmt.Printf("Error opening editor: %v\n", err)
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
	// Fallback to common editors
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
