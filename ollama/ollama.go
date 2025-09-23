package ollama

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"diny/config"
)

const server = "http://127.0.0.1:11434"
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
	prompt := `You are an expert at writing Git commit messages. Generate a clear commit message based on the provided git diff. 

Rules:
- Use imperative mood
- Write only the commit message, no code snippets or explanations
- Focus on what changed and why`

	if userConfig.Length == config.Short {
		prompt += "\n- Keep it under 50 characters"
	} else if userConfig.Length == config.Normal {
		prompt += "\n- Keep it concise but descriptive, around 50-72 characters"
	} else if userConfig.Length == config.Long {
		prompt += "\n- Write a detailed message explaining the changes"
	}

	if userConfig.Style == config.Conventional {
		prompt += "\n- Use conventional commit format with type prefix like feat:, fix:, docs:"
	} else if userConfig.Style == config.Gitmoji {
		prompt += "\n- Start with an appropriate emoji"
	} else if userConfig.Style == config.Simple {
		prompt += "\n- Use simple, clear language without prefixes"
	}

	if userConfig.Tone == config.Casual {
		prompt += "\n- Use casual language"
	} else if userConfig.Tone == config.Friendly {
		prompt += "\n- Use warm, approachable language"
	} else {
		prompt += "\n- Use professional language"
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
