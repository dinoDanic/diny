/* Copyright © 2025 dinoDanic dino.danic@gmail.com */
package config

import (
	_ "embed"
	"fmt"
	"os"
	"path/filepath"
	"slices"

	"gopkg.in/yaml.v3"
)

//go:embed defaults.yaml
var defaultConfigTemplate string

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

type Config struct {
	Theme  string       `yaml:"theme"`
	Commit CommitConfig `yaml:"commit"`
}

type CommitConfig struct {
	Conventional       bool     `yaml:"conventional"`
	ConventionalFormat []string `yaml:"conventional_format"`
	Emoji              bool     `yaml:"emoji"`
	Tone               Tone     `yaml:"tone"`
	Length             Length   `yaml:"length"`
}

var cfg *Config

func loadDefaultConfig() (*Config, error) {
	var defaultCfg Config
	if err := yaml.Unmarshal([]byte(defaultConfigTemplate), &defaultCfg); err != nil {
		return nil, fmt.Errorf("failed to parse embedded defaults: %w", err)
	}
	return &defaultCfg, nil
}

func getConfigPath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	return filepath.Join(home, ".config", "diny", "config.yaml")
}

func createDefaultConfig(configPath string) error {
	dir := filepath.Dir(configPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	if err := os.WriteFile(configPath, []byte(defaultConfigTemplate), 0644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}

func Load(cfgFile string) (*Config, error) {
	configPath := cfgFile
	if configPath == "" {
		configPath = getConfigPath()
		if configPath == "" {
			return loadDefaultConfig()
		}
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		if os.IsNotExist(err) {
			if err := createDefaultConfig(configPath); err != nil {
				fmt.Printf("Using default configuration (couldn't create config file: %v)\n", err)
				return loadDefaultConfig()
			}
			fmt.Printf("Created default config at: %s\n", configPath)
			data, err = os.ReadFile(configPath)
			if err != nil {
				return loadDefaultConfig()
			}
		} else {
			return nil, fmt.Errorf("failed to read config file: %w", err)
		}
	}

	// Parse YAML
	var loadedCfg Config
	if err := yaml.Unmarshal(data, &loadedCfg); err != nil {
		return nil, fmt.Errorf("invalid config file: %w", err)
	}

	// Validate configuration
	if err := loadedCfg.Validate(); err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}

	return &loadedCfg, nil
}

func (c *Config) Validate() error {
	// Validate tone
	validTones := []Tone{Professional, Casual, Friendly}
	if !slices.Contains(validTones, c.Commit.Tone) {
		return fmt.Errorf("invalid tone '%s', must be one of: professional, casual, friendly", c.Commit.Tone)
	}

	// Validate length
	validLengths := []Length{Short, Normal, Long}
	if !slices.Contains(validLengths, c.Commit.Length) {
		return fmt.Errorf("invalid length '%s', must be one of: short, normal, long", c.Commit.Length)
	}

	// Validate theme (basic check - just ensure it's not empty)
	if c.Theme == "" {
		c.Theme = "catppuccin"
	}

	return nil
}

// Get returns the loaded configuration or defaults if not loaded
func Get() *Config {
	if cfg == nil {
		defaultCfg, err := loadDefaultConfig()
		if err != nil {
			return &Config{
				Theme: "catppuccin",
				Commit: CommitConfig{
					Conventional:       false,
					ConventionalFormat: []string{"feat", "fix", "docs", "chore", "style", "refactor", "test"},
					Emoji:              false,
					Tone:               Casual,
					Length:             Short,
				},
			}
		}
		return defaultCfg
	}
	return cfg
}

func Set(c *Config) {
	cfg = c
}
