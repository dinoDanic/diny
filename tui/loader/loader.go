package loader

import (
	"math/rand"
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/dinoDanic/diny/ui"
)

var brailleSpinner = spinner.Spinner{
	Frames: []string{"◐", "◓", "◑", "◒"},
	FPS:    time.Second / 6,
}

var InitMessages = []string{
	"waking up...",
	"stretching...",
	"booting...",
	"caffeinating...",
}

var GeneratingMessages = []string{
	"thinkering...",
	"vibing...",
	"cooking...",
	"brewing...",
	"conjuring...",
	"pondering...",
	"manifesting...",
	"hallucinating (nicely)...",
}

var CommittingMessages = []string{
	"shipping...",
	"yolo-ing...",
	"sending it...",
	"no going back now...",
}

var VariantMessages = []string{
	"cooking up options...",
	"brainstorming...",
	"generating variants...",
	"thinking of alternatives...",
}

type Model struct {
	Tick    tea.Cmd
	spinner spinner.Model
	message string
}

func New(messages []string) Model {
	s := spinner.New()
	s.Spinner = brailleSpinner
	msg := messages[rand.Intn(len(messages))]
	return Model{
		Tick:    s.Tick,
		spinner: s,
		message: msg,
	}
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd
	m.spinner, cmd = m.spinner.Update(msg)
	return m, cmd
}

func (m Model) View() string {
	t := ui.GetCurrentTheme()
	style := lipgloss.NewStyle().Foreground(t.PrimaryForeground)
	return style.Render(m.spinner.View() + " " + m.message)
}
