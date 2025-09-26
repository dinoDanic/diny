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

	// Core task
	b.WriteString("Write ONLY a git commit message for the provided diff.\n")

	// Global constraints
	b.WriteString("\nRules:\n")
	b.WriteString("- Output only the commit message (no pre/post text).\n")
	b.WriteString("- Do not echo the diff or include code.\n")
	b.WriteString("- No explanations, comments, or markdown.\n")
	b.WriteString("- Use imperative mood and focus on WHAT changed and WHY it matters (not HOW).\n")
	b.WriteString("- Avoid vague terms (e.g., 'various', 'some'); be specific and concise.\n")

	// Optional formats
	if userConfig.UseConventional {
		b.WriteString("\nFormat: type(scope): subject  — types: feat, fix, docs, style, refactor, test, chore, perf.\n")
	}
	if userConfig.UseEmoji {
		b.WriteString("\nOptional emoji prefixes: 🚀 feat  🐛 fix  📚 docs  🎨 style  ♻️ refactor  ✅ test  🔧 chore  ⚡ perf.\n")
	}

	// Tone hint
	switch userConfig.Tone {
	case config.Professional:
		b.WriteString("\nTone: professional.\n")
	case config.Casual:
		b.WriteString("\nTone: casual.\n")
	case config.Friendly:
		b.WriteString("\nTone: friendly.\n")
	}

	// Length/structure
	b.WriteString("\nStructure:\n")
	switch userConfig.Length {
	case config.Short:
		b.WriteString("- Subject only, ≤50 chars, imperative verb first word.\n")
	case config.Normal:
		b.WriteString("- Subject ≤50 chars (imperative). If needed, add 1–4 bullets starting with '- ' for key WHY.\n")
	case config.Long:
		b.WriteString("- Subject ≤50 chars (imperative). Then 2–6 terse bullets for context/impact (no code, no diff).\n")
	default:
		b.WriteString("- Subject ≤50 chars (imperative). Add up to 3 bullets only if essential.\n")
	}

	// Edge handling
	b.WriteString("\nEdge cases:\n")
	b.WriteString("- If diff is large/noisy, summarize main intent and primary impact.\n")
	b.WriteString("- If changes are mechanical (format/rename), state it plainly.\n")
	b.WriteString("- Never include secrets, paths, or identifiers from the diff.\n")

	// Explicit diff label
	b.WriteString("\ngit diff:\n")

	return b.String()
}

func BuildTimelinePrompt(userConfig config.UserConfig) string {
	var b strings.Builder

	// Role & Goal
	b.WriteString("You analyze git commits and produce a client-facing daily work log.\n")
	b.WriteString("Start with ONE summary sentence (past tense), then bullet points of completed actions.\n\n")

	// Hard rules
	b.WriteString("Rules:\n")
	b.WriteString("- Output plain text only (no code fences, no markdown headings).\n")
	b.WriteString("- First line: single-sentence summary (≤20 words), past tense, no prefix labels.\n")
	b.WriteString("- Then a blank line, then bullets of what was done (no explanations or methodology).\n")
	b.WriteString("- Do NOT echo full commit messages or diffs; synthesize clean outcomes.\n")
	b.WriteString("- Use past tense, action verbs (Added, Fixed, Refactored, Created, Updated).\n")
	b.WriteString("- Keep each bullet ≤16 words. One line per bullet.\n")
	b.WriteString("- Optional compact tags: [feat], [fix], [docs], [refactor], [test], [chore], [perf].\n")
	b.WriteString("- Be concrete; quantify when obvious. No guesses.\n\n")

	// Tone
	switch userConfig.Tone {
	case config.Professional:
		b.WriteString("Tone: professional and concise.\n")
	case config.Casual:
		b.WriteString("Tone: casual but precise.\n")
	case config.Friendly:
		b.WriteString("Tone: friendly and clear.\n")
	default:
		b.WriteString("Tone: professional and concise.\n")
	}

	// Length / Structure
	b.WriteString("\nStructure:\n")
	switch userConfig.Length {
	case config.Short:
		b.WriteString("- 1 summary sentence, then 3–5 bullets.\n")
	case config.Normal:
		b.WriteString("- 1 summary sentence, then 5–9 bullets.\n")
	case config.Long:
		b.WriteString("- 1 summary sentence, then 8–12 bullets; group implicitly by tag if helpful.\n")
	default:
		b.WriteString("- 1 summary sentence, then 5–9 bullets.\n")
	}

	// Content guidance
	b.WriteString("\nFocus on extracting:\n")
	b.WriteString("- Concrete outcomes (features, endpoints, refactors, config/schema changes).\n")
	b.WriteString("- Scope hints (APIs, modules) when obvious.\n")
	b.WriteString("- Clear user/business impact only if explicit; otherwise omit.\n")

	// Edge cases
	b.WriteString("\nEdge cases:\n")
	b.WriteString("- If commits are noisy/repetitive, collapse to one bullet per theme.\n")
	b.WriteString("- If data is sparse, produce best-effort summary + 3–4 bullets.\n")

	// Few-shot style anchors (no headings, just examples of format)
	b.WriteString("\nExamples:\n")
	b.WriteString("Delivered timeline logging and simplified commit workflow.\n\n")
	b.WriteString("- [feat] Added timeline API with background DB logging (after()).\n")
	b.WriteString("- [refactor] Simplified commit API; removed crypto; trimmed telemetry fields.\n")
	b.WriteString("- [chore] Updated model config; adjusted service port.\n")
	b.WriteString("- [sql] Created minimal commits table for telemetry.\n")

	return b.String()
}
