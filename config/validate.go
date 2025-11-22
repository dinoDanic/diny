package config

import (
	"fmt"
	"slices"
)

// Validate checks if the config has all required fields with valid values
func (c *Config) Validate() error {
	// Theme is required
	if c.Theme == "" {
		return fmt.Errorf("theme is required")
	}

	// Tone is required and must be valid
	validTones := []Tone{Professional, Casual, Friendly}
	if c.Commit.Tone == "" {
		return fmt.Errorf("tone is required")
	}
	if !slices.Contains(validTones, c.Commit.Tone) {
		return fmt.Errorf("invalid tone '%s', must be one of: professional, casual, friendly", c.Commit.Tone)
	}

	// Length is required and must be valid
	validLengths := []Length{Short, Normal, Long}
	if c.Commit.Length == "" {
		return fmt.Errorf("length is required")
	}
	if !slices.Contains(validLengths, c.Commit.Length) {
		return fmt.Errorf("invalid length '%s', must be one of: short, normal, long", c.Commit.Length)
	}

	// ConventionalFormat is required
	if len(c.Commit.ConventionalFormat) == 0 {
		return fmt.Errorf("conventional_format is required")
	}

	return nil
}
