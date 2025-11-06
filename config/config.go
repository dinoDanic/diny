package config

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/charmbracelet/huh"
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
	UseLocalAPI     bool   `json:"useLocalAPI"`
	OllamaURL       string `json:"ollamaURL,omitempty"`
	OllamaModel     string `json:"ollamaModel,omitempty"`
	BackendURL      string `json:"backendURL,omitempty"`
}

func Load() (*UserConfig, error) {
	return LoadMerged()
}

// LoadMerged loads config with precedence: local > global > defaults
func LoadMerged() (*UserConfig, error) {
	config := getDefaultUserConfig()

	globalPath, err := GetGlobalConfigPath()
	if err == nil {
		globalConfig, err := tryLoadConfig(globalPath)
		if err == nil && globalConfig != nil {
			config = mergeConfig(config, globalConfig)
		}
	}

	localPath, err := GetLocalConfigPath()
	if err == nil {
		localConfig, err := tryLoadConfig(localPath)
		if err == nil && localConfig != nil {
			config = mergeConfig(config, localConfig)
		}
	}

	if !isValidConfig(config) {
		return config, nil // Return with defaults if invalid
	}

	return config, nil
}

// LoadGlobal loads only the global config
func LoadGlobal() (*UserConfig, error) {
	globalPath, err := GetGlobalConfigPath()
	if err != nil {
		return nil, err
	}

	config, err := tryLoadConfig(globalPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return handleConfigError(globalPath, err)
	}

	return config, nil
}

// LoadLocal loads only the local config
func LoadLocal() (*UserConfig, error) {
	localPath, err := GetLocalConfigPath()
	if err != nil {
		return nil, err
	}

	config, err := tryLoadConfig(localPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return handleConfigError(localPath, err)
	}

	return config, nil
}

// mergeConfig merges src into dst (src values override dst)
func mergeConfig(dst, src *UserConfig) *UserConfig {
	result := *dst

	result.UseConventional = src.UseConventional
	result.UseEmoji = src.UseEmoji
	result.UseLocalAPI = src.UseLocalAPI

	if src.Tone != "" {
		result.Tone = src.Tone
	}

	if src.Length != "" {
		result.Length = src.Length
	}

	if src.OllamaURL != "" {
		result.OllamaURL = src.OllamaURL
	}

	if src.OllamaModel != "" {
		result.OllamaModel = src.OllamaModel
	}

	if src.BackendURL != "" {
		result.BackendURL = src.BackendURL
	}

	return &result
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

func handleInvalidConfig(configPath string) (*UserConfig, error) {
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

	err = huh.NewConfirm().
		Title("Use local Ollama API?").
		Description("Connect to local Ollama instance at http://localhost:11434 (requires Ollama installed)").
		Value(&config.UseLocalAPI).
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
	return SaveLocal(config)
}

// SaveGlobal saves config to global location
func SaveGlobal(config UserConfig) error {
	configPath, err := GetGlobalConfigPath()
	if err != nil {
		return fmt.Errorf("failed to get global config path: %w", err)
	}

	return saveConfigToPath(config, configPath)
}

// SaveLocal saves config to local (repository) location
func SaveLocal(config UserConfig) error {
	configPath, err := GetLocalConfigPath()
	if err != nil {
		return fmt.Errorf("not in a git repository: %w", err)
	}

	return saveConfigToPath(config, configPath)
}

// saveConfigToPath saves config to a specific path
func saveConfigToPath(config UserConfig, configPath string) error {
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
	content := fmt.Sprintf("• Emoji: %t\n• Conventional: %t\n• Tone: %s\n• Length: %s\n• Local API: %t",
		userConfig.UseEmoji,
		userConfig.UseConventional,
		userConfig.Tone,
		userConfig.Length,
		userConfig.UseLocalAPI)

	if userConfig.UseLocalAPI {
		if userConfig.OllamaURL != "" {
			content += fmt.Sprintf("\n• Ollama URL: %s", userConfig.OllamaURL)
		}
		if userConfig.OllamaModel != "" {
			content += fmt.Sprintf("\n• Ollama Model: %s", userConfig.OllamaModel)
		}
	} else {
		if userConfig.BackendURL != "" {
			content += fmt.Sprintf("\n• Backend URL: %s", userConfig.BackendURL)
		}
	}

	ui.Box(ui.BoxOptions{Title: "Configuration", Message: content})
}

// GetConfigSummary returns a single-line summary of the configuration
func GetConfigSummary(userConfig UserConfig) string {
	summary := fmt.Sprintf("emoji:%t conv:%t tone:%s len:%s",
		userConfig.UseEmoji,
		userConfig.UseConventional,
		userConfig.Tone,
		userConfig.Length)

	configService := GetService()
	apiConfig := configService.GetAPIConfig()

	if apiConfig.Provider == LocalOllama {
		summary += fmt.Sprintf(" api:ollama(%s)", apiConfig.Model)
	} else {
		summary += " api:cloud"
	}

	return summary
}

func PrintEffectiveConfiguration(userConfig UserConfig) {
	configService := GetService()
	apiConfig := configService.GetAPIConfig()

	content := "Active Settings\n\n"
	content += fmt.Sprintf("• Emoji: %t\n", userConfig.UseEmoji)
	content += fmt.Sprintf("• Conventional: %t\n", userConfig.UseConventional)
	content += fmt.Sprintf("• Tone: %s\n", userConfig.Tone)
	content += fmt.Sprintf("• Length: %s\n", userConfig.Length)
	content += fmt.Sprintf("• Local API: %t\n\n", userConfig.UseLocalAPI)

	content += "Effective API Configuration\n\n"
	content += fmt.Sprintf("• Provider: %s\n", apiConfig.Provider)
	content += fmt.Sprintf("• URL: %s\n", apiConfig.BaseURL)

	ollamaURLSource := getConfigSource("DINY_OLLAMA_URL", userConfig.OllamaURL, "http://127.0.0.1:11434")
	content += fmt.Sprintf("  └─ Source: %s\n", ollamaURLSource)

	if apiConfig.Provider == LocalOllama && apiConfig.Model != "" {
		content += fmt.Sprintf("• Model: %s\n", apiConfig.Model)
		modelSource := getConfigSource("DINY_OLLAMA_MODEL", userConfig.OllamaModel, "llama3.2")
		content += fmt.Sprintf("  └─ Source: %s\n", modelSource)
	}

	if apiConfig.Provider == CloudBackend {
		backendSource := getConfigSource("DINY_BACKEND_URL", userConfig.BackendURL, "https://diny-cli.vercel.app")
		content += fmt.Sprintf("  └─ Source: %s\n", backendSource)
	}

	content += "\nConfiguration Keys\n\n"
	content += "JSON:       ollamaURL, ollamaModel, backendURL\n"
	content += "Env Vars:   DINY_OLLAMA_URL, DINY_OLLAMA_MODEL, DINY_BACKEND_URL\n"
	content += "Precedence: env vars > local JSON > global JSON > defaults"

	ui.Box(ui.BoxOptions{Title: "Effective Configuration", Message: content})
}

// GetConfigSource returns a string describing where a config value comes from
func GetConfigSource(envVar string, jsonValue string, defaultValue string) string {
	if envValue := os.Getenv(envVar); envValue != "" {
		return fmt.Sprintf("env var %s", envVar)
	}
	if jsonValue != "" {
		return "config file"
	}
	return "default"
}

func getConfigSource(envVar string, jsonValue string, defaultValue string) string {
	return GetConfigSource(envVar, jsonValue, defaultValue)
}
