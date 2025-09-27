package groq

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/dinoDanic/diny/config"
	"github.com/dinoDanic/diny/server"
)

type commitData struct {
	CommitMessage string `json:"commitMessage"`
}

type commitResp struct {
	Error *string     `json:"error,omitempty"`
	Data  *commitData `json:"data,omitempty"`
}

func CreateCommitMessageWithGroq(gitDiff string, userConfig *config.UserConfig) (string, error) {
	payload := map[string]interface{}{
		"gitDiff": gitDiff,
		"version": "1.0.0", // TODO: Get actual version
		"name":    "diny",  // TODO: Get actual name
	}

	if userConfig != nil {
		payload["userConfig"] = *userConfig
	} else {
		payload["userConfig"] = nil
	}

	buf, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("marshal payload: %w", err)
	}

	req, err := http.NewRequestWithContext(context.Background(),
		http.MethodPost,
		server.ServerConfig.BaseURL+"/api/commit",
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

	var out commitResp

	if err := json.Unmarshal(body, &out); err != nil {
		return "", fmt.Errorf("decode response: %w", err)
	}

	if out.Error != nil {
		return "", fmt.Errorf("commit generation failed: %s", *out.Error)
	}

	if out.Data == nil {
		return "", fmt.Errorf("no data in response")
	}

	if out.Data.CommitMessage == "" {
		return "", fmt.Errorf("empty commit message from server")
	}

	return out.Data.CommitMessage, nil
}
