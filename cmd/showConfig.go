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

// showUserConfig displays the current configuration or prompts to create one
func showUserConfig() {
	// Check if we're in a git repository and get config
	gitRoot, gitErr := git.FindGitRoot()
	if gitErr != nil {
		fmt.Println("❌ Error: Not in a git repository")
		fmt.Println("Please run this command from within a git repository.")
		os.Exit(1)
	}

	// Check if config file exists
	configExists := configFileExists(gitRoot)

	if !configExists {
		fmt.Println("🔧 No configuration found!")
		fmt.Println("Diny needs to be configured before use.")
		fmt.Println()

		// Ask user if they want to create config
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
			fmt.Println("👋 Configuration setup cancelled.")
			fmt.Println("Run 'diny init' when you're ready to configure.")
			return
		}

		// Run the init command logic
		runInitSetup()
		return
	}

	// Config exists, load and display it
	userConfig := config.Load()
	displayConfig(userConfig)
}

// configFileExists checks if the config file exists
func configFileExists(gitRoot string) bool {
	configPath := filepath.Join(gitRoot, ".git", "diny-config.json")
	_, err := os.Stat(configPath)
	return err == nil
}

// runInitSetup runs the init command setup logic
func runInitSetup() {
	fmt.Println("🚀 Starting Diny configuration setup...")
	fmt.Println()

	userConfig := RunConfigurationSetup()

	// Save the configuration
	err := config.Save(userConfig)
	if err != nil {
		fmt.Printf("Error saving configuration: %v\n", err)
		os.Exit(1)
	}

	fmt.Println()
	fmt.Println("🎉 Configuration saved successfully!")
	displayConfig(userConfig)
}

// displayConfig shows the current configuration
func displayConfig(userConfig config.UserConfig) {
	fmt.Println()
	fmt.Println("⚙️  Diny Configuration")
	fmt.Println("=======================")
	fmt.Printf("📁 Location: .git/diny-config.json\n")
	fmt.Println()
	fmt.Printf("🎨 Use Emoji: %t\n", userConfig.UseEmoji)
	fmt.Printf("📋 Conventional: %t\n", userConfig.UseConventional)
	fmt.Printf("💬 Tone: %s\n", userConfig.Tone)
	fmt.Printf("📏 Length: %s\n", userConfig.Length)
	fmt.Println()
	fmt.Println("💡 To modify configuration, run: diny init")
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
