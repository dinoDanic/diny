package app

import (
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	"github.com/dinoDanic/diny/config"
	"github.com/dinoDanic/diny/git"
)

type state int

const (
	stateWelcome state = iota
	stateGenerating
	stateReady
	stateFeedback
	stateEditing
	stateHelp
	stateCommitting
	stateSuccess
	stateNoStaged
	stateError
)

// Messages

type repoInfoMsg struct {
	repoName    string
	branchName  string
	gitUserName string
}

type stagedFilesMsg struct {
	files []git.StagedFile
}

type welcomeTimerDoneMsg struct{}

type diffAndCommitMsg struct {
	diff          string
	commitMessage string
}

type commitDoneMsg struct {
	hash string
	push bool
}

type errMsg struct {
	err error
}

type draftSavedMsg struct{}
type copiedMsg struct{}

type editorFinishedMsg struct {
	newMessage string
}

// Model

type model struct {
	cfg           *config.Config
	version       string
	state         state
	width, height int

	// Context (loaded async in Init)
	repoName    string
	branchName  string
	gitUserName string
	stagedFiles []git.StagedFile

	// Commit state
	diff             string
	commitMessage    string
	previousMessages []string

	// Components
	spinner   spinner.Model
	viewport  viewport.Model
	textinput textinput.Model
	textarea  textarea.Model

	// Status
	statusMessage string
	statusIsError bool
	err           error

	// Welcome timing
	welcomeReady bool // min display timer elapsed
	dataReady    bool // staged files + repo info loaded
	repoLoaded   bool
	filesLoaded  bool

	// Commit options
	pendingPush     bool
	pendingNoVerify bool
}

func newModel(cfg *config.Config, version string) model {
	s := spinner.New()
	s.Spinner = spinner.Dot

	ti := textinput.New()
	ti.Placeholder = "Describe what to change..."
	ti.CharLimit = 200

	ta := textarea.New()
	ta.SetHeight(8)

	return model{
		cfg:     cfg,
		version: version,
		state:   stateWelcome,
		spinner: s,
	}
}
