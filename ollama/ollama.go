package ollama

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
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

type UserConfig struct {
	Language string `json:"language"` // Options: "English", "Spanish", "French", "German", "Portuguese", "Italian", "Chinese", "Japanese"
	Style    string `json:"style"`    // Options: "conventional", "gitmoji", "simple"
	Tone     string `json:"tone"`     // Options: "professional", "casual", "friendly"
}

var defaultUserConfig = UserConfig{
	Language: "English",
	Style:    "conventional",
	Tone:     "professional",
}

func findGitRoot() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	for {
		gitPath := filepath.Join(dir, ".git")
		if _, err := os.Stat(gitPath); err == nil {
			return dir, nil
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			return "", fmt.Errorf("not in a git repository")
		}
		dir = parent
	}
}

func loadUserConfig() UserConfig {
	gitRoot, err := findGitRoot()
	if err != nil {
		return defaultUserConfig
	}

	configPath := filepath.Join(gitRoot, ".diny.json")
	data, err := os.ReadFile(configPath)
	if err != nil {
		return defaultUserConfig
	}

	var config UserConfig
	if err := json.Unmarshal(data, &config); err != nil {
		return defaultUserConfig
	}

	return config
}

func buildSystemPrompt(userConfig UserConfig) string {
	prompt := `You are an expert at writing Git commit messages. Generate a concise, clear commit message based on the provided git diff. The message should:
- Be in imperative mood (e.g., "Add feature" not "Added feature")
- Be under 50 characters for the subject line
- Focus on the "what" and "why", not the "how"`

	if userConfig.Style == "conventional" {
		prompt += "\n- Follow conventional commit format with type prefix (feat:, fix:, docs:, etc.)"
	} else if userConfig.Style == "gitmoji" {
		prompt += "\n- Include appropriate emoji at the beginning (✨ for new features, 🐛 for bug fixes, etc.)"
	} else if userConfig.Style == "simple" {
		prompt += "\n- Use simple, descriptive format without prefixes"
	}

	if userConfig.Tone == "casual" {
		prompt += "\n- Use a casual, friendly tone"
	} else if userConfig.Tone == "friendly" {
		prompt += "\n- Use a warm, approachable tone"
	} else {
		prompt += "\n- Use a professional tone"
	}

	prompt += fmt.Sprintf("\n\nPlease respond in %s language.\n\nHere is the git diff:\n\n", userConfig.Language)

	return prompt
}

func Main(gitdiff string) (string, error) {
	userConfig := loadUserConfig()
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
