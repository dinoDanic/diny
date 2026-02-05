package cmd

import (
	"github.com/dinoDanic/diny/changelog"
	"github.com/spf13/cobra"
)

var changelogCmd = &cobra.Command{
	Use:   "changelog",
	Short: "View GitHub release changelogs",
	Long: `View and interact with GitHub release changelogs.

This command fetches releases from the GitHub repository,
displays them in an interactive list, and allows you to:
- View changelog details
- Open releases in your browser
- Copy changelog content to clipboard
- Browse through multiple releases`,
	Run: func(cmd *cobra.Command, args []string) {
		changelog.Main(AppConfig)
	},
}

func init() {
	rootCmd.AddCommand(changelogCmd)
}
