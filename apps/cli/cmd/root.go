/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/dinoDanic/diny/cli/config"
	"github.com/spf13/cobra"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "diny",
	Short: "Privacy-first AI-powered git commit message generator",
	Long: `Diny is a privacy-focused CLI tool that helps you write better git commit messages.

It analyzes your staged changes and generates conventional, meaningful commit
messages while keeping your data private—no tracking, no user identification.

Commands:
  commit - Generate commit messages from staged changes
  theme  - Manage UI themes

Examples:
  diny commit                           # Interactive commit message generation
  diny commit --print                   # Non-interactive mode
  diny theme list                       # Show available themes
  diny theme set catppuccin            # Set your preferred theme`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
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
	// Call config.Load before running any command
	cobra.OnInitialize(func() {
		cfg, err := config.Load(cfgFile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error loading config: %v\n", err)
			os.Exit(1)
		}
		config.Set(cfg)
	})

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.config/diny/config.yaml)")
}
