package config

import (
	"fmt"
	"slices"
)

func (c *Config) Validate() error {
	if c.Theme == "" {
		return fmt.Errorf("theme is required")
	}

	validTones := []Tone{Professional, Casual, Friendly}
	if c.Commit.Tone == "" {
		return fmt.Errorf("tone is required")
	}
	if !slices.Contains(validTones, c.Commit.Tone) {
		return fmt.Errorf("invalid tone '%s', must be one of: professional, casual, friendly", c.Commit.Tone)
	}

	validLengths := []Length{Short, Normal, Long}
	if c.Commit.Length == "" {
		return fmt.Errorf("length is required")
	}
	if !slices.Contains(validLengths, c.Commit.Length) {
		return fmt.Errorf("invalid length '%s', must be one of: short, normal, long", c.Commit.Length)
	}

	if len(c.Commit.ConventionalFormat) == 0 {
		return fmt.Errorf("conventional_format is required")
	}

	if c.Commit.EmojiMap == nil {
		return fmt.Errorf("emoji_map is required")
	}

	return nil
}
