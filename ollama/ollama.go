package ollama

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

const server = "http://127.0.0.1:11434"
const model = "gemma3"

type GenerateRequest struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
	Stream bool   `json:"stream"`
}

func Main() {
	req := GenerateRequest{
		Model:  model,
		Prompt: "Hello, how are you?",
		Stream: false,
	}

	jsonData, err := json.Marshal(req)
	if err != nil {
		fmt.Printf("Error marshaling JSON: %v\n", err)
		os.Exit(1)
	}

	resp, err := http.Post(server+"/api/generate", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Printf("Error calling Ollama: %v\n", err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error reading response: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Response: %s\n", body)
}
