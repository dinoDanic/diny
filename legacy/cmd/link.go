/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/dinoDanic/diny/link"
	"github.com/spf13/cobra"
)

// linkCmd represents the link command
var linkCmd = &cobra.Command{
	Use:   "link [tool]",
	Short: "Link diny with other developer tools",
	Long: `The 'link' command integrates diny with external developer tools 
such as LazyGit. It automatically updates the tool's configuration 
so you can trigger 'diny commit' directly from within that tool.

Examples:
  # Link diny with LazyGit
  diny link lazygit

This will add a custom command to LazyGit's config.yml that lets you 
generate commit messages with diny straight from the UI.`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		tool := args[0]

		switch tool {
		case "lazygit":
			if err := link.LinkLazyGit(); err != nil {
				fmt.Fprintf(os.Stderr, "Error: %v\n", err)
				os.Exit(1)
			}
		default:
			fmt.Fprintf(os.Stderr, "Error: unknown tool '%s'\n", tool)
			fmt.Fprintf(os.Stderr, "Available tools: lazygit\n")
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(linkCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// linkCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// linkCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
