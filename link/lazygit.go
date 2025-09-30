package link

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/charmbracelet/huh"
	"github.com/dinoDanic/diny/ui"
	"gopkg.in/yaml.v3"
)

type LazyGitConfig struct {
	CustomCommands []CustomCommand `yaml:"customCommands,omitempty"`
}

type CustomCommand struct {
	Key         string `yaml:"key"`
	Description string `yaml:"description"`
	Command     string `yaml:"command"`
	Context     string `yaml:"context"`
	Output      string `yaml:"output,omitempty"`
}

func LinkLazyGit() error {
	var keyBinding string

	fmt.Println()
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Choose a key binding for diny commit").
				Description("This key will trigger diny commit in LazyGit, default C").
				Value(&keyBinding).
				Placeholder("C"),
		),
	)

	err := form.Run()
	if err != nil {
		return fmt.Errorf("failed to get key binding: %w", err)
	}

	if keyBinding == "" {
		keyBinding = "C"
	}

	configPath, err := getLazyGitConfigPath()
	if err != nil {
		return err
	}

	if err := os.MkdirAll(filepath.Dir(configPath), 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	var config LazyGitConfig

	data, err := os.ReadFile(configPath)
	if err != nil {
		if !os.IsNotExist(err) {
			return fmt.Errorf("failed to read config file: %w", err)
		}
	} else {
		if err := yaml.Unmarshal(data, &config); err != nil {
			return fmt.Errorf("failed to parse config file: %w", err)
		}
	}

	dinyCommand := CustomCommand{
		Key:         keyBinding,
		Description: "Commit 🦕",
		Command:     "diny commit",
		Context:     "files",
		Output:      "terminal",
	}

	found := false
	for i, cmd := range config.CustomCommands {
		if cmd.Command == "diny commit" {
			config.CustomCommands[i] = dinyCommand
			found = true
			break
		}
	}

	if !found {
		config.CustomCommands = append(config.CustomCommands, dinyCommand)
	}

	yamlData, err := yaml.Marshal(&config)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	if err := os.WriteFile(configPath, yamlData, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	ui.RenderSuccess("Successfully linked diny with LazyGit!")

	return nil
}

func getLazyGitConfigPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get home directory: %w", err)
	}

	return filepath.Join(home, ".config", "lazygit", "config.yml"), nil

}
