/*
Copyright © 2025 NAME HERE dino.danic@gmail.com
*/
package git

import (
	"os/exec"
	"path/filepath"
	"strings"
)

func GetGitName() string {
	cmd := exec.Command("git", "config", "user.name")
	output, err := cmd.Output()
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(output))
}

func GetRepoName() string {
	cmd := exec.Command("git", "config", "--get", "remote.origin.url")
	output, err := cmd.Output()
	if err != nil {
		cwd, cwdErr := exec.Command("pwd").Output()
		if cwdErr != nil {
			return ""
		}
		return filepath.Base(strings.TrimSpace(string(cwd)))
	}

	url := strings.TrimSpace(string(output))
	url = strings.TrimSuffix(url, ".git")

	parts := strings.Split(url, "/")
	if len(parts) > 0 {
		return parts[len(parts)-1]
	}

	return ""
}
