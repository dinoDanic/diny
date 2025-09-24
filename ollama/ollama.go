package ollama

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
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

// MainStream generates a commit message with live streaming output
func MainStream(prompt string) (string, error) {
	req := GenerateRequest{
		Model:  model,
		Prompt: prompt,
		Stream: true,
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

	fmt.Println() // New line after streaming is done

	if err := scanner.Err(); err != nil {
		return "", fmt.Errorf("error reading stream: %v", err)
	}

	return fullResponse.String(), nil
}

func Main(prompt string) (string, error) {
	req := GenerateRequest{
		Model:  model,
		Prompt: prompt,
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
