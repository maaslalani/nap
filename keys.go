package main

import "github.com/charmbracelet/bubbles/key"

// KeyMap is the mappings of actions to key bindings.
type KeyMap struct {
	Quit            key.Binding
	Search          key.Binding
	ToggleHelp      key.Binding
	NewSnippet      key.Binding
	MoveSnippetUp   key.Binding
	MoveSnippetDown key.Binding
	DeleteSnippet   key.Binding
	EditSnippet     key.Binding
	CopySnippet     key.Binding
	PasteSnippet    key.Binding
	SetFolder       key.Binding
	RenameSnippet   key.Binding
	TagSnippet      key.Binding
	SetLanguage     key.Binding
	Confirm         key.Binding
	Cancel          key.Binding
	NextPane        key.Binding
	PreviousPane    key.Binding
	ChangeFolder    key.Binding
}

// DefaultKeyMap is the default key map for the application.
var DefaultKeyMap = KeyMap{
	Quit:            key.NewBinding(key.WithKeys("q", "ctrl+c"), key.WithHelp("q", "exit")),
	Search:          key.NewBinding(key.WithKeys("/"), key.WithHelp("/", "search")),
	ToggleHelp:      key.NewBinding(key.WithKeys("?"), key.WithHelp("?", "help")),
	NewSnippet:      key.NewBinding(key.WithKeys("n"), key.WithHelp("n", "new")),
	MoveSnippetDown: key.NewBinding(key.WithKeys("J"), key.WithHelp("J", "move snippet down")),
	MoveSnippetUp:   key.NewBinding(key.WithKeys("K"), key.WithHelp("K", "move snippet up")),
	DeleteSnippet:   key.NewBinding(key.WithKeys("x"), key.WithHelp("x", "delete")),
	EditSnippet:     key.NewBinding(key.WithKeys("e"), key.WithHelp("e", "edit")),
	CopySnippet:     key.NewBinding(key.WithKeys("c"), key.WithHelp("c", "copy")),
	PasteSnippet:    key.NewBinding(key.WithKeys("p"), key.WithHelp("p", "paste")),
	RenameSnippet:   key.NewBinding(key.WithKeys("r"), key.WithHelp("r", "rename snippet")),
	SetFolder:       key.NewBinding(key.WithKeys("R"), key.WithHelp("R", "rename folder")),
	SetLanguage:     key.NewBinding(key.WithKeys("L"), key.WithHelp("L", "set file type")),
	TagSnippet:      key.NewBinding(key.WithKeys("t"), key.WithHelp("t", "tag"), key.WithDisabled()),
	Confirm:         key.NewBinding(key.WithKeys("y"), key.WithHelp("y", "confirm")),
	Cancel:          key.NewBinding(key.WithKeys("N", "esc"), key.WithHelp("N", "cancel")),
	NextPane:        key.NewBinding(key.WithKeys("tab", "right", "l"), key.WithHelp("tab", "navigate")),
	PreviousPane:    key.NewBinding(key.WithKeys("shift+tab", "left", "h"), key.WithHelp("shift+tab", "navigate")),
	ChangeFolder:    key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", "change folder"), key.WithDisabled()),
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
		{k.NewSnippet, k.EditSnippet, k.PasteSnippet, k.CopySnippet, k.DeleteSnippet},
		{k.MoveSnippetDown, k.MoveSnippetUp},
		{k.RenameSnippet, k.SetFolder, k.TagSnippet, k.SetLanguage},
		{k.NextPane, k.PreviousPane},
		{k.Search, k.ToggleHelp, k.Quit},
	}
}
