package launcher

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/dinoDanic/diny/config"
	"github.com/dinoDanic/diny/internal/tui/commitui"
)

const leftPanelWidth = 22

type item struct {
	title       string
	description string
}

// subDoneMsg is sent when the sub-panel signals quit (instead of tea.Quit).
type subDoneMsg struct{}

type model struct {
	items       []item
	cursor      int
	activeIndex int // -1 = no sub-panel active
	subModel    tea.Model
	cfg         *config.Config
	width       int
	height      int
}

func newModel(cfg *config.Config) model {
	return model{
		items: []item{
			{title: "Commit", description: "Generate AI commit message from staged changes"},
		},
		cursor:      0,
		activeIndex: -1,
		cfg:         cfg,
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

// wrapSubCmd intercepts tea.QuitMsg produced by sub-model commands and
// replaces it with subDoneMsg so the launcher can return to the menu.
func wrapSubCmd(cmd tea.Cmd) tea.Cmd {
	if cmd == nil {
		return nil
	}
	return func() tea.Msg {
		msg := cmd()
		if _, ok := msg.(tea.QuitMsg); ok {
			return subDoneMsg{}
		}
		return msg
	}
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		if m.subModel != nil {
			sub, cmd := m.subModel.Update(tea.WindowSizeMsg{Width: m.rightWidth(), Height: msg.Height})
			m.subModel = sub
			return m, wrapSubCmd(cmd)
		}
		return m, nil

	case subDoneMsg:
		m.subModel = nil
		m.activeIndex = -1
		return m, nil

	case tea.KeyMsg:
		if m.subModel != nil {
			sub, cmd := m.subModel.Update(msg)
			m.subModel = sub
			return m, wrapSubCmd(cmd)
		}
		switch msg.String() {
		case "q", "ctrl+c", "esc":
			return m, tea.Quit
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(m.items)-1 {
				m.cursor++
			}
		case "enter":
			return m.activateItem(m.cursor)
		}
		return m, nil
	}

	// Forward all other msgs (spinner ticks, API responses, etc.) to sub-model.
	if m.subModel != nil {
		sub, cmd := m.subModel.Update(msg)
		m.subModel = sub
		return m, wrapSubCmd(cmd)
	}

	return m, nil
}

func (m model) activateItem(index int) (tea.Model, tea.Cmd) {
	var sub tea.Model
	switch index {
	case 0:
		sub = commitui.New(m.cfg, false)
	default:
		return m, nil
	}

	// Send initial window size before Init so the sub-model knows its dimensions.
	sub, sizeCmd := sub.Update(tea.WindowSizeMsg{Width: m.rightWidth(), Height: m.height})
	m.subModel = sub
	m.activeIndex = index

	initCmd := m.subModel.Init()
	return m, tea.Batch(wrapSubCmd(sizeCmd), wrapSubCmd(initCmd))
}

func (m model) rightWidth() int {
	rw := m.width - leftPanelWidth - 1
	if rw < 20 {
		rw = 20
	}
	return rw
}
