package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestAIMode_Constants(t *testing.T) {
	t.Run("valid AI modes exist", func(t *testing.T) {
		if AIRemote != "remote" {
			t.Errorf("expected AIRemote to be 'remote', got '%s'", AIRemote)
		}
		if AILocal != "local" {
			t.Errorf("expected AILocal to be 'local', got '%s'", AILocal)
		}
		if AICustom != "custom" {
			t.Errorf("expected AICustom to be 'custom', got '%s'", AICustom)
		}
	})
}

func TestConfig_AIFields(t *testing.T) {
	t.Run("config has AI section", func(t *testing.T) {
		cfg := Config{
			Theme: "catppuccin",
			AI: AIConfig{
				Mode: AIRemote,
			},
			Commit: CommitConfig{
				Tone:               Casual,
				Length:             Short,
				ConventionalFormat: []string{"feat"},
				EmojiMap:           map[string]string{"feat": "x"},
			},
		}
		if cfg.AI.Mode != AIRemote {
			t.Errorf("expected AI.Mode to be 'remote', got '%s'", cfg.AI.Mode)
		}
	})

	t.Run("local config fields", func(t *testing.T) {
		cfg := AIConfig{
			Mode:     AILocal,
			LocalURL: "http://localhost:11434",
		}
		if cfg.LocalURL != "http://localhost:11434" {
			t.Errorf("expected LocalURL, got '%s'", cfg.LocalURL)
		}
	})

	t.Run("custom config fields", func(t *testing.T) {
		cfg := AIConfig{
			Mode:   AICustom,
			APIURL: "https://api.example.com/v1/chat",
			APIKey: "sk-test-key",
			Model:  "gpt-4",
		}
		if cfg.APIURL != "https://api.example.com/v1/chat" {
			t.Errorf("expected APIURL, got '%s'", cfg.APIURL)
		}
		if cfg.APIKey != "sk-test-key" {
			t.Errorf("expected APIKey, got '%s'", cfg.APIKey)
		}
		if cfg.Model != "gpt-4" {
			t.Errorf("expected Model, got '%s'", cfg.Model)
		}
	})
}

func TestConfig_AIValidation(t *testing.T) {
	baseCommit := CommitConfig{
		Tone:               Casual,
		Length:             Short,
		ConventionalFormat: []string{"feat"},
		EmojiMap:           map[string]string{"feat": "x"},
	}

	t.Run("remote mode requires no extra fields", func(t *testing.T) {
		cfg := Config{
			Theme:  "catppuccin",
			AI:     AIConfig{Mode: AIRemote},
			Commit: baseCommit,
		}
		if err := cfg.Validate(); err != nil {
			t.Errorf("remote mode should validate without extra fields: %v", err)
		}
	})

	t.Run("local mode requires local_url", func(t *testing.T) {
		cfg := Config{
			Theme:  "catppuccin",
			AI:     AIConfig{Mode: AILocal},
			Commit: baseCommit,
		}
		err := cfg.Validate()
		if err == nil {
			t.Error("local mode without local_url should fail validation")
		}
	})

	t.Run("local mode with local_url passes", func(t *testing.T) {
		cfg := Config{
			Theme: "catppuccin",
			AI: AIConfig{
				Mode:     AILocal,
				LocalURL: "http://localhost:11434",
			},
			Commit: baseCommit,
		}
		if err := cfg.Validate(); err != nil {
			t.Errorf("local mode with local_url should pass: %v", err)
		}
	})

	t.Run("custom mode requires api_url and api_key", func(t *testing.T) {
		cfg := Config{
			Theme:  "catppuccin",
			AI:     AIConfig{Mode: AICustom},
			Commit: baseCommit,
		}
		err := cfg.Validate()
		if err == nil {
			t.Error("custom mode without api_url/api_key should fail")
		}
	})

	t.Run("custom mode with api_url and api_key passes", func(t *testing.T) {
		cfg := Config{
			Theme: "catppuccin",
			AI: AIConfig{
				Mode:   AICustom,
				APIURL: "https://api.example.com/v1/chat",
				APIKey: "sk-test-key",
			},
			Commit: baseCommit,
		}
		if err := cfg.Validate(); err != nil {
			t.Errorf("custom mode with api_url and api_key should pass: %v", err)
		}
	})

	t.Run("invalid AI mode fails", func(t *testing.T) {
		cfg := Config{
			Theme:  "catppuccin",
			AI:     AIConfig{Mode: "invalid"},
			Commit: baseCommit,
		}
		err := cfg.Validate()
		if err == nil {
			t.Error("invalid AI mode should fail validation")
		}
	})

	t.Run("empty AI mode defaults to remote in validation", func(t *testing.T) {
		cfg := Config{
			Theme:  "catppuccin",
			AI:     AIConfig{},
			Commit: baseCommit,
		}
		// Empty mode should be treated as remote (backward compatible)
		if err := cfg.Validate(); err != nil {
			t.Errorf("empty AI mode should pass (defaults to remote): %v", err)
		}
	})
}

func TestConfig_YAML_AISection(t *testing.T) {
	t.Run("AI config loads from YAML", func(t *testing.T) {
		yamlContent := `
theme: catppuccin
ai:
  mode: local
  local_url: "http://localhost:11434"
commit:
  conventional: false
  conventional_format: ['feat', 'fix']
  emoji: false
  emoji_map:
    feat: x
  tone: casual
  length: short
  custom_instructions: ""
  hash_after_commit: false
`
		tmpDir := t.TempDir()
		cfgPath := filepath.Join(tmpDir, "config.yaml")
		if err := os.WriteFile(cfgPath, []byte(yamlContent), 0644); err != nil {
			t.Fatalf("failed to write temp config: %v", err)
		}

		cfg, err := Load(cfgPath)
		if err != nil {
			t.Fatalf("failed to load config: %v", err)
		}

		if cfg.AI.Mode != AILocal {
			t.Errorf("expected AI.Mode='local', got '%s'", cfg.AI.Mode)
		}
		if cfg.AI.LocalURL != "http://localhost:11434" {
			t.Errorf("expected LocalURL, got '%s'", cfg.AI.LocalURL)
		}
	})

	t.Run("config without AI section defaults to remote", func(t *testing.T) {
		yamlContent := `
theme: catppuccin
commit:
  conventional: false
  conventional_format: ['feat', 'fix']
  emoji: false
  emoji_map:
    feat: x
  tone: casual
  length: short
  custom_instructions: ""
  hash_after_commit: false
`
		tmpDir := t.TempDir()
		cfgPath := filepath.Join(tmpDir, "config.yaml")
		if err := os.WriteFile(cfgPath, []byte(yamlContent), 0644); err != nil {
			t.Fatalf("failed to write temp config: %v", err)
		}

		cfg, err := Load(cfgPath)
		if err != nil {
			t.Fatalf("failed to load config: %v", err)
		}

		// Empty mode is treated as remote
		if cfg.AI.Mode != "" && cfg.AI.Mode != AIRemote {
			t.Errorf("expected AI.Mode to be empty or 'remote', got '%s'", cfg.AI.Mode)
		}
	})
}

func TestLocalConfig_AIFields(t *testing.T) {
	t.Run("local config overlay has AI fields", func(t *testing.T) {
		mode := AILocal
		localCfg := LocalConfig{
			AI: LocalAIConfig{
				Mode:     &mode,
				LocalURL: "http://localhost:11434",
			},
		}

		if *localCfg.AI.Mode != AILocal {
			t.Errorf("expected AI.Mode pointer to AILocal")
		}
		if localCfg.AI.LocalURL != "http://localhost:11434" {
			t.Errorf("expected LocalURL")
		}
	})
}

func TestMergeConfig_AIFields(t *testing.T) {
	t.Run("overlay AI mode overrides base", func(t *testing.T) {
		base := &Config{
			Theme: "catppuccin",
			AI:    AIConfig{Mode: AIRemote},
			Commit: CommitConfig{
				Tone:               Casual,
				Length:             Short,
				ConventionalFormat: []string{"feat"},
				EmojiMap:           map[string]string{"feat": "x"},
			},
		}

		mode := AILocal
		overlay := &LocalConfig{
			AI: LocalAIConfig{
				Mode:     &mode,
				LocalURL: "http://localhost:11434",
			},
		}

		merged := mergeConfigWithLocal(base, overlay)
		if merged.AI.Mode != AILocal {
			t.Errorf("expected merged AI.Mode='local', got '%s'", merged.AI.Mode)
		}
		if merged.AI.LocalURL != "http://localhost:11434" {
			t.Errorf("expected merged AI.LocalURL='http://localhost:11434', got '%s'", merged.AI.LocalURL)
		}
	})

	t.Run("overlay without AI preserves base", func(t *testing.T) {
		base := &Config{
			Theme: "catppuccin",
			AI: AIConfig{
				Mode:   AICustom,
				APIURL: "https://api.example.com",
				APIKey: "sk-key",
			},
			Commit: CommitConfig{
				Tone:               Casual,
				Length:             Short,
				ConventionalFormat: []string{"feat"},
				EmojiMap:           map[string]string{"feat": "x"},
			},
		}

		overlay := &LocalConfig{}

		merged := mergeConfigWithLocal(base, overlay)
		if merged.AI.Mode != AICustom {
			t.Errorf("expected AI.Mode preserved as 'custom', got '%s'", merged.AI.Mode)
		}
		if merged.AI.APIURL != "https://api.example.com" {
			t.Errorf("expected APIURL preserved")
		}
		if merged.AI.APIKey != "sk-key" {
			t.Errorf("expected APIKey preserved")
		}
	})
}
