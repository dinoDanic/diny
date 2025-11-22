/*
Copyright © 2025 dinoDanic dino.danic@gmail.com
*/
package cmd

import (
	"github.com/dinoDanic/diny/ui"
	"github.com/spf13/cobra"
)

var themeCmd = &cobra.Command{
	Use:   "theme",
	Short: "List available UI themes",
	Long: `Display all available color themes with previews.

To change your theme, edit ~/.config/diny/config.yaml and set the 'theme' field.`,
	Run: func(cmd *cobra.Command, args []string) {
		ui.PrintThemeList()
	},
}

func init() {
	rootCmd.AddCommand(themeCmd)
}
