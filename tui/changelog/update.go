package changelog

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/dinoDanic/diny/tui/loader"
)

func (m model) Init() tea.Cmd {
	return tea.Batch(
		loadRepoInfo(),
		tea.WindowSize(),
	)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		return m, nil

	case repoInfoMsg:
		m.repoName = msg.repoName
		m.branchName = msg.branchName
		return m, nil

	case refsLoadedMsg:
		m.tags = msg.tags
		m.commits = msg.commits
		// Validate we have enough refs
		if m.mode == "tag" && len(m.tags) < 2 {
			m.err = fmt.Errorf("at least two tags are required; found %d", len(m.tags))
			m.state = stateError
			return m, nil
		}
		if m.mode == "commit" && len(m.commits) < 2 {
			m.err = fmt.Errorf("at least two commits are required; found %d", len(m.commits))
			m.state = stateError
			return m, nil
		}
		m.listCursor = 0
		m.listOffset = 0
		m.state = stateSelectNewerRef
		return m, nil

	case changelogReadyMsg:
		m.result = msg.result
		m.prompt = msg.prompt
		m.state = stateResults
		return m, nil

	case noCommitsMsg:
		m.state = stateNoCommits
		return m, nil

	case copiedMsg:
		m.statusMessage = "Copied!"
		m.statusIsError = false
		return m, nil

	case savedMsg:
		m.statusMessage = "Saved: " + msg.filePath
		m.statusIsError = false
		return m, nil

	case errMsg:
		m.err = msg.err
		m.state = stateError
		return m, nil

	case tea.KeyMsg:
		return m.handleKey(msg)
	}

	// Update loader for spinner states
	switch m.state {
	case stateLoadingRefs, stateGenerating, stateRegenerating:
		var cmd tea.Cmd
		m.loader, cmd = m.loader.Update(msg)
		return m, cmd
	}

	return m, nil
}

func (m model) handleKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	key := msg.String()

	switch m.state {
	case stateModeSelect:
		switch key {
		case "up", "k":
			if m.modeCursor > 0 {
				m.modeCursor--
			}
		case "down", "j":
			if m.modeCursor < len(modeMenuItems)-1 {
				m.modeCursor++
			}
		case "enter":
			m.mode = modeMenuItems[m.modeCursor].value
			m.state = stateLoadingRefs
			m.loader = loader.New(loader.GeneratingMessages)
			return m, tea.Batch(m.loader.Tick, loadRefs(m.mode))
		case "q", "ctrl+c":
			return m, tea.Quit
		}
		return m, nil

	case stateSelectNewerRef:
		items := m.currentListLabels()
		switch key {
		case "up", "k":
			if m.listCursor > 0 {
				m.listCursor--
				if m.listCursor < m.listOffset {
					m.listOffset = m.listCursor
				}
			}
		case "down", "j":
			if m.listCursor < len(items)-1 {
				m.listCursor++
				if m.listCursor >= m.listOffset+listPageSize {
					m.listOffset = m.listCursor - listPageSize + 1
				}
			}
		case "enter":
			m.newerRef = m.selectedValue()
			if m.mode == "tag" {
				newerIdx := indexOf(m.tags, m.newerRef)
				m.olderTags = m.tags[newerIdx+1:]
				if len(m.olderTags) == 0 {
					m.err = fmt.Errorf("no older tags available before %s", m.newerRef)
					m.state = stateError
					return m, nil
				}
			}
			m.listCursor = 0
			m.listOffset = 0
			m.state = stateSelectOlderRef
		case "esc":
			m.state = stateModeSelect
		case "q", "ctrl+c":
			return m, tea.Quit
		}
		return m, nil

	case stateSelectOlderRef:
		items := m.currentListLabels()
		switch key {
		case "up", "k":
			if m.listCursor > 0 {
				m.listCursor--
				if m.listCursor < m.listOffset {
					m.listOffset = m.listCursor
				}
			}
		case "down", "j":
			if m.listCursor < len(items)-1 {
				m.listCursor++
				if m.listCursor >= m.listOffset+listPageSize {
					m.listOffset = m.listCursor - listPageSize + 1
				}
			}
		case "enter":
			m.olderRef = m.selectedValue()
			m.rangeLabel = fmt.Sprintf("%s → %s", m.olderRef, m.newerRef)
			m.state = stateGenerating
			m.loader = loader.New(loader.GeneratingMessages)
			return m, tea.Batch(m.loader.Tick, doGenerate(m.olderRef, m.newerRef, m.cfg))
		case "esc":
			m.listCursor = 0
			m.listOffset = 0
			m.state = stateSelectNewerRef
		case "q", "ctrl+c":
			return m, tea.Quit
		}
		return m, nil

	case stateResults:
		switch key {
		case "c":
			return m, doCopy(m.result)
		case "s":
			return m, doSave(m.result, m.rangeLabel)
		case "r":
			m.previousResults = append(m.previousResults, m.result)
			m.state = stateRegenerating
			m.loader = loader.New(loader.GeneratingMessages)
			return m, tea.Batch(m.loader.Tick, doRegenerate(m.prompt, m.cfg, m.previousResults))
		case "n":
			return m.resetToModeSelect()
		case "q", "ctrl+c":
			return m, tea.Quit
		}
		return m, nil

	case stateNoCommits:
		switch key {
		case "n":
			return m.resetToModeSelect()
		case "q", "ctrl+c":
			return m, tea.Quit
		}
		return m, nil

	case stateError:
		switch key {
		case "q", "ctrl+c":
			return m, tea.Quit
		}
		return m, nil
	}

	if key == "ctrl+c" {
		return m, tea.Quit
	}
	return m, nil
}

// selectedValue returns the ref value (tag name or commit SHA) at the current cursor.
func (m model) selectedValue() string {
	switch m.state {
	case stateSelectNewerRef:
		if m.mode == "tag" {
			if m.listCursor < len(m.tags) {
				return m.tags[m.listCursor]
			}
		} else {
			if m.listCursor < len(m.commits) {
				return m.commits[m.listCursor].SHA
			}
		}
	case stateSelectOlderRef:
		if m.mode == "tag" {
			if m.listCursor < len(m.olderTags) {
				return m.olderTags[m.listCursor]
			}
		} else {
			if m.listCursor < len(m.commits) {
				return m.commits[m.listCursor].SHA
			}
		}
	}
	return ""
}

func (m model) resetToModeSelect() (tea.Model, tea.Cmd) {
	m.state = stateModeSelect
	m.modeCursor = 0
	m.listCursor = 0
	m.listOffset = 0
	m.mode = ""
	m.newerRef = ""
	m.olderRef = ""
	m.rangeLabel = ""
	m.result = ""
	m.prompt = ""
	m.previousResults = nil
	m.statusMessage = ""
	m.olderTags = nil
	return m, nil
}

func indexOf(slice []string, target string) int {
	for i, s := range slice {
		if s == target {
			return i
		}
	}
	return -1
}
