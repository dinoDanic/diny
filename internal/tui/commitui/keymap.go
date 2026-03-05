package commitui

import "github.com/charmbracelet/bubbles/key"

type keyMap struct {
	Commit      key.Binding
	NoVerify    key.Binding
	Push        key.Binding
	Regenerate  key.Binding
	Feedback    key.Binding
	Edit        key.Binding
	Editor      key.Binding
	SaveDraft   key.Binding
	Variants    key.Binding
	Copy        key.Binding
	Help        key.Binding
	Quit        key.Binding
	Submit      key.Binding
	Cancel      key.Binding
}

var keys = keyMap{
	Commit:     key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", "commit")),
	NoVerify:   key.NewBinding(key.WithKeys("n"), key.WithHelp("n", "no-verify")),
	Push:       key.NewBinding(key.WithKeys("p"), key.WithHelp("p", "push")),
	Regenerate: key.NewBinding(key.WithKeys("r"), key.WithHelp("r", "regen")),
	Feedback:   key.NewBinding(key.WithKeys("f"), key.WithHelp("f", "feedback")),
	Edit:       key.NewBinding(key.WithKeys("e"), key.WithHelp("e", "edit")),
	Editor:     key.NewBinding(key.WithKeys("E"), key.WithHelp("E", "$EDITOR")),
	SaveDraft:  key.NewBinding(key.WithKeys("s"), key.WithHelp("s", "draft")),
	Variants:   key.NewBinding(key.WithKeys("v"), key.WithHelp("v", "variants")),
	Copy:       key.NewBinding(key.WithKeys("y"), key.WithHelp("y", "copy")),
	Help:       key.NewBinding(key.WithKeys("?"), key.WithHelp("?", "help")),
	Quit:       key.NewBinding(key.WithKeys("q", "ctrl+c"), key.WithHelp("q", "quit")),
	Submit:     key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", "submit")),
	Cancel:     key.NewBinding(key.WithKeys("esc"), key.WithHelp("esc", "cancel")),
}
