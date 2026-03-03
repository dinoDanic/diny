package git

import (
	"fmt"
	"os/exec"
	"strings"
)

type CommitInfo struct {
	SHA     string
	Message string
}

func GetTags() ([]string, error) {
	cmd := exec.Command("git", "tag", "-l", "--sort=-version:refname")
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to list tags: %w", err)
	}

	raw := strings.TrimSpace(string(output))
	if raw == "" {
		return []string{}, nil
	}

	tags := strings.Split(raw, "\n")
	var result []string
	for _, t := range tags {
		if t = strings.TrimSpace(t); t != "" {
			result = append(result, t)
		}
	}
	return result, nil
}

func GetRecentCommits(limit int) ([]CommitInfo, error) {
	cmd := exec.Command("git", "log",
		fmt.Sprintf("--pretty=format:%%h|||%%s"),
		"--no-merges",
		fmt.Sprintf("-n%d", limit),
	)
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to get recent commits: %w", err)
	}

	raw := strings.TrimSpace(string(output))
	if raw == "" {
		return []CommitInfo{}, nil
	}

	lines := strings.Split(raw, "\n")
	var commits []CommitInfo
	for _, line := range lines {
		parts := strings.SplitN(line, "|||", 2)
		if len(parts) != 2 {
			continue
		}
		commits = append(commits, CommitInfo{
			SHA:     strings.TrimSpace(parts[0]),
			Message: strings.TrimSpace(parts[1]),
		})
	}
	return commits, nil
}

func GetDiffBetweenRefs(ref1, ref2 string) (string, error) {
	cmd := exec.Command("git", "diff", ref1+"..."+ref2,
		"-U3", "--no-color", "--ignore-all-space", "--ignore-blank-lines",
		":(exclude)*.lock", ":(exclude)*package-lock.json", ":(exclude)*yarn.lock",
		":(exclude)node_modules/", ":(exclude)dist/", ":(exclude)build/",
	)
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to get diff between %s and %s: %w", ref1, ref2, err)
	}
	return string(output), nil
}

func GetCommitsBetweenRefs(ref1, ref2 string) ([]string, error) {
	cmd := exec.Command("git", "log", ref1+".."+ref2,
		"--pretty=format:%s",
		"--no-merges",
	)
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to get commits between %s and %s: %w", ref1, ref2, err)
	}

	raw := strings.TrimSpace(string(output))
	if raw == "" {
		return []string{}, nil
	}

	lines := strings.Split(raw, "\n")
	var commits []string
	for _, line := range lines {
		if line = strings.TrimSpace(line); line != "" {
			commits = append(commits, line)
		}
	}
	return commits, nil
}
