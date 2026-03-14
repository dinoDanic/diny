package app

import (
	"bytes"
	"fmt"
	"os/exec"
	"regexp"
	"strings"
	"time"

	"github.com/atotto/clipboard"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/dinoDanic/diny/commit"
	"github.com/dinoDanic/diny/config"
	"github.com/dinoDanic/diny/git"
)

type commitProgressMsg struct {
	line string
}

type lineWriter struct {
	ch  chan<- string
	buf []byte
	all []byte
}

func (w *lineWriter) Write(p []byte) (int, error) {
	w.all = append(w.all, p...)
	w.buf = append(w.buf, p...)
	for {
		i := bytes.IndexByte(w.buf, '\n')
		if i < 0 {
			break
		}
		line := strings.TrimSpace(stripAnsi(string(w.buf[:i])))
		w.buf = w.buf[i+1:]
		if line != "" {
			select {
			case w.ch <- line:
			default:
			}
		}
	}
	return len(p), nil
}

var ansiEscape = regexp.MustCompile(`\x1b\[[0-9;]*[mKGJ]`)

func stripAnsi(s string) string {
	return ansiEscape.ReplaceAllString(s, "")
}

func waitForCommitLine(ch <-chan string) tea.Cmd {
	return func() tea.Msg {
		line, ok := <-ch
		if !ok {
			return nil
		}
		return commitProgressMsg{line: line}
	}
}

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

func doCommit(message string, push bool, noVerify bool, amend bool, cfg *config.Config, progressCh chan string) tea.Cmd {
	return func() tea.Msg {
		defer close(progressCh)

		// Build args
		var args []string
		if amend {
			args = []string{"commit", "--amend", "-m", message}
		} else {
			args = []string{"commit", "-m", message}
		}
		if noVerify {
			args = append([]string{args[0], "--no-verify"}, args[1:]...)
		}

		lw := &lineWriter{ch: progressCh}
		cmd := exec.Command("git", args...)
		cmd.Stdout = lw
		cmd.Stderr = lw
		if err := cmd.Run(); err != nil {
			return errMsg{err: fmt.Errorf("commit failed: %s", strings.TrimSpace(string(lw.all)))}
		}

		if amend {
			return commitDoneMsg{hash: "", push: false}
		}

		// hash after commit
		var hash string
		if cfg != nil && cfg.Commit.HashAfterCommit {
			if out, err := exec.Command("git", "rev-parse", "--short", "HEAD").Output(); err == nil {
				hash = strings.TrimSpace(string(out))
				_ = clipboard.WriteAll(hash)
			}
		}

		// optional push
		if push {
			pushLw := &lineWriter{ch: progressCh}
			pushCmd := exec.Command("git", "push")
			pushCmd.Stdout = pushLw
			pushCmd.Stderr = pushLw
			if err := pushCmd.Run(); err != nil {
				return errMsg{err: fmt.Errorf("committed but push failed: %s", strings.TrimSpace(string(pushLw.all)))}
			}
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

func loadUnstagedFiles() tea.Cmd {
	return func() tea.Msg {
		files, err := git.GetUnstagedFiles()
		if err != nil {
			return errMsg{err: fmt.Errorf("failed to get unstaged files: %w", err)}
		}
		return unstagedFilesMsg{files: files}
	}
}

func doStageFiles(paths []string) tea.Cmd {
	return func() tea.Msg {
		args := append([]string{"add", "--"}, paths...)
		cmd := exec.Command("git", args...)
		if out, err := cmd.CombinedOutput(); err != nil {
			return errMsg{err: fmt.Errorf("git add failed: %s", strings.TrimSpace(string(out)))}
		}
		files, err := git.GetStagedFiles()
		if err != nil {
			return errMsg{err: fmt.Errorf("failed to get staged files: %w", err)}
		}
		return stagedFilesMsg{files: files}
	}
}

func doGenerateVariants(diff string, cfg *config.Config, previousMessages []string, current string) tea.Cmd {
	return func() tea.Msg {
		const n = 3
		type result struct {
			msg string
			err error
		}

		modifiedDiff := diff
		allPrev := append(previousMessages, current)
		modifiedDiff += "\n\nPrevious commit messages that were not satisfactory:\n"
		for i, m := range allPrev {
			modifiedDiff += fmt.Sprintf("%d. %s\n", i+1, m)
		}
		modifiedDiff += "\nPlease generate a different commit message that avoids the style and approach of the previous ones."

		ch := make(chan result, n)
		for range n {
			go func() {
				msg, err := commit.CreateCommitMessage(modifiedDiff, cfg)
				ch <- result{msg, err}
			}()
		}

		var variants []string
		var lastErr error
		for range n {
			r := <-ch
			if r.err != nil {
				lastErr = r.err
			} else {
				variants = append(variants, r.msg)
			}
		}
		if len(variants) == 0 {
			return errMsg{err: fmt.Errorf("all variants failed: %w", lastErr)}
		}
		for len(variants) < n {
			variants = append(variants, variants[0])
		}
		return variantsReadyMsg{variants: variants}
	}
}

func loadAllFiles() tea.Cmd {
	return func() tea.Msg {
		staged, err := git.GetStagedFiles()
		if err != nil {
			return errMsg{err: fmt.Errorf("failed to get staged files: %w", err)}
		}
		unstaged, err := git.GetUnstagedFiles()
		if err != nil {
			return errMsg{err: fmt.Errorf("failed to get unstaged files: %w", err)}
		}
		var entries []fileEntry
		for _, f := range staged {
			entries = append(entries, fileEntry{path: f.Path, status: f.Status, currentStaged: true, wantStaged: true})
		}
		for _, f := range unstaged {
			entries = append(entries, fileEntry{path: f.Path, status: f.Status, currentStaged: false, wantStaged: false})
		}
		return allFilesMsg{entries: entries}
	}
}

func doApplyFilePicker(entries []fileEntry) tea.Cmd {
	return func() tea.Msg {
		for _, e := range entries {
			if e.wantStaged == e.currentStaged {
				continue
			}
			var cmd *exec.Cmd
			if e.wantStaged {
				cmd = exec.Command("git", "add", "--", e.path)
			} else {
				cmd = exec.Command("git", "restore", "--staged", "--", e.path)
			}
			if out, err := cmd.CombinedOutput(); err != nil {
				return errMsg{err: fmt.Errorf("git operation failed: %s", strings.TrimSpace(string(out)))}
			}
		}
		files, err := git.GetStagedFiles()
		if err != nil {
			return errMsg{err: fmt.Errorf("failed to get staged files: %w", err)}
		}
		return filePickerDoneMsg{files: files}
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
