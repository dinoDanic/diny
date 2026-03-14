package changelog

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
		return repoInfoMsg{repoName: repoName, branchName: branchName}
	}
}

func loadRefs(mode string) tea.Cmd {
	return func() tea.Msg {
		if mode == "tag" {
			tags, err := git.GetTags()
			if err != nil {
				return errMsg{err: fmt.Errorf("failed to load tags: %w", err)}
			}
			return refsLoadedMsg{tags: tags}
		}
		commits, err := git.GetRecentCommits(50)
		if err != nil {
			return errMsg{err: fmt.Errorf("failed to load commits: %w", err)}
		}
		return refsLoadedMsg{commits: commits}
	}
}

func doGenerate(olderRef, newerRef string, cfg *config.Config) tea.Cmd {
	return func() tea.Msg {
		commits, err := git.GetCommitsBetweenRefs(olderRef, newerRef)
		if err != nil {
			return errMsg{err: fmt.Errorf("failed to get commits: %w", err)}
		}
		if len(commits) == 0 {
			return noCommitsMsg{}
		}

		diff, err := git.GetDiffBetweenRefs(olderRef, newerRef)
		if err != nil {
			return errMsg{err: fmt.Errorf("failed to get diff: %w", err)}
		}

		repoName := git.GetRepoName()
		gitName := git.GetGitName()
		prompt := buildChangelogPrompt(repoName, gitName, olderRef, newerRef, commits, diff)

		result, err := groq.CreateChangelogWithGroq(prompt, cfg)
		if err != nil {
			return errMsg{err: fmt.Errorf("failed to generate changelog: %w", err)}
		}
		return changelogReadyMsg{result: result, prompt: prompt}
	}
}

func doRegenerate(prompt string, cfg *config.Config, previousResults []string) tea.Cmd {
	return func() tea.Msg {
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

		result, err := groq.CreateChangelogWithGroq(modifiedPrompt, cfg)
		if err != nil {
			return errMsg{err: fmt.Errorf("failed to regenerate changelog: %w", err)}
		}
		return changelogReadyMsg{result: result, prompt: prompt}
	}
}

func doCopy(result string) tea.Cmd {
	return func() tea.Msg {
		if err := clipboard.WriteAll(result); err != nil {
			return errMsg{err: fmt.Errorf("failed to copy to clipboard: %w", err)}
		}
		return copiedMsg{}
	}
}

func doSave(result, rangeLabel string) tea.Cmd {
	return func() tea.Msg {
		filePath, err := saveChangelog(result, rangeLabel)
		if err != nil {
			return errMsg{err: err}
		}
		return savedMsg{filePath: filePath}
	}
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
		repoName, gitName, olderRef, newerRef,
		len(commits), strings.Join(commitLines, "\n"),
		diffSummary,
	)
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
