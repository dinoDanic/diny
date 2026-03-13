package configtui

import (
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/dinoDanic/diny/config"
	"github.com/dinoDanic/diny/ui"
)

type state int

const (
	statePicker  state = iota
	stateMenu
	stateEditing
)

type fieldKind int

const (
	kindBool   fieldKind = iota
	kindSelect
	kindText
)

type field struct {
	label   string
	kind    fieldKind
	value   string
	options []string
}

// Messages

type savedMsg struct{}
type saveAndQuitMsg struct{}

type errMsg struct{ err error }

// Picker option

type pickerOption struct {
	label      string
	configType string
	configPath string
}

// Model

type model struct {
	version, repoName, branchName string
	width, height                 int

	configPath string
	configType string
	cfg        *config.Config

	state  state
	cursor int
	fields []field

	// statePicker
	pickerCursor  int
	pickerOptions []pickerOption

	// stateEditing
	activeField  int
	optionCursor int
	textinput    textinput.Model

	statusMessage string
	statusIsError bool
}

func newModel(version, repoName, branchName, configPath, configType string, cfg *config.Config) model {
	ti := textinput.New()
	ti.CharLimit = 500

	m := model{
		version:    version,
		repoName:   repoName,
		branchName: branchName,
		configPath: configPath,
		configType: configType,
		cfg:        cfg,
		textinput:  ti,
	}

	if configPath == "" {
		// In git repo: show picker
		m.state = statePicker
		m.pickerOptions = []pickerOption{
			{
				label:      "Global  (~/.config/diny/config.yaml)",
				configType: "global",
				configPath: config.GetConfigPath(),
			},
			{
				label:      "Versioned  (.diny.yaml)",
				configType: "versioned",
				configPath: config.GetVersionedProjectConfigPath(),
			},
			{
				label:      "Local  (<gitdir>/diny/config.yaml)",
				configType: "local",
				configPath: config.GetLocalProjectConfigPath(),
			},
		}
	} else {
		// Not in git repo or configPath already known: go straight to menu
		m.state = stateMenu
		if cfg != nil {
			m.fields = buildFields(cfg)
		}
	}

	return m
}

func buildFields(cfg *config.Config) []field {
	return []field{
		{
			label:   "theme",
			kind:    kindSelect,
			value:   cfg.Theme,
			options: ui.GetAvailableThemes(),
		},
		{
			label:   "tone",
			kind:    kindSelect,
			value:   string(cfg.Commit.Tone),
			options: []string{"professional", "casual", "friendly"},
		},
		{
			label:   "length",
			kind:    kindSelect,
			value:   string(cfg.Commit.Length),
			options: []string{"short", "normal", "long"},
		},
		{
			label:  "conventional",
			kind:   kindBool,
			value:  boolStr(cfg.Commit.Conventional),
		},
		{
			label:  "emoji",
			kind:   kindBool,
			value:  boolStr(cfg.Commit.Emoji),
		},
		{
			label:  "hash_after_commit",
			kind:   kindBool,
			value:  boolStr(cfg.Commit.HashAfterCommit),
		},
		{
			label: "custom_instructions",
			kind:  kindText,
			value: cfg.Commit.CustomInstructions,
		},
	}
}

func boolStr(b bool) string {
	if b {
		return "true"
	}
	return "false"
}

func toggleBool(val string) string {
	if val == "true" {
		return "false"
	}
	return "true"
}

func applyFieldsToConfig(fields []field, cfg *config.Config) {
	for _, f := range fields {
		switch f.label {
		case "theme":
			cfg.Theme = f.value
		case "tone":
			cfg.Commit.Tone = config.Tone(f.value)
		case "length":
			cfg.Commit.Length = config.Length(f.value)
		case "conventional":
			cfg.Commit.Conventional = f.value == "true"
		case "emoji":
			cfg.Commit.Emoji = f.value == "true"
		case "hash_after_commit":
			cfg.Commit.HashAfterCommit = f.value == "true"
		case "custom_instructions":
			cfg.Commit.CustomInstructions = f.value
		}
	}
}
