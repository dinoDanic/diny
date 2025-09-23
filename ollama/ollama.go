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
	prompt := `You are an expert at writing Git commit messages. Generate a clear commit message based on the provided git diff in English language. 

Rules:
- Use imperative mood
- Write only the commit message, no code snippets or explanations
- Focus on what changed and why`

	switch userConfig.Length {
	case config.Short:
		prompt += "\n- Keep it under 50 characters"
	case config.Normal:
		prompt += "\n- Keep it concise but descriptive, around 50-72 characters"
	case config.Long:
		prompt += "\n- Write a detailed message explaining the changes"
	default:
		panic(fmt.Sprintf("unhandled Length value: %v", userConfig.Length))
	}

	// Handle conventional commits as boolean (can be mixed with other styles)
	if userConfig.Conventional {
		prompt += "\n- Use conventional commit format with type prefix like feat:, fix:, docs:"
	}

	// Handle other style options
	switch userConfig.Style {
	case config.Gitmoji:
		prompt += "\n- Start with an appropriate emoji"
	case config.Simple:
		prompt += "\n- Use simple, clear language without prefixes"
	default:
		panic(fmt.Sprintf("unhandled Style value: %v", userConfig.Style))
	}

	switch userConfig.Tone {
	case config.Professional:
		prompt += "\n- Use professional language"
	case config.Casual:
		prompt += "\n- Use casual language"
	case config.Friendly:
		prompt += "\n- Use warm, approachable language"
	default:
		panic(fmt.Sprintf("unhandled Tone value: %v", userConfig.Tone))
	}

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

	fmt.Printf("loading..")

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
