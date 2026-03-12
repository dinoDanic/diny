package configtui

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/dinoDanic/diny/config"
)

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case tea.KeyMsg:
		return m.handleKey(msg)

	case savedMsg:
		m.statusMessage = "Saved!"
		m.statusIsError = false
		return m, nil

	case saveAndQuitMsg:
		return m, tea.Quit

	case errMsg:
		m.statusMessage = "Error: " + msg.err.Error()
		m.statusIsError = true
		return m, nil
	}

	return m, nil
}

func (m model) handleKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch m.state {
	case statePicker:
		return m.handlePickerKey(msg)
	case stateMenu:
		return m.handleMenuKey(msg)
	case stateEditing:
		return m.handleEditingKey(msg)
	}
	return m, nil
}

func (m model) handlePickerKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "up", "k":
		if m.pickerCursor > 0 {
			m.pickerCursor--
		}
	case "down", "j":
		if m.pickerCursor < len(m.pickerOptions)-1 {
			m.pickerCursor++
		}
	case "enter":
		opt := m.pickerOptions[m.pickerCursor]
		m.configPath = opt.configPath
		m.configType = opt.configType

		// Ensure the config file exists
		switch m.configType {
		case "versioned":
			config.CreateVersionedProjectConfigIfNeeded()
		case "local":
			config.CreateLocalProjectConfigIfNeeded()
		}

		// Load the merged effective config for display
		cfg, _, err := config.LoadWithProjectOverride("")
		if err != nil {
			cfg, _ = config.Load("")
		}
		if cfg == nil {
			cfg = &config.Config{}
		}
		m.cfg = cfg
		m.fields = buildFields(cfg)
		m.state = stateMenu
		m.cursor = 0
	case "q", "ctrl+c":
		return m, tea.Quit
	}
	return m, nil
}

func (m model) handleMenuKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "up", "k":
		if m.cursor > 0 {
			m.cursor--
		}
	case "down", "j":
		if m.cursor < len(m.fields)-1 {
			m.cursor++
		}
	case "enter":
		f := m.fields[m.cursor]
		switch f.kind {
		case kindBool:
			m.fields[m.cursor].value = toggleBool(f.value)
			applyFieldsToConfig(m.fields, m.cfg)
			return m, doSaveConfig(m.cfg, m.configPath, m.configType)
		case kindSelect:
			m.activeField = m.cursor
			// Position option cursor at current value
			for i, opt := range f.options {
				if opt == f.value {
					m.optionCursor = i
					break
				}
			}
			m.state = stateEditing
		case kindText:
			m.activeField = m.cursor
			m.textinput = textinput.New()
			m.textinput.SetValue(f.value)
			m.textinput.CharLimit = 500
			m.textinput.Focus()
			m.state = stateEditing
			return m, m.textinput.Cursor.BlinkCmd()
		}
	case "s":
		applyFieldsToConfig(m.fields, m.cfg)
		return m, doSaveConfig(m.cfg, m.configPath, m.configType)
	case "w":
		applyFieldsToConfig(m.fields, m.cfg)
		return m, doSaveAndQuit(m.cfg, m.configPath, m.configType)
	case "q", "ctrl+c":
		return m, tea.Quit
	}
	return m, nil
}

func (m model) handleEditingKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	f := m.fields[m.activeField]
	switch f.kind {
	case kindSelect:
		switch msg.String() {
		case "up", "k":
			if m.optionCursor > 0 {
				m.optionCursor--
			}
		case "down", "j":
			if m.optionCursor < len(f.options)-1 {
				m.optionCursor++
			}
		case "enter":
			m.fields[m.activeField].value = f.options[m.optionCursor]
			applyFieldsToConfig(m.fields, m.cfg)
			m.state = stateMenu
			return m, doSaveConfig(m.cfg, m.configPath, m.configType)
		case "esc":
			m.state = stateMenu
		}

	case kindText:
		switch msg.String() {
		case "enter":
			m.fields[m.activeField].value = m.textinput.Value()
			applyFieldsToConfig(m.fields, m.cfg)
			m.state = stateMenu
			return m, doSaveConfig(m.cfg, m.configPath, m.configType)
		case "esc":
			m.state = stateMenu
		default:
			var cmd tea.Cmd
			m.textinput, cmd = m.textinput.Update(msg)
			return m, cmd
		}
	}

	return m, nil
}
