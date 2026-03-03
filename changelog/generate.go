package changelog

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

func GenerateByTag(cfg *config.Config) error {
	tags, err := git.GetTags()
	if err != nil {
		ui.Box(ui.BoxOptions{Message: fmt.Sprintf("Failed to list tags: %v", err), Variant: ui.Error})
		return err
	}

	if len(tags) < 2 {
		ui.Box(ui.BoxOptions{
			Message: "At least two tags are required to generate a changelog. No tags found or only one tag exists.",
			Variant: ui.Warning,
		})
		return nil
	}

	newerTag, err := selectTagPrompt("Select newer tag (to)", tags)
	if err != nil {
		return nil // user exited
	}

	newerIdx := indexOf(tags, newerTag)
	olderTags := tags[newerIdx+1:]
	if len(olderTags) == 0 {
		ui.Box(ui.BoxOptions{
			Message: fmt.Sprintf("No older tags available before %s.", newerTag),
			Variant: ui.Warning,
		})
		return nil
	}

	olderTag, err := selectTagPrompt("Select older tag (from)", olderTags)
	if err != nil {
		return nil
	}

	return generateChangelog(olderTag, newerTag, cfg)
}

func GenerateByCommit(cfg *config.Config) error {
	commits, err := git.GetRecentCommits(50)
	if err != nil {
		ui.Box(ui.BoxOptions{Message: fmt.Sprintf("Failed to get recent commits: %v", err), Variant: ui.Error})
		return err
	}

	if len(commits) < 2 {
		ui.Box(ui.BoxOptions{
			Message: "At least two commits are required to generate a changelog.",
			Variant: ui.Warning,
		})
		return nil
	}

	options := make([]huh.Option[string], len(commits))
	for i, c := range commits {
		label := fmt.Sprintf("%s  %s", c.SHA, c.Message)
		options[i] = huh.NewOption(label, c.SHA)
	}

	newerSHA, err := selectRefFromOptions("Select newer commit (to)", options)
	if err != nil {
		return nil
	}

	olderSHA, err := selectRefFromOptions("Select older commit (from)", options)
	if err != nil {
		return nil
	}

	return generateChangelog(olderSHA, newerSHA, cfg)
}

func generateChangelog(olderRef, newerRef string, cfg *config.Config) error {
	commits, err := git.GetCommitsBetweenRefs(olderRef, newerRef)
	if err != nil {
		ui.Box(ui.BoxOptions{Message: fmt.Sprintf("Failed to get commits: %v", err), Variant: ui.Error})
		return err
	}

	rangeLabel := fmt.Sprintf("%s → %s", olderRef, newerRef)

	if len(commits) == 0 {
		ui.Box(ui.BoxOptions{
			Message: fmt.Sprintf("No commits found between %s.", rangeLabel),
			Variant: ui.Warning,
		})
		return nil
	}

	diff, err := git.GetDiffBetweenRefs(olderRef, newerRef)
	if err != nil {
		ui.Box(ui.BoxOptions{Message: fmt.Sprintf("Failed to get diff: %v", err), Variant: ui.Error})
		return err
	}

	repoName := git.GetRepoName()
	gitName := git.GetGitName()

	prompt := buildChangelogPrompt(repoName, gitName, olderRef, newerRef, commits, diff)

	var result string
	spinErr := ui.WithSpinner("Generating changelog...", func() error {
		var genErr error
		result, genErr = groq.CreateChangelogWithGroq(prompt, cfg)
		return genErr
	})

	if spinErr != nil {
		ui.Box(ui.BoxOptions{Message: fmt.Sprintf("Failed to generate changelog: %v", spinErr), Variant: ui.Error})
		return spinErr
	}

	ui.Box(ui.BoxOptions{Title: fmt.Sprintf("Changelog: %s", rangeLabel), Message: result, Variant: ui.Primary})

	handleGenerateFlow(result, rangeLabel, olderRef, newerRef, prompt, cfg, []string{})
	return nil
}

func buildChangelogPrompt(repoName, gitName, olderRef, newerRef string, commits []string, diff string) string {
	commitLines := make([]string, len(commits))
	for i, c := range commits {
		commitLines[i] = "- " + c
	}

	diffSummary := diff
	if len(diffSummary) > 4000 {
		diffSummary = diffSummary[:4000] + "\n... (diff truncated)"
	}

	return fmt.Sprintf(`Generate a changelog for the following repository changes.

Repository: %s
Author: %s
Range: %s → %s

Commits (%d):
%s

Diff Summary:
%s

Generate a well-structured markdown changelog with sections:
## What's Changed
Use sub-bullets for individual changes. Group logically if possible.
## Bug Fixes (if any)
## Breaking Changes (if any)
Keep it concise and human-readable.`,
		repoName,
		gitName,
		olderRef,
		newerRef,
		len(commits),
		strings.Join(commitLines, "\n"),
		diffSummary,
	)
}

func handleGenerateFlow(result, rangeLabel, olderRef, newerRef, prompt string, cfg *config.Config, previousResults []string) {
	var choice string

	err := huh.NewSelect[string]().
		Title("What would you like to do next?").
		Description("Select an option using arrow keys or j,k and press Enter").
		Options(
			huh.NewOption("Copy to clipboard", "copy"),
			huh.NewOption("Save as CHANGELOG.md", "save"),
			huh.NewOption("Regenerate", "regenerate"),
			huh.NewOption("Exit", "exit"),
		).
		Value(&choice).
		Height(6).
		WithTheme(ui.GetHuhPrimaryTheme()).
		Run()

	if err != nil {
		return
	}

	switch choice {
	case "copy":
		if err := clipboard.WriteAll(result); err != nil {
			ui.Box(ui.BoxOptions{Message: fmt.Sprintf("Failed to copy to clipboard: %v", err), Variant: ui.Error})
		} else {
			ui.Box(ui.BoxOptions{Message: "Changelog copied to clipboard!", Variant: ui.Success})
		}

	case "save":
		filePath, saveErr := saveChangelog(result, rangeLabel)
		if saveErr != nil {
			ui.Box(ui.BoxOptions{Message: fmt.Sprintf("Failed to save changelog: %v", saveErr), Variant: ui.Error})
			handleGenerateFlow(result, rangeLabel, olderRef, newerRef, prompt, cfg, previousResults)
			return
		}
		ui.Box(ui.BoxOptions{Message: fmt.Sprintf("Changelog saved!\n\n%s", filePath), Variant: ui.Success})

	case "regenerate":
		modifiedPrompt := prompt
		if len(previousResults) > 0 {
			modifiedPrompt += "\n\nPrevious changelogs that were not satisfactory:\n"
			for i, prev := range previousResults {
				modifiedPrompt += fmt.Sprintf("%d. %s\n", i+1, prev)
			}
			modifiedPrompt += "\nPlease generate a different changelog with a fresh perspective."
		} else {
			modifiedPrompt += "\n\nPlease provide an alternative changelog with a different approach."
		}

		var newResult string
		spinErr := ui.WithSpinner("Regenerating changelog...", func() error {
			var genErr error
			newResult, genErr = groq.CreateChangelogWithGroq(modifiedPrompt, cfg)
			return genErr
		})
		if spinErr != nil {
			ui.Box(ui.BoxOptions{Message: fmt.Sprintf("Error: %v", spinErr), Variant: ui.Error})
			return
		}

		ui.Box(ui.BoxOptions{Title: fmt.Sprintf("Changelog: %s", rangeLabel), Message: newResult, Variant: ui.Primary})
		updated := append(previousResults, result)
		handleGenerateFlow(newResult, rangeLabel, olderRef, newerRef, prompt, cfg, updated)

	case "exit":
		return
	}
}

func saveChangelog(content, rangeLabel string) (string, error) {
	gitDir, err := git.FindGitDir()
	if err != nil {
		return "", fmt.Errorf("failed to find git repository: %v", err)
	}

	changelogDir := filepath.Join(gitDir, "diny", "changelog")
	if err := os.MkdirAll(changelogDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create changelog directory: %v", err)
	}

	timestamp := time.Now().Format("2006-01-02-150405")
	sanitized := strings.NewReplacer(" ", "-", ":", "-", "→", "-", "/", "-").Replace(rangeLabel)
	fileName := fmt.Sprintf("changelog-%s-%s.md", sanitized, timestamp)
	filePath := filepath.Join(changelogDir, fileName)

	header := fmt.Sprintf("# Changelog: %s\n\nGenerated: %s\n\n", rangeLabel, time.Now().Format("2006-01-02 15:04:05"))
	if err := os.WriteFile(filePath, []byte(header+content+"\n"), 0644); err != nil {
		return "", fmt.Errorf("failed to write changelog file: %v", err)
	}

	return filePath, nil
}

func selectTagPrompt(title string, tags []string) (string, error) {
	options := make([]huh.Option[string], len(tags))
	for i, t := range tags {
		options[i] = huh.NewOption(t, t)
	}

	var selected string
	err := huh.NewSelect[string]().
		Title(title).
		Description("Select using arrow keys or j,k and press Enter").
		Options(options...).
		Value(&selected).
		Height(10).
		WithTheme(ui.GetHuhPrimaryTheme()).
		Run()

	if err != nil {
		return "", err
	}
	return selected, nil
}

func selectRefFromOptions(title string, options []huh.Option[string]) (string, error) {
	var selected string
	err := huh.NewSelect[string]().
		Title(title).
		Description("Select using arrow keys or j,k and press Enter").
		Options(options...).
		Value(&selected).
		Height(10).
		WithTheme(ui.GetHuhPrimaryTheme()).
		Run()

	if err != nil {
		return "", err
	}
	return selected, nil
}

func indexOf(slice []string, target string) int {
	for i, s := range slice {
		if s == target {
			return i
		}
	}
	return -1
}
