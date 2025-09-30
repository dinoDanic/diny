/*
Copyright © 2025 NAME HERE dino.danic@gmail.com
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/dinoDanic/diny/ui"
	"github.com/dinoDanic/diny/update"
	"github.com/spf13/cobra"
)

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update diny to the latest version",
	Long: `Update diny to the latest version.

This command will:
- Detect your installation method (Homebrew, Scoop, Winget)
- Update using the appropriate package manager
- Show manual instructions if installed manually

Examples:
  diny update
  diny update --force    # Force update even if already latest
`,
	Run: func(cmd *cobra.Command, args []string) {
		force, _ := cmd.Flags().GetBool("force")
		runUpdate(force)
	},
}

func init() {
	rootCmd.AddCommand(updateCmd)
	updateCmd.Flags().BoolP("force", "f", false, "Force update even if already on latest version")
}

func runUpdate(force bool) {
	ui.RenderTitle("Checking for diny updates...")

	checker := update.NewUpdateChecker(Version)
	latestVersion, err := checker.GetLatestVersion()
	if err != nil {
		ui.RenderError(fmt.Sprintf("Failed to check for updates: %v", err))
		os.Exit(1)
	}

	if !force && !checker.CompareVersions(Version, latestVersion) {
		ui.RenderSuccess(fmt.Sprintf("You're already on the latest version (%s)", Version))
		return
	}

	method := checker.DetectInstallMethod()
	ui.RenderBox("Update Available", fmt.Sprintf("Updating from %s to %s...\nInstalled via: %s", Version, latestVersion, method))

	if err := checker.PerformUpdate(); err != nil {
		ui.RenderError(fmt.Sprintf("Update failed: %v", err))
		os.Exit(1)
	}
}
