package timeline

import (
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/dinoDanic/diny/config"
	"github.com/dinoDanic/diny/tui/loader"
)

type state int

const (
	stateDateSelect state = iota
	stateEnterDate
	stateEnterStartDate
	stateEnterEndDate
	statePickDate
	stateFetching
	stateResults
	stateFeedbackInput
	stateRegenerating
	stateNoCommits
	stateError
)

// Messages

type repoInfoMsg struct {
	repoName   string
	branchName string
}

type analysisReadyMsg struct {
	commits    []string
	analysis   string
	fullPrompt string
}

type noCommitsMsg struct{}

type copiedMsg struct{}

type savedMsg struct {
	filePath string
}

type errMsg struct {
	err error
}

// Model

type model struct {
	cfg     *config.Config
	version string
	state   state
	width   int

	repoName   string
	branchName string

	dateChoice string // "today" / "date" / "range"
	dateCursor int    // selected row in date-select menu

	startDate string
	endDate   string
	dateRange string

	commits          []string
	analysis         string
	previousAnalyses []string
	fullPrompt       string

	loader    loader.Model
	textinput textinput.Model
	picker    datePicker

	statusMessage string
	statusIsError bool

	err error
}

func newModel(cfg *config.Config, version string) model {
	ti := textinput.New()
	ti.Focus()
	return model{
		cfg:       cfg,
		version:   version,
		state:     stateDateSelect,
		dateCursor: 0,
		loader:    loader.New(loader.GeneratingMessages),
		textinput: ti,
	}
}
