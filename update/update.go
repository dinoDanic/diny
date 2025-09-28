/*
Copyright © 2025 NAME HERE dino.danic@gmail.com
*/
package update

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/dinoDanic/diny/ui"
)

type GitHubRelease struct {
	TagName string `json:"tag_name"`
}

type UpdateChecker struct {
	currentVersion string
	repoURL        string
}

func NewUpdateChecker(currentVersion string) *UpdateChecker {
	return &UpdateChecker{
		currentVersion: currentVersion,
		repoURL:        "https://api.github.com/repos/dinoDanic/diny/releases/latest",
	}
}

func (uc *UpdateChecker) getLatestVersion() (string, error) {
	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Get(uc.repoURL)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return "", fmt.Errorf("GitHub API returned status %d", resp.StatusCode)
	}

	var release GitHubRelease
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return "", err
	}

	return release.TagName, nil
}

func (uc *UpdateChecker) compareVersions(current, latest string) bool {
	current = strings.TrimPrefix(current, "v")
	latest = strings.TrimPrefix(latest, "v")

	if current == "dev" {
		return false
	}

	current = uc.normalizeVersion(current)
	latest = uc.normalizeVersion(latest)

	currentParts := strings.Split(current, ".")
	latestParts := strings.Split(latest, ".")

	maxLen := len(currentParts)
	if len(latestParts) > maxLen {
		maxLen = len(latestParts)
	}

	for i := 0; i < maxLen; i++ {
		var currentPart, latestPart int

		if i < len(currentParts) {
			currentPart, _ = strconv.Atoi(currentParts[i])
		}
		if i < len(latestParts) {
			latestPart, _ = strconv.Atoi(latestParts[i])
		}

		if latestPart > currentPart {
			return true
		} else if latestPart < currentPart {
			return false
		}
	}

	return false
}

func (uc *UpdateChecker) normalizeVersion(version string) string {
	if idx := strings.Index(version, "-"); idx != -1 {
		return version[:idx]
	}
	return version
}

func (uc *UpdateChecker) CheckForUpdate() {
	latestVersion, err := uc.getLatestVersion()

	if err != nil {
		return
	}

	if uc.compareVersions(uc.currentVersion, latestVersion) {
		uc.printUpdateNotification(latestVersion)
	}
}

func (uc *UpdateChecker) printUpdateNotification(version string) {
	content := fmt.Sprintf("New version %s available!\n\nPlease update with: diny update\nOr visit: https://github.com/dinoDanic/diny\n\nUpdate is crucial due to early development stage.", version)

	ui.RenderWarning(content)
}

func (uc *UpdateChecker) GetLatestVersion() (string, error) {
	return uc.getLatestVersion()
}

func (uc *UpdateChecker) CompareVersions(current, latest string) bool {
	return uc.compareVersions(current, latest)
}
