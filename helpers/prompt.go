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
	b.WriteString("You analyze a sequence of git commits and produce a client-facing daily work log.\n")
	b.WriteString("Write only what was accomplished, as concise past-tense bullet points suitable for copy-paste.\n")
	b.WriteString("Do not explain methodology or rationale; just state the completed actions.\n\n")

	// Hard rules
	b.WriteString("Rules:\n")
	b.WriteString("- Output plain text only (no code fences, no markdown headings, no bold).\n")
	b.WriteString("- Do NOT write an intro line like 'Today’s changes...' or any conclusions.\n")
	b.WriteString("- Do NOT echo full commit messages or diffs; synthesize them into clean bullets.\n")
	b.WriteString("- Use past tense, action-oriented phrasing (e.g., 'Added', 'Refactored', 'Fixed').\n")
	b.WriteString("- Keep each bullet ≤ 16 words. Prefer one line per bullet.\n")
	b.WriteString("- If helpful, prefix with a compact tag like [feat], [fix], [refactor], [chore], [docs], [perf], [test].\n")
	b.WriteString("- Quantify when obvious (files/modules/areas touched), but avoid guessing.\n")
	b.WriteString("- No tables. No numbered lists. Bullets only.\n")
	b.WriteString("- If input is too vague, produce 3–5 best-effort bullets and say 'Consolidated minor tweaks' as one bullet.\n\n")

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

	// Length / Structure -> bullets only, no overview
	b.WriteString("\nStructure:\n")
	switch userConfig.Length {
	case config.Short:
		b.WriteString("- 3–5 bullets. No header or footer.\n")
	case config.Normal:
		b.WriteString("- 5–9 bullets. No header or footer.\n")
	case config.Long:
		b.WriteString("- 8–12 bullets, grouped implicitly by tag if applicable. No header or footer.\n")
	default:
		b.WriteString("- 5–9 bullets. No header or footer.\n")
	}

	// Content guidance
	b.WriteString("\nFocus on extracting:\n")
	b.WriteString("- Concrete outcomes per commit/theme (features added, endpoints created, refactors, config changes).\n")
	b.WriteString("- Scope hints (API names, endpoints, modules) when obvious from context.\n")
	b.WriteString("- Brief metrics if explicit (counts of files, tokens, or latency improvements).\n")

	// Edge cases
	b.WriteString("\nEdge cases:\n")
	b.WriteString("- If commits are noisy/repetitive, collapse them into one bullet per theme.\n")
	b.WriteString("- If timeframe mixes tasks, still output a single bullet list with clear actions.\n")

	// Few-shot examples to anchor the style
	b.WriteString("\nExamples:\n")
	b.WriteString("- [feat] Added timeline API with background DB logging (after()).\n")
	b.WriteString("- [refactor] Simplified commit API; removed crypto and trimmed telemetry fields.\n")
	b.WriteString("- [chore] Updated model config; aligned service port.\n")
	b.WriteString("- [sql] Created minimal commits table and timeline schema.\n")
	b.WriteString("- [feat] Replaced prompt truncation with explicit system prompt.\n")

	return b.String()
}
