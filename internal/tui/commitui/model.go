package commitui

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
	stateLoading state = iota
	stateReady
	stateError
	stateFeedback
	stateEditing
	stateVariants
	stateHelp
)

type model struct {
	cfg      *config.Config
	noVerify bool

	state        state
	width        int
	height       int
	stackedLayout bool

	stagedFiles    []git.StagedFile
	diff           string
	commitMessage  string
	previousMessages []string

	viewport viewport.Model
	spinner  spinner.Model
	textinput textinput.Model
	textarea  textarea.Model

	variants       []string
	variantCursor  int

	statusMessage string
	statusIsError bool

	err error
}

// Messages
type stagedFilesMsg struct {
	files []git.StagedFile
}

type commitMessageMsg struct {
	message string
	diff    string
}

type errorMsg struct {
	err error
}

type commitDoneMsg struct {
	pushed bool
}

type regenerateMsg struct {
	message string
}

type feedbackMsg struct {
	message string
}

type editorMsg struct {
	message string
}

type draftMsg struct{}

type copyMsg struct{}

type variantsMsg struct {
	variants []string
}

func newModel(cfg *config.Config, noVerify bool) model {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = getSpinnerStyle()

	ti := textinput.New()
	ti.Placeholder = "Describe what to change..."
	ti.CharLimit = 200

	ta := textarea.New()
	ta.CharLimit = 0

	return model{
		cfg:      cfg,
		noVerify: noVerify,
		state:    stateLoading,
		spinner:  s,
		textinput: ti,
		textarea:  ta,
	}
}
