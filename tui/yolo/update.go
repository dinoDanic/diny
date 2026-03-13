package yolo

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/dinoDanic/diny/tui/loader"
)

func (m model) Init() tea.Cmd {
	return tea.Batch(
		m.loader.Tick,
		loadRepoInfo(),
		doStageAll(),
		tea.WindowSize(),
	)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		return m, nil

	case tea.KeyMsg:
		if m.state == stateError {
			if msg.String() == "q" || msg.String() == "ctrl+c" {
				return m, tea.Quit
			}
		}
		if msg.String() == "ctrl+c" {
			return m, tea.Quit
		}
		return m, nil

	case repoInfoMsg:
		m.repoName = msg.repoName
		m.branchName = msg.branchName
		return m, nil

	case stageDoneMsg:
		m.stagedFiles = msg.files
		m.state = stateGenerating
		m.loader = loader.New(loader.GeneratingMessages)
		return m, tea.Batch(m.loader.Tick, loadDiffAndGenerate(m.cfg))

	case nothingToCommitMsg:
		m.state = stateNothingToCommit
		return m, tea.Quit

	case generateDoneMsg:
		m.commitMessage = msg.commitMessage
		m.state = stateCommitting
		m.loader = loader.New(loader.CommittingMessages)
		return m, tea.Batch(m.loader.Tick, doCommitAndPush(m.commitMessage, m.cfg))

	case commitDoneMsg:
		m.hash = msg.hash
		m.state = stateSuccess
		return m, tea.Quit

	case errMsg:
		m.err = msg.err
		m.state = stateError
		return m, nil
	}

	// Update loader for spinner states
	switch m.state {
	case stateStaging, stateGenerating, stateCommitting:
		var cmd tea.Cmd
		m.loader, cmd = m.loader.Update(msg)
		return m, cmd
	}

	return m, nil
}
