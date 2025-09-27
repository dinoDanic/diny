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
)

func Main() {

	// Show date selection menu
	choice := timelinePrompt("📅 Choose timeline for commit analysis:")
	fmt.Println()

	var timelineCommits []string
	var dateRange string
	var err error

	switch choice {
	case "today":
		fmt.Println("🦕 Analyzing today's commits...")
		timelineCommits, err = git.GetCommitsToday()
		dateRange = "today"
	case "date":
		selectedDate := dateInputPrompt("Enter date (DD MM YYYY):")
		fmt.Printf("🦕 Analyzing commits from %s...\n", selectedDate)
		timelineCommits, err = git.GetCommitsByDate(selectedDate)
		dateRange = selectedDate
	case "range":
		startDate := dateInputPrompt("Enter start date (DD MM YYYY):")
		endDate := dateInputPrompt("Enter end date (DD MM YYYY):")
		fmt.Printf("🦕 Analyzing commits from %s to %s...\n", startDate, endDate)
		timelineCommits, err = git.GetCommitsByDateRange(startDate+" 00:00:00", endDate+" 23:59:59")
		dateRange = fmt.Sprintf("%s to %s", startDate, endDate)
	}

	if err != nil {
		fmt.Printf("❌ Failed to get timeline commits: %v\n", err)
		os.Exit(1)
	}

	if len(timelineCommits) == 0 {
		fmt.Printf("🦕 No commits found for the selected period (%s).\n", dateRange)
		return
	}

	fmt.Printf("🦕 Found %d commits from %s:\n\n", len(timelineCommits), dateRange)

	// Display the commit messages
	for i, commit := range timelineCommits {
		fmt.Printf("%d. %s\n", i+1, commit)
	}

	userConfig, err := config.Load()

	fmt.Println()
	if err == nil && userConfig != nil {
		config.PrintConfiguration(*userConfig)
	}

	prompt := fmt.Sprintf("Timeline: %s\nCommits:\n%s", dateRange, strings.Join(timelineCommits, "\n"))

	fmt.Println()
	fmt.Println("🦕 Generating timeline analysis...")
	fmt.Println()

	// Call Groq API for analysis
	analysis, err := groq.CreateTimelineWithGroq(prompt, userConfig)
	if err != nil {
		fmt.Printf("❌ Failed to generate analysis: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("\n📊 Timeline Analysis:\n%s\n", analysis)
}

func timelinePrompt(message string) string {
	var choice string

	err := huh.NewSelect[string]().
		Title(message).
		Options(
			huh.NewOption("Today", "today"),
			huh.NewOption("Pick specific date", "date"),
			huh.NewOption("Choose date range", "range"),
		).
		Value(&choice).
		Run()

	if err != nil {
		fmt.Printf("Error running prompt: %v\n", err)
		os.Exit(1)
	}

	return choice
}

func dateInputPrompt(message string) string {
	var input string

	err := huh.NewInput().
		Title(message).
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
		fmt.Printf("Error running prompt: %v\n", err)
		os.Exit(1)
	}

	return input
}
