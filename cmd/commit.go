package cmd

import (
	"github.com/dinoDanic/diny/commit"

	"github.com/spf13/cobra"
)

// commitCmd represents the commit command
var commitCmd = &cobra.Command{
	Use:   "commit",
	Short: "Generate clean, conventional commit messages from staged changes",
	Long: `Diny reads your staged changes and proposes concise, conventional-friendly
commit messages. Use it to keep a tidy, consistent history—interactively or in
scripts.

Examples:
  diny commit                           # Interactive mode with options
  diny commit --print                   # Print message to stdout only
  diny commit --print | git commit -F - # Generate and commit directly
  diny commit --print | pbcopy          # Copy to clipboard (macOS)
  diny commit --print | xclip -sel clip # Copy to clipboard (Linux)
  diny commit --print | clip            # Copy to clipboard (Windows)`,
	Run: func(cmd *cobra.Command, args []string) {
		commit.Main(cmd, args)
	},
}

func init() {
	rootCmd.AddCommand(commitCmd)

	commitCmd.Flags().BoolP("print", "p", false, "Print commit message to stdout (no interactive UI)")
}
