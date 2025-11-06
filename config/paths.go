package config

import (
	"os"
	"path/filepath"
	"runtime"
)

// GetGlobalConfigPath returns the path to the global config file
// following XDG Base Directory specification on Linux
// Can be overridden with DINY_CONFIG_PATH environment variable
func GetGlobalConfigPath() (string, error) {
	if customPath := os.Getenv("DINY_CONFIG_PATH"); customPath != "" {
		configDir := filepath.Dir(customPath)
		if err := os.MkdirAll(configDir, 0755); err != nil {
			return "", err
		}
		return customPath, nil
	}

	var configDir string

	switch runtime.GOOS {
	case "windows":
		configDir = os.Getenv("APPDATA")
		if configDir == "" {
			configDir = filepath.Join(os.Getenv("USERPROFILE"), "AppData", "Roaming")
		}
		configDir = filepath.Join(configDir, "diny")

	case "darwin":
		home, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		configDir = filepath.Join(home, "Library", "Application Support", "diny")

	default:
		configDir = os.Getenv("XDG_CONFIG_HOME")
		if configDir == "" {
			home, err := os.UserHomeDir()
			if err != nil {
				return "", err
			}
			configDir = filepath.Join(home, ".config", "diny")
		} else {
			configDir = filepath.Join(configDir, "diny")
		}
	}

	if err := os.MkdirAll(configDir, 0755); err != nil {
		return "", err
	}

	return filepath.Join(configDir, "config.json"), nil
}

// GetLocalConfigPath returns the path to the local (repository) config file
func GetLocalConfigPath() (string, error) {
	gitRoot, err := findGitRoot()
	if err != nil {
		return "", err
	}
	return filepath.Join(gitRoot, ".git", "diny-config.json"), nil
}

// findGitRoot finds the root of the git repository
func findGitRoot() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	for {
		gitDir := filepath.Join(dir, ".git")
		if info, err := os.Stat(gitDir); err == nil && info.IsDir() {
			return dir, nil
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			return "", os.ErrNotExist
		}
		dir = parent
	}
}
