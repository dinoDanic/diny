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
	Type          string         `json:"type"`
	UserPrompt    string         `json:"userPrompt"`
	Version       string         `json:"version"`
	Name          string         `json:"name"`
	Email         string         `json:"email"`
	RepoName      string         `json:"repoName"`
	Config        *config.Config `json:"config"`
	System        string         `json:"system,omitempty"`
	PreviousPlans [][]SplitGroup `json:"previousPlans,omitempty"`
	Feedback      string         `json:"feedback,omitempty"`
}

type SplitGroup struct {
	Order   int      `json:"order"`
	Type    string   `json:"type"`
	Message string   `json:"message"`
	Files   []string `json:"files"`
}

type responseData struct {
	Message string       `json:"message"`
	Groups  []SplitGroup `json:"groups"`
}

type response struct {
	Error *string       `json:"error,omitempty"`
	Data  *responseData `json:"data,omitempty"`
}

type RequestExtras struct {
	PreviousPlans [][]SplitGroup
	Feedback      string
}

func doRequest(reqType string, userPrompt string, cfg *config.Config, extras *RequestExtras) (*responseData, error) {
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
	if extras != nil {
		payload.PreviousPlans = extras.PreviousPlans
		payload.Feedback = extras.Feedback
	}

	buf, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("marshal payload: %w", err)
	}

	req, err := http.NewRequestWithContext(context.Background(),
		http.MethodPost,
		server.ServerConfig.BaseURL+"/api/requests",
		bytes.NewReader(buf),
	)
	if err != nil {
		return nil, fmt.Errorf("new request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 60 * time.Second}
	res, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("do request: %w", err)
	}
	defer res.Body.Close()

	body, _ := io.ReadAll(res.Body)

	var out response
	if err := json.Unmarshal(body, &out); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}

	if out.Error != nil {
		return nil, fmt.Errorf("%s", *out.Error)
	}

	if out.Data == nil {
		return nil, fmt.Errorf("no data in response")
	}

	return out.Data, nil
}

func sendRequest(reqType string, userPrompt string, cfg *config.Config) (string, error) {
	data, err := doRequest(reqType, userPrompt, cfg, nil)
	if err != nil {
		return "", err
	}
	if data.Message == "" {
		return "", fmt.Errorf("empty message from server")
	}
	return data.Message, nil
}

func CreateCommitMessageWithGroq(gitDiff string, cfg *config.Config) (string, error) {
	return sendRequest("commit", gitDiff, cfg)
}

func CreateTimelineWithGroq(prompt string, cfg *config.Config) (string, error) {
	return sendRequest("timeline", prompt, cfg)
}

func CreateSplitPlanWithGroq(gitDiff string, cfg *config.Config, extras *RequestExtras) ([]SplitGroup, error) {
	data, err := doRequest("split", gitDiff, cfg, extras)
	if err != nil {
		return nil, err
	}
	if len(data.Groups) == 0 {
		return nil, fmt.Errorf("empty split plan from server")
	}
	return data.Groups, nil
}
