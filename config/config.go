/* Copyright © 2025 dinoDanic dino.danic@gmail.com */
package config

import (
	_ "embed"
	"fmt"
	"os"
	"path/filepath"
	"strings"

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
	Theme  string       `yaml:"theme" json:"Theme"`
	Commit CommitConfig `yaml:"commit" json:"Commit"`
}

type CommitConfig struct {
	Conventional       bool              `yaml:"conventional" json:"Conventional"`
	ConventionalFormat []string          `yaml:"conventional_format" json:"ConventionalFormat"`
	Emoji              bool              `yaml:"emoji" json:"Emoji"`
	EmojiMap           map[string]string `yaml:"emoji_map" json:"EmojiMap"`
	Tone               Tone              `yaml:"tone" json:"Tone"`
	Length             Length            `yaml:"length" json:"Length"`
	CustomInstructions string            `yaml:"custom_instructions" json:"CustomInstructions"`
	HashAfterCommit    bool              `yaml:"hash_after_commit" json:"HashAfterCommit"`
}

func loadDefaultConfig() (*Config, error) {
	var defaultCfg Config
	if err := yaml.Unmarshal([]byte(defaultConfigTemplate), &defaultCfg); err != nil {
		return nil, fmt.Errorf("failed to parse embedded defaults: %w", err)
	}
	return &defaultCfg, nil
}

func GetConfigPath() string {
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
		configPath = GetConfigPath()
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
			data, err = os.ReadFile(configPath)
			if err != nil {
				return loadDefaultConfig()
			}
		} else {
			return nil, fmt.Errorf("failed to read config file: %w", err)
		}
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("invalid config file: %w", err)
	}

	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("config validation failed: %w", err)
	}

	return &cfg, nil
}

type LoadResult struct {
	Config        *Config
	RecoveryMsg   string
	ValidationErr string
}

func LoadOrRecover(cfgFile string) (*LoadResult, error) {
	configPath := cfgFile
	if configPath == "" {
		configPath = GetConfigPath()
	}

	cfg, err := Load(cfgFile)
	if err == nil {
		return &LoadResult{Config: cfg}, nil
	}

	if configPath != "" {
		if _, statErr := os.Stat(configPath); statErr == nil {
			validationErr := err.Error()

			backupPath := getBackupPath(configPath)
			if renameErr := os.Rename(configPath, backupPath); renameErr != nil {
				return nil, fmt.Errorf("failed to backup config: %w", renameErr)
			}

			if createErr := createDefaultConfig(configPath); createErr != nil {
				return nil, fmt.Errorf("failed to create new config: %w", createErr)
			}

			newCfg, loadErr := Load(cfgFile)
			if loadErr != nil {
				return nil, fmt.Errorf("failed to load new config: %w", loadErr)
			}

			return &LoadResult{
				Config:        newCfg,
				RecoveryMsg:   "Invalid config backed up, new default created",
				ValidationErr: validationErr,
			}, nil
		}
	}

	return nil, err
}

func getBackupPath(configPath string) string {
	dir := filepath.Dir(configPath)
	ext := filepath.Ext(configPath)
	base := strings.TrimSuffix(filepath.Base(configPath), ext)

	for i := 1; ; i++ {
		backupName := fmt.Sprintf("%s.backup%d%s", base, i, ext)
		backupPath := filepath.Join(dir, backupName)
		if _, err := os.Stat(backupPath); os.IsNotExist(err) {
			return backupPath
		}
	}
}
