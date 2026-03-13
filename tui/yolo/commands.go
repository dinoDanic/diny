package yolo

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/dinoDanic/diny/commit"
	"github.com/dinoDanic/diny/config"
	"github.com/dinoDanic/diny/git"
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

func doStageAll() tea.Cmd {
	return func() tea.Msg {
		if err := git.AddAll(); err != nil {
			return errMsg{err: fmt.Errorf("failed to stage changes: %w", err)}
		}
		files, err := git.GetStagedFiles()
		if err != nil {
			return errMsg{err: fmt.Errorf("failed to get staged files: %w", err)}
		}
		return stageDoneMsg{files: files}
	}
}

func loadDiffAndGenerate(cfg *config.Config) tea.Cmd {
	return func() tea.Msg {
		diff, err := git.GetGitDiff()
		if err != nil {
			return errMsg{err: fmt.Errorf("failed to get git diff: %w", err)}
		}
		if diff == "" {
			return nothingToCommitMsg{}
		}

		msg, err := commit.CreateCommitMessage(diff, cfg)
		if err != nil {
			return errMsg{err: fmt.Errorf("failed to generate commit message: %w", err)}
		}
		return generateDoneMsg{commitMessage: msg}
	}
}

func doCommitAndPush(message string, cfg *config.Config) tea.Cmd {
	return func() tea.Msg {
		hash, err := commit.TryCommit(message, true, true, cfg)
		if err != nil {
			return errMsg{err: err}
		}
		return commitDoneMsg{hash: hash}
	}
}
