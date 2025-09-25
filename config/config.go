package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/dinoDanic/diny/git"
)

// Custom types for type safety
type Tone string
type Length string

// Tone constants
const (
	Professional Tone = "professional"
	Casual       Tone = "casual"
	Friendly     Tone = "friendly"
)

// Length constants
const (
	Short  Length = "short"
	Normal Length = "normal"
	Long   Length = "long"
)

type UserConfig struct {
	UseEmoji        bool   `json:"use_emoji"`
	UseConventional bool   `json:"use_conventional"`
	Tone            Tone   `json:"tone"`
	Length          Length `json:"length"`
}

var DefaultUserConfig = UserConfig{
	UseEmoji:        false,
	UseConventional: false,
	Tone:            Casual,
	Length:          Short,
}

func Load() UserConfig {
	gitRoot, err := git.FindGitRoot()
	if err != nil {
		return DefaultUserConfig
	}

	configPath := filepath.Join(gitRoot, ".git", "diny-config.json")
	data, err := os.ReadFile(configPath)
	if err != nil {
		return DefaultUserConfig
	}

	var config UserConfig
	if err := json.Unmarshal(data, &config); err != nil {
		return DefaultUserConfig
	}

	return config
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
