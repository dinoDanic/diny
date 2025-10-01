/*
Copyright © 2025 NAME HERE dino.danic@gmail.com
*/
package update

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
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

	return current != latest
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
	method := uc.DetectInstallMethod()

	var updateCmd string
	switch method {
	case "brew":
		updateCmd = "brew upgrade dinoDanic/tap/diny"
	case "scoop":
		updateCmd = "scoop update diny"
	case "winget":
		updateCmd = "winget upgrade dinoDanic.diny"
	default:
		updateCmd = "download from https://github.com/dinoDanic/diny/releases"
	}

	content := fmt.Sprintf("New version %s available!\n\nUpdate with: diny update\nOr manually: %s\n\nUpdate is crucial due to early development stage.", version, updateCmd)

	ui.Box(ui.BoxOptions{Message: content, Variant: ui.Warning})
}

func (uc *UpdateChecker) GetLatestVersion() (string, error) {
	return uc.getLatestVersion()
}

func (uc *UpdateChecker) CompareVersions(current, latest string) bool {
	return uc.compareVersions(current, latest)
}

func (uc *UpdateChecker) DetectInstallMethod() string {
	execPath, err := os.Executable()
	if err != nil {
		return "manual"
	}

	execPath, err = filepath.EvalSymlinks(execPath)
	if err != nil {
		execPath, _ = os.Executable()
	}

	switch runtime.GOOS {
	case "darwin", "linux":
		if strings.Contains(execPath, "/Cellar/") || strings.Contains(execPath, "/homebrew/") || strings.Contains(execPath, "/linuxbrew/") {
			return "brew"
		}

		cmd := exec.Command("brew", "list", "diny")
		if cmd.Run() == nil {
			return "brew"
		}

	case "windows":
		if strings.Contains(execPath, "\\scoop\\apps\\") {
			return "scoop"
		}

		cmd := exec.Command("winget", "list", "--id", "dinoDanic.diny")
		cmd.Stdout = nil
		cmd.Stderr = nil
		if cmd.Run() == nil {
			return "winget"
		}
	}

	return "manual"
}

func (uc *UpdateChecker) PerformUpdate() error {
	method := uc.DetectInstallMethod()

	ui.RenderTitle(fmt.Sprintf("Updating diny using %s...", method))

	switch method {
	case "brew":
		cmd := exec.Command("brew", "upgrade", "dinoDanic/tap/diny")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("brew upgrade failed: %w", err)
		}
		ui.Box(ui.BoxOptions{Message: "Successfully updated diny via Homebrew!", Variant: ui.Success})
		return nil

	case "scoop":
		cmd := exec.Command("scoop", "update", "diny")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("scoop update failed: %w", err)
		}
		ui.Box(ui.BoxOptions{Message: "Successfully updated diny via Scoop!", Variant: ui.Success})
		return nil

	case "winget":
		cmd := exec.Command("winget", "upgrade", "dinoDanic.diny")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("winget upgrade failed: %w", err)
		}
		ui.Box(ui.BoxOptions{Message: "Successfully updated diny via Winget!", Variant: ui.Success})
		return nil

	default:
		return fmt.Errorf("manual installation detected\n\nPlease download the latest version from:\nhttps://github.com/dinoDanic/diny/releases")
	}
}
