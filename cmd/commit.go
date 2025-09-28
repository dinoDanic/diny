package cmd

import (
	"github.com/dinoDanic/diny/commit"

	"github.com/spf13/cobra"
)

// commitCmd represents the commit command
var commitCmd = &cobra.Command{
	Use:   "commit",
	Short: "Generate commit messages from staged changes",
	Long: `Diny analyzes your staged git changes and generates clear,
well-formatted commit messages using AI. 

This helps you keep a clean, consistent commit history with less effort.

Examples:
  diny commit
  diny commit --lang hr
  diny commit --style conventional`,
	Run: func(cmd *cobra.Command, args []string) {
		commit.Main(cmd, args)
	},
}

func init() {
	rootCmd.AddCommand(commitCmd)
}
