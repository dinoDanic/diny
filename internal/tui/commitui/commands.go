package commitui

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"sync"

	"github.com/atotto/clipboard"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/dinoDanic/diny/commit"
	"github.com/dinoDanic/diny/config"
	"github.com/dinoDanic/diny/git"
)

func loadStagedFiles() tea.Msg {
	files, err := git.GetStagedFiles()
	if err != nil {
		return errorMsg{err: fmt.Errorf("failed to get staged files: %w", err)}
	}
	return stagedFilesMsg{files: files}
}

func loadDiffAndGenerate(cfg *config.Config) tea.Cmd {
	return func() tea.Msg {
		diff, err := git.GetGitDiff()
		if err != nil {
			return errorMsg{err: fmt.Errorf("failed to get git diff: %w", err)}
		}
		if len(diff) == 0 {
			return errorMsg{err: fmt.Errorf("no staged changes found")}
		}

		msg, err := commit.CreateCommitMessage(diff, cfg)
		if err != nil {
			return errorMsg{err: fmt.Errorf("failed to generate commit message: %w", err)}
		}

		return commitMessageMsg{message: msg, diff: diff}
	}
}

func doCommit(msg string, push, noVerify bool, cfg *config.Config) tea.Cmd {
	return func() tea.Msg {
		err := commit.TryCommit(msg, push, noVerify, cfg)
		if err != nil {
			return errorMsg{err: err}
		}
		return commitDoneMsg{pushed: push}
	}
}

func doRegenerate(diff string, cfg *config.Config, previousMessages []string, currentMessage string) tea.Cmd {
	return func() tea.Msg {
		modifiedPrompt := diff
		allPrevious := append(previousMessages, currentMessage)
		if len(allPrevious) > 0 {
			modifiedPrompt += "\n\nPrevious commit messages that were not satisfactory:\n"
			for i, msg := range allPrevious {
				modifiedPrompt += fmt.Sprintf("%d. %s\n", i+1, msg)
			}
			modifiedPrompt += "\nPlease generate a different commit message that avoids the style and approach of the previous ones."
		}

		msg, err := commit.CreateCommitMessage(modifiedPrompt, cfg)
		if err != nil {
			return errorMsg{err: fmt.Errorf("failed to regenerate: %w", err)}
		}
		return regenerateMsg{message: msg}
	}
}

func doFeedback(diff string, cfg *config.Config, currentMessage, feedback string) tea.Cmd {
	return func() tea.Msg {
		modifiedPrompt := diff + fmt.Sprintf("\n\nCurrent commit message:\n%s\n\nUser feedback: %s\n\nPlease generate a new commit message that addresses the user's feedback.", currentMessage, feedback)

		msg, err := commit.CreateCommitMessage(modifiedPrompt, cfg)
		if err != nil {
			return errorMsg{err: fmt.Errorf("failed to refine: %w", err)}
		}
		return feedbackMsg{message: msg}
	}
}

func doOpenEditor(message string) tea.Cmd {
	editor := git.GetGitEditor()

	tmpFile, err := os.CreateTemp("", "diny-commit-*.txt")
	if err != nil {
		return func() tea.Msg {
			return errorMsg{err: fmt.Errorf("failed to create temp file: %w", err)}
		}
	}

	if _, err := tmpFile.WriteString(message); err != nil {
		os.Remove(tmpFile.Name())
		return func() tea.Msg {
			return errorMsg{err: fmt.Errorf("failed to write temp file: %w", err)}
		}
	}
	tmpFile.Close()

	editorArgs := strings.Fields(editor)
	editorCmd := editorArgs[0]
	args := append(editorArgs[1:], tmpFile.Name())

	c := exec.Command(editorCmd, args...)
	return tea.ExecProcess(c, func(err error) tea.Msg {
		defer os.Remove(tmpFile.Name())
		if err != nil {
			return errorMsg{err: fmt.Errorf("editor exited with error: %w", err)}
		}
		content, readErr := os.ReadFile(tmpFile.Name())
		if readErr != nil {
			return errorMsg{err: fmt.Errorf("failed to read edited file: %w", readErr)}
		}
		return editorMsg{message: strings.TrimSpace(string(content))}
	})
}

func doSaveDraft(message string) tea.Cmd {
	return func() tea.Msg {
		if err := commit.SaveDraft(message); err != nil {
			return errorMsg{err: fmt.Errorf("failed to save draft: %w", err)}
		}
		return draftMsg{}
	}
}

func doCopy(message string) tea.Cmd {
	return func() tea.Msg {
		if err := clipboard.WriteAll(message); err != nil {
			return errorMsg{err: fmt.Errorf("failed to copy: %w", err)}
		}
		return copyMsg{}
	}
}

func doGenerateVariants(diff string, cfg *config.Config, currentMessage string) tea.Cmd {
	return func() tea.Msg {
		var wg sync.WaitGroup
		results := make([]string, 3)
		errs := make([]error, 3)

		for i := 0; i < 3; i++ {
			wg.Add(1)
			go func(idx int) {
				defer wg.Done()
				prompt := diff + fmt.Sprintf("\n\nGenerate a unique alternative (variant %d of 3). Avoid this message: %s\nBe creative with a different angle or focus.", idx+1, currentMessage)
				msg, err := commit.CreateCommitMessage(prompt, cfg)
				results[idx] = msg
				errs[idx] = err
			}(i)
		}

		wg.Wait()

		var variants []string
		for i, r := range results {
			if errs[i] == nil && r != "" {
				variants = append(variants, r)
			}
		}

		if len(variants) == 0 {
			return errorMsg{err: fmt.Errorf("failed to generate any variants")}
		}

		return variantsMsg{variants: variants}
	}
}
