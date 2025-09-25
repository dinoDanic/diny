/*
Copyright © 2025 NAME HERE dino.danic@gmail.com
*/
package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

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
