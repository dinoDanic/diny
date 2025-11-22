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

type TimelineRequest struct {
	Prompt   string         `json:"prompt"`
	Version  string         `json:"version"`
	Name     string         `json:"name"`
	RepoName string         `json:"repoName"`
	Config   *config.Config `json:"config"`
	System   string         `json:"system,omitempty"`
}

type timelineData struct {
	Message string `json:"message"`
}

type timelineResp struct {
	Error *string       `json:"error,omitempty"`
	Data  *timelineData `json:"data,omitempty"`
}

func CreateTimelineWithGroq(prompt string, cfg *config.Config) (string, error) {
	gitName := git.GetGitName()
	repoName := git.GetRepoName()

	payload := TimelineRequest{
		Config:   cfg,
		Version:  version.Get(),
		Prompt:   prompt,
		Name:     gitName,
		RepoName: repoName,
		System:   runtime.GOOS,
	}

	buf, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("marshal payload: %w", err)
	}

	req, err := http.NewRequestWithContext(context.Background(),
		http.MethodPost,
		server.ServerConfig.BaseURL+"/api/v2/timeline",
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

	var out timelineResp
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
