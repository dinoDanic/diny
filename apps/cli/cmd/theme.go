/*
Copyright © 2025 dinoDanic dino.danic@gmail.com
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/charmbracelet/huh"
	"github.com/dinoDanic/diny/cli/ui"
	"github.com/spf13/cobra"
)

var themeCmd = &cobra.Command{
	Use:   "theme",
	Short: "Change the UI theme",
	Long:  `Select from available color themes for the diny UI.`,
	Run: func(cmd *cobra.Command, args []string) {
		currentTheme := ui.LoadTheme()
		if currentTheme == "" {
			currentTheme = "catppuccin"
		}

		var selectedTheme string

		darkThemes := []huh.Option[string]{
			huh.NewOption("Catppuccin Mocha (dark)", "catppuccin"),
			huh.NewOption("Tokyo Night (dark)", "tokyo"),
			huh.NewOption("Nord (dark)", "nord"),
			huh.NewOption("Dracula (dark)", "dracula"),
			huh.NewOption("Gruvbox Dark (dark)", "gruvbox-dark"),
			huh.NewOption("One Dark (dark)", "onedark"),
			huh.NewOption("Monokai (dark)", "monokai"),
			huh.NewOption("Solarized Dark (dark)", "solarized-dark"),
			huh.NewOption("Everforest Dark (dark)", "everforest-dark"),
			huh.NewOption("Flexoki Dark (dark)", "flexoki-dark"),
		}

		lightThemes := []huh.Option[string]{
			huh.NewOption("Solarized Light (light)", "solarized-light"),
			huh.NewOption("GitHub Light (light)", "github-light"),
			huh.NewOption("Gruvbox Light (light)", "gruvbox-light"),
			huh.NewOption("Flexoki Light (light)", "flexoki-light"),
		}

		allThemes := append(darkThemes, lightThemes...)

		err := huh.NewSelect[string]().
			Title("Select Theme").
			Description(fmt.Sprintf("Current: %s", ui.GetCurrentTheme().Name)).
			Options(allThemes...).
			Value(&selectedTheme).
			WithTheme(ui.GetHuhPrimaryTheme()).
			Run()

		if err != nil {
			ui.ErrorMsg(fmt.Sprintf("Error: %v", err))
			os.Exit(1)
		}

		if selectedTheme == currentTheme {
			ui.InfoMsg(fmt.Sprintf("Already using %s theme", ui.GetCurrentTheme().Name))
			return
		}

		if !ui.SetTheme(selectedTheme) {
			ui.ErrorMsg("Invalid theme selected")
			os.Exit(1)
		}

		if err := ui.SaveTheme(selectedTheme); err != nil {
			ui.ErrorMsg(fmt.Sprintf("Failed to save theme: %v", err))
			os.Exit(1)
		}

		ui.SuccessMsg(fmt.Sprintf("Theme changed to: %s ✓", ui.GetCurrentTheme().Name))
	},
}

var themeListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all available themes",
	Long:  `Display all available color themes with previews.`,
	Run: func(cmd *cobra.Command, args []string) {
		ui.PrintThemeList()
	},
}

func init() {
	rootCmd.AddCommand(themeCmd)
	themeCmd.AddCommand(themeListCmd)
}
