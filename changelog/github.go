package changelog

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type GitHubRelease struct {
	TagName     string    `json:"tag_name"`
	Name        string    `json:"name"`
	PublishedAt time.Time `json:"published_at"`
	Body        string    `json:"body"`
	HTMLURL     string    `json:"html_url"`
	Prerelease  bool      `json:"prerelease"`
	Draft       bool      `json:"draft"`
}

type ReleaseClient struct {
	repoURL string
	timeout time.Duration
}

func NewReleaseClient() *ReleaseClient {
	return &ReleaseClient{
		repoURL: "https://api.github.com/repos/dinoDanic/diny/releases",
		timeout: 10 * time.Second,
	}
}

func (c *ReleaseClient) FetchReleases() ([]GitHubRelease, error) {
	client := &http.Client{
		Timeout: c.timeout,
	}

	req, err := http.NewRequest("GET", c.repoURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Accept", "application/vnd.github.v3+json")

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch releases: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == 403 {
		return nil, fmt.Errorf("GitHub API rate limit exceeded. Please try again later")
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GitHub API returned status %d", resp.StatusCode)
	}

	var releases []GitHubRelease
	if err := json.NewDecoder(resp.Body).Decode(&releases); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	filtered := make([]GitHubRelease, 0, len(releases))
	for _, release := range releases {
		if !release.Draft {
			filtered = append(filtered, release)
		}
	}

	if len(filtered) > 30 {
		filtered = filtered[:30]
	}

	return filtered, nil
}
