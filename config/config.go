/* Copyright © 2025 dinoDanic dino.danic@gmail.com */
package config

import (
	_ "embed"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/dinoDanic/diny/git"
	"gopkg.in/yaml.v3"
)

//go:embed defaults.yaml
var defaultConfigTemplate string

type Tone string
type Length string
type AIMode string

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

const (
	AIRemote    AIMode = "remote"
	AILocal     AIMode = "local"
	AICustom    AIMode = "custom"
	AIAnthropic AIMode = "anthropic"
)

type AIConfig struct {
	Mode     AIMode `yaml:"mode" json:"Mode"`
	LocalURL string `yaml:"local_url" json:"LocalURL"`
	APIURL   string `yaml:"api_url" json:"APIURL"`
	APIKey   string `yaml:"api_key" json:"APIKey"`
	Model    string `yaml:"model" json:"Model"`
}

type Config struct {
	Theme  string       `yaml:"theme" json:"Theme"`
	AI     AIConfig     `yaml:"ai" json:"AI"`
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

type LocalAIConfig struct {
	Mode     *AIMode `yaml:"mode,omitempty"`
	LocalURL string  `yaml:"local_url,omitempty"`
	APIURL   string  `yaml:"api_url,omitempty"`
	APIKey   string  `yaml:"api_key,omitempty"`
	Model    string  `yaml:"model,omitempty"`
}

type LocalConfig struct {
	Theme  string            `yaml:"theme,omitempty"`
	AI     LocalAIConfig     `yaml:"ai,omitempty"`
	Commit LocalCommitConfig `yaml:"commit,omitempty"`
}

type LocalCommitConfig struct {
	Conventional       *bool             `yaml:"conventional,omitempty"`
	ConventionalFormat []string          `yaml:"conventional_format,omitempty"`
	Emoji              *bool             `yaml:"emoji,omitempty"`
	EmojiMap           map[string]string `yaml:"emoji_map,omitempty"`
	Tone               Tone              `yaml:"tone,omitempty"`
	Length             Length            `yaml:"length,omitempty"`
	CustomInstructions string            `yaml:"custom_instructions,omitempty"`
	HashAfterCommit    *bool             `yaml:"hash_after_commit,omitempty"`
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

func GetVersionedProjectConfigPath() string {
	repoRoot, err := git.FindGitRoot()
	if err != nil {
		return ""
	}
	return filepath.Join(repoRoot, ".diny.yaml")
}

func GetLocalProjectConfigPath() string {
	gitDir, err := git.FindGitDir()
	if err != nil {
		return ""
	}
	return filepath.Join(gitDir, "diny", "config.yaml")
}

func loadLocalConfig(path string) (*LocalConfig, error) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil, nil 
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("unreadable: %w", err)
	}

	var cfg LocalConfig
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("invalid YAML: %w", err)
	}

	return &cfg, nil
}

func mergeConfigWithLocal(base *Config, overlay *LocalConfig) *Config {
	merged := &Config{
		Theme: base.Theme,
		AI: AIConfig{
			Mode:     base.AI.Mode,
			LocalURL: base.AI.LocalURL,
			APIURL:   base.AI.APIURL,
			APIKey:   base.AI.APIKey,
			Model:    base.AI.Model,
		},
		Commit: CommitConfig{
			Conventional:       base.Commit.Conventional,
			ConventionalFormat: make([]string, len(base.Commit.ConventionalFormat)),
			Emoji:              base.Commit.Emoji,
			EmojiMap:           make(map[string]string),
			Tone:               base.Commit.Tone,
			Length:             base.Commit.Length,
			CustomInstructions: base.Commit.CustomInstructions,
			HashAfterCommit:    base.Commit.HashAfterCommit,
		},
	}

	copy(merged.Commit.ConventionalFormat, base.Commit.ConventionalFormat)
	for k, v := range base.Commit.EmojiMap {
		merged.Commit.EmojiMap[k] = v
	}

	if overlay.Theme != "" {
		merged.Theme = overlay.Theme
	}

	if overlay.AI.Mode != nil {
		merged.AI.Mode = *overlay.AI.Mode
	}
	if overlay.AI.LocalURL != "" {
		merged.AI.LocalURL = overlay.AI.LocalURL
	}
	if overlay.AI.APIURL != "" {
		merged.AI.APIURL = overlay.AI.APIURL
	}
	if overlay.AI.APIKey != "" {
		merged.AI.APIKey = overlay.AI.APIKey
	}
	if overlay.AI.Model != "" {
		merged.AI.Model = overlay.AI.Model
	}

	if overlay.Commit.Conventional != nil {
		merged.Commit.Conventional = *overlay.Commit.Conventional
	}
	if overlay.Commit.Emoji != nil {
		merged.Commit.Emoji = *overlay.Commit.Emoji
	}
	if overlay.Commit.HashAfterCommit != nil {
		merged.Commit.HashAfterCommit = *overlay.Commit.HashAfterCommit
	}

	if overlay.Commit.Tone != "" {
		merged.Commit.Tone = overlay.Commit.Tone
	}
	if overlay.Commit.Length != "" {
		merged.Commit.Length = overlay.Commit.Length
	}
	if overlay.Commit.CustomInstructions != "" {
		merged.Commit.CustomInstructions = overlay.Commit.CustomInstructions
	}

	if len(overlay.Commit.ConventionalFormat) > 0 {
		merged.Commit.ConventionalFormat = make([]string, len(overlay.Commit.ConventionalFormat))
		copy(merged.Commit.ConventionalFormat, overlay.Commit.ConventionalFormat)
	}

	if len(overlay.Commit.EmojiMap) > 0 {
		for k, v := range overlay.Commit.EmojiMap {
			merged.Commit.EmojiMap[k] = v
		}
	}

	return merged
}

func LoadWithProjectOverride(globalCfgFile string) (*Config, string, error) {
	globalConfig, err := Load(globalCfgFile)
	if err != nil {
		return nil, "", fmt.Errorf("failed to load global config: %w", err)
	}

	sources := []string{"global"}
	result := globalConfig

	versionedPath := GetVersionedProjectConfigPath()
	if versionedPath != "" {
		versionedCfg, err := loadLocalConfig(versionedPath)
		if err != nil {
			sources = append(sources, "versioned-error")
		} else if versionedCfg != nil {
			result = mergeConfigWithLocal(result, versionedCfg)
			sources = append(sources, "versioned")
		}
	}

	localPath := GetLocalProjectConfigPath()
	if localPath != "" {
		localCfg, err := loadLocalConfig(localPath)
		if err != nil {
			sources = append(sources, "local-error")
		} else if localCfg != nil {
			result = mergeConfigWithLocal(result, localCfg)
			sources = append(sources, "local")
		}
	}

	if err := result.Validate(); err != nil {
		return globalConfig, "global (merged config invalid)", nil
	}

	sourceStr := strings.Join(sources, " + ")
	return result, sourceStr, nil
}

func LoadOrRecoverWithProject(cfgFile string) (*LoadResult, error) {
	configPath := cfgFile
	if configPath == "" {
		configPath = GetConfigPath()
	}

	cfg, source, err := LoadWithProjectOverride(cfgFile)
	if err == nil {
		result := &LoadResult{Config: cfg}

		if strings.Contains(source, "versioned-error") {
			result.RecoveryMsg = "Versioned project config (.diny.yaml) has errors, skipping"
		} else if strings.Contains(source, "local-error") {
			result.RecoveryMsg = "Local project config (<gitdir>/diny/config.yaml) has errors, skipping"
		} else if source == "global (merged config invalid)" {
			result.RecoveryMsg = "Merged config is invalid, using global config only"
		}

		return result, nil
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

			newCfg, source, loadErr := LoadWithProjectOverride(cfgFile)
			if loadErr != nil {
				return nil, fmt.Errorf("failed to load new config: %w", loadErr)
			}

			recoveryMsg := "Invalid global config backed up, new default created"
			if strings.Contains(source, "versioned") || strings.Contains(source, "local") {
				recoveryMsg += " (project config applied)"
			}

			return &LoadResult{
				Config:        newCfg,
				RecoveryMsg:   recoveryMsg,
				ValidationErr: validationErr,
			}, nil
		}
	}

	return nil, err
}

func createVersionedProjectConfigIfNeeded() error {
	path := GetVersionedProjectConfigPath()
	if path == "" {
		return fmt.Errorf("not in a git repository")
	}

	if _, err := os.Stat(path); err == nil {
		return nil 
	}

	template := `# Diny Project Configuration (Versioned)
# This file can be committed to version control and shared with your team
# It overlays on top of global config (~/.config/diny/config.yaml)
# Only specify the settings you want to override from the global config
# Learn more: https://github.com/dinoDanic/diny

# UI theme (catppuccin, tokyonight, nord, dracula, gruvbox, etc.)
# theme: catppuccin

# AI generation settings
# ai:
#   # Generation mode: remote (default), local (e.g. Ollama), custom (OpenAI-compatible), or anthropic (Claude)
#   mode: remote
#
#   # Local AI server URL (required when mode: local)
#   # local_url: "http://localhost:11434"
#
#   # Custom API URL (required when mode: custom, must be OpenAI-compatible)
#   # Note: Anthropic/Claude API is NOT compatible — use an OpenAI-compatible endpoint
#   # api_url: "https://api.openai.com/v1/chat/completions"
#
#   # API key (required when mode: custom or anthropic)
#   # api_key: ""
#
#   # Model name for local, custom, or anthropic mode
#   # model: "llama3"

# Commit configuration
# commit:
#   # Use conventional commit format (feat, fix, docs, etc.)
#   conventional: false
#
#   # Conventional commit types (only used if conventional: true)
#   conventional_format: ['feat', 'fix', 'docs', 'chore', 'style', 'refactor', 'test', 'perf']
#
#   # Add emoji prefix to commits
#   emoji: false
#
#   # Emoji mappings for each commit type (only used if emoji: true)
#   emoji_map:
#     feat: 🚀
#     fix: 🐛
#     docs: 📚
#     chore: 🔧
#     style: 💄
#     refactor: ♻️
#     test: ✅
#     perf: ⚡
#
#   # Commit message tone: professional, casual, or friendly
#   tone: casual
#
#   # Commit message length: short, normal, or long
#   length: short
#
#   # Custom instructions for AI (e.g., "Include JIRA ticket from branch name")
#   custom_instructions: ""
#
#   # Show/copy commit hash after committing
#   hash_after_commit: false
`

	if err := os.WriteFile(path, []byte(template), 0644); err != nil {
		return fmt.Errorf("failed to create versioned project config: %w", err)
	}

	return nil
}

func createLocalProjectConfigIfNeeded() error {
	path := GetLocalProjectConfigPath()
	if path == "" {
		return fmt.Errorf("not in a git repository")
	}

	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	if _, err := os.Stat(path); err == nil {
		return nil // Already exists
	}

	template := `# Diny Local Project Configuration (Not Versioned)
# This file is stored under <gitdir>/diny/ and will never be committed
# Use this for personal overrides on top of team config (.diny.yaml) and global config
# It has highest priority: local > versioned (.diny.yaml) > global
# Only specify the settings you want to override
# Learn more: https://github.com/dinoDanic/diny

# UI theme (catppuccin, tokyonight, nord, dracula, gruvbox, etc.)
# theme: dracula

# AI generation settings
# ai:
#   # Generation mode: remote (default), local (e.g. Ollama), custom (OpenAI-compatible), or anthropic (Claude)
#   mode: local
#
#   # Local AI server URL (required when mode: local)
#   # local_url: "http://localhost:11434"
#
#   # Custom API URL (required when mode: custom, must be OpenAI-compatible)
#   # Note: Anthropic/Claude API is NOT compatible — use an OpenAI-compatible endpoint
#   # api_url: "https://api.openai.com/v1/chat/completions"
#
#   # API key (required when mode: custom or anthropic)
#   # api_key: "sk-your-key-here"
#
#   # Model name for local, custom, or anthropic mode
#   # model: "llama3"

# Commit configuration
# commit:
#   # Use conventional commit format (feat, fix, docs, etc.)
#   conventional: false
#
#   # Conventional commit types (only used if conventional: true)
#   conventional_format: ['feat', 'fix', 'docs', 'chore', 'style', 'refactor', 'test', 'perf']
#
#   # Add emoji prefix to commits
#   emoji: false
#
#   # Emoji mappings for each commit type (only used if emoji: true)
#   emoji_map:
#     feat: 🚀
#     fix: 🐛
#     docs: 📚
#     chore: 🔧
#     style: 💄
#     refactor: ♻️
#     test: ✅
#     perf: ⚡
#
#   # Commit message tone: professional, casual, or friendly
#   tone: casual
#
#   # Commit message length: short, normal, or long
#   length: short
#
#   # Custom instructions for AI (e.g., "Include JIRA ticket from branch name")
#   custom_instructions: ""
#
#   # Show/copy commit hash after committing
#   hash_after_commit: false
`

	if err := os.WriteFile(path, []byte(template), 0644); err != nil {
		return fmt.Errorf("failed to create local project config: %w", err)
	}

	return nil
}

func CreateVersionedProjectConfigIfNeeded() error {
	return createVersionedProjectConfigIfNeeded()
}

func CreateLocalProjectConfigIfNeeded() error {
	return createLocalProjectConfigIfNeeded()
}
