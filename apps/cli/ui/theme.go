package ui

import (
	"github.com/dinoDanic/diny/cli/config"
	"github.com/dinoDanic/diny/cli/ui/themes"
	"github.com/spf13/viper"
)

var currentTheme *themes.Theme

func init() {
	cfg := config.Get()
	if cfg.Theme != "" {
		SetTheme(cfg.Theme)
	} else {
		currentTheme = themes.Catppuccin()
	}
}

func SetTheme(name string) bool {
	var theme *themes.Theme
	switch name {
	case "catppuccin":
		theme = themes.Catppuccin()
	case "tokyo":
		theme = themes.Tokyo()
	case "nord":
		theme = themes.Nord()
	case "dracula":
		theme = themes.Dracula()
	case "gruvbox-dark":
		theme = themes.GruvboxDark()
	case "onedark":
		theme = themes.OneDark()
	case "monokai":
		theme = themes.Monokai()
	case "solarized-dark":
		theme = themes.SolarizedDark()
	case "solarized-light":
		theme = themes.SolarizedLight()
	case "github-light":
		theme = themes.GithubLight()
	case "gruvbox-light":
		theme = themes.GruvboxLight()
	case "everforest-dark":
		theme = &themes.EverforestDark
	case "flexoki-dark":
		theme = themes.FlexokiDark()
	case "flexoki-light":
		theme = themes.FlexokiLight()
	default:
		return false
	}
	currentTheme = theme
	return true
}

func GetCurrentTheme() *themes.Theme {
	return currentTheme
}

func GetAvailableThemes() []string {
	return []string{
		"catppuccin",
		"tokyo",
		"nord",
		"dracula",
		"gruvbox-dark",
		"onedark",
		"monokai",
		"solarized-dark",
		"everforest-dark",
		"flexoki-dark",
		"solarized-light",
		"github-light",
		"gruvbox-light",
		"flexoki-light",
	}
}

func GetDarkThemes() []string {
	return []string{
		"catppuccin",
		"tokyo",
		"nord",
		"dracula",
		"gruvbox-dark",
		"onedark",
		"monokai",
		"solarized-dark",
		"everforest-dark",
		"flexoki-dark",
	}
}

func GetLightThemes() []string {
	return []string{
		"solarized-light",
		"github-light",
		"gruvbox-light",
		"flexoki-light",
	}
}

func SaveTheme(themeName string) error {
	viper.Set("theme", themeName)
	return viper.WriteConfig()
}

func LoadTheme() string {
	return viper.GetString("theme")
}
