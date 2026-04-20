package prompts

import (
	"fmt"
	"os/exec"
	"runtime"
	"time"

	"github.com/dinoDanic/diny/config"
	tuiprompts "github.com/dinoDanic/diny/tui/prompts"
)

// ShowStar runs the star TUI and mutates state.
// Returns "starred" | "already_given" | "dismissed"; "" on error.
func ShowStar(state *State, cfg *config.Config) string {
	res := tuiprompts.RunStar(cfg)
	if res.Outcome == "" {
		return ""
	}

	now := time.Now()
	switch res.Outcome {
	case "starred":
		state.Prompts.Star.Status = StatusStarred
		state.Prompts.Star.AnsweredAt = &now
		openBrowser(GitHubRepoURL)
	case "already_given":
		state.Prompts.Star.Status = StatusAlreadyGiven
		state.Prompts.Star.AnsweredAt = &now
	case "dismissed":
		state.Prompts.Star.Status = StatusDismissed
		state.Prompts.Star.AnsweredAt = &now
	}

	return res.Outcome
}

// openBrowser opens the given URL in the default browser.
// If the command fails or is not found, it prints the URL as a fallback.
func openBrowser(url string) {
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "darwin":
		cmd = exec.Command("open", url)
	case "windows":
		cmd = exec.Command("rundll32", "url.dll,FileProtocolHandler", url)
	default: // linux, freebsd, etc.
		cmd = exec.Command("xdg-open", url)
	}

	if err := cmd.Start(); err != nil {
		fmt.Printf("\nOpen this URL to star: %s\n", url)
	}
}
