/*
Copyright © 2025 NAME HERE dino.danic@gmail.com
*/
package helpers

import (
	"github.com/dinoDanic/diny/config"
)

func BuildSystemPrompt(userConfig config.UserConfig) string {
	prompt := "Write a git commit message. "

	if userConfig.UseConventional {
		prompt += "Use conventional format: type(scope): description. "
	}

	if userConfig.UseEmoji {
		prompt += "Start with appropriate emoji. "
	}

	// Add length instruction
	switch userConfig.Length {
	case config.Short:
		prompt += "Keep it short, subject only. "
	case config.Normal:
		prompt += "Subject + optional 1-4 body bullets. "
	case config.Long:
		prompt += "Subject + detailed 2-6 body bullets. "
	}

	prompt += "Git diff:\n\n"
	return prompt
}
