package app

import (
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	"github.com/dinoDanic/diny/commit"
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
	stateFilePicker
	stateSplitGenerating
	stateSplitPlan
	stateSplitCommitting
	stateSplitSuccess
	stateSplitFailure
	stateSplitFeedback
)

type fileEntry struct {
	path          string
	status        string
	currentStaged bool
	wantStaged    bool
}

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

type splitEditorFinishedMsg struct {
	groupIdx   int
	newMessage string
}

type variantsReadyMsg struct {
	variants []string
}

type allFilesMsg struct {
	entries []fileEntry
}

type filePickerDoneMsg struct {
	files []git.StagedFile
}

type splitPlanReadyMsg struct {
	plan []commit.SplitGroup
}

type splitCommitDoneMsg struct {
	hashes []string
}

// splitCommitFailureMsg reports a structured failure after one or more
// successful commits have already landed.
type splitCommitFailureMsg struct {
	committedHashes []string       // hashes from groups that succeeded (indexed 0..K-1)
	failedIndex     int            // index of the failing group in the plan
	failedStderr    string         // stderr from git for the failing group
	remainingFiles  []string       // files from groups that never ran
	failedFiles     []string       // files from the failing group (left staged)
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
	commitProgress   string
	commitOutputCh   <-chan string

	// Components
	loader   loader.Model
	viewport viewport.Model
	textinput textinput.Model
	textarea  textarea.Model

	// Status
	statusMessage string
	statusIsError bool
	err           error
	currentTip    string

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

	// File picker (stateFilePicker)
	fileEntries      []fileEntry
	filePickerCursor int

	// Split plan (stateSplitPlan / stateSplitCommitting / stateSplitSuccess)
	splitPlan     []commit.SplitGroup
	splitCursor   int
	splitExpanded map[int]bool
	splitHashes   []string
	splitFailure  *splitCommitFailureMsg
	splitPushed   bool

	// Split move mode — reassigning a file from its current group to another
	splitMoveMode     bool // cursor is on a file in splitCursor's group
	splitMoveFileIdx  int  // index into splitPlan[splitCursor].Files
	splitMovePickDest bool // second step: choose destination group (>9 groups or via enter)
	splitMoveDestIdx  int  // cursor into destination group list

	// Split regeneration — session history and overlay
	splitPrevPlans    [][]commit.SplitGroup // all rejected plans in this session
	splitRegenerating bool                  // spinner overlay on plan view

	// CLI flags forwarded from cobra
	cliNoVerify bool
	cliPush     bool
	cliPrint    bool
}

func newModel(cfg *config.Config, version string, opts Options) model {
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
		currentTip:        randomTip(),
		cliNoVerify:       opts.NoVerify,
		cliPush:           opts.Push,
		cliPrint:          opts.Print,
	}
}
