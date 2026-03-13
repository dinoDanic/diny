package app

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/dinoDanic/diny/config"
	"github.com/dinoDanic/diny/git"
	"github.com/dinoDanic/diny/tui/loader"
)

var conventionalTypes = []string{"feat", "fix", "docs", "style", "refactor", "perf", "test", "chore"}

func (m model) Init() tea.Cmd {
	return tea.Batch(
		m.loader.Tick,
		loadRepoInfo(),
		loadStagedFiles(),
		startWelcomeTimer(),
	)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case tea.KeyMsg:
		return m.handleKey(msg)

	case repoInfoMsg:
		m.repoName = msg.repoName
		m.branchName = msg.branchName
		m.gitUserName = msg.gitUserName
		m.repoLoaded = true
		return m.checkWelcomeDone()

	case stagedFilesMsg:
		m.stagedFiles = msg.files
		m.filesLoaded = true
		return m.checkWelcomeDone()

	case unstagedFilesMsg:
		m.unstagedFiles = msg.files
		m.fileSelected = make([]bool, len(msg.files))
		m.fileCursor = 0
		return m, nil

	case welcomeTimerDoneMsg:
		m.welcomeReady = true
		return m.checkWelcomeDone()

	case diffAndCommitMsg:
		m.diff = msg.diff
		m.commitMessage = msg.commitMessage
		m.messageHistoryIdx = -1
		m.savedMessage = ""
		m.state = stateReady
		m.statusMessage = ""
		m.statusIsError = false
		return m, nil

	case restoreStagedDoneMsg:
		m.stagedFiles = msg.files
		if len(m.stagedFiles) == 0 {
			m.state = stateNoStaged
			return m, loadUnstagedFiles()
		}
		m.state = stateReady
		return m, nil

	case commitDoneMsg:
		m.state = stateSuccess
		m.statusMessage = "Committed!"
		if msg.hash != "" {
			m.statusMessage += " (" + msg.hash + ")"
		}
		if msg.push {
			m.statusMessage += " Pushed!"
		}
		return m, tea.Quit

	case draftSavedMsg:
		m.statusMessage = "Draft saved!"
		m.statusIsError = false
		return m, nil

	case copiedMsg:
		m.statusMessage = "Copied to clipboard!"
		m.statusIsError = false
		return m, nil

	case errMsg:
		if m.state == stateCommitting {
			m.state = stateReady
			m.statusMessage = msg.err.Error()
			m.statusIsError = true
			return m, nil
		}
		m.state = stateError
		m.err = msg.err
		return m, nil

	case editorFinishedMsg:
		if msg.newMessage != "" && msg.newMessage != m.commitMessage {
			m.commitMessage = msg.newMessage
		}
		m.state = stateReady
		return m, nil

	case variantsReadyMsg:
		m.variants = msg.variants
		m.variantCursor = 0
		m.state = stateVariantPicking
		return m, nil
	}

	// Update sub-components
	var cmd tea.Cmd
	switch m.state {
	case stateWelcome, stateGenerating, stateCommitting:
		m.loader, cmd = m.loader.Update(msg)
		return m, cmd
	case stateFeedback:
		m.textinput, cmd = m.textinput.Update(msg)
		return m, cmd
	case stateEditing:
		m.textarea, cmd = m.textarea.Update(msg)
		return m, cmd
	case stateDiffView:
		m.viewport, cmd = m.viewport.Update(msg)
		return m, cmd
	}

	return m, nil
}

func (m model) handleKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch m.state {
	case stateReady:
		return m.handleReadyKey(msg)
	case stateFeedback:
		return m.handleFeedbackKey(msg)
	case stateEditing:
		return m.handleEditingKey(msg)
	case stateHelp:
		return m.handleHelpKey(msg)
	case stateNoStaged:
		return m.handleNoStagedKey(msg)
	case stateVariantPicking:
		return m.handleVariantPickingKey(msg)
	case stateDiffView:
		return m.handleDiffViewKey(msg)
	case stateTypePicker:
		return m.handleTypePickerKey(msg)
	case stateUnstaging:
		return m.handleUnstagingKey(msg)
	case stateError, stateSuccess:
		if msg.String() == "q" || msg.String() == "ctrl+c" {
			return m, tea.Quit
		}
	case stateWelcome, stateGenerating, stateCommitting:
		if msg.String() == "ctrl+c" {
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m model) handleReadyKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch {
	case msg.String() == "enter":
		m.state = stateCommitting
		m.loader = loader.New(loader.CommittingMessages)
		return m, doCommit(m.commitMessage, false, false, m.pendingAmend, m.cfg)
	case msg.String() == "n":
		m.state = stateCommitting
		m.loader = loader.New(loader.CommittingMessages)
		return m, doCommit(m.commitMessage, false, true, m.pendingAmend, m.cfg)
	case msg.String() == "p":
		m.state = stateCommitting
		m.loader = loader.New(loader.CommittingMessages)
		return m, doCommit(m.commitMessage, true, false, m.pendingAmend, m.cfg)
	case msg.String() == "r":
		m.state = stateGenerating
		m.loader = loader.New(loader.GeneratingMessages)
		prev := m.previousMessages
		m.previousMessages = append(m.previousMessages, m.commitMessage)
		return m, tea.Batch(m.loader.Tick, doRegenerate(m.diff, m.cfg, prev, m.commitMessage))
	case msg.String() == "v":
		m.state = stateGenerating
		m.loader = loader.New(loader.VariantMessages)
		return m, tea.Batch(m.loader.Tick, doGenerateVariants(m.diff, m.cfg, m.previousMessages, m.commitMessage))
	case msg.String() == "f":
		m.state = stateFeedback
		m.textinput = textinput.New()
		m.textinput.Placeholder = "Describe what to change..."
		m.textinput.CharLimit = 200
		m.textinput.Width = 60
		m.textinput.Focus()
		return m, m.textinput.Cursor.BlinkCmd()
	case msg.String() == "e":
		m.state = stateEditing
		m.textarea = textarea.New()
		m.textarea.SetValue(m.commitMessage)
		m.textarea.SetHeight(8)
		m.textarea.SetWidth(60)
		m.textarea.Focus()
		return m, m.textarea.Cursor.BlinkCmd()
	case msg.String() == "E":
		return m.openExternalEditor()
	case msg.String() == "A":
		m.pendingAmend = true
		m.state = stateGenerating
		m.loader = loader.New(loader.GeneratingMessages)
		return m, tea.Batch(m.loader.Tick, loadDiffAndGenerate(m.cfg))
	case msg.String() == "d":
		vp := viewport.New(m.width-6, m.height-8)
		vp.SetContent(m.diff)
		m.viewport = vp
		m.state = stateDiffView
		return m, nil
	case msg.String() == "[":
		if len(m.previousMessages) == 0 {
			return m, nil
		}
		if m.messageHistoryIdx == -1 {
			m.savedMessage = m.commitMessage
			m.messageHistoryIdx = len(m.previousMessages) - 1
		} else if m.messageHistoryIdx > 0 {
			m.messageHistoryIdx--
		}
		m.commitMessage = m.previousMessages[m.messageHistoryIdx]
		return m, nil
	case msg.String() == "]":
		if m.messageHistoryIdx == -1 {
			return m, nil
		}
		if m.messageHistoryIdx >= len(m.previousMessages)-1 {
			m.messageHistoryIdx = -1
			m.commitMessage = m.savedMessage
		} else {
			m.messageHistoryIdx++
			m.commitMessage = m.previousMessages[m.messageHistoryIdx]
		}
		return m, nil
	case msg.String() == "t":
		if !m.cfg.Commit.Conventional {
			m.statusMessage = "Enable conventional commits in config"
			m.statusIsError = false
			return m, nil
		}
		m.typeCursor = 0
		m.state = stateTypePicker
		return m, nil
	case msg.String() == "L":
		switch m.cfg.Commit.Length {
		case config.Short:
			m.cfg.Commit.Length = config.Normal
		case config.Normal:
			m.cfg.Commit.Length = config.Long
		default:
			m.cfg.Commit.Length = config.Short
		}
		m.statusMessage = fmt.Sprintf("Length: %s", m.cfg.Commit.Length)
		m.statusIsError = false
		m.state = stateGenerating
		m.loader = loader.New(loader.GeneratingMessages)
		prev := m.previousMessages
		m.previousMessages = append(m.previousMessages, m.commitMessage)
		return m, tea.Batch(m.loader.Tick, doRegenerate(m.diff, m.cfg, prev, m.commitMessage))
	case msg.String() == "M":
		m.cfg.Commit.Emoji = !m.cfg.Commit.Emoji
		emojiStatus := "off"
		if m.cfg.Commit.Emoji {
			emojiStatus = "on"
		}
		m.statusMessage = fmt.Sprintf("Emoji: %s", emojiStatus)
		m.statusIsError = false
		m.state = stateGenerating
		m.loader = loader.New(loader.GeneratingMessages)
		prev := m.previousMessages
		m.previousMessages = append(m.previousMessages, m.commitMessage)
		return m, tea.Batch(m.loader.Tick, doRegenerate(m.diff, m.cfg, prev, m.commitMessage))
	case msg.String() == "x":
		m.unstageCursor = 0
		m.unstageSelected = make([]bool, len(m.stagedFiles))
		m.state = stateUnstaging
		return m, nil
	case msg.String() == "s":
		return m, doSaveDraft(m.commitMessage)
	case msg.String() == "y":
		return m, doCopy(m.commitMessage)
	case msg.String() == "?":
		m.state = stateHelp
		return m, nil
	case msg.String() == "q" || msg.String() == "ctrl+c":
		return m, tea.Quit
	}
	return m, nil
}

func (m model) handleFeedbackKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "enter":
		feedback := strings.TrimSpace(m.textinput.Value())
		if feedback == "" {
			m.state = stateReady
			return m, nil
		}
		m.state = stateGenerating
		m.loader = loader.New(loader.GeneratingMessages)
		m.previousMessages = append(m.previousMessages, m.commitMessage)
		return m, tea.Batch(m.loader.Tick, doFeedback(m.diff, m.cfg, m.commitMessage, feedback))
	case "esc":
		m.state = stateReady
		return m, nil
	}

	var cmd tea.Cmd
	m.textinput, cmd = m.textinput.Update(msg)
	return m, cmd
}

func (m model) handleEditingKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc":
		newMsg := strings.TrimSpace(m.textarea.Value())
		if newMsg != "" {
			m.commitMessage = newMsg
		}
		m.state = stateReady
		return m, nil
	}

	var cmd tea.Cmd
	m.textarea, cmd = m.textarea.Update(msg)
	return m, cmd
}

func (m model) handleNoStagedKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "q", "ctrl+c":
		return m, tea.Quit
	case "up", "k":
		if m.fileCursor > 0 {
			m.fileCursor--
		}
	case "down", "j":
		if m.fileCursor < len(m.unstagedFiles)-1 {
			m.fileCursor++
		}
	case " ":
		if len(m.fileSelected) > m.fileCursor {
			m.fileSelected[m.fileCursor] = !m.fileSelected[m.fileCursor]
		}
	case "a":
		allSelected := true
		for _, s := range m.fileSelected {
			if !s {
				allSelected = false
				break
			}
		}
		for i := range m.fileSelected {
			m.fileSelected[i] = !allSelected
		}
	case "enter":
		var paths []string
		for i, s := range m.fileSelected {
			if s {
				paths = append(paths, m.unstagedFiles[i].Path)
			}
		}
		if len(paths) == 0 {
			return m, nil
		}
		return m, doStageFiles(paths)
	}
	return m, nil
}

func (m model) handleVariantPickingKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "up", "k":
		if m.variantCursor > 0 {
			m.variantCursor--
		}
	case "down", "j":
		if m.variantCursor < len(m.variants)-1 {
			m.variantCursor++
		}
	case "1", "2", "3":
		idx := int(msg.String()[0] - '1')
		if idx < len(m.variants) {
			m.variantCursor = idx
			return m.selectVariant()
		}
	case "enter":
		return m.selectVariant()
	case "esc", "q":
		m.state = stateReady
		return m, nil
	case "ctrl+c":
		return m, tea.Quit
	}
	return m, nil
}

func (m model) selectVariant() (model, tea.Cmd) {
	m.previousMessages = append(m.previousMessages, m.commitMessage)
	m.commitMessage = m.variants[m.variantCursor]
	m.variants = nil
	m.state = stateReady
	m.statusMessage = "Variant selected"
	m.statusIsError = false
	return m, nil
}

func (m model) handleDiffViewKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc", "d", "q":
		m.state = stateReady
		return m, nil
	case "ctrl+c":
		return m, tea.Quit
	}
	var cmd tea.Cmd
	m.viewport, cmd = m.viewport.Update(msg)
	return m, cmd
}

func (m model) handleTypePickerKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "up", "k":
		if m.typeCursor > 0 {
			m.typeCursor--
		}
	case "down", "j":
		if m.typeCursor < len(conventionalTypes)-1 {
			m.typeCursor++
		}
	case "1", "2", "3", "4", "5", "6", "7", "8":
		idx := int(msg.String()[0] - '1')
		if idx < len(conventionalTypes) {
			m.typeCursor = idx
			return m.selectType()
		}
	case "enter":
		return m.selectType()
	case "esc", "q":
		m.state = stateReady
		return m, nil
	case "ctrl+c":
		return m, tea.Quit
	}
	return m, nil
}

func (m model) selectType() (model, tea.Cmd) {
	selected := conventionalTypes[m.typeCursor]
	m.previousMessages = append(m.previousMessages, m.commitMessage)
	m.state = stateGenerating
	m.loader = loader.New(loader.GeneratingMessages)
	return m, tea.Batch(m.loader.Tick, doFeedback(m.diff, m.cfg, m.commitMessage, "Force type prefix: "+selected))
}

func (m model) handleUnstagingKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "q", "ctrl+c":
		return m, tea.Quit
	case "esc":
		m.state = stateReady
		return m, nil
	case "up", "k":
		if m.unstageCursor > 0 {
			m.unstageCursor--
		}
	case "down", "j":
		if m.unstageCursor < len(m.stagedFiles)-1 {
			m.unstageCursor++
		}
	case " ":
		if m.unstageCursor < len(m.unstageSelected) {
			m.unstageSelected[m.unstageCursor] = !m.unstageSelected[m.unstageCursor]
		}
	case "enter":
		var paths []string
		for i, sel := range m.unstageSelected {
			if sel && i < len(m.stagedFiles) {
				paths = append(paths, m.stagedFiles[i].Path)
			}
		}
		if len(paths) == 0 {
			return m, nil
		}
		return m, doUnstageFiles(paths)
	}
	return m, nil
}

func (m model) handleHelpKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	m.state = stateReady
	return m, nil
}

func (m model) checkWelcomeDone() (tea.Model, tea.Cmd) {
	m.dataReady = m.repoLoaded && m.filesLoaded

	if !m.dataReady || !m.welcomeReady {
		return m, nil
	}

	if len(m.stagedFiles) == 0 {
		m.state = stateNoStaged
		return m, loadUnstagedFiles()
	}

	m.state = stateGenerating
	m.loader = loader.New(loader.GeneratingMessages)
	return m, tea.Batch(m.loader.Tick, loadDiffAndGenerate(m.cfg))
}

func (m model) openExternalEditor() (tea.Model, tea.Cmd) {
	editor := git.GetGitEditor()
	editorArgs := strings.Fields(editor)

	tmpFile, err := os.CreateTemp("", "diny-commit-*.txt")
	if err != nil {
		m.statusMessage = "Failed to create temp file"
		m.statusIsError = true
		return m, nil
	}

	if _, err := tmpFile.WriteString(m.commitMessage); err != nil {
		os.Remove(tmpFile.Name())
		m.statusMessage = "Failed to write temp file"
		m.statusIsError = true
		return m, nil
	}
	tmpFile.Close()
	tmpPath := tmpFile.Name()

	args := append(editorArgs[1:], tmpPath)
	c := exec.Command(editorArgs[0], args...)

	return m, tea.ExecProcess(c, func(err error) tea.Msg {
		defer os.Remove(tmpPath)
		if err != nil {
			return editorFinishedMsg{newMessage: ""}
		}
		content, readErr := os.ReadFile(tmpPath)
		if readErr != nil {
			return editorFinishedMsg{newMessage: ""}
		}
		return editorFinishedMsg{newMessage: strings.TrimSpace(string(content))}
	})
}
