package commitui

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
)

func (m model) Init() tea.Cmd {
	return tea.Batch(
		m.spinner.Tick,
		tea.SetWindowTitle("diny"),
		loadStagedFiles,
		loadDiffAndGenerate(m.cfg),
	)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.stackedLayout = msg.Width < 60
		m = m.updateViewportSize()
		return m, nil

	case spinner.TickMsg:
		if m.state == stateLoading {
			var cmd tea.Cmd
			m.spinner, cmd = m.spinner.Update(msg)
			return m, cmd
		}
		return m, nil

	case stagedFilesMsg:
		m.stagedFiles = msg.files
		return m, nil

	case commitMessageMsg:
		m.commitMessage = msg.message
		m.diff = msg.diff
		m.state = stateReady
		m.viewport.SetContent(m.commitMessage)
		return m, nil

	case errorMsg:
		if m.state == stateLoading {
			m.state = stateError
			m.err = msg.err
			return m, nil
		}
		m.statusMessage = msg.err.Error()
		m.statusIsError = true
		if m.state == stateFeedback || m.state == stateVariants {
			m.state = stateReady
		}
		return m, nil

	case commitDoneMsg:
		if msg.pushed {
			m.statusMessage = "Committed and pushed!"
		} else {
			m.statusMessage = "Committed!"
		}
		m.statusIsError = false
		return m, tea.Quit

	case regenerateMsg:
		m.previousMessages = append(m.previousMessages, m.commitMessage)
		m.commitMessage = msg.message
		m.viewport.SetContent(m.commitMessage)
		m.viewport.GotoTop()
		m.statusMessage = "Message regenerated"
		m.statusIsError = false
		m.state = stateReady
		return m, nil

	case feedbackMsg:
		m.previousMessages = append(m.previousMessages, m.commitMessage)
		m.commitMessage = msg.message
		m.viewport.SetContent(m.commitMessage)
		m.viewport.GotoTop()
		m.statusMessage = "Message refined with feedback"
		m.statusIsError = false
		m.state = stateReady
		return m, nil

	case editorMsg:
		if msg.message != "" && msg.message != m.commitMessage {
			m.commitMessage = msg.message
			m.viewport.SetContent(m.commitMessage)
			m.viewport.GotoTop()
			m.statusMessage = "Message updated from editor"
		}
		m.statusIsError = false
		m.state = stateReady
		return m, nil

	case draftMsg:
		m.statusMessage = "Draft saved!"
		m.statusIsError = false
		return m, nil

	case copyMsg:
		m.statusMessage = "Copied to clipboard!"
		m.statusIsError = false
		return m, nil

	case variantsMsg:
		m.variants = msg.variants
		m.variantCursor = 0
		m.state = stateVariants
		m.statusMessage = ""
		return m, nil

	case tea.KeyMsg:
		return m.handleKey(msg)
	}

	return m, nil
}

func (m model) handleKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	// Global quit
	if key.Matches(msg, keys.Quit) && m.state != stateFeedback && m.state != stateEditing {
		return m, tea.Quit
	}

	switch m.state {
	case stateError:
		return m, tea.Quit

	case stateReady:
		return m.handleReadyKey(msg)

	case stateFeedback:
		return m.handleFeedbackKey(msg)

	case stateEditing:
		return m.handleEditingKey(msg)

	case stateVariants:
		return m.handleVariantsKey(msg)

	case stateHelp:
		m.state = stateReady
		return m, nil
	}

	return m, nil
}

func (m model) handleReadyKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch {
	case key.Matches(msg, keys.Commit):
		m.state = stateLoading
		m.statusMessage = "Committing..."
		return m, doCommit(m.commitMessage, false, m.noVerify, m.cfg)

	case key.Matches(msg, keys.NoVerify):
		m.state = stateLoading
		m.statusMessage = "Committing (no-verify)..."
		return m, doCommit(m.commitMessage, false, true, m.cfg)

	case key.Matches(msg, keys.Push):
		m.state = stateLoading
		m.statusMessage = "Committing and pushing..."
		return m, doCommit(m.commitMessage, true, m.noVerify, m.cfg)

	case key.Matches(msg, keys.Regenerate):
		m.state = stateLoading
		m.statusMessage = "Regenerating..."
		return m, tea.Batch(m.spinner.Tick, doRegenerate(m.diff, m.cfg, m.previousMessages, m.commitMessage))

	case key.Matches(msg, keys.Feedback):
		m.state = stateFeedback
		m.textinput.Focus()
		m.textinput.Reset()
		m.statusMessage = ""
		return m, m.textinput.Cursor.BlinkCmd()

	case key.Matches(msg, keys.Edit):
		m.state = stateEditing
		m.textarea.SetValue(m.commitMessage)
		m.textarea.Focus()
		return m, m.textarea.Cursor.BlinkCmd()

	case key.Matches(msg, keys.Editor):
		return m, doOpenEditor(m.commitMessage)

	case key.Matches(msg, keys.SaveDraft):
		return m, doSaveDraft(m.commitMessage)

	case key.Matches(msg, keys.Variants):
		m.state = stateLoading
		m.statusMessage = "Generating variants..."
		return m, tea.Batch(m.spinner.Tick, doGenerateVariants(m.diff, m.cfg, m.commitMessage))

	case key.Matches(msg, keys.Copy):
		return m, doCopy(m.commitMessage)

	case key.Matches(msg, keys.Help):
		m.state = stateHelp
		return m, nil

	default:
		var cmd tea.Cmd
		m.viewport, cmd = m.viewport.Update(msg)
		return m, cmd
	}
}

func (m model) handleFeedbackKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.Type {
	case tea.KeyEnter:
		feedback := m.textinput.Value()
		if feedback == "" {
			m.state = stateReady
			return m, nil
		}
		m.state = stateLoading
		m.statusMessage = "Refining with feedback..."
		m.textinput.Blur()
		return m, tea.Batch(m.spinner.Tick, doFeedback(m.diff, m.cfg, m.commitMessage, feedback))

	case tea.KeyEscape, tea.KeyCtrlC:
		m.state = stateReady
		m.textinput.Blur()
		return m, nil

	default:
		var cmd tea.Cmd
		m.textinput, cmd = m.textinput.Update(msg)
		return m, cmd
	}
}

func (m model) handleEditingKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.Type {
	case tea.KeyEscape:
		edited := m.textarea.Value()
		if edited != "" && edited != m.commitMessage {
			m.commitMessage = edited
			m.viewport.SetContent(m.commitMessage)
			m.viewport.GotoTop()
			m.statusMessage = "Message updated"
			m.statusIsError = false
		}
		m.state = stateReady
		m.textarea.Blur()
		return m, nil

	case tea.KeyCtrlC:
		m.state = stateReady
		m.textarea.Blur()
		m.statusMessage = "Edit cancelled"
		m.statusIsError = false
		return m, nil

	default:
		var cmd tea.Cmd
		m.textarea, cmd = m.textarea.Update(msg)
		return m, cmd
	}
}

func (m model) handleVariantsKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch {
	case msg.String() == "j" || msg.String() == "down":
		if m.variantCursor < len(m.variants)-1 {
			m.variantCursor++
		}
		return m, nil

	case msg.String() == "k" || msg.String() == "up":
		if m.variantCursor > 0 {
			m.variantCursor--
		}
		return m, nil

	case msg.Type == tea.KeyEnter:
		m.commitMessage = m.variants[m.variantCursor]
		m.viewport.SetContent(m.commitMessage)
		m.viewport.GotoTop()
		m.statusMessage = "Variant selected"
		m.statusIsError = false
		m.state = stateReady
		return m, nil

	case msg.Type == tea.KeyEscape:
		m.state = stateReady
		m.statusMessage = ""
		return m, nil
	}

	return m, nil
}

func (m model) updateViewportSize() model {
	_, bodyHeight := m.getLayoutDimensions()

	if m.stackedLayout {
		m.viewport = viewport.New(m.width-2, bodyHeight-3)
	} else {
		rightWidth := m.getRightPaneWidth()
		m.viewport = viewport.New(rightWidth-2, bodyHeight)
	}
	m.viewport.SetContent(m.commitMessage)

	m.textarea.SetWidth(m.getRightPaneWidth() - 4)
	m.textarea.SetHeight(bodyHeight - 2)

	return m
}

func (m model) getLayoutDimensions() (int, int) {
	headerHeight := 2
	footerHeight := 2
	statusHeight := 1
	bodyHeight := m.height - headerHeight - footerHeight - statusHeight
	if bodyHeight < 3 {
		bodyHeight = 3
	}
	return m.width, bodyHeight
}

func (m model) getLeftPaneWidth() int {
	if m.stackedLayout {
		return m.width
	}
	w := int(float64(m.width) * 0.35)
	if w < 25 {
		w = 25
	}
	if w > 50 {
		w = 50
	}
	return w
}

func (m model) getRightPaneWidth() int {
	if m.stackedLayout {
		return m.width
	}
	return m.width - m.getLeftPaneWidth() - 1
}
