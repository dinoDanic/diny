/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/charmbracelet/huh"
	"github.com/dinoDanic/diny/config"
	"github.com/dinoDanic/diny/git"
	"github.com/dinoDanic/diny/ui"
	"github.com/spf13/cobra"
)

// showConfigCmd represents the showConfig command
var showConfigCmd = &cobra.Command{
	Use:   "config",
	Short: "Show current Diny configuration",
	Long: `Display the current Diny configuration settings for commit message generation.

If no configuration exists, you'll be prompted to create one through the interactive setup.

The configuration includes:
- Emoji: Whether to use emoji prefixes in commit messages
- Conventional: Whether to use Conventional Commits format
- Tone: Professional, casual, or friendly language style
- Length: Short, normal, or detailed commit message length

Configuration is stored in .git/diny-config.json in your git repository.`,
	Run: func(cmd *cobra.Command, args []string) {
		showUserConfig()
	},
}

func showUserConfig() {
	gitRoot, gitErr := git.FindGitRoot()
	if gitErr != nil {
		fmt.Println("❌ Error: Not in a git repository")
		fmt.Println("Please run this command from within a git repository.")
		os.Exit(1)
	}

	configExists := configFileExists(gitRoot)

	if !configExists {
		fmt.Println()

		var createConfig bool
		err := huh.NewConfirm().
			Title("Would you like to create a configuration now?").
			Description("This will start the interactive setup process").
			Affirmative("Yes, let's configure Diny").
			Negative("No, exit").
			Value(&createConfig).
			Run()

		if err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}

		if !createConfig {
			ui.RenderTitle("Configuration setup cancelled.")
			return
		}

		runInitSetup()
		return
	}

	userConfig, err := config.Load()
	if err != nil {
		ui.RenderError("Error loading configuration")
		os.Exit(1)
	}
	if userConfig != nil {
		fmt.Println()
		config.PrintConfiguration(*userConfig)
	}
}

func configFileExists(gitRoot string) bool {
	configPath := filepath.Join(gitRoot, ".git", "diny-config.json")
	_, err := os.Stat(configPath)
	return err == nil
}

func runInitSetup() {
	ui.RenderTitle("Starting Diny configuration setup...")

	userConfig := RunConfigurationSetup()

	err := config.Save(userConfig)
	if err != nil {
		fmt.Printf("Error saving configuration: %v\n", err)
		os.Exit(1)
	}

	config.PrintConfiguration(userConfig)
}

func init() {
	rootCmd.AddCommand(showConfigCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// showConfigCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// showConfigCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
