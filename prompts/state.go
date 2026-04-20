package prompts

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

// Status values for rating prompt.
const (
	StatusPending   = "pending"
	StatusRated     = "rated"
	StatusDismissed = "dismissed"
)

// Status values for star prompt (superset of rating statuses).
const (
	StatusStarred      = "starred"
	StatusAlreadyGiven = "already_given"
	// StatusDismissed is shared with rating.
)

type RatingState struct {
	Status     string     `yaml:"status"`
	Value      int        `yaml:"value,omitempty"`
	AnsweredAt *time.Time `yaml:"answered_at,omitempty"`
}

type StarState struct {
	Status     string     `yaml:"status"`
	AnsweredAt *time.Time `yaml:"answered_at,omitempty"`
}

type PromptsState struct {
	CommitCount int         `yaml:"commit_count"`
	Rating      RatingState `yaml:"rating"`
	Star        StarState   `yaml:"star"`
}

type State struct {
	Prompts PromptsState `yaml:"prompts"`
}

func defaultState() *State {
	return &State{
		Prompts: PromptsState{
			CommitCount: 0,
			Rating:      RatingState{Status: StatusPending},
			Star:        StarState{Status: StatusPending},
		},
	}
}

func GetStatePath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	return filepath.Join(home, ".config", "diny", "state.yaml")
}

func LoadState() *State {
	path := GetStatePath()
	if path == "" {
		return defaultState()
	}

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return defaultState()
		}
		// Corrupted or unreadable — back up and recreate.
		backupAndRecreate(path)
		return defaultState()
	}

	var s State
	if err := yaml.Unmarshal(data, &s); err != nil {
		backupAndRecreate(path)
		return defaultState()
	}

	// Normalize missing status fields to pending.
	if s.Prompts.Rating.Status == "" {
		s.Prompts.Rating.Status = StatusPending
	}
	if s.Prompts.Star.Status == "" {
		s.Prompts.Star.Status = StatusPending
	}

	return &s
}

func SaveState(s *State) error {
	path := GetStatePath()
	if path == "" {
		return fmt.Errorf("cannot determine state file path")
	}

	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create state directory: %w", err)
	}

	data, err := yaml.Marshal(s)
	if err != nil {
		return fmt.Errorf("failed to marshal state: %w", err)
	}

	// Atomic write: temp file + rename.
	tmp := path + ".tmp"
	if err := os.WriteFile(tmp, data, 0644); err != nil {
		return fmt.Errorf("failed to write state temp file: %w", err)
	}
	if err := os.Rename(tmp, path); err != nil {
		os.Remove(tmp)
		return fmt.Errorf("failed to rename state file: %w", err)
	}

	return nil
}

func backupAndRecreate(path string) {
	backupPath := getBackupPath(path)
	_ = os.Rename(path, backupPath)
	s := defaultState()
	data, err := yaml.Marshal(s)
	if err != nil {
		return
	}
	_ = os.WriteFile(path, data, 0644)
}

func getBackupPath(path string) string {
	dir := filepath.Dir(path)
	ext := filepath.Ext(path)
	base := strings.TrimSuffix(filepath.Base(path), ext)

	for i := 1; ; i++ {
		bp := filepath.Join(dir, fmt.Sprintf("%s.backup%d%s", base, i, ext))
		if _, err := os.Stat(bp); os.IsNotExist(err) {
			return bp
		}
	}
}
