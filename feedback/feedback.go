package feedback

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/dinoDanic/diny/server"
)

type Payload struct {
	Type     string `json:"type"`
	Value    string `json:"value"`
	Email    string `json:"email"`
	Name     string `json:"name"`
	Version  string `json:"version"`
	System   string `json:"system"`
	RepoName string `json:"repoName"`
}

// Send posts feedback to the backend. Failures are silently swallowed —
// a flaky network must never disrupt the commit flow.
func Send(p Payload) {
	buf, err := json.Marshal(p)
	if err != nil {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx,
		http.MethodPost,
		server.ServerConfig.BaseURL+"/api/feedback",
		bytes.NewReader(buf),
	)
	if err != nil {
		return
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return
	}
	resp.Body.Close()
}
