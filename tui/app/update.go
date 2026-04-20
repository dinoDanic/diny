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
	"github.com/dinoDanic/diny/commit"
	"github.com/dinoDanic/diny/config"
	"github.com/dinoDanic/diny/git"
	"github.com/dinoDanic/diny/tui/loader"
)

func clonePlan(plan []commit.SplitGroup) []commit.SplitGroup {
	out := make([]commit.SplitGroup, len(plan))
	for i, g := range plan {
		files := make([]string, len(g.Files))
		copy(files, g.Files)
		out[i] = commit.SplitGroup{
			Order:   g.Order,
			Type:    g.Type,
			Message: g.Message,
			Files:   files,
		}
	}
	return out
}

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

	case allFilesMsg:
		m.fileEntries = msg.entries
		m.filePickerCursor = 0
		return m, nil

	case filePickerDoneMsg:
		m.stagedFiles = msg.files
		if len(m.stagedFiles) == 0 {
			m.state = stateNoStaged
			return m, loadUnstagedFiles()
		}
		m.state = stateGenerating
		m.loader = loader.New(loader.GeneratingMessages)
		return m, tea.Batch(m.loader.Tick, loadDiffAndGenerate(m.cfg))

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

	case splitEditorFinishedMsg:
		if msg.newMessage != "" && msg.groupIdx >= 0 && msg.groupIdx < len(m.splitPlan) {
			m.splitPlan[msg.groupIdx].Message = msg.newMessage
		}
		return m, nil

	case variantsReadyMsg:
		m.variants = msg.variants
		m.variantCursor = 0
		m.state = stateVariantPicking
		return m, nil

	case commitProgressMsg:
		m.commitProgress = msg.line
		return m, waitForCommitLine(m.commitOutputCh)

	case splitPlanReadyMsg:
		m.splitPlan = msg.plan
		m.splitCursor = 0
		m.splitExpanded = map[int]bool{}
		m.splitRegenerating = false
		m.state = stateSplitPlan
		m.statusMessage = ""
		m.statusIsError = false
		return m, nil

	case splitCommitDoneMsg:
		m.splitHashes = msg.hashes
		m.splitPushed = m.cliPush
		m.state = stateSplitSuccess
		return m, tea.Quit

	case splitCommitFailureMsg:
		failure := msg
		m.splitFailure = &failure
		m.splitHashes = msg.committedHashes
		m.state = stateSplitFailure
		return m, nil
	}

	// Update sub-components
	var cmd tea.Cmd
	switch m.state {
	case stateWelcome, stateGenerating, stateCommitting, stateSplitGenerating, stateSplitCommitting:
		m.loader, cmd = m.loader.Update(msg)
		return m, cmd
	case stateSplitPlan:
		if m.splitRegenerating {
			m.loader, cmd = m.loader.Update(msg)
			return m, cmd
		}
	case stateFeedback, stateSplitFeedback:
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
	case stateFilePicker:
		return m.handleFilePickerKey(msg)
	case stateSplitPlan:
		return m.handleSplitPlanKey(msg)
	case stateSplitFeedback:
		return m.handleSplitFeedbackKey(msg)
	case stateSplitSuccess:
		if msg.String() == "q" || msg.String() == "ctrl+c" || msg.String() == "enter" {
			return m, tea.Quit
		}
	case stateSplitFailure:
		if msg.String() == "q" || msg.String() == "ctrl+c" || msg.String() == "enter" {
			return m, tea.Quit
		}
	case stateError, stateSuccess:
		if msg.String() == "q" || msg.String() == "ctrl+c" {
			return m, tea.Quit
		}
	case stateWelcome, stateGenerating, stateCommitting, stateSplitGenerating, stateSplitCommitting:
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
		m.commitProgress = ""
		ch := make(chan string, 20)
		m.commitOutputCh = ch
		return m, tea.Batch(doCommit(m.commitMessage, false, false, m.pendingAmend, m.cfg, ch), waitForCommitLine(ch), m.loader.Tick)
	case msg.String() == "n":
		m.state = stateCommitting
		m.loader = loader.New(loader.CommittingMessages)
		m.commitProgress = ""
		ch := make(chan string, 20)
		m.commitOutputCh = ch
		return m, tea.Batch(doCommit(m.commitMessage, false, true, m.pendingAmend, m.cfg, ch), waitForCommitLine(ch), m.loader.Tick)
	case msg.String() == "p":
		m.state = stateCommitting
		m.loader = loader.New(loader.CommittingMessages)
		m.commitProgress = ""
		ch := make(chan string, 20)
		m.commitOutputCh = ch
		return m, tea.Batch(doCommit(m.commitMessage, true, false, m.pendingAmend, m.cfg, ch), waitForCommitLine(ch), m.loader.Tick)
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
		m.state = stateFilePicker
		return m, loadAllFiles()
	case msg.String() == "s":
		return m, doSaveDraft(m.commitMessage)
	case msg.String() == "S":
		if m.cliPrint {
			m.statusMessage = "--print is incompatible with split; rerun without --print"
			m.statusIsError = true
			return m, nil
		}
		if len(m.stagedFiles) == 0 {
			m.statusMessage = "No staged files to split"
			m.statusIsError = true
			return m, nil
		}
		m.state = stateSplitGenerating
		m.loader = loader.New(loader.GeneratingMessages)
		return m, tea.Batch(m.loader.Tick, loadSplitPlan(m.diff, m.cfg, m.stagedFiles, nil, ""))
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

func (m model) handleFilePickerKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "q", "ctrl+c":
		return m, tea.Quit
	case "esc":
		m.state = stateReady
		return m, nil
	case "up", "k":
		if m.filePickerCursor > 0 {
			m.filePickerCursor--
		}
	case "down", "j":
		if m.filePickerCursor < len(m.fileEntries)-1 {
			m.filePickerCursor++
		}
	case " ":
		if m.filePickerCursor < len(m.fileEntries) {
			m.fileEntries[m.filePickerCursor].wantStaged = !m.fileEntries[m.filePickerCursor].wantStaged
		}
	case "a":
		anyPending := false
		for _, e := range m.fileEntries {
			if e.wantStaged != e.currentStaged {
				anyPending = true
				break
			}
		}
		for i := range m.fileEntries {
			if anyPending {
				m.fileEntries[i].wantStaged = m.fileEntries[i].currentStaged
			} else {
				m.fileEntries[i].wantStaged = !m.fileEntries[i].currentStaged
			}
		}
	case "enter":
		hasPending := false
		for _, e := range m.fileEntries {
			if e.wantStaged != e.currentStaged {
				hasPending = true
				break
			}
		}
		if !hasPending {
			m.state = stateReady
			return m, nil
		}
		return m, doApplyFilePicker(m.fileEntries)
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

func (m model) handleSplitFeedbackKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "enter":
		feedback := strings.TrimSpace(m.textinput.Value())
		if feedback == "" {
			m.state = stateSplitPlan
			return m, nil
		}
		prev := append([][]commit.SplitGroup(nil), m.splitPrevPlans...)
		prev = append(prev, clonePlan(m.splitPlan))
		m.splitPrevPlans = prev
		m.splitRegenerating = true
		m.state = stateSplitPlan
		m.loader = loader.New(loader.GeneratingMessages)
		return m, tea.Batch(m.loader.Tick, loadSplitPlan(m.diff, m.cfg, m.stagedFiles, prev, feedback))
	case "esc":
		m.state = stateSplitPlan
		return m, nil
	case "ctrl+c":
		return m, tea.Quit
	}
	var cmd tea.Cmd
	m.textinput, cmd = m.textinput.Update(msg)
	return m, cmd
}

func (m model) handleSplitPlanKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	if m.splitMoveMode {
		return m.handleSplitMoveKey(msg)
	}

	switch msg.String() {
	case "q", "esc":
		m.state = stateReady
		m.splitPlan = nil
		m.splitExpanded = nil
		m.splitCursor = 0
		m.statusMessage = ""
		m.statusIsError = false
		return m, nil
	case "ctrl+c":
		return m, tea.Quit
	case "up", "k":
		if m.splitCursor > 0 {
			m.splitCursor--
		}
		return m, nil
	case "down", "j":
		if m.splitCursor < len(m.splitPlan)-1 {
			m.splitCursor++
		}
		return m, nil
	case "enter", " ":
		if m.splitExpanded == nil {
			m.splitExpanded = map[int]bool{}
		}
		m.splitExpanded[m.splitCursor] = !m.splitExpanded[m.splitCursor]
		return m, nil
	case "e":
		if m.splitCursor < 0 || m.splitCursor >= len(m.splitPlan) {
			return m, nil
		}
		return m.openSplitGroupEditor(m.splitCursor)
	case "r":
		if m.splitRegenerating {
			return m, nil
		}
		prev := append([][]commit.SplitGroup(nil), m.splitPrevPlans...)
		prev = append(prev, clonePlan(m.splitPlan))
		m.splitPrevPlans = prev
		m.splitRegenerating = true
		m.loader = loader.New(loader.GeneratingMessages)
		return m, tea.Batch(m.loader.Tick, loadSplitPlan(m.diff, m.cfg, m.stagedFiles, prev, ""))
	case "f":
		m.state = stateSplitFeedback
		m.textinput = textinput.New()
		m.textinput.Placeholder = "What's wrong with this plan? (e.g. 'merge the refactor with the fix')"
		m.textinput.CharLimit = 400
		m.textinput.Width = 80
		m.textinput.Focus()
		return m, m.textinput.Cursor.BlinkCmd()
	case "m":
		if m.splitCursor < 0 || m.splitCursor >= len(m.splitPlan) {
			return m, nil
		}
		if len(m.splitPlan[m.splitCursor].Files) == 0 {
			return m, nil
		}
		m.splitMoveMode = true
		m.splitMoveFileIdx = 0
		m.splitMovePickDest = false
		m.splitMoveDestIdx = 0
		if m.splitExpanded == nil {
			m.splitExpanded = map[int]bool{}
		}
		m.splitExpanded[m.splitCursor] = true
		m.statusMessage = ""
		m.statusIsError = false
		return m, nil
	case "c":
		m.state = stateSplitCommitting
		m.loader = loader.New(loader.CommittingMessages)
		return m, tea.Batch(m.loader.Tick, doExecuteSplit(m.splitPlan, m.cliNoVerify, m.cliPush, m.cfg))
	}
	return m, nil
}

func (m model) handleSplitMoveKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	if m.splitCursor < 0 || m.splitCursor >= len(m.splitPlan) {
		m = m.exitSplitMoveMode()
		return m, nil
	}
	srcGroup := m.splitPlan[m.splitCursor]

	if m.splitMovePickDest {
		return m.handleSplitMoveDestKey(msg)
	}

	switch msg.String() {
	case "esc", "q":
		m = m.exitSplitMoveMode()
		return m, nil
	case "ctrl+c":
		return m, tea.Quit
	case "up", "k":
		if m.splitMoveFileIdx > 0 {
			m.splitMoveFileIdx--
		}
		return m, nil
	case "down", "j":
		if m.splitMoveFileIdx < len(srcGroup.Files)-1 {
			m.splitMoveFileIdx++
		}
		return m, nil
	case "enter":
		// Enter destination picker (arrow selection)
		m.splitMovePickDest = true
		m.splitMoveDestIdx = 0
		// Skip over source group — can't move a file to its own group
		if m.splitMoveDestIdx == m.splitCursor && len(m.splitPlan) > 1 {
			m.splitMoveDestIdx = 1
			if m.splitCursor == 0 {
				m.splitMoveDestIdx = 1
			} else {
				m.splitMoveDestIdx = 0
			}
		}
		return m, nil
	}

	// Digit shortcuts 1..9 for ≤9 groups
	if len(m.splitPlan) <= 9 {
		s := msg.String()
		if len(s) == 1 && s[0] >= '1' && s[0] <= '9' {
			dest := int(s[0]-'1') // zero-based
			if dest >= 0 && dest < len(m.splitPlan) && dest != m.splitCursor {
				return m.applySplitMove(dest)
			}
		}
	}
	return m, nil
}

func (m model) handleSplitMoveDestKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc":
		m.splitMovePickDest = false
		return m, nil
	case "q":
		m = m.exitSplitMoveMode()
		return m, nil
	case "ctrl+c":
		return m, tea.Quit
	case "up", "k":
		m.splitMoveDestIdx = m.prevOtherGroup(m.splitMoveDestIdx)
		return m, nil
	case "down", "j":
		m.splitMoveDestIdx = m.nextOtherGroup(m.splitMoveDestIdx)
		return m, nil
	case "enter":
		dest := m.splitMoveDestIdx
		if dest == m.splitCursor || dest < 0 || dest >= len(m.splitPlan) {
			return m, nil
		}
		return m.applySplitMove(dest)
	}
	return m, nil
}

func (m model) nextOtherGroup(i int) int {
	n := len(m.splitPlan)
	if n <= 1 {
		return i
	}
	for step := 1; step <= n; step++ {
		cand := (i + step) % n
		if cand != m.splitCursor {
			return cand
		}
	}
	return i
}

func (m model) prevOtherGroup(i int) int {
	n := len(m.splitPlan)
	if n <= 1 {
		return i
	}
	for step := 1; step <= n; step++ {
		cand := (i - step + n) % n
		if cand != m.splitCursor {
			return cand
		}
	}
	return i
}

func (m model) applySplitMove(dest int) (tea.Model, tea.Cmd) {
	src := m.splitCursor
	if src == dest || src < 0 || src >= len(m.splitPlan) || dest < 0 || dest >= len(m.splitPlan) {
		m = m.exitSplitMoveMode()
		return m, nil
	}
	srcGroup := &m.splitPlan[src]
	if m.splitMoveFileIdx < 0 || m.splitMoveFileIdx >= len(srcGroup.Files) {
		m = m.exitSplitMoveMode()
		return m, nil
	}
	file := srcGroup.Files[m.splitMoveFileIdx]
	srcGroup.Files = append(srcGroup.Files[:m.splitMoveFileIdx], srcGroup.Files[m.splitMoveFileIdx+1:]...)
	m.splitPlan[dest].Files = append(m.splitPlan[dest].Files, file)

	// If source group is now empty, remove it and adjust cursor.
	if len(srcGroup.Files) == 0 {
		m.splitPlan = append(m.splitPlan[:src], m.splitPlan[src+1:]...)
		// Renumber and fix orders.
		for i := range m.splitPlan {
			m.splitPlan[i].Order = i + 1
		}
		// Adjust cursor to the destination's new index.
		newDest := dest
		if dest > src {
			newDest = dest - 1
		}
		m.splitCursor = newDest
		// Rebuild expanded map by index shift.
		newExp := map[int]bool{}
		for k, v := range m.splitExpanded {
			if k == src {
				continue
			}
			nk := k
			if k > src {
				nk = k - 1
			}
			newExp[nk] = v
		}
		m.splitExpanded = newExp
	}
	m.statusMessage = fmt.Sprintf("moved %s", file)
	m.statusIsError = false
	m = m.exitSplitMoveMode()
	return m, nil
}

func (m model) exitSplitMoveMode() model {
	m.splitMoveMode = false
	m.splitMovePickDest = false
	m.splitMoveFileIdx = 0
	m.splitMoveDestIdx = 0
	return m
}

func (m model) openSplitGroupEditor(idx int) (tea.Model, tea.Cmd) {
	if idx < 0 || idx >= len(m.splitPlan) {
		return m, nil
	}
	editor := git.GetGitEditor()
	editorArgs := strings.Fields(editor)

	tmpFile, err := os.CreateTemp("", "diny-split-*.txt")
	if err != nil {
		m.statusMessage = "Failed to create temp file"
		m.statusIsError = true
		return m, nil
	}
	if _, err := tmpFile.WriteString(m.splitPlan[idx].Message); err != nil {
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
			return splitEditorFinishedMsg{groupIdx: idx, newMessage: ""}
		}
		content, readErr := os.ReadFile(tmpPath)
		if readErr != nil {
			return splitEditorFinishedMsg{groupIdx: idx, newMessage: ""}
		}
		return splitEditorFinishedMsg{groupIdx: idx, newMessage: strings.TrimSpace(string(content))}
	})
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
