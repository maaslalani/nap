package main

import "github.com/charmbracelet/bubbles/key"

// KeyMap is the mappings of actions to key bindings.
type KeyMap struct {
	Quit          key.Binding
	Search        key.Binding
	ToggleHelp    key.Binding
	NewSnippet    key.Binding
	DeleteSnippet key.Binding
	EditSnippet   key.Binding
	CopySnippet   key.Binding
	PasteSnippet  key.Binding
	SetFolder     key.Binding
	RenameSnippet key.Binding
	TagSnippet    key.Binding
	SetLanguage   key.Binding
	Confirm       key.Binding
	Cancel        key.Binding
	NextPane      key.Binding
	PreviousPane  key.Binding
}

// DefaultKeyMap is the default key map for the application.
var DefaultKeyMap = KeyMap{
	Quit:          key.NewBinding(key.WithKeys("q", "ctrl+c"), key.WithHelp("q", "exit")),
	Search:        key.NewBinding(key.WithKeys("/"), key.WithHelp("/", "search")),
	ToggleHelp:    key.NewBinding(key.WithKeys("?"), key.WithHelp("?", "help")),
	NewSnippet:    key.NewBinding(key.WithKeys("n"), key.WithHelp("n", "new")),
	DeleteSnippet: key.NewBinding(key.WithKeys("x"), key.WithHelp("x", "delete")),
	EditSnippet:   key.NewBinding(key.WithKeys("e"), key.WithHelp("e", "edit")),
	CopySnippet:   key.NewBinding(key.WithKeys("c"), key.WithHelp("c", "copy")),
	PasteSnippet:  key.NewBinding(key.WithKeys("p"), key.WithHelp("p", "paste")),
	RenameSnippet: key.NewBinding(key.WithKeys("r"), key.WithHelp("r", "rename")),
	SetFolder:     key.NewBinding(key.WithKeys("f"), key.WithHelp("f", "folder")),
	SetLanguage:   key.NewBinding(key.WithKeys("L"), key.WithHelp("L", "language")),
	TagSnippet:    key.NewBinding(key.WithKeys("t"), key.WithHelp("t", "tag")),
	Confirm:       key.NewBinding(key.WithKeys("y"), key.WithHelp("y", "confirm")),
	Cancel:        key.NewBinding(key.WithKeys("N", "esc"), key.WithHelp("N", "cancel")),
	NextPane:      key.NewBinding(key.WithKeys("tab"), key.WithHelp("tab", "navigate")),
	PreviousPane:  key.NewBinding(key.WithKeys("shift+tab"), key.WithHelp("shift+tab", "navigate")),
}

// ShortHelp returns a quick help menu.
func (k KeyMap) ShortHelp() []key.Binding {
	return []key.Binding{
		k.NextPane,
		k.Search,
		k.EditSnippet,
		k.DeleteSnippet,
		k.CopySnippet,
		k.NewSnippet,
		k.ToggleHelp,
	}
}

// FullHelp returns all help options in a more detailed view.
func (k KeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.NextPane, k.PreviousPane, k.ToggleHelp},
		{k.Search, k.CopySnippet, k.DeleteSnippet},
		{k.EditSnippet, k.PasteSnippet, k.RenameSnippet, k.SetFolder, k.SetLanguage},
	}
}
