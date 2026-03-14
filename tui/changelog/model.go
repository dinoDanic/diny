package changelog

import (
	"github.com/dinoDanic/diny/config"
	"github.com/dinoDanic/diny/git"
	"github.com/dinoDanic/diny/tui/loader"
)

type state int

const (
	stateModeSelect state = iota
	stateLoadingRefs
	stateSelectNewerRef
	stateSelectOlderRef
	stateGenerating
	stateResults
	stateRegenerating
	stateNoCommits
	stateError
)

const listPageSize = 10

// Messages

type repoInfoMsg struct {
	repoName   string
	branchName string
}

type refsLoadedMsg struct {
	tags    []string
	commits []git.CommitInfo
}

type changelogReadyMsg struct {
	result string
	prompt string
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

	mode       string // "tag" / "commit"
	modeCursor int

	listCursor int
	listOffset int

	tags      []string
	commits   []git.CommitInfo
	olderTags []string

	newerRef   string
	olderRef   string
	rangeLabel string

	result          string
	prompt          string
	previousResults []string

	loader loader.Model

	statusMessage string
	statusIsError bool

	err error
}

func newModel(cfg *config.Config, version string) model {
	return model{
		cfg:     cfg,
		version: version,
		state:   stateModeSelect,
		loader:  loader.New(loader.GeneratingMessages),
	}
}
