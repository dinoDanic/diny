package cmd

import (
	"fmt"
	"os"

	"github.com/dinoDanic/diny/config"
	"github.com/dinoDanic/diny/ui"
	"github.com/spf13/cobra"
)

var showConfigCmd = &cobra.Command{
	Use:   "config",
	Short: "Show current Diny configuration",
	Long: `Display the current Diny configuration settings for commit message generation.

Shows the effective configuration after merging global and local configs,
with environment variable overrides applied.

The configuration includes:
- Emoji: Whether to use emoji prefixes in commit messages
- Conventional: Whether to use Conventional Commits format
- Tone: Professional, casual, or friendly language style
- Length: Short, normal, or detailed commit message length
- API settings: URLs and models being used

Configuration precedence: env vars > local JSON > global JSON > defaults`,
	Run: func(cmd *cobra.Command, args []string) {
		verbose, _ := cmd.Flags().GetBool("verbose")
		showUserConfig(verbose)
	},
}

func showUserConfig(verbose bool) {
	configService := config.GetService()
	if err := configService.LoadUserConfig(); err != nil {
		ui.Box(ui.BoxOptions{Message: "Error loading configuration", Variant: ui.Error})
		os.Exit(1)
	}

	userConfig := configService.GetUserConfig()

	if verbose {
		config.PrintEffectiveConfiguration(*userConfig)
		fmt.Println()
		showConfigPaths()
	} else {
		config.PrintConfiguration(*userConfig)
	}
}

func showConfigPaths() {
	ui.RenderTitle("Configuration File Locations")

	globalPath, _ := config.GetGlobalConfigPath()
	localPath, _ := config.GetLocalConfigPath()

	fmt.Printf("Global: %s\n", globalPath)
	if customPath := os.Getenv("DINY_CONFIG_PATH"); customPath != "" {
		fmt.Println("        [overridden by DINY_CONFIG_PATH]")
	}
	if _, err := os.Stat(globalPath); err == nil {
		fmt.Println("        [exists]")
	} else {
		fmt.Println("        [not found]")
	}

	fmt.Printf("\nLocal:  %s\n", localPath)
	if _, err := os.Stat(localPath); err == nil {
		fmt.Println("        [exists]")
	} else {
		fmt.Println("        [not found - using global or defaults]")
	}
}

func runInitSetup() {
	ui.RenderTitle("Starting Diny configuration setup...")

	userConfig := RunConfigurationSetup(nil)

	err := config.Save(userConfig)
	if err != nil {
		fmt.Printf("Error saving configuration: %v\n", err)
		os.Exit(1)
	}

	config.PrintConfiguration(userConfig)
}

func init() {
	rootCmd.AddCommand(showConfigCmd)

	showConfigCmd.Flags().BoolP("verbose", "v", false, "Show effective configuration with source precedence")
}
