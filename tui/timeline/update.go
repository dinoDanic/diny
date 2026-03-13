package timeline

import (
	"github.com/charmbracelet/bubbles/textinput"
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

	case analysisReadyMsg:
		if msg.commits != nil {
			m.commits = msg.commits
		}
		m.analysis = msg.analysis
		m.fullPrompt = msg.fullPrompt
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
	case stateFetching, stateRegenerating:
		var cmd tea.Cmd
		m.loader, cmd = m.loader.Update(msg)
		return m, cmd
	}

	// Update textinput for input states
	switch m.state {
	case stateEnterDate, stateEnterStartDate, stateEnterEndDate, stateFeedbackInput:
		var cmd tea.Cmd
		m.textinput, cmd = m.textinput.Update(msg)
		return m, cmd
	}

	return m, nil
}

func (m model) handleKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	key := msg.String()

	switch m.state {
	case stateDateSelect:
		switch key {
		case "up", "k":
			if m.dateCursor > 0 {
				m.dateCursor--
			}
			return m, nil
		case "down", "j":
			if m.dateCursor < len(dateMenuItems)-1 {
				m.dateCursor++
			}
			return m, nil
		case "enter":
			return m.confirmDateChoice()
		case "q", "ctrl+c":
			return m, tea.Quit
		}

	case stateEnterDate:
		switch key {
		case "enter":
			val := m.textinput.Value()
			m.startDate = val
			m.dateRange = val
			m.state = stateFetching
			m.loader = loader.New(loader.GeneratingMessages)
			return m, tea.Batch(m.loader.Tick, fetchAndGenerate(m.dateChoice, m.startDate, m.endDate, m.dateRange, m.cfg))
		case "esc":
			m.state = stateDateSelect
			return m, nil
		default:
			var cmd tea.Cmd
			m.textinput, cmd = m.textinput.Update(msg)
			return m, cmd
		}

	case stateEnterStartDate:
		switch key {
		case "enter":
			m.startDate = m.textinput.Value()
			m.state = stateEnterEndDate
			ti := textinput.New()
			ti.Placeholder = "YYYY-MM-DD"
			ti.Focus()
			m.textinput = ti
			return m, nil
		case "esc":
			m.state = stateDateSelect
			return m, nil
		default:
			var cmd tea.Cmd
			m.textinput, cmd = m.textinput.Update(msg)
			return m, cmd
		}

	case stateEnterEndDate:
		switch key {
		case "enter":
			m.endDate = m.textinput.Value()
			m.dateRange = m.startDate + " to " + m.endDate
			m.state = stateFetching
			m.loader = loader.New(loader.GeneratingMessages)
			return m, tea.Batch(m.loader.Tick, fetchAndGenerate(m.dateChoice, m.startDate, m.endDate, m.dateRange, m.cfg))
		case "esc":
			m.state = stateEnterStartDate
			ti := textinput.New()
			ti.Placeholder = "YYYY-MM-DD"
			ti.SetValue(m.startDate)
			ti.Focus()
			m.textinput = ti
			return m, nil
		default:
			var cmd tea.Cmd
			m.textinput, cmd = m.textinput.Update(msg)
			return m, cmd
		}

	case stateResults:
		switch key {
		case "c":
			return m, doCopy(m.analysis)
		case "s":
			return m, doSave(m.analysis, m.dateRange)
		case "r":
			m.previousAnalyses = append(m.previousAnalyses, m.analysis)
			m.state = stateRegenerating
			m.loader = loader.New(loader.GeneratingMessages)
			return m, tea.Batch(m.loader.Tick, doRegenerate(m.fullPrompt, m.cfg, m.previousAnalyses))
		case "f":
			ti := textinput.New()
			ti.Placeholder = "e.g., focus more on patterns, include statistics..."
			ti.CharLimit = 200
			ti.Focus()
			m.textinput = ti
			m.state = stateFeedbackInput
			return m, nil
		case "n":
			return m.resetToDateSelect()
		case "q", "ctrl+c":
			return m, tea.Quit
		}

	case stateFeedbackInput:
		switch key {
		case "enter":
			feedback := m.textinput.Value()
			m.previousAnalyses = append(m.previousAnalyses, m.analysis)
			m.state = stateRegenerating
			m.loader = loader.New(loader.GeneratingMessages)
			return m, tea.Batch(m.loader.Tick, doFeedback(m.fullPrompt, m.analysis, feedback, m.cfg))
		case "esc":
			m.state = stateResults
			return m, nil
		default:
			var cmd tea.Cmd
			m.textinput, cmd = m.textinput.Update(msg)
			return m, cmd
		}

	case stateNoCommits:
		switch key {
		case "n":
			return m.resetToDateSelect()
		case "q", "ctrl+c":
			return m, tea.Quit
		}

	case stateError:
		switch key {
		case "q", "ctrl+c":
			return m, tea.Quit
		}
	}

	if key == "ctrl+c" {
		return m, tea.Quit
	}

	return m, nil
}

func (m model) confirmDateChoice() (tea.Model, tea.Cmd) {
	switch m.dateCursor {
	case 0: // Today
		m.dateChoice = "today"
		m.dateRange = "today"
		m.state = stateFetching
		m.loader = loader.New(loader.GeneratingMessages)
		return m, tea.Batch(m.loader.Tick, fetchAndGenerate("today", "", "", "today", m.cfg))
	case 1: // Specific date
		m.dateChoice = "date"
		m.state = stateEnterDate
		ti := textinput.New()
		ti.Placeholder = "YYYY-MM-DD"
		ti.Focus()
		m.textinput = ti
		return m, nil
	case 2: // Date range
		m.dateChoice = "range"
		m.state = stateEnterStartDate
		ti := textinput.New()
		ti.Placeholder = "YYYY-MM-DD"
		ti.Focus()
		m.textinput = ti
		return m, nil
	}
	return m, nil
}

func (m model) resetToDateSelect() (tea.Model, tea.Cmd) {
	m.state = stateDateSelect
	m.dateCursor = 0
	m.dateChoice = ""
	m.startDate = ""
	m.endDate = ""
	m.dateRange = ""
	m.commits = nil
	m.analysis = ""
	m.previousAnalyses = nil
	m.fullPrompt = ""
	m.statusMessage = ""
	return m, nil
}
