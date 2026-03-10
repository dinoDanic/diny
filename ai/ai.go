package ai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/dinoDanic/diny/config"
	"github.com/dinoDanic/diny/server"
)

// options holds optional overrides for testing and flexibility.
type options struct {
	baseURL string
}

// Option configures the AI generation behavior.
type Option func(*options)

// WithBaseURL overrides the remote server URL (useful for testing).
func WithBaseURL(url string) Option {
	return func(o *options) {
		o.baseURL = url
	}
}

func applyOptions(opts []Option) options {
	var o options
	for _, fn := range opts {
		fn(&o)
	}
	return o
}

// GenerateCommitMessage generates a commit message from a git diff using the configured AI mode.
func GenerateCommitMessage(gitDiff string, cfg *config.Config, opts ...Option) (string, error) {
	mode := cfg.EffectiveAIMode()

	switch mode {
	case config.AILocal:
		prompt := buildCommitPrompt(gitDiff, cfg)
		return requestLocal(cfg.AI.LocalURL, cfg.AI.Model, prompt)
	case config.AICustom:
		prompt := buildCommitPrompt(gitDiff, cfg)
		return requestCustom(cfg.AI.APIURL, cfg.AI.APIKey, cfg.AI.Model, prompt)
	default:
		return requestRemoteCommit(gitDiff, cfg, applyOptions(opts))
	}
}

// GenerateTimeline generates a timeline analysis using the configured AI mode.
func GenerateTimeline(prompt string, cfg *config.Config, opts ...Option) (string, error) {
	mode := cfg.EffectiveAIMode()

	switch mode {
	case config.AILocal:
		return requestLocal(cfg.AI.LocalURL, cfg.AI.Model, prompt)
	case config.AICustom:
		return requestCustom(cfg.AI.APIURL, cfg.AI.APIKey, cfg.AI.Model, prompt)
	default:
		return requestRemoteTimeline(prompt, cfg, applyOptions(opts))
	}
}

// GenerateChangelog generates a changelog using the configured AI mode.
func GenerateChangelog(prompt string, cfg *config.Config, opts ...Option) (string, error) {
	return GenerateTimeline(prompt, cfg, opts...)
}

// --- Remote (default diny server) ---

type remoteCommitRequest struct {
	GitDiff string         `json:"gitDiff"`
	Version string         `json:"version"`
	Name    string         `json:"name"`
	Config  *config.Config `json:"config"`
	System  string         `json:"system,omitempty"`
}

type remoteCommitData struct {
	CommitMessage string `json:"commitMessage"`
}

type remoteCommitResp struct {
	Error *string           `json:"error,omitempty"`
	Data  *remoteCommitData `json:"data,omitempty"`
}

func requestRemoteCommit(gitDiff string, cfg *config.Config, o options) (string, error) {
	baseURL := server.ServerConfig.BaseURL
	if o.baseURL != "" {
		baseURL = o.baseURL
	}

	payload := remoteCommitRequest{
		Config:  cfg,
		GitDiff: gitDiff,
	}

	body, err := doPost(baseURL+"/api/v2/commit", nil, payload, 30*time.Second)
	if err != nil {
		return "", err
	}

	var out remoteCommitResp
	if err := json.Unmarshal(body, &out); err != nil {
		return "", fmt.Errorf("decode response: %w", err)
	}
	if out.Error != nil {
		return "", fmt.Errorf("%s", *out.Error)
	}
	if out.Data == nil || out.Data.CommitMessage == "" {
		return "", fmt.Errorf("empty commit message from server")
	}

	return out.Data.CommitMessage, nil
}

type remoteTimelineRequest struct {
	Prompt string         `json:"prompt"`
	Config *config.Config `json:"config"`
}

type remoteTimelineData struct {
	Message string `json:"message"`
}

type remoteTimelineResp struct {
	Error *string             `json:"error,omitempty"`
	Data  *remoteTimelineData `json:"data,omitempty"`
}

func requestRemoteTimeline(prompt string, cfg *config.Config, o options) (string, error) {
	baseURL := server.ServerConfig.BaseURL
	if o.baseURL != "" {
		baseURL = o.baseURL
	}

	payload := remoteTimelineRequest{
		Prompt: prompt,
		Config: cfg,
	}

	body, err := doPost(baseURL+"/api/v2/timeline", nil, payload, 30*time.Second)
	if err != nil {
		return "", err
	}

	var out remoteTimelineResp
	if err := json.Unmarshal(body, &out); err != nil {
		return "", fmt.Errorf("decode response: %w", err)
	}
	if out.Error != nil {
		return "", fmt.Errorf("%s", *out.Error)
	}
	if out.Data == nil || out.Data.Message == "" {
		return "", fmt.Errorf("empty message from server")
	}

	return out.Data.Message, nil
}

// --- Local (Ollama-compatible) ---

type localChatRequest struct {
	Model    string             `json:"model"`
	Messages []localChatMessage `json:"messages"`
	Stream   bool               `json:"stream"`
}

type localChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type localChatResponse struct {
	Message *localChatMessage `json:"message,omitempty"`
}

func requestLocal(localURL, model, prompt string) (string, error) {
	if model == "" {
		model = "llama3"
	}

	payload := localChatRequest{
		Model: model,
		Messages: []localChatMessage{
			{Role: "user", Content: prompt},
		},
		Stream: false,
	}

	client := newLocalHTTPClient(localURL)

	buf, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("marshal payload: %w", err)
	}

	req, err := http.NewRequestWithContext(context.Background(),
		http.MethodPost, localURL+"/api/chat", bytes.NewReader(buf))
	if err != nil {
		return "", fmt.Errorf("new request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("local request: %w", err)
	}
	defer res.Body.Close()

	body, _ := io.ReadAll(res.Body)

	var out localChatResponse
	if err := json.Unmarshal(body, &out); err != nil {
		return "", fmt.Errorf("decode local response: %w", err)
	}

	if out.Message == nil || out.Message.Content == "" {
		return "", fmt.Errorf("empty response from local AI")
	}

	return strings.TrimSpace(out.Message.Content), nil
}

// --- Custom API (OpenAI-compatible) ---

type customChatRequest struct {
	Model    string               `json:"model"`
	Messages []customChatMessage  `json:"messages"`
}

type customChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type customChatChoice struct {
	Message customChatMessage `json:"message"`
}

type customChatResponse struct {
	Choices []customChatChoice `json:"choices"`
}

func requestCustom(apiURL, apiKey, model, prompt string) (string, error) {
	if model == "" {
		model = "gpt-4"
	}

	payload := customChatRequest{
		Model: model,
		Messages: []customChatMessage{
			{Role: "user", Content: prompt},
		},
	}

	headers := map[string]string{
		"Authorization": "Bearer " + apiKey,
	}

	body, err := doPost(apiURL, headers, payload, 60*time.Second)
	if err != nil {
		return "", err
	}

	var out customChatResponse
	if err := json.Unmarshal(body, &out); err != nil {
		return "", fmt.Errorf("decode custom response: %w", err)
	}

	if len(out.Choices) == 0 || out.Choices[0].Message.Content == "" {
		return "", fmt.Errorf("empty response from custom AI")
	}

	return strings.TrimSpace(out.Choices[0].Message.Content), nil
}

// --- Local-mode network restriction ---

// newLocalHTTPClient returns an HTTP client that only allows requests
// to the configured local URL. All other requests are blocked.
func newLocalHTTPClient(allowedURL string) *http.Client {
	parsed, _ := url.Parse(allowedURL)
	allowedHost := ""
	if parsed != nil {
		allowedHost = parsed.Host
	}

	return &http.Client{
		Timeout: 60 * time.Second,
		Transport: &localOnlyTransport{
			allowedHost: allowedHost,
			inner:       http.DefaultTransport,
		},
	}
}

type localOnlyTransport struct {
	allowedHost string
	inner       http.RoundTripper
}

func (t *localOnlyTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	if req.URL.Host != t.allowedHost {
		return nil, fmt.Errorf("request to %s blocked: local mode only allows requests to %s", req.URL.Host, t.allowedHost)
	}
	return t.inner.RoundTrip(req)
}

// --- Shared helpers ---

func doPost(url string, headers map[string]string, payload any, timeout time.Duration) ([]byte, error) {
	buf, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("marshal payload: %w", err)
	}

	req, err := http.NewRequestWithContext(context.Background(),
		http.MethodPost, url, bytes.NewReader(buf))
	if err != nil {
		return nil, fmt.Errorf("new request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	for k, v := range headers {
		req.Header.Set(k, v)
	}

	client := &http.Client{Timeout: timeout}
	res, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("do request: %w", err)
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("read response: %w", err)
	}

	if res.StatusCode < 200 || res.StatusCode >= 300 {
		detail := string(body)
		if len(detail) > 200 {
			detail = detail[:200]
		}
		if detail == "" {
			detail = "(empty body)"
		}
		return nil, fmt.Errorf("HTTP %d from %s: %s", res.StatusCode, url, detail)
	}

	if len(body) == 0 {
		return nil, fmt.Errorf("empty response body from %s", url)
	}

	return body, nil
}

// buildCommitPrompt creates a prompt for local/custom AI models from the git diff and config.
func buildCommitPrompt(gitDiff string, cfg *config.Config) string {
	var sb strings.Builder

	sb.WriteString("Generate a git commit message for the following diff.\n\n")

	sb.WriteString(fmt.Sprintf("Tone: %s\n", cfg.Commit.Tone))
	sb.WriteString(fmt.Sprintf("Length: %s\n", cfg.Commit.Length))

	if cfg.Commit.Conventional {
		sb.WriteString(fmt.Sprintf("Use conventional commit format with types: %s\n",
			strings.Join(cfg.Commit.ConventionalFormat, ", ")))
	}

	if cfg.Commit.Emoji {
		sb.WriteString("Include an emoji prefix.\n")
	}

	switch cfg.Commit.Length {
	case config.Short:
		sb.WriteString("Keep it to a subject line only, 1-2 sentences max, <=60 chars, imperative verb first. No bullets.\n")
	case config.Normal:
		sb.WriteString("Subject <=70 chars (imperative). If needed, add 1-3 terse bullets for WHY/impact.\n")
	case config.Long:
		sb.WriteString("Subject <=80 chars (imperative). Then 2-6 terse bullets for context/impact.\n")
	}

	if cfg.Commit.CustomInstructions != "" {
		sb.WriteString(fmt.Sprintf("\nAdditional instructions: %s\n", cfg.Commit.CustomInstructions))
	}

	sb.WriteString(fmt.Sprintf("\nDiff:\n%s\n", gitDiff))
	sb.WriteString("\nRespond with ONLY the commit message, no explanations or extra text.")

	return sb.String()
}
