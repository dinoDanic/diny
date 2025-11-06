package config

import (
	"os"
	"sync"
)

type APIProvider string

const (
	CloudBackend APIProvider = "cloud"
	LocalOllama  APIProvider = "ollama"
)

type APIConfig struct {
	Provider    APIProvider
	BaseURL     string
	Model       string
}

type Service struct {
	userConfig *UserConfig
	apiConfig  *APIConfig
	mu         sync.RWMutex
}

var (
	globalService *Service
	once          sync.Once
)

func GetService() *Service {
	once.Do(func() {
		globalService = &Service{}
	})
	return globalService
}

func (s *Service) LoadUserConfig() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	userConfig, err := Load()
	if err != nil {
		return err
	}

	if userConfig == nil {
		userConfig = getDefaultUserConfig()
	}

	s.userConfig = userConfig
	s.apiConfig = s.buildAPIConfig(userConfig)

	return nil
}

func (s *Service) GetUserConfig() *UserConfig {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if s.userConfig == nil {
		return getDefaultUserConfig()
	}

	return s.userConfig
}

func (s *Service) GetAPIConfig() *APIConfig {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if s.apiConfig == nil {
		return s.buildAPIConfig(getDefaultUserConfig())
	}

	return s.apiConfig
}

func (s *Service) buildAPIConfig(userConfig *UserConfig) *APIConfig {
	if userConfig.UseLocalAPI {
		baseURL := getEnvOrDefault("DINY_OLLAMA_URL", userConfig.OllamaURL)
		if baseURL == "" {
			baseURL = "http://127.0.0.1:11434"
		}

		model := getEnvOrDefault("DINY_OLLAMA_MODEL", userConfig.OllamaModel)
		if model == "" {
			model = "llama3.2"
		}

		return &APIConfig{
			Provider: LocalOllama,
			BaseURL:  baseURL,
			Model:    model,
		}
	}

	baseURL := getEnvOrDefault("DINY_BACKEND_URL", userConfig.BackendURL)
	if baseURL == "" {
		baseURL = "https://diny-cli.vercel.app"
	}

	return &APIConfig{
		Provider: CloudBackend,
		BaseURL:  baseURL,
		Model:    "",
	}
}

func getDefaultUserConfig() *UserConfig {
	return &UserConfig{
		UseConventional: false,
		UseEmoji:        false,
		Tone:            Casual,
		Length:          Normal,
		UseLocalAPI:     false,
	}
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func (s *Service) IsUsingLocalAPI() bool {
	config := s.GetAPIConfig()
	return config.Provider == LocalOllama
}

func (s *Service) GetAPIBaseURL() string {
	config := s.GetAPIConfig()
	return config.BaseURL
}

func (s *Service) GetModel() string {
	config := s.GetAPIConfig()
	return config.Model
}
