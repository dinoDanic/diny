package config

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/dinoDanic/diny/git"
)

type GitInfo struct {
	RepoName  string `json:"repoName"`
	RepoOwner string `json:"repoOwner"`
	RepoURL   string `json:"repoURL"`
}

func GetGitInfo() (*GitInfo, error) {
	gitRoot, err := git.FindGitRoot()
	if err != nil {
		return nil, fmt.Errorf("not in a git repository: %w", err)
	}

	configPath := filepath.Join(gitRoot, ".git", "config")
	file, err := os.Open(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open git config: %w", err)
	}
	defer file.Close()

	var originURL string
	scanner := bufio.NewScanner(file)
	inOriginSection := false

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		if line == `[remote "origin"]` {
			inOriginSection = true
			continue
		}

		if inOriginSection && strings.HasPrefix(line, "[") && line != `[remote "origin"]` {
			inOriginSection = false
			continue
		}

		if inOriginSection && strings.HasPrefix(line, "url = ") {
			originURL = strings.TrimPrefix(line, "url = ")
			break
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("failed to read git config: %w", err)
	}

	if originURL == "" {
		return nil, fmt.Errorf("no origin remote found in git config")
	}

	repoOwner, repoName, err := parseGitURL(originURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse git URL: %w", err)
	}

	return &GitInfo{
		RepoName:  repoName,
		RepoOwner: repoOwner,
		RepoURL:   originURL,
	}, nil
}

func parseGitURL(url string) (owner, repo string, err error) {
	url = strings.TrimSuffix(url, ".git")

	httpsRegex := regexp.MustCompile(`https?://[^/]+/([^/]+)/([^/]+)`)
	if matches := httpsRegex.FindStringSubmatch(url); len(matches) == 3 {
		return matches[1], matches[2], nil
	}

	sshRegex := regexp.MustCompile(`git@[^:]+:([^/]+)/([^/]+)`)
	if matches := sshRegex.FindStringSubmatch(url); len(matches) == 3 {
		return matches[1], matches[2], nil
	}

	gitRegex := regexp.MustCompile(`git://[^/]+/([^/]+)/([^/]+)`)
	if matches := gitRegex.FindStringSubmatch(url); len(matches) == 3 {
		return matches[1], matches[2], nil
	}

	return "", "", fmt.Errorf("unsupported git URL format: %s", url)
}
