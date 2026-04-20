package prompts

import (
	"fmt"
	"os/exec"
	"runtime"
	"time"

	"github.com/dinoDanic/diny/ui"
)

// ShowStar displays the GitHub star prompt and handles the user's keypress.
// Returns the outcome string for the backend: "starred", "already_given", "dismissed",
// or "" on read error.
func ShowStar(state *State) string {
	ui.Box("Diny is free & open source. A GitHub star helps a lot!",
		"[1] Star it  [2] Don't ask again  [3] Already starred  [0] Close")

	for {
		r := readSingleRune()
		if r == 0 {
			return "" // read error
		}

		now := time.Now()

		switch r {
		case '1':
			state.Prompts.Star.Status = StatusStarred
			state.Prompts.Star.AnsweredAt = &now
			openBrowser(GitHubRepoURL)
			fmt.Println()
			ui.Success("Thanks for starring!")
			return "starred"

		case '2', '0':
			state.Prompts.Star.Status = StatusDismissed
			state.Prompts.Star.AnsweredAt = &now
			return "dismissed"

		case '3':
			state.Prompts.Star.Status = StatusAlreadyGiven
			state.Prompts.Star.AnsweredAt = &now
			fmt.Println()
			ui.Success("Thanks!")
			return "already_given"

		default:
			continue
		}
	}
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
