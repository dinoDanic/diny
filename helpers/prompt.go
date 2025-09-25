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

func BuildTimelinePrompt(userConfig config.UserConfig) string {
	var b strings.Builder

	// Role & Goal
	b.WriteString("You analyze a sequence of git commits and produce a clear timeline summary.\n")
	b.WriteString("Your job is to summarize WHAT changed and WHY it mattered, not HOW it was implemented.\n\n")

	// Hard rules
	b.WriteString("Rules:\n")
	b.WriteString("- Output plain text only (no code fences, no markdown headings).\n")
	b.WriteString("- Do NOT echo diffs or full commit messages; summarize themes and impact.\n")
	b.WriteString("- Be concrete and specific (quantify when possible: files touched, modules, scope).\n")
	b.WriteString("- Group changes by category (feat, fix, docs, refactor, test, chore, perf) and/or by area/scope if available.\n")
	b.WriteString("- Prefer ISO dates for ranges (YYYY-MM-DD) when referencing time.\n")
	b.WriteString("- Keep it skimmable and consistent with the requested tone and length.\n")

	// Optional emoji legend
	if userConfig.UseEmoji {
		b.WriteString("\nIf helpful, prefix bullets with these emojis:\n")
		b.WriteString("🚀 feat  🐛 fix  📚 docs  🎨 style  ♻️ refactor  ✅ test  🔧 chore  ⚡ perf\n")
		b.WriteString("Use them only at the start of bullets; avoid overuse.\n")
	}

	// Tone
	switch userConfig.Tone {
	case config.Professional:
		b.WriteString("\nTone: professional, neutral, analytical.\n")
	case config.Casual:
		b.WriteString("\nTone: casual and concise, still precise.\n")
	case config.Friendly:
		b.WriteString("\nTone: friendly, encouraging, and clear.\n")
	}

	// Length / Structure
	b.WriteString("\nStructure:\n")
	switch userConfig.Length {
	case config.Short:
		b.WriteString("- 1–2 line high-level summary.\n")
		b.WriteString("- 3–5 short bullets of key themes (group by category or area).\n")
	case config.Normal:
		b.WriteString("- 2–4 line overview summarizing the period, focus, and impact.\n")
		b.WriteString("- 5–9 bullets grouped by category/area (each bullet <= 20 words).\n")
		b.WriteString("- Finish with a 1-line outlook/next steps if discernible.\n")
	case config.Long:
		b.WriteString("- 1 short overview paragraph (3–6 lines) with time span and main objectives.\n")
		b.WriteString("- Thematic sections as bullet blocks (feat/fix/refactor/test/etc.), 2–5 bullets per block.\n")
		b.WriteString("- Add a compact metrics block if derivable (e.g., commits: N, dominant category, notable scopes).\n")
		b.WriteString("- Close with risks/debt and next focus (if inferable from commits).\n")
	default:
		b.WriteString("- 2–4 line overview, then 5–9 bullets by category.\n")
	}

	// Content guidance: what to extract
	b.WriteString("\nFocus on extracting:\n")
	b.WriteString("- Objectives/themes (e.g., 'auth revamp', 'performance in X', 'bugfixes in payments').\n")
	b.WriteString("- Impact/benefit (e.g., 'reduced cold start', 'fewer 500s', 'unlock feature Y').\n")
	b.WriteString("- Notable scopes/modules/packages.\n")
	b.WriteString("- Any visible progression (early groundwork → refactor → feature; or bug surge → stabilization).\n")
	b.WriteString("- Risks/tech debt signals (temporary workarounds, TODOs, follow-ups) if apparent.\n")

	// Formatting constraints
	b.WriteString("\nFormatting constraints:\n")
	b.WriteString("- No tables. No numbered lists unless explicitly helpful.\n")
	b.WriteString("- Start bullets with a category keyword or emoji (if enabled), then a crisp summary.\n")
	b.WriteString("- Keep phrasing action-oriented; avoid vague words like 'various' or 'some'.\n")

	// Safety rails / edge cases
	b.WriteString("\nEdge cases:\n")
	b.WriteString("- If commits are noisy or repetitive, state that and collapse into a single summarized bullet per theme.\n")
	b.WriteString("- If timeframe is mixed, call out shifts in focus (e.g., 'Early week: tests; Late week: perf').\n")
	b.WriteString("- If insufficient data, say so briefly and provide the best-effort high-level view.\n")

	return b.String()
}
