package ai

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/dinoDanic/diny/config"
)

func newTestConfig(mode config.AIMode) *config.Config {
	return &config.Config{
		Theme: "catppuccin",
		AI: config.AIConfig{
			Mode: mode,
		},
		Commit: config.CommitConfig{
			Tone:               config.Casual,
			Length:              config.Short,
			ConventionalFormat: []string{"feat"},
			EmojiMap:           map[string]string{"feat": "x"},
		},
	}
}

func TestGenerateCommitMessage_Remote(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/api/v2/commit" {
			t.Errorf("expected /api/v2/commit, got %s", r.URL.Path)
		}

		resp := map[string]any{
			"data": map[string]any{
				"commitMessage": "feat: add new feature",
			},
		}
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	cfg := newTestConfig(config.AIRemote)

	msg, err := GenerateCommitMessage("diff content", cfg, WithBaseURL(server.URL))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if msg != "feat: add new feature" {
		t.Errorf("expected 'feat: add new feature', got '%s'", msg)
	}
}

func TestGenerateCommitMessage_Local(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/chat" {
			t.Errorf("expected /api/chat, got %s", r.URL.Path)
		}

		var req map[string]any
		json.NewDecoder(r.Body).Decode(&req)

		// Verify model is passed if set
		if _, ok := req["model"]; !ok {
			t.Error("expected model in request")
		}

		resp := map[string]any{
			"message": map[string]any{
				"content": "fix: resolve null pointer",
			},
		}
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	cfg := newTestConfig(config.AILocal)
	cfg.AI.LocalURL = server.URL
	cfg.AI.Model = "llama3"

	msg, err := GenerateCommitMessage("diff content", cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if msg != "fix: resolve null pointer" {
		t.Errorf("expected 'fix: resolve null pointer', got '%s'", msg)
	}
}

func TestGenerateCommitMessage_Custom(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify API key is sent
		authHeader := r.Header.Get("Authorization")
		if authHeader != "Bearer sk-test-key" {
			t.Errorf("expected Authorization 'Bearer sk-test-key', got '%s'", authHeader)
		}

		resp := map[string]any{
			"choices": []map[string]any{
				{
					"message": map[string]any{
						"content": "docs: update readme",
					},
				},
			},
		}
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	cfg := newTestConfig(config.AICustom)
	cfg.AI.APIURL = server.URL
	cfg.AI.APIKey = "sk-test-key"
	cfg.AI.Model = "gpt-4"

	msg, err := GenerateCommitMessage("diff content", cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if msg != "docs: update readme" {
		t.Errorf("expected 'docs: update readme', got '%s'", msg)
	}
}

func TestGenerateCommitMessage_EmptyModeDefaultsToRemote(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := map[string]any{
			"data": map[string]any{
				"commitMessage": "chore: cleanup",
			},
		}
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	cfg := newTestConfig("")

	msg, err := GenerateCommitMessage("diff content", cfg, WithBaseURL(server.URL))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if msg != "chore: cleanup" {
		t.Errorf("expected 'chore: cleanup', got '%s'", msg)
	}
}

func TestGenerateTimeline_Remote(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v2/timeline" {
			t.Errorf("expected /api/v2/timeline, got %s", r.URL.Path)
		}
		resp := map[string]any{
			"data": map[string]any{
				"message": "timeline analysis result",
			},
		}
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	cfg := newTestConfig(config.AIRemote)

	msg, err := GenerateTimeline("prompt", cfg, WithBaseURL(server.URL))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if msg != "timeline analysis result" {
		t.Errorf("expected 'timeline analysis result', got '%s'", msg)
	}
}

func TestGenerateTimeline_Local(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := map[string]any{
			"message": map[string]any{
				"content": "local timeline analysis",
			},
		}
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	cfg := newTestConfig(config.AILocal)
	cfg.AI.LocalURL = server.URL
	cfg.AI.Model = "llama3"

	msg, err := GenerateTimeline("prompt", cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if msg != "local timeline analysis" {
		t.Errorf("expected 'local timeline analysis', got '%s'", msg)
	}
}

func TestLocalMode_BlocksExternalRequests(t *testing.T) {
	cfg := newTestConfig(config.AILocal)
	cfg.AI.LocalURL = "http://localhost:11434"

	client := newLocalHTTPClient(cfg.AI.LocalURL)

	// Try to make a request to an external URL - should be blocked
	req, _ := http.NewRequest("GET", "https://example.com", nil)
	_, err := client.Do(req)
	if err == nil {
		t.Error("expected error when making external request in local mode")
	}
	if !strings.Contains(err.Error(), "blocked") {
		t.Errorf("expected 'blocked' in error message, got: %v", err)
	}
}

func TestLocalMode_AllowsLocalRequests(t *testing.T) {
	localServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer localServer.Close()

	client := newLocalHTTPClient(localServer.URL)

	req, _ := http.NewRequest("GET", localServer.URL+"/test", nil)
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("expected local request to succeed: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected 200, got %d", resp.StatusCode)
	}
}

func TestDoPost_HTTPError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(`{"error":"invalid api key"}`))
	}))
	defer server.Close()

	_, err := doPost(server.URL, nil, map[string]string{"hello": "world"}, 5*time.Second)
	if err == nil {
		t.Fatal("expected error for 401 response")
	}
	if !strings.Contains(err.Error(), "HTTP 401") {
		t.Errorf("expected 'HTTP 401' in error, got: %v", err)
	}
	if !strings.Contains(err.Error(), "invalid api key") {
		t.Errorf("expected error body in message, got: %v", err)
	}
}

func TestDoPost_EmptyBody(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		// write nothing
	}))
	defer server.Close()

	_, err := doPost(server.URL, nil, map[string]string{"hello": "world"}, 5*time.Second)
	if err == nil {
		t.Fatal("expected error for empty body")
	}
	if !strings.Contains(err.Error(), "empty response body") {
		t.Errorf("expected 'empty response body' in error, got: %v", err)
	}
}

func TestBuildCommitPrompt(t *testing.T) {
	cfg := newTestConfig(config.AILocal)
	cfg.Commit.Tone = config.Professional
	cfg.Commit.Length = config.Normal
	cfg.Commit.Conventional = true
	cfg.Commit.ConventionalFormat = []string{"feat", "fix", "docs"}
	cfg.Commit.CustomInstructions = "Include JIRA ticket"

	prompt := buildCommitPrompt("diff content", cfg)

	if !strings.Contains(prompt, "diff content") {
		t.Error("prompt should contain the diff")
	}
	if !strings.Contains(prompt, "professional") {
		t.Error("prompt should contain tone")
	}
	if !strings.Contains(prompt, "normal") {
		t.Error("prompt should contain length")
	}
	if !strings.Contains(prompt, "conventional") {
		t.Error("prompt should mention conventional format")
	}
	if !strings.Contains(prompt, "Include JIRA ticket") {
		t.Error("prompt should contain custom instructions")
	}
}
