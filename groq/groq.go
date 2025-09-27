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

type commitResp struct {
	Message string      `json:"message"`
	Note    string      `json:"note"`
	Error   string      `json:"error"`
	Details interface{} `json:"details"`
}

func CreateCommitMessageWithGroq(gitDiff string, userConfig *config.UserConfig) (string, string, error) {
	payload := map[string]interface{}{
		"git_diff": gitDiff,
	}

	if userConfig != nil {
		payload["user_config"] = *userConfig
	}

	buf, err := json.Marshal(payload)
	if err != nil {
		return "", "", fmt.Errorf("marshal payload: %w", err)
	}

	req, err := http.NewRequestWithContext(context.Background(),
		http.MethodPost,
		server.ServerConfig.BaseURL+"/api/commit",
		bytes.NewReader(buf),
	)
	if err != nil {
		return "", "", fmt.Errorf("new request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 30 * time.Second}
	res, err := client.Do(req)

	if err != nil {
		return "", "", fmt.Errorf("do request: %w", err)
	}

	defer res.Body.Close()

	body, _ := io.ReadAll(res.Body)

	if res.StatusCode != http.StatusOK {
		var e commitResp
		_ = json.Unmarshal(body, &e)
		if e.Error != "" {
			return "", "", fmt.Errorf("proxy %d: %s", res.StatusCode, e.Error)
		}
		return "", "", fmt.Errorf("proxy %d: %s", res.StatusCode, string(body))
	}

	var out commitResp
	if err := json.Unmarshal(body, &out); err != nil {
		return "", "", fmt.Errorf("decode response: %w", err)
	}
	if out.Message == "" {
		return "", "", fmt.Errorf("empty message from proxy")
	}

	return out.Message, out.Note, nil
}
