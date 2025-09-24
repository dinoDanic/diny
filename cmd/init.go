/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/charmbracelet/huh"
	"github.com/dinoDanic/diny/config"
	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize Diny configuration with an interactive setup",
	Long: `Initialize Diny configuration with an interactive setup.

This command will guide you through configuring your commit message preferences:
- Emoji: Add emoji prefixes to commit messages
- Format: Conventional commits or free-form messages
- Tone: Professional, casual, or friendly
- Length: Short, normal, or detailed messages

The configuration will be saved to .git/diny-config.json in your git repository.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Start with default configuration values
		userConfig := config.DefaultUserConfig

		// Emoji confirmation - use proper confirm component
		err := huh.NewConfirm().
			Title("Use emoji prefixes in commit messages?").
			Description("Add emojis like ✨ feat: or 🐛 fix: to commit messages").
			Affirmative("Yes").
			Negative("No").
			Value(&userConfig.UseEmoji).
			Run()

		if err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("UseEmoji: %t\n", userConfig.UseEmoji)

		// Conventional commits confirmation - use proper confirm component
		err = huh.NewConfirm().
			Title("Use Conventional Commits format?").
			Description("Format: type(scope): description").
			Affirmative("Yes").
			Negative("No").
			Value(&userConfig.UseConventional).
			Run()

		if err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("UseConventional: %t\n", userConfig.UseConventional)

		// Tone selection - use actual config types
		err = huh.NewSelect[config.Tone]().
			Title("Choose your commit message tone").
			Options(
				huh.NewOption("Professional - formal and matter-of-fact", config.Professional),
				huh.NewOption("Casual - light but clear", config.Casual),
				huh.NewOption("Friendly - warm and approachable", config.Friendly),
			).
			Value(&userConfig.Tone).
			Run()

		if err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Tone: %s\n", userConfig.Tone)

		// Length selection - use actual config types
		err = huh.NewSelect[config.Length]().
			Title("Choose your commit message length").
			Options(
				huh.NewOption("Short - subject only (no body)", config.Short),
				huh.NewOption("Normal - subject + optional body (1-4 bullets)", config.Normal),
				huh.NewOption("Long - subject + detailed body (2-6 bullets)", config.Long),
			).
			Value(&userConfig.Length).
			Run()

		if err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Length: %s\n", userConfig.Length)

		// Save the configuration
		err = config.Save(userConfig)
		if err != nil {
			fmt.Printf("Error saving configuration: %v\n", err)
			os.Exit(1)
		}

		// Print the configuration for testing
		fmt.Println()
		fmt.Println("🎉 Configuration Complete!")
		fmt.Println("===========================")
		fmt.Printf("Use Emoji: %t\n", userConfig.UseEmoji)
		fmt.Printf("Use Conventional: %t\n", userConfig.UseConventional)
		fmt.Printf("Tone: %s\n", userConfig.Tone)
		fmt.Printf("Length: %s\n", userConfig.Length)
		fmt.Printf("\nConfiguration saved to .git/diny-config.json\n")
	},
}

func init() {
	rootCmd.AddCommand(initCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// initCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// initCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
