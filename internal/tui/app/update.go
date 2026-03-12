package app

import (
	"os"
	"os/exec"
	"strings"

	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/dinoDanic/diny/git"
	"github.com/dinoDanic/diny/internal/tui/loader"
)

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
		m.state = stateReady
		m.statusMessage = ""
		m.statusIsError = false
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
		return m, doCommit(m.commitMessage, false, false, m.cfg)
	case msg.String() == "n":
		m.state = stateCommitting
		m.loader = loader.New(loader.CommittingMessages)
		return m, doCommit(m.commitMessage, false, true, m.cfg)
	case msg.String() == "p":
		m.state = stateCommitting
		m.loader = loader.New(loader.CommittingMessages)
		return m, doCommit(m.commitMessage, true, false, m.cfg)
	case msg.String() == "r":
		m.state = stateGenerating
		m.loader = loader.New(loader.GeneratingMessages)
		prev := m.previousMessages
		m.previousMessages = append(m.previousMessages, m.commitMessage)
		return m, tea.Batch(m.loader.Tick, doRegenerate(m.diff, m.cfg, prev, m.commitMessage))
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
