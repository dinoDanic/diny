package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// Custom types for type safety
type Style string
type Tone string
type Length string

// Style constants (for non-conventional styles)
const (
	Gitmoji Style = "gitmoji"
	Simple  Style = "simple"
)

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
	Style        Style  `json:"style"`
	Conventional bool   `json:"conventional"`
	Tone         Tone   `json:"tone"`
	Length       Length `json:"length"`
}

var defaultUserConfig = UserConfig{
	Style:        Simple,
	Conventional: true,
	Tone:         Casual,
	Length:       Short,
}

func findGitRoot() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	for {
		gitPath := filepath.Join(dir, ".git")
		if _, err := os.Stat(gitPath); err == nil {
			return dir, nil
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			return "", fmt.Errorf("not in a git repository")
		}
		dir = parent
	}
}

func Load() UserConfig {
	gitRoot, err := findGitRoot()
	if err != nil {
		return defaultUserConfig
	}

	configPath := filepath.Join(gitRoot, ".diny.json")
	data, err := os.ReadFile(configPath)
	if err != nil {
		return defaultUserConfig
	}

	var config UserConfig
	if err := json.Unmarshal(data, &config); err != nil {
		return defaultUserConfig
	}

	return config
}
