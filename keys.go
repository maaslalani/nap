package main

import "github.com/charmbracelet/bubbles/key"

// KeyMap is the mappings of actions to key bindings.
type KeyMap struct {
	Quit          key.Binding
	ToggleHelp    key.Binding
	NewSnippet    key.Binding
	DeleteSnippet key.Binding
	EditSnippet   key.Binding
	Confirm       key.Binding
	Cancel        key.Binding
	NextPane      key.Binding
	PreviousPane  key.Binding
}

// DefaultKeyMap is the default key map for the application.
var DefaultKeyMap = KeyMap{
	Quit:          key.NewBinding(key.WithKeys("q", "ctrl+c"), key.WithHelp("q", "exit")),
	ToggleHelp:    key.NewBinding(key.WithKeys("?"), key.WithHelp("?", "help")),
	NewSnippet:    key.NewBinding(key.WithKeys("n"), key.WithHelp("n", "new")),
	DeleteSnippet: key.NewBinding(key.WithKeys("x"), key.WithHelp("x", "delete")),
	EditSnippet:   key.NewBinding(key.WithKeys("e"), key.WithHelp("e", "edit")),
	Confirm:       key.NewBinding(key.WithKeys("y"), key.WithHelp("y", "confirm")),
	Cancel:        key.NewBinding(key.WithKeys("N"), key.WithHelp("N", "cancel")),
	NextPane:      key.NewBinding(key.WithKeys("tab", "l"), key.WithHelp("tab", "focus next")),
	PreviousPane:  key.NewBinding(key.WithKeys("shift+tab", "h"), key.WithHelp("shift+tab", "focus previous")),
}
