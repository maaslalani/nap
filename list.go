package main

import (
	"fmt"
	"io"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/dustin/go-humanize"
)

// FilterValue is the snippet filter value that can be used when searching.
func (s Snippet) FilterValue() string {
	return s.Folder + "/" + s.Name + "\n" + "+" + strings.Join(s.Tags, "+") + "\n" + s.Language
}

// snippetDelegate represents the snippet list item.
type snippetDelegate struct{}

// Height is the number of lines the snippet list item takes up.
func (d snippetDelegate) Height() int {
	return 2
}

// Spacing is the number of lines to insert between list items.
func (d snippetDelegate) Spacing() int {
	return 1
}

// Update is called when the list is updated.
// We use this to update the snippet code view.
func (d snippetDelegate) Update(msg tea.Msg, m *list.Model) tea.Cmd {
	return func() tea.Msg {
		if m.SelectedItem() == nil {
			return nil
		}
		return updateContentMsg(m.SelectedItem().(Snippet))
	}
}

// Render renders the list item for the snippet which includes the title,
// folder, and date.
func (d snippetDelegate) Render(w io.Writer, m list.Model, index int, item list.Item) {
	if item == nil {
		return
	}
	s, ok := item.(Snippet)
	if !ok {
		return
	}
	if index == m.Index() {
		fmt.Fprintln(w, "  "+DefaultStyles.Snippets.Focused.SelectedTitle.Render(s.Name))
		fmt.Fprint(w, "  "+DefaultStyles.Snippets.Focused.SelectedSubtitle.Render(humanize.Time(s.Date)))
		return
	}
	fmt.Fprintln(w, "  "+DefaultStyles.Snippets.Focused.UnselectedTitle.Render(s.Name))
	fmt.Fprint(w, "  "+DefaultStyles.Snippets.Focused.UnselectedSubtitle.Render(humanize.Time(s.Date)))
}

// Folder represents a group of snippets in a directory.
type Folder string

// FilterValue is the searchable value for the folder.
func (f Folder) FilterValue() string {
	return string(f)
}

// folderDelegate represents a folder list item.
type folderDelegate struct{}

// Height is the number of lines the folder list item takes up.
func (d folderDelegate) Height() int {
	return 1
}

// Spacing is the number of lines to insert between folder items.
func (d folderDelegate) Spacing() int {
	return 0
}

// Update is what is called when the folder selection is updated.
// TODO: Update the filter search for the snippets with the folder name.
func (d folderDelegate) Update(msg tea.Msg, m *list.Model) tea.Cmd {
	return nil
}

// Render renders a folder list item.
func (d folderDelegate) Render(w io.Writer, m list.Model, index int, item list.Item) {
	f, ok := item.(Folder)
	if !ok {
		return
	}
	fmt.Fprint(w, "  ")
	if index == m.Index() {
		fmt.Fprint(w, DefaultStyles.Folders.Focused.Selected.Render(string(f)))
		return
	}
	fmt.Fprint(w, DefaultStyles.Folders.Focused.Unselected.Render(string(f)))
}
