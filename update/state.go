package update

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"
)

const (
	stateFileName  = "update-state.json"
	snoozeDuration = 48 * time.Hour
	fetchThrottle  = 15 * time.Minute
)

type state struct {
	DismissedAt         time.Time `json:"dismissed_at,omitempty"`
	LastCheckedAt       time.Time `json:"last_checked_at,omitempty"`
	CachedLatestVersion string    `json:"cached_latest_version,omitempty"`
}

func stateFilePath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".config", "diny", stateFileName), nil
}

func loadState() state {
	var s state
	path, err := stateFilePath()
	if err != nil {
		return s
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return s
	}
	_ = json.Unmarshal(data, &s)
	return s
}

func saveState(s state) error {
	path, err := stateFilePath()
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}
	data, err := json.Marshal(s)
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}

func isSnoozed() bool {
	s := loadState()
	if s.DismissedAt.IsZero() {
		return false
	}
	return time.Since(s.DismissedAt) < snoozeDuration
}

func dismiss() error {
	s := loadState()
	s.DismissedAt = time.Now().UTC()
	return saveState(s)
}

func cachedVersion() (version string, fresh bool) {
	s := loadState()
	if s.CachedLatestVersion == "" || s.LastCheckedAt.IsZero() {
		return "", false
	}
	return s.CachedLatestVersion, time.Since(s.LastCheckedAt) < fetchThrottle
}

func saveFetchedVersion(v string) error {
	s := loadState()
	s.CachedLatestVersion = v
	s.LastCheckedAt = time.Now().UTC()
	return saveState(s)
}
