package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/dinoDanic/diny/git"
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
	UseEmoji        bool   `json:"use_emoji"`
	UseConventional bool   `json:"use_conventional"`
	Tone            Tone   `json:"tone"`
	Length          Length `json:"length"`
}

func Load() (*UserConfig, error) {
	gitRoot, err := git.FindGitRoot()
	if err != nil {
		return nil, err
	}

	configPath := filepath.Join(gitRoot, ".git", "diny-config.json")
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
	fmt.Println("⚙️  Configuration:")
	fmt.Printf("   • Emoji: %t\n", userConfig.UseEmoji)
	fmt.Printf("   • Conventional: %t\n", userConfig.UseConventional)
	fmt.Printf("   • Tone: %s\n", userConfig.Tone)
	fmt.Printf("   • Length: %s\n", userConfig.Length)
}
