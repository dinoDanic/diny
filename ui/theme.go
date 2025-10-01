package ui

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/dinoDanic/diny/git"
	"github.com/dinoDanic/diny/ui/themes"
)

var currentTheme *themes.Theme

func init() {
	loadedTheme := LoadTheme()
	if loadedTheme != "" {
		SetTheme(loadedTheme)
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
		"solarized-light",
		"github-light",
		"gruvbox-light",
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
	}
}

func GetLightThemes() []string {
	return []string{
		"solarized-light",
		"github-light",
		"gruvbox-light",
	}
}

func SaveTheme(themeName string) error {
	gitRoot, err := git.FindGitRoot()
	if err != nil {
		return err
	}

	dinyDir := filepath.Join(gitRoot, ".git", "diny")
	if err := os.MkdirAll(dinyDir, 0755); err != nil {
		return err
	}

	themePath := filepath.Join(dinyDir, "theme")
	return os.WriteFile(themePath, []byte(themeName), 0644)
}

func LoadTheme() string {
	gitRoot, err := git.FindGitRoot()
	if err != nil {
		return ""
	}

	themePath := filepath.Join(gitRoot, ".git", "diny", "theme")
	data, err := os.ReadFile(themePath)
	if err != nil {
		return ""
	}

	return strings.TrimSpace(string(data))
}
