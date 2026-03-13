package yolo

import (
	"github.com/dinoDanic/diny/config"
	"github.com/dinoDanic/diny/git"
	"github.com/dinoDanic/diny/tui/loader"
)

type state int

const (
	stateStaging state = iota
	stateGenerating
	stateCommitting
	stateSuccess
	stateNothingToCommit
	stateError
)

// Messages

type repoInfoMsg struct {
	repoName   string
	branchName string
}

type stageDoneMsg struct {
	files []git.StagedFile
}

type generateDoneMsg struct {
	commitMessage string
}

type nothingToCommitMsg struct{}

type commitDoneMsg struct {
	hash string
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

	repoName     string
	branchName   string
	stagedFiles  []git.StagedFile
	commitMessage string
	hash         string
	err          error

	loader   loader.Model
	quitting bool
}

func newModel(cfg *config.Config, version string) model {
	return model{
		cfg:     cfg,
		version: version,
		state:   stateStaging,
		loader:  loader.New(loader.InitMessages),
	}
}
