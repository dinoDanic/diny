package timeline

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/atotto/clipboard"
	"github.com/charmbracelet/huh"
	"github.com/dinoDanic/diny/config"
	"github.com/dinoDanic/diny/git"
	"github.com/dinoDanic/diny/groq"
	"github.com/dinoDanic/diny/ui"
)

func Main(cfg *config.Config) {
	choice := timelinePrompt("Choose timeline for commit analysis:")

	var timelineCommits []string
	var dateRange string
	var err error

	switch choice {
	case "today":
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

	prompt := fmt.Sprintf("Timeline: %s\nCommits:\n%s", dateRange, strings.Join(timelineCommits, "n"))

	var analysis string
	err = ui.WithSpinner("Generating timeline analysis...", func() error {
		var genErr error
		analysis, genErr = groq.CreateTimelineWithGroq(prompt, cfg)
		return genErr
	})

	if err != nil {
		ui.Box(ui.BoxOptions{Message: fmt.Sprintf("Failed to generate analysis: %v", err), Variant: ui.Error})
		os.Exit(1)
	}

	ui.Box(ui.BoxOptions{Title: "Timeline Analysis", Message: analysis})

	HandleTimelineFlow(analysis, prompt, cfg, dateRange, []string{})
}

func timelinePrompt(message string) string {
	var choice string

	err := huh.NewSelect[string]().
		Title(message).
		Description("Select an option using arrow keys or j,k and press Enter").
		Options(
			huh.NewOption("Today", "today"),
			huh.NewOption("Pick specific date", "date"),
			huh.NewOption("Choose date range", "range"),
		).
		Value(&choice).
		Height(5).
		WithTheme(ui.GetHuhPrimaryTheme()).
		Run()

	if err != nil {
		ui.Box(ui.BoxOptions{Message: fmt.Sprintf("Error running prompt: %v", err), Variant: ui.Error})
		os.Exit(1)
	}

	return choice
}

func dateInputPrompt(message string) string {
	now := time.Now()
	day := now.Day()
	month := int(now.Month())
	year := now.Year()

	err := huh.NewSelect[int]().
		Title(message).
		Description("Day").
		Options(generateDayOptions()...).
		Value(&day).
		Height(7).
		WithTheme(ui.GetHuhPrimaryTheme()).
		Run()
	if err != nil {
		ui.Box(ui.BoxOptions{Message: fmt.Sprintf("Error selecting day: %v", err), Variant: ui.Error})
		os.Exit(1)
	}

	err = huh.NewSelect[int]().
		Title(message).
		Description("Month").
		Options(generateMonthOptions()...).
		Value(&month).
		Height(7).
		WithTheme(ui.GetHuhPrimaryTheme()).
		Run()
	if err != nil {
		ui.Box(ui.BoxOptions{Message: fmt.Sprintf("Error selecting month: %v", err), Variant: ui.Error})
		os.Exit(1)
	}

	err = huh.NewSelect[int]().
		Title(message).
		Description("Year").
		Options(generateYearOptions()...).
		Value(&year).
		Height(7).
		WithTheme(ui.GetHuhPrimaryTheme()).
		Run()
	if err != nil {
		ui.Box(ui.BoxOptions{Message: fmt.Sprintf("Error selecting year: %v", err), Variant: ui.Error})
		os.Exit(1)
	}

	return fmt.Sprintf("%02d %02d %d", day, month, year)
}

func generateDayOptions() []huh.Option[int] {
	options := make([]huh.Option[int], 31)
	for i := 1; i <= 31; i++ {
		options[i-1] = huh.NewOption(fmt.Sprintf("%02d", i), i)
	}
	return options
}

func generateMonthOptions() []huh.Option[int] {
	months := []string{"January", "February", "March", "April", "May", "June", "July", "August", "September", "October", "November", "December"}
	options := make([]huh.Option[int], 12)
	for i := 1; i <= 12; i++ {
		options[i-1] = huh.NewOption(fmt.Sprintf("%02d - %s", i, months[i-1]), i)
	}
	return options
}

func generateYearOptions() []huh.Option[int] {
	currentYear := time.Now().Year()
	options := make([]huh.Option[int], 10)
	for i := 0; i < 10; i++ {
		year := currentYear - i
		options[i] = huh.NewOption(fmt.Sprintf("%d", year), year)
	}
	return options
}

func HandleTimelineFlow(analysis, fullPrompt string, cfg *config.Config, dateRange string, previousAnalyses []string) {
	choice := timelineChoicePrompt()

	switch choice {
	case "copy":
		if err := clipboard.WriteAll(analysis); err != nil {
			ui.Box(ui.BoxOptions{Message: fmt.Sprintf("Failed to copy to clipboard: %v", err), Variant: ui.Error})
			HandleTimelineFlow(analysis, fullPrompt, cfg, dateRange, previousAnalyses)
			return
		}
		ui.Box(ui.BoxOptions{Message: "Analysis copied to clipboard!", Variant: ui.Success})
		fmt.Println()
	case "save":
		filePath, err := saveTimelineAnalysis(analysis, dateRange)
		if err != nil {
			ui.Box(ui.BoxOptions{Message: fmt.Sprintf("Failed to save analysis: %v", err), Variant: ui.Error})
			HandleTimelineFlow(analysis, fullPrompt, cfg, dateRange, previousAnalyses)
			return
		}
		ui.Box(ui.BoxOptions{Message: fmt.Sprintf("Analysis saved!\n\n%s", filePath), Variant: ui.Success})
		fmt.Println()
	case "regenerate":
		modifiedPrompt := fullPrompt
		if len(previousAnalyses) > 0 {
			modifiedPrompt += "\n\nPrevious analyses that were not satisfactory:\n"
			for i, msg := range previousAnalyses {
				modifiedPrompt += fmt.Sprintf("%d. %s\n", i+1, msg)
			}
			modifiedPrompt += "\nPlease generate a different analysis with a different perspective or focus."
		} else {
			modifiedPrompt += "\n\nPlease provide an alternative analysis with a different approach or focus."
		}

		var newAnalysis string
		err := ui.WithSpinner("Generating alternative analysis...", func() error {
			var genErr error
			newAnalysis, genErr = groq.CreateTimelineWithGroq(modifiedPrompt, cfg)
			return genErr
		})
		if err != nil {
			ui.Box(ui.BoxOptions{Message: fmt.Sprintf("Error: %v", err), Variant: ui.Error})
			os.Exit(1)
		}

		ui.Box(ui.BoxOptions{Title: "Timeline Analysis", Message: newAnalysis})
		updatedHistory := append(previousAnalyses, analysis)
		HandleTimelineFlow(newAnalysis, fullPrompt, cfg, dateRange, updatedHistory)
	case "custom":
		customInput := customTimelineInputPrompt("What changes would you like to see in the analysis?")

		modifiedPrompt := fullPrompt + fmt.Sprintf("\n\nCurrent analysis:\n%s\n\nUser feedback: %s\n\nPlease generate a new analysis that addresses the user's feedback.", analysis, customInput)

		var newAnalysis string
		err := ui.WithSpinner("Refining analysis with your feedback...", func() error {
			var genErr error
			newAnalysis, genErr = groq.CreateTimelineWithGroq(modifiedPrompt, cfg)
			return genErr
		})
		if err != nil {
			ui.Box(ui.BoxOptions{Message: fmt.Sprintf("Error: %v", err), Variant: ui.Error})
			os.Exit(1)
		}

		ui.Box(ui.BoxOptions{Title: "Timeline Analysis", Message: newAnalysis})
		updatedHistory := append(previousAnalyses, analysis)
		HandleTimelineFlow(newAnalysis, fullPrompt, cfg, dateRange, updatedHistory)
	case "new":
		Main(cfg)
	case "exit":
		ui.RenderTitle("Bye!")
		os.Exit(0)
	}
}

func timelineChoicePrompt() string {
	var choice string

	err := huh.NewSelect[string]().
		Title("What would you like to do next?").
		Description("Select an option using arrow keys or j,k and press Enter").
		Options(
			huh.NewOption("Copy to clipboard", "copy"),
			huh.NewOption("Save analysis to file", "save"),
			huh.NewOption("Generate different analysis", "regenerate"),
			huh.NewOption("Refine with feedback", "custom"),
			huh.NewOption("Analyze different period", "new"),
			huh.NewOption("Exit", "exit"),
		).
		Value(&choice).
		Height(8).
		WithTheme(ui.GetHuhPrimaryTheme()).
		Run()

	if err != nil {
		ui.Box(ui.BoxOptions{Message: fmt.Sprintf("Error running prompt: %v", err), Variant: ui.Error})
		os.Exit(1)
	}

	return choice
}

func customTimelineInputPrompt(message string) string {
	var input string

	err := huh.NewInput().
		Title(message).
		Description("Provide specific feedback to improve the analysis").
		Placeholder("e.g., focus more on patterns, include statistics, be more detailed...").
		CharLimit(200).
		Value(&input).
		WithTheme(ui.GetHuhPrimaryTheme()).
		Run()

	if err != nil {
		ui.Box(ui.BoxOptions{Message: fmt.Sprintf("Error running prompt: %v", err), Variant: ui.Error})
		os.Exit(1)
	}

	return input
}

func saveTimelineAnalysis(analysis, dateRange string) (string, error) {
	repoRoot, err := git.FindGitRoot()
	if err != nil {
		return "", fmt.Errorf("failed to find git repository: %v", err)
	}

	timelineDir := filepath.Join(repoRoot, ".git", "diny", "timeline")
	if err := os.MkdirAll(timelineDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create timeline directory: %v", err)
	}

	timestamp := time.Now().Format("2006-01-02-150405")
	sanitizedRange := strings.ReplaceAll(dateRange, " ", "-")
	sanitizedRange = strings.ReplaceAll(sanitizedRange, ":", "-")
	fileName := fmt.Sprintf("diny-timeline-%s-%s.md", sanitizedRange, timestamp)
	filePath := filepath.Join(timelineDir, fileName)

	content := fmt.Sprintf("# Timeline Analysis: %s\n\nGenerated: %s\n\n%s\n", dateRange, time.Now().Format("2006-01-02 15:04:05"), analysis)

	if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
		return "", fmt.Errorf("failed to write analysis file: %v", err)
	}

	return filePath, nil
}
