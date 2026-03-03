package changelog

import (
	"fmt"
	"os"

	"github.com/charmbracelet/huh"
	"github.com/dinoDanic/diny/config"
	"github.com/dinoDanic/diny/ui"
)

func Main(cfg *config.Config) {
	var mode string

	err := huh.NewSelect[string]().
		Title("Generate a changelog for...").
		Description("Select using arrow keys or j,k and press Enter").
		Options(
			huh.NewOption("Between two tags", "tag"),
			huh.NewOption("Between two commits", "commit"),
		).
		Value(&mode).
		Height(4).
		WithTheme(ui.GetHuhPrimaryTheme()).
		Run()

	if err != nil {
		os.Exit(0)
	}

	switch mode {
	case "tag":
		if err := GenerateByTag(cfg); err != nil {
			ui.Box(ui.BoxOptions{Message: fmt.Sprintf("Error: %v", err), Variant: ui.Error})
		}
	case "commit":
		if err := GenerateByCommit(cfg); err != nil {
			ui.Box(ui.BoxOptions{Message: fmt.Sprintf("Error: %v", err), Variant: ui.Error})
		}
	}
}
