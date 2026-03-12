package cmd

import (
	"github.com/dinoDanic/diny/commit"
	"github.com/dinoDanic/diny/config"
	"github.com/dinoDanic/diny/internal/tui/app"

	"github.com/spf13/cobra"
)

var length string
var custom string

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
  diny commit --no-werify               # skip hooks
  diny commit --length short            # Force short commit message length
  diny commit --length normal           # Force normal commit message length
  diny commit --length long             # Force long commit message length
  diny commit --print | git commit -F - # Generate and commit directly
  diny commit --print | pbcopy          # Copy to clipboard (macOS)
  diny commit --print | xclip -sel clip # Copy to clipboard (Linux)
  diny commit --print | clip            # Copy to clipboard (Windows)
  diny commit --custom "include jira ticket from branch name"`,
	Run: func(cmd *cobra.Command, args []string) {
		tuiMode, _ := cmd.Flags().GetBool("tui")
		if tuiMode {
			app.Run(AppConfig, Version)
			return
		}
		if length != "" {
			AppConfig.Commit.Length = config.Length(length)
		}
		if custom != "" {
			AppConfig.Commit.CustomInstructions = custom
		}
		commit.Main(cmd, args, AppConfig)
	},
}

func init() {
	rootCmd.AddCommand(commitCmd)

	commitCmd.Flags().BoolP("print", "p", false, "Print commit message to stdout (no interactive UI)")
	commitCmd.Flags().BoolP("no-verify", "n", false, "Skip pre-commit and commit-msg hooks")
	commitCmd.Flags().StringVarP(&length, "length", "l", "", "Override commit message length: short | normal | long")
	commitCmd.Flags().StringVarP(&custom, "custom", "c", "", "One-off custom instruction for the AI")
	commitCmd.Flags().BoolP("tui", "t", false, "Interactive TUI mode")
}
