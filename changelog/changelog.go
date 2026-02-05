package changelog

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/atotto/clipboard"
	"github.com/charmbracelet/huh"
	"github.com/dinoDanic/diny/config"
	"github.com/dinoDanic/diny/ui"
)

func Main(cfg *config.Config) {
	var releases []GitHubRelease

	err := ui.WithSpinner("Fetching releases from GitHub...", func() error {
		client := NewReleaseClient()
		var fetchErr error
		releases, fetchErr = client.FetchReleases()
		return fetchErr
	})

	if err != nil {
		ui.Box(ui.BoxOptions{
			Variant: ui.Error,
			Title:   "Error",
			Message: fmt.Sprintf("Failed to fetch releases: %v", err),
		})
		os.Exit(1)
	}

	if len(releases) == 0 {
		ui.Box(ui.BoxOptions{
			Variant: ui.Warning,
			Title:   "No Releases Found",
			Message: "No releases are available for this project.",
		})
		os.Exit(0)
	}

	selectedRelease := selectReleasePrompt(releases)

	displayChangelog(selectedRelease)

	handleChangelogFlow(selectedRelease, releases, cfg)
}

func selectReleasePrompt(releases []GitHubRelease) GitHubRelease {
	options := make([]huh.Option[int], len(releases))

	for i, release := range releases {
		label := fmt.Sprintf("%s (%s)",
			release.TagName,
			release.PublishedAt.Format("Jan 02, 2006"),
		)
		if release.Prerelease {
			label += " [pre-release]"
		}
		options[i] = huh.NewOption(label, i)
	}

	var selectedIndex int

	err := huh.NewSelect[int]().
		Title("Select a release to view changelog").
		Description("Use arrow keys or j,k to navigate, Enter to select").
		Options(options...).
		Value(&selectedIndex).
		Height(12).
		WithTheme(ui.GetHuhPrimaryTheme()).
		Run()

	if err != nil {
		os.Exit(0)
	}

	return releases[selectedIndex]
}

func displayChangelog(release GitHubRelease) {
	cleanedBody := cleanChangelogBody(release.Body)

	ui.Box(ui.BoxOptions{
		Variant: ui.Primary,
		Title:   fmt.Sprintf("Changelog for %s", release.TagName),
		Message: cleanedBody,
	})
}

func cleanChangelogBody(body string) string {
	lines := strings.Split(body, "\n")
	cleaned := make([]string, 0, len(lines))

	for _, line := range lines {
		if strings.Contains(line, "Released by GoReleaser") {
			break
		}
		cleaned = append(cleaned, line)
	}

	result := strings.TrimSpace(strings.Join(cleaned, "\n"))
	if result == "" {
		return "No changelog content available."
	}

	return result
}

func handleChangelogFlow(release GitHubRelease, allReleases []GitHubRelease, cfg *config.Config) {
	for {
		action := changelogActionPrompt()

		switch action {
		case "browser":
			openInBrowser(release.HTMLURL)
		case "clipboard":
			copyToClipboard(release.Body)
		case "another":
			Main(cfg)
			return
		case "exit":
			os.Exit(0)
		}
	}
}

func changelogActionPrompt() string {
	options := []huh.Option[string]{
		huh.NewOption("View in browser", "browser"),
		huh.NewOption("Copy to clipboard", "clipboard"),
		huh.NewOption("View another release", "another"),
		huh.NewOption("Exit", "exit"),
	}

	var choice string

	err := huh.NewSelect[string]().
		Title("What would you like to do next?").
		Description("Select an option using arrow keys or j,k and press Enter").
		Options(options...).
		Value(&choice).
		Height(6).
		WithTheme(ui.GetHuhPrimaryTheme()).
		Run()

	if err != nil {
		os.Exit(0)
	}

	return choice
}

func openInBrowser(url string) {
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "darwin":
		cmd = exec.Command("open", url)
	case "linux":
		cmd = exec.Command("xdg-open", url)
	case "windows":
		cmd = exec.Command("cmd", "/c", "start", url)
	default:
		ui.Box(ui.BoxOptions{
			Variant: ui.Warning,
			Title:   "Unsupported Platform",
			Message: fmt.Sprintf("Cannot open browser on %s. Visit: %s", runtime.GOOS, url),
		})
		return
	}

	err := cmd.Start()
	if err != nil {
		ui.Box(ui.BoxOptions{
			Variant: ui.Warning,
			Title:   "Browser Open Failed",
			Message: fmt.Sprintf("Could not open browser: %v\nVisit: %s", err, url),
		})
		return
	}

	ui.Box(ui.BoxOptions{
		Variant: ui.Success,
		Title:   "Opened in Browser",
		Message: "Release page opened in your default browser.",
	})
}

func copyToClipboard(content string) {
	err := clipboard.WriteAll(content)
	if err != nil {
		ui.Box(ui.BoxOptions{
			Variant: ui.Error,
			Title:   "Clipboard Error",
			Message: fmt.Sprintf("Failed to copy to clipboard: %v", err),
		})
		return
	}

	ui.Box(ui.BoxOptions{
		Variant: ui.Success,
		Title:   "Copied to Clipboard",
		Message: "Changelog has been copied to your clipboard.",
	})
}
