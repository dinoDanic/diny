/*
Copyright © 2025 NAME HERE dino.danic@gmail.com
*/
package helpers

import (
	"strings"

	"github.com/dinoDanic/diny/config"
)

func BuildSystemPrompt(userConfig config.UserConfig) string {
	var b strings.Builder

	b.WriteString("Write ONLY a git commit message for the provided diff.\n\n")
	b.WriteString("Rules:\n")
	b.WriteString("- Output only the message (no pre/post text)\n")
	b.WriteString("- Don’t echo the diff\n")
	b.WriteString("- No explanations, comments, or markdown\n")
	b.WriteString("- Emphasize WHY and WHAT, not HOW\n")

	if userConfig.UseConventional {
		b.WriteString("\nFormat: type(scope): subject. Types: feat, fix, docs, style, refactor, test, chore, perf\n")
	}
	if userConfig.UseEmoji {
		b.WriteString("\nPrefix emoji: 🚀 feat, 🐛 fix, 📚 docs, 🎨 style, ♻️ refactor, ✅ test, 🔧 chore, ⚡ perf\n")
	}

	switch userConfig.Tone {
	case config.Professional:
		b.WriteString("\nTone: professional\n")
	case config.Casual:
		b.WriteString("\nTone: casual\n")
	case config.Friendly:
		b.WriteString("\nTone: friendly\n")
	}

	switch userConfig.Length {
	case config.Short:
		b.WriteString("\nStructure: subject only (<=50 chars)\n")
	case config.Normal:
		b.WriteString("\nStructure: subject (<=50, imperative) + 1–4 bullets starting with '-'\n")
	case config.Long:
		b.WriteString("\nStructure: subject (<=50, imperative) + 2–6 bullets w/ context & impact\n")
	}

	return b.String()
}
