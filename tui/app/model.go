package app

import (
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	"github.com/dinoDanic/diny/config"
	"github.com/dinoDanic/diny/git"
	"github.com/dinoDanic/diny/tui/loader"
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
	stateVariantPicking
	stateDiffView
	stateTypePicker
	stateUnstaging
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

type unstagedFilesMsg struct {
	files []git.StagedFile
}

type errMsg struct {
	err error
}

type draftSavedMsg struct{}
type copiedMsg struct{}

type editorFinishedMsg struct {
	newMessage string
}

type variantsReadyMsg struct {
	variants []string
}

type restoreStagedDoneMsg struct {
	files []git.StagedFile
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
	loader   loader.Model
	viewport viewport.Model
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
	pendingAmend    bool

	// History navigation ([/] keys)
	messageHistoryIdx int    // -1 = current; >=0 = index into previousMessages
	savedMessage      string // preserved current message when browsing history

	// File picker (stateNoStaged)
	unstagedFiles []git.StagedFile
	fileCursor    int
	fileSelected  []bool

	// Variant picker (stateVariantPicking)
	variants      []string
	variantCursor int

	// Type picker (stateTypePicker)
	typeCursor int

	// Unstaging (stateUnstaging)
	unstageCursor   int
	unstageSelected []bool
}

func newModel(cfg *config.Config, version string) model {
	ti := textinput.New()
	ti.Placeholder = "Describe what to change..."
	ti.CharLimit = 200

	ta := textarea.New()
	ta.SetHeight(8)

	return model{
		cfg:               cfg,
		version:           version,
		state:             stateWelcome,
		loader:            loader.New(loader.InitMessages),
		messageHistoryIdx: -1,
	}
}
