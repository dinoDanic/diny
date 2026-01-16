/*
Copyright © 2025 NAME HERE dino.danic@gmail.com
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/dinoDanic/diny/config"
	"github.com/dinoDanic/diny/ui"
	"github.com/dinoDanic/diny/update"
	"github.com/spf13/cobra"
)

var AppConfig *config.Config

// Version will be set at build time via ldflags
var Version = "dev"

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "diny",
	Short: "Generate meaningful commit messages from git diff",
	Long: `diny is a simple CLI tool that analyzes your git diff 
and generates commit messages. 

It helps you maintain clean, consistent commit history without 
spending time manually writing messages.
`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		fmt.Println()

		if cmd.Name() == "theme" {
			return
		}

		result, err := config.LoadOrRecoverWithProject("")
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

		AppConfig = result.Config

		if AppConfig != nil && AppConfig.Theme != "" {
			ui.SetTheme(AppConfig.Theme)
		}
	},
	PersistentPostRun: func(cmd *cobra.Command, args []string) {
		checker := update.NewUpdateChecker(Version)
		checker.CheckForUpdate()
		fmt.Println()
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Set version for the root command
	rootCmd.Version = Version

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.diny.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
