package config

import (
	"fmt"
	"slices"
)

func (c *Config) Validate() error {
	if c.Theme == "" {
		return fmt.Errorf("theme is required")
	}

	if err := c.validateAI(); err != nil {
		return err
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

func (c *Config) validateAI() error {
	validModes := []AIMode{AIRemote, AILocal, AICustom}

	// Empty mode is treated as remote (backward compatible)
	if c.AI.Mode == "" {
		return nil
	}

	if !slices.Contains(validModes, c.AI.Mode) {
		return fmt.Errorf("invalid ai.mode '%s', must be one of: remote, local, custom", c.AI.Mode)
	}

	switch c.AI.Mode {
	case AILocal:
		if c.AI.LocalURL == "" {
			return fmt.Errorf("ai.local_url is required when ai.mode is 'local'")
		}
	case AICustom:
		if c.AI.APIURL == "" {
			return fmt.Errorf("ai.api_url is required when ai.mode is 'custom'")
		}
		if c.AI.APIKey == "" {
			return fmt.Errorf("ai.api_key is required when ai.mode is 'custom'")
		}
	}

	return nil
}

// EffectiveAIMode returns the AI mode, defaulting to AIRemote if empty.
func (c *Config) EffectiveAIMode() AIMode {
	if c.AI.Mode == "" {
		return AIRemote
	}
	return c.AI.Mode
}
