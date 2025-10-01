package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/charmbracelet/huh"
	"github.com/dinoDanic/diny/git"
	"github.com/dinoDanic/diny/ui"
)

type Tone string
type Length string

const (
	Professional Tone = "professional"
	Casual       Tone = "casual"
	Friendly     Tone = "friendly"
)

const (
	Short  Length = "short"
	Normal Length = "normal"
	Long   Length = "long"
)

type UserConfig struct {
	UseConventional bool   `json:"useConventional"`
	UseEmoji        bool   `json:"useEmoji"`
	Tone            Tone   `json:"tone"`
	Length          Length `json:"length"`
}

func Load() (*UserConfig, error) {
	gitRoot, err := git.FindGitRoot()
	if err != nil {
		return nil, err
	}

	configPath := filepath.Join(gitRoot, ".git", "diny-config.json")

	config, err := tryLoadConfig(configPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return handleConfigError(configPath, err)
	}

	if !isValidConfig(config) {
		return handleInvalidConfig(configPath, config)
	}

	return config, nil
}

func tryLoadConfig(configPath string) (*UserConfig, error) {
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	var config UserConfig
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	return &config, nil
}

func isValidConfig(config *UserConfig) bool {
	validTones := []Tone{Professional, Casual, Friendly}
	if !contains(validTones, config.Tone) {
		return false
	}

	validLengths := []Length{Short, Normal, Long}
	if !contains(validLengths, config.Length) {
		return false
	}

	return true
}

func contains[T comparable](slice []T, item T) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

func handleConfigError(configPath string, err error) (*UserConfig, error) {
	ui.Box(ui.BoxOptions{Message: fmt.Sprintf("Configuration file is corrupted: %v", err), Variant: ui.Error})
	return promptConfigAction(configPath)
}

func handleInvalidConfig(configPath string, config *UserConfig) (*UserConfig, error) {
	ui.Box(ui.BoxOptions{Message: "Invalid configuration values detected!", Variant: ui.Warning})
	return promptConfigAction(configPath)
}

func promptConfigAction(configPath string) (*UserConfig, error) {
	var choice string

	err := huh.NewSelect[string]().
		Title("🦕 What would you like to do?").
		Description("Choose how to handle the invalid configuration").
		Options(
			huh.NewOption("Fix configuration interactively", "fix"),
			huh.NewOption("Use defaults", "defaults"),
		).
		Value(&choice).
		Height(4).
		Run()

	if err != nil {
		return nil, fmt.Errorf("error in prompt: %w", err)
	}

	switch choice {
	case "fix":
		return promptUserForValidConfig(configPath)
	case "defaults":
		if err := os.Remove(configPath); err != nil {
			ui.Box(ui.BoxOptions{Message: fmt.Sprintf("Could not delete config file: %v", err), Variant: ui.Warning})
		} else {
			ui.RenderTitle("Using defaults..")
		}
		return nil, nil
	default:
		return nil, fmt.Errorf("invalid choice")
	}
}

func promptUserForValidConfig(configPath string) (*UserConfig, error) {
	ui.RenderTitle("Let's create a new configuration")

	config := &UserConfig{}

	err := huh.NewConfirm().
		Title("Use conventional commit format?").
		Description("e.g., 'feat: add new feature' or 'fix: resolve bug'").
		Value(&config.UseConventional).
		Run()
	if err != nil {
		return nil, err
	}

	err = huh.NewConfirm().
		Title("Include emojis in commit messages?").
		Description("e.g., '✨ Add new feature' or '🐛 Fix bug'").
		Value(&config.UseEmoji).
		Run()
	if err != nil {
		return nil, err
	}

	var toneStr string
	err = huh.NewSelect[string]().
		Title("Choose commit message tone").
		Options(
			huh.NewOption("Professional - Formal and precise", "professional"),
			huh.NewOption("Casual - Relaxed and conversational", "casual"),
			huh.NewOption("Friendly - Warm and approachable", "friendly"),
		).
		Value(&toneStr).
		Run()
	if err != nil {
		return nil, err
	}
	config.Tone = Tone(toneStr)

	var lengthStr string
	err = huh.NewSelect[string]().
		Title("Choose commit message length").
		Options(
			huh.NewOption("Short - Concise and brief", "short"),
			huh.NewOption("Normal - Balanced detail", "normal"),
			huh.NewOption("Long - Detailed and descriptive", "long"),
		).
		Value(&lengthStr).
		Run()
	if err != nil {
		return nil, err
	}
	config.Length = Length(lengthStr)

	if err := Save(*config); err != nil {
		return nil, fmt.Errorf("failed to save new config: %w", err)
	}

	return config, nil
}

func Save(config UserConfig) error {
	gitRoot, err := git.FindGitRoot()
	if err != nil {
		return fmt.Errorf("not in a git repository: %w", err)
	}

	configPath := filepath.Join(gitRoot, ".git", "diny-config.json")

	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	err = os.WriteFile(configPath, data, 0644)
	if err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}

func PrintConfiguration(userConfig UserConfig) {
	content := fmt.Sprintf("• Emoji: %t\n• Conventional: %t\n• Tone: %s\n• Length: %s",
		userConfig.UseEmoji,
		userConfig.UseConventional,
		userConfig.Tone,
		userConfig.Length)
	ui.Box(ui.BoxOptions{Title: "Configuration", Message: content})
}
