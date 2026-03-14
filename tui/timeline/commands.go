package timeline

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/atotto/clipboard"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/dinoDanic/diny/config"
	"github.com/dinoDanic/diny/git"
	"github.com/dinoDanic/diny/groq"
)

func loadRepoInfo() tea.Cmd {
	return func() tea.Msg {
		repoName := git.GetRepoName()
		branchName, _ := git.GetCurrentBranch()
		return repoInfoMsg{
			repoName:   repoName,
			branchName: branchName,
		}
	}
}

func fetchAndGenerate(dateChoice, startDate, endDate, dateRange string, cfg *config.Config) tea.Cmd {
	return func() tea.Msg {
		var commits []string
		var err error

		switch dateChoice {
		case "today":
			commits, err = git.GetCommitsToday()
		case "date":
			commits, err = git.GetCommitsByDate(startDate)
		case "range":
			commits, err = git.GetCommitsByDateRange(startDate+" 00:00:00", endDate+" 23:59:59")
		}

		if err != nil {
			return errMsg{err: fmt.Errorf("failed to get commits: %w", err)}
		}

		if len(commits) == 0 {
			return noCommitsMsg{}
		}

		prompt := fmt.Sprintf("Timeline: %s\nCommits:\n%s", dateRange, strings.Join(commits, "\n"))
		analysis, err := groq.CreateTimelineWithGroq(prompt, cfg)
		if err != nil {
			return errMsg{err: fmt.Errorf("failed to generate analysis: %w", err)}
		}

		return analysisReadyMsg{
			commits:    commits,
			analysis:   analysis,
			fullPrompt: prompt,
		}
	}
}

func doRegenerate(fullPrompt string, cfg *config.Config, previousAnalyses []string) tea.Cmd {
	return func() tea.Msg {
		modifiedPrompt := fullPrompt
		if len(previousAnalyses) > 0 {
			modifiedPrompt += "\n\nPrevious analyses that were not satisfactory:\n"
			for i, a := range previousAnalyses {
				modifiedPrompt += fmt.Sprintf("%d. %s\n", i+1, a)
			}
			modifiedPrompt += "\nPlease generate a different analysis with a different perspective or focus."
		} else {
			modifiedPrompt += "\n\nPlease provide an alternative analysis with a different approach or focus."
		}

		analysis, err := groq.CreateTimelineWithGroq(modifiedPrompt, cfg)
		if err != nil {
			return errMsg{err: fmt.Errorf("failed to regenerate analysis: %w", err)}
		}

		return analysisReadyMsg{
			commits:    nil,
			analysis:   analysis,
			fullPrompt: fullPrompt,
		}
	}
}

func doFeedback(fullPrompt, current, feedback string, cfg *config.Config) tea.Cmd {
	return func() tea.Msg {
		modifiedPrompt := fullPrompt + fmt.Sprintf(
			"\n\nCurrent analysis:\n%s\n\nUser feedback: %s\n\nPlease generate a new analysis that addresses the user's feedback.",
			current, feedback,
		)

		analysis, err := groq.CreateTimelineWithGroq(modifiedPrompt, cfg)
		if err != nil {
			return errMsg{err: fmt.Errorf("failed to refine analysis: %w", err)}
		}

		return analysisReadyMsg{
			commits:    nil,
			analysis:   analysis,
			fullPrompt: fullPrompt,
		}
	}
}

func doCopy(analysis string) tea.Cmd {
	return func() tea.Msg {
		if err := clipboard.WriteAll(analysis); err != nil {
			return errMsg{err: fmt.Errorf("failed to copy to clipboard: %w", err)}
		}
		return copiedMsg{}
	}
}

func doSave(analysis, dateRange string) tea.Cmd {
	return func() tea.Msg {
		filePath, err := saveTimelineAnalysis(analysis, dateRange)
		if err != nil {
			return errMsg{err: err}
		}
		return savedMsg{filePath: filePath}
	}
}

func saveTimelineAnalysis(analysis, dateRange string) (string, error) {
	gitDir, err := git.FindGitDir()
	if err != nil {
		return "", fmt.Errorf("failed to find git repository: %v", err)
	}

	timelineDir := filepath.Join(gitDir, "diny", "timeline")
	if err := os.MkdirAll(timelineDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create timeline directory: %v", err)
	}

	timestamp := time.Now().Format("2006-01-02-150405")
	sanitizedRange := strings.ReplaceAll(dateRange, " ", "-")
	sanitizedRange = strings.ReplaceAll(sanitizedRange, ":", "-")
	fileName := fmt.Sprintf("diny-timeline-%s-%s.md", sanitizedRange, timestamp)
	filePath := filepath.Join(timelineDir, fileName)

	content := fmt.Sprintf("# Timeline Analysis: %s\n\nGenerated: %s\n\n%s\n",
		dateRange, time.Now().Format("2006-01-02 15:04:05"), analysis)

	if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
		return "", fmt.Errorf("failed to write analysis file: %v", err)
	}

	return filePath, nil
}
