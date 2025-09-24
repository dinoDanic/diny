package ollama

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/dinoDanic/diny/config"
)

// const server = "http://127.0.0.1:11434"
const server = "http://167.235.150.40"

const model = "qwen2.5:7b-instruct"

type GenerateRequest struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
	Stream bool   `json:"stream"`
}

type GenerateResponse struct {
	Response string `json:"response"`
	Done     bool   `json:"done"`
}

func buildSystemPrompt(userConfig config.UserConfig) string {
	// Base header
	prompt := `You are an expert that writes high-quality Git commit messages.

You will be given a unified git diff. Produce ONLY a commit message. Do not include code, explanations, or any text outside the specified format. Be specific to changed files and behaviors.

General rules (always apply):
- Use imperative mood (e.g., "add", "fix", "refactor", not "added" or "adds").
- The first line is a concise subject ≤ 72 chars (hard limit).
- The body (if present) explains WHY and key WHAT in bullets, wrapped ≈72 chars/line.
- Prefer concrete nouns and scopes over vague terms (avoid "update", "improve" alone).
- Infer a logical scope from paths (e.g., auth, ui, css, assets, build, api).
- Never include trailing punctuation in the subject.
- If multiple areas changed, pick the primary type/scope for the subject; list others in the body bullets.
- If nothing meaningful changed, say "chore: reformat" and stop.

Formatting modes:`

	// Length
	switch userConfig.Length {
	case config.Short:
		prompt += `
- Length: "short" → subject only (no body) unless a breaking change.`
	case config.Normal:
		prompt += `
- Length: "normal" → subject + optional body (1–4 bullets).`
	case config.Long:
		prompt += `
- Length: "long" → subject + body (2–6 bullets) explaining rationale and impact.`
	default:
		panic(fmt.Sprintf("unhandled Length value: %v", userConfig.Length))
	}

	// Conventional toggle
	if userConfig.Conventional {
		prompt += `
- Use Conventional Commits in the subject: <type>(<scope>): <subject>
  - type ∈ {feat, fix, refactor, perf, docs, style, test, build, chore}
  - scope: infer from paths (e.g., auth, ui, css, assets)
  - examples:
    - feat(auth): add OAuth login flow
    - style(ui): tweak button padding and hover states`
	} else {
		prompt += `
- No prefix required in subject (but still use imperative mood and ≤72 chars).`
	}

	// Style
	switch userConfig.Style {
	case config.Gitmoji:
		prompt += `
- Start the subject with a single appropriate emoji (gitmoji style), then a space.
  - examples: ✨ feat(...): ..., 🐛 fix(...): ..., ♻️ refactor(...): ...
  - Only use ONE emoji at the very start.`
	case config.Simple:
		prompt += `
- Use simple, clear language. NO emojis, NO decorative text.`
	default:
		panic(fmt.Sprintf("unhandled Style value: %v", userConfig.Style))
	}

	// Tone
	switch userConfig.Tone {
	case config.Professional:
		prompt += `
- Tone: professional and matter-of-fact.`
	case config.Casual:
		prompt += `
- Tone: light but clear (still professional).`
	case config.Friendly:
		prompt += `
- Tone: warm and approachable (still concise).`
	default:
		panic(fmt.Sprintf("unhandled Tone value: %v", userConfig.Tone))
	}

	// Extra guardrails + examples
	prompt += `
Hard constraints:
- Do NOT include code fences, backticks, or extra commentary.
- Do NOT mention "this diff" / "AI" / "model".
- Do NOT output more than one commit message.
- Wrap lines ~72 chars; no lines > 80 chars.

Body guidance (when a body is required by length or complexity):
- Start with a blank line after the subject.
- Use bullet points ("- ") focusing on WHY and key WHAT.
- Mention notable file/area changes (e.g., "split layout into Header/Content", "remove focus-visible CSS").
- If assets were added (images, icons), include one bullet (“add C2ButtonIcon SVG and login visual”).
- If only formatting removed (e.g., CSS focus/hover), prefer "style:" or "chore:" appropriately.

Examples:
- Good:
  feat(auth): redesign login layout with image variant
  - introduce variant-based API (image|icon)
  - split layout into Header and Content components
  - update login page to use AuthLayout image variant
  - add C2ButtonIcon and login visuals
  - remove unused focus-visible and hover styles

- Bad (do NOT write):
  "update stuff", "improve UI", "misc changes", "wip", emojis mid-body, explanations about the diff.

Output delimiters:
Begin your answer with <<<COMMIT and end with COMMIT>>>. Print nothing else.`

	prompt += "\n\nHere is the git diff:\n\n"
	return prompt
}

func Main(gitdiff string) (string, error) {
	userConfig := config.Load()
	systemPrompt := buildSystemPrompt(userConfig)

	req := GenerateRequest{
		Model:  model,
		Prompt: systemPrompt + gitdiff,
		Stream: false,
	}

	jsonData, err := json.Marshal(req)
	if err != nil {
		return "", fmt.Errorf("error marshaling JSON: %v", err)
	}

	fmt.Println("🐢 My tiny server is thinking hard, thanks for your patience!")

	resp, err := http.Post(server+"/api/generate", "application/json", bytes.NewBuffer(jsonData))

	if err != nil {
		return "", fmt.Errorf("error calling Ollama: %v", err)
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading response: %v", err)
	}

	var generateResp GenerateResponse
	err = json.Unmarshal(body, &generateResp)
	if err != nil {
		return "", fmt.Errorf("error parsing response: %v", err)
	}

	return generateResp.Response, nil
}
