package cmd

import (
	"github.com/dinoDanic/diny/commit"

	"github.com/spf13/cobra"
)

// commitCmd represents the commit command
var commitCmd = &cobra.Command{
	Use:   "commit",
	Short: "Generate AI-powered commit messages from staged changes",
	Long: `Diny analyzes your staged git changes and generates clear,
well-formatted commit messages using advanced AI models. 

This helps you maintain a clean, consistent commit history with minimal effort
while following best practices and conventional commit standards.

Features:
  • Smart analysis of code changes and context
  • Multiple language support for international teams
  • Conventional commit format support
  • Custom commit message styles and templates
  • Integration with popular git workflows

Examples:
  diny commit                           # Generate commit with default settings
  diny commit --lang hr                 # Generate in Croatian language
  diny commit --style conventional      # Use conventional commit format
  diny commit --template feature        # Use feature template
  diny commit --auto-stage              # Auto-stage changes before committing`,
	Run: func(cmd *cobra.Command, args []string) {
		commit.Main(cmd, args)
	},
}

func init() {
	rootCmd.AddCommand(commitCmd)

	// Add language flag for internationalization
	commitCmd.Flags().StringP("lang", "l", "en", "Language for commit message (en, hr, es, fr, de)")

	// Add style flags for different commit formats
	commitCmd.Flags().StringP("style", "s", "default", "Commit message style (default, conventional, angular, gitmoji)")

	// Add template flag for predefined commit types
	commitCmd.Flags().StringP("template", "t", "", "Use predefined template (feature, bugfix, refactor, docs, test)")

	// Add auto-stage flag
	commitCmd.Flags().BoolP("auto-stage", "a", false, "Automatically stage all changes before committing")

	// Add dry-run flag
	commitCmd.Flags().BoolP("dry-run", "d", false, "Show generated commit message without committing")

	// Add word limit flag
	commitCmd.Flags().IntP("max-length", "m", 72, "Maximum commit message length (50-100 characters)")
}
