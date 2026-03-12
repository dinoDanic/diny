package app

import (
	"fmt"
	"time"

	"github.com/atotto/clipboard"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/dinoDanic/diny/commit"
	"github.com/dinoDanic/diny/config"
	"github.com/dinoDanic/diny/git"
)

func loadRepoInfo() tea.Cmd {
	return func() tea.Msg {
		repoName := git.GetRepoName()
		branchName, _ := git.GetCurrentBranch()
		gitUserName := git.GetGitName()
		return repoInfoMsg{
			repoName:    repoName,
			branchName:  branchName,
			gitUserName: gitUserName,
		}
	}
}

func loadStagedFiles() tea.Cmd {
	return func() tea.Msg {
		files, err := git.GetStagedFiles()
		if err != nil {
			return errMsg{err: fmt.Errorf("failed to get staged files: %w", err)}
		}
		return stagedFilesMsg{files: files}
	}
}

func startWelcomeTimer() tea.Cmd {
	return tea.Tick(800*time.Millisecond, func(t time.Time) tea.Msg {
		return welcomeTimerDoneMsg{}
	})
}

func loadDiffAndGenerate(cfg *config.Config) tea.Cmd {
	return func() tea.Msg {
		diff, err := git.GetGitDiff()
		if err != nil {
			return errMsg{err: fmt.Errorf("failed to get git diff: %w", err)}
		}
		if diff == "" {
			return errMsg{err: fmt.Errorf("no diff found for staged changes")}
		}

		msg, err := commit.CreateCommitMessage(diff, cfg)
		if err != nil {
			return errMsg{err: fmt.Errorf("failed to generate commit message: %w", err)}
		}
		return diffAndCommitMsg{diff: diff, commitMessage: msg}
	}
}

func doCommit(message string, push bool, noVerify bool, cfg *config.Config) tea.Cmd {
	return func() tea.Msg {
		hash, err := commit.TryCommit(message, push, noVerify, cfg)
		if err != nil {
			return errMsg{err: err}
		}
		return commitDoneMsg{hash: hash, push: push}
	}
}

func doRegenerate(diff string, cfg *config.Config, previousMessages []string, current string) tea.Cmd {
	return func() tea.Msg {
		modifiedDiff := diff
		allPrev := append(previousMessages, current)
		if len(allPrev) > 0 {
			modifiedDiff += "\n\nPrevious commit messages that were not satisfactory:\n"
			for i, msg := range allPrev {
				modifiedDiff += fmt.Sprintf("%d. %s\n", i+1, msg)
			}
			modifiedDiff += "\nPlease generate a different commit message that avoids the style and approach of the previous ones."
		}

		msg, err := commit.CreateCommitMessage(modifiedDiff, cfg)
		if err != nil {
			return errMsg{err: fmt.Errorf("failed to regenerate: %w", err)}
		}
		return diffAndCommitMsg{diff: diff, commitMessage: msg}
	}
}

func doFeedback(diff string, cfg *config.Config, current string, feedback string) tea.Cmd {
	return func() tea.Msg {
		modifiedDiff := diff + fmt.Sprintf("\n\nCurrent commit message:\n%s\n\nUser feedback: %s\n\nPlease generate a new commit message that addresses the user's feedback.", current, feedback)

		msg, err := commit.CreateCommitMessage(modifiedDiff, cfg)
		if err != nil {
			return errMsg{err: fmt.Errorf("failed to refine: %w", err)}
		}
		return diffAndCommitMsg{diff: diff, commitMessage: msg}
	}
}

func doSaveDraft(message string) tea.Cmd {
	return func() tea.Msg {
		if err := commit.SaveDraft(message); err != nil {
			return errMsg{err: fmt.Errorf("failed to save draft: %w", err)}
		}
		return draftSavedMsg{}
	}
}

func doCopy(message string) tea.Cmd {
	return func() tea.Msg {
		if err := clipboard.WriteAll(message); err != nil {
			return errMsg{err: fmt.Errorf("failed to copy: %w", err)}
		}
		return copiedMsg{}
	}
}
