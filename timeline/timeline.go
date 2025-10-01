package timeline

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/charmbracelet/huh"
	"github.com/dinoDanic/diny/config"
	"github.com/dinoDanic/diny/git"
	"github.com/dinoDanic/diny/groq"
	"github.com/dinoDanic/diny/ui"
)

func Main() {
	choice := timelinePrompt("Choose timeline for commit analysis:")

	var timelineCommits []string
	var dateRange string
	var err error

	switch choice {
	case "today":
		ui.RenderTitle("Analyzing today's commits...")
		timelineCommits, err = git.GetCommitsToday()
		dateRange = "today"
	case "date":
		selectedDate := dateInputPrompt("Enter date (DD MM YYYY):")
		ui.RenderTitle(fmt.Sprintf("Analyzing commits from %s...", selectedDate))
		timelineCommits, err = git.GetCommitsByDate(selectedDate)
		dateRange = selectedDate
	case "range":
		startDate := dateInputPrompt("Enter start date (DD MM YYYY):")
		endDate := dateInputPrompt("Enter end date (DD MM YYYY):")
		ui.RenderTitle(fmt.Sprintf("Analyzing commits from %s to %s...", startDate, endDate))
		timelineCommits, err = git.GetCommitsByDateRange(startDate+" 00:00:00", endDate+" 23:59:59")
		dateRange = fmt.Sprintf("%s to %s", startDate, endDate)
	}

	if err != nil {
		ui.Box(ui.BoxOptions{Message: fmt.Sprintf("Failed to get timeline commits: %v", err), Variant: ui.Error})
		os.Exit(1)
	}

	if len(timelineCommits) == 0 {
		ui.Box(ui.BoxOptions{Message: fmt.Sprintf("No commits found for the selected period (%s).", dateRange), Variant: ui.Warning})
		return
	}

	commitList := ""
	for i, commit := range timelineCommits {
		commitList += fmt.Sprintf("%d. %s\n", i+1, commit)
	}
	ui.Box(ui.BoxOptions{Title: fmt.Sprintf("Found %d commits from %s", len(timelineCommits), dateRange), Message: strings.TrimSpace(commitList)})

	userConfig, err := config.Load()
	prompt := fmt.Sprintf("Timeline: %s\nCommits:\n%s", dateRange, strings.Join(timelineCommits, "n"))

	var analysis string
	err = ui.WithSpinner("Generating timeline analysis...", func() error {
		var genErr error
		analysis, genErr = groq.CreateTimelineWithGroq(prompt, userConfig)
		return genErr
	})

	if err != nil {
		ui.Box(ui.BoxOptions{Message: fmt.Sprintf("Failed to generate analysis: %v", err), Variant: ui.Error})
		os.Exit(1)
	}

	ui.Box(ui.BoxOptions{Title: "Timeline Analysis", Message: analysis})
}

func timelinePrompt(message string) string {
	var choice string

	err := huh.NewSelect[string]().
		Title("🦕 "+message).
		Description("Select an option using arrow keys or j,k and press Enter").
		Options(
			huh.NewOption("Today", "today"),
			huh.NewOption("Pick specific date", "date"),
			huh.NewOption("Choose date range", "range"),
		).
		Value(&choice).
		Height(5).
		Run()

	if err != nil {
		ui.Box(ui.BoxOptions{Message: fmt.Sprintf("Error running prompt: %v", err), Variant: ui.Error})
		os.Exit(1)
	}

	return choice
}

func dateInputPrompt(message string) string {
	var input string

	err := huh.NewInput().
		Title("🦕 " + message).
		Description("Use format: DD MM YYYY (e.g., 15 01 2025)").
		Placeholder("15 01 2025").
		Validate(func(s string) error {
			_, err := time.Parse("02 01 2006", s)
			if err != nil {
				return fmt.Errorf("invalid date format, use DD MM YYYY")
			}
			return nil
		}).
		Value(&input).
		Run()

	if err != nil {
		ui.Box(ui.BoxOptions{Message: fmt.Sprintf("Error running prompt: %v", err), Variant: ui.Error})
		os.Exit(1)
	}

	return input
}
