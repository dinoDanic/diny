package groq

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"runtime"
	"time"

	"github.com/dinoDanic/diny/config"
	"github.com/dinoDanic/diny/git"
	"github.com/dinoDanic/diny/server"
	"github.com/dinoDanic/diny/version"
)

type Request struct {
	Type       string         `json:"type"`
	UserPrompt string         `json:"userPrompt"`
	Version    string         `json:"version"`
	Name       string         `json:"name"`
	Email      string         `json:"email"`
	RepoName   string         `json:"repoName"`
	Config     *config.Config `json:"config"`
	System     string         `json:"system,omitempty"`
}

type responseData struct {
	Message string `json:"message"`
}

type response struct {
	Error *string       `json:"error,omitempty"`
	Data  *responseData `json:"data,omitempty"`
}

func sendRequest(reqType string, userPrompt string, cfg *config.Config) (string, error) {
	payload := Request{
		Type:       reqType,
		Config:     cfg,
		Version:    version.Get(),
		UserPrompt: userPrompt,
		Name:       git.GetGitName(),
		Email:      git.GetGitEmail(),
		RepoName:   git.GetRepoName(),
		System:     runtime.GOOS,
	}

	buf, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("marshal payload: %w", err)
	}

	req, err := http.NewRequestWithContext(context.Background(),
		http.MethodPost,
		server.ServerConfig.BaseURL+"/api/requests",
		bytes.NewReader(buf),
	)
	if err != nil {
		return "", fmt.Errorf("new request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 30 * time.Second}
	res, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("do request: %w", err)
	}
	defer res.Body.Close()

	body, _ := io.ReadAll(res.Body)

	var out response
	if err := json.Unmarshal(body, &out); err != nil {
		return "", fmt.Errorf("decode response: %w", err)
	}

	if out.Error != nil {
		return "", fmt.Errorf("%s", *out.Error)
	}

	if out.Data == nil {
		return "", fmt.Errorf("no data in response")
	}

	if out.Data.Message == "" {
		return "", fmt.Errorf("empty message from server")
	}

	return out.Data.Message, nil
}

func CreateCommitMessageWithGroq(gitDiff string, cfg *config.Config) (string, error) {
	return sendRequest("commit", gitDiff, cfg)
}

func CreateTimelineWithGroq(prompt string, cfg *config.Config) (string, error) {
	return sendRequest("timeline", prompt, cfg)
}
