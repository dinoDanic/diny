package ollama

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/dinoDanic/diny/config"
)

type GenerateRequest struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
	Stream bool   `json:"stream"`
}

type GenerateResponse struct {
	Response string `json:"response"`
	Done     bool   `json:"done"`
}

func MainStream(prompt string) (string, error) {
	configService := config.GetService()
	apiConfig := configService.GetAPIConfig()

	req := GenerateRequest{
		Model:  apiConfig.Model,
		Prompt: prompt,
		Stream: true,
	}

	jsonData, err := json.Marshal(req)
	if err != nil {
		return "", fmt.Errorf("error marshaling JSON: %v", err)
	}

	resp, err := http.Post(apiConfig.BaseURL+"/api/generate", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("error calling Ollama: %v", err)
	}
	defer resp.Body.Close()

	var fullResponse strings.Builder
	scanner := bufio.NewScanner(resp.Body)

	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}

		var streamResp GenerateResponse
		if err := json.Unmarshal([]byte(line), &streamResp); err != nil {
			continue // Skip invalid JSON lines
		}

		// Print each chunk as it comes
		fmt.Print(streamResp.Response)
		fullResponse.WriteString(streamResp.Response)

		if streamResp.Done {
			break
		}
	}

	if err := scanner.Err(); err != nil {
		return "", fmt.Errorf("error reading stream: %v", err)
	}

	return fullResponse.String(), nil
}

func Main(prompt string) (string, error) {
	configService := config.GetService()
	apiConfig := configService.GetAPIConfig()

	req := GenerateRequest{
		Model:  apiConfig.Model,
		Prompt: prompt,
		Stream: false,
	}

	jsonData, err := json.Marshal(req)
	if err != nil {
		return "", fmt.Errorf("error marshaling JSON: %v", err)
	}

	fmt.Println("My tiny server is thinking hard, thanks for your patience!")

	resp, err := http.Post(apiConfig.BaseURL+"/api/generate", "application/json", bytes.NewBuffer(jsonData))

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

func CheckHealth(baseURL string) error {
	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	resp, err := client.Get(baseURL + "/api/tags")
	if err != nil {
		errorMsg := fmt.Sprintf("Cannot connect to Ollama at %s\n", baseURL)
		errorMsg += fmt.Sprintf("Error: %v\n\n", err)
		errorMsg += "Is Ollama running?\n"
		errorMsg += "  Local:  ollama serve\n"
		errorMsg += "  Docker: docker run -d -p 11434:11434 ollama/ollama\n"
		errorMsg += "          (or check: docker ps | grep ollama)\n\n"
		errorMsg += "Need to install Ollama? See: ollama/README.md"
		return fmt.Errorf(errorMsg)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Ollama responded with status %d at %s\nExpected 200 OK", resp.StatusCode, baseURL)
	}

	return nil
}

func CheckModelExists(baseURL, modelName string) error {
	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	resp, err := client.Get(baseURL + "/api/tags")
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to fetch models list: status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response: %v", err)
	}

	type ModelInfo struct {
		Name string `json:"name"`
	}

	type TagsResponse struct {
		Models []ModelInfo `json:"models"`
	}

	var tags TagsResponse
	if err := json.Unmarshal(body, &tags); err != nil {
		return fmt.Errorf("failed to parse response: %v", err)
	}

	for _, model := range tags.Models {
		if model.Name == modelName || model.Name == modelName+":latest" {
			return nil
		}
	}

	var modelNames []string
	for _, m := range tags.Models {
		modelNames = append(modelNames, m.Name)
	}

	availableList := "none"
	if len(modelNames) > 0 {
		availableList = strings.Join(modelNames, ", ")
	}

	errorMsg := fmt.Sprintf("Model '%s' not found\nAvailable models: %s\n\nTo pull this model:", modelName, availableList)
	errorMsg += fmt.Sprintf("\n  Local:  ollama pull %s", modelName)
	errorMsg += fmt.Sprintf("\n  Docker: docker exec -it ollama ollama pull %s", modelName)
	errorMsg += "\n\nOr choose a different model from the available list"

	return fmt.Errorf(errorMsg)
}

func CreateCommitMessage(gitDiff string, userConfig *config.UserConfig) (string, error) {
	prompt := buildCommitPrompt(gitDiff, userConfig)
	return Main(prompt)
}

func buildCommitPrompt(gitDiff string, userConfig *config.UserConfig) string {
	prompt := "You are a commit message generator. Generate a clear, concise commit message based on the following git diff.\n\n"

	if userConfig.UseConventional {
		prompt += "Use Conventional Commits format (type(scope): description).\n"
	}

	if userConfig.UseEmoji {
		prompt += "Include appropriate emoji prefixes.\n"
	}

	prompt += fmt.Sprintf("Tone: %s\n", userConfig.Tone)
	prompt += fmt.Sprintf("Length: %s\n", userConfig.Length)

	prompt += "\nGit diff:\n" + gitDiff + "\n\n"
	prompt += "Generate only the commit message, no explanations:"

	return prompt
}
