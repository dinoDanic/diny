/*
Copyright © 2025 dinoDanic dino.danic@gmail.com
*/
package cmd

import (
	"github.com/dinoDanic/diny/cli/ui"
	"github.com/spf13/cobra"
)

var demoCmd = &cobra.Command{
	Use:   "demo",
	Short: "Demonstrate UI components",
	Long:  `Display examples of all available UI components and styles.`,
	Run: func(cmd *cobra.Command, args []string) {
		ui.RenderTitle("UI Components Demo")

		// Old style (still works)
		ui.Box(ui.BoxOptions{
			Title:   "Old Style",
			Message: "This is the traditional Box() function with BoxOptions struct",
			Variant: ui.Primary,
		})

		// New convenience functions with title and message
		ui.PrimaryBox("Primary Box", "This is a primary box using PrimaryBox() convenience function")
		ui.SuccessBox("Success Box", "This is a success box using SuccessBox() convenience function")
		ui.ErrorBox("Error Box", "This is an error box using ErrorBox() convenience function")
		ui.WarningBox("Warning Box", "This is a warning box using WarningBox() convenience function")

		// New convenience functions with message only
		ui.InfoMsg("This is an info message using InfoMsg()")
		ui.SuccessMsg("This is a success message using SuccessMsg()")
		ui.ErrorMsg("This is an error message using ErrorMsg()")
		ui.WarningMsg("This is a warning message using WarningMsg()")

		// Title only
		ui.RenderTitle("This is just a title using RenderTitle()")
	},
}

func init() {
	rootCmd.AddCommand(demoCmd)
}
