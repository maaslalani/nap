package main

import (
	"fmt"
	"io"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/dustin/go-humanize"
)

func (s Snippet) FilterValue() string {
	return s.Folder + "/" + s.Title + "\n" + "+" + strings.Join(s.Tags, "+") + "\n" + s.Language
}

type snippetDelegate struct{}

func (d snippetDelegate) Height() int {
	return 2
}

func (d snippetDelegate) Spacing() int {
	return 1
}

func (d snippetDelegate) Update(msg tea.Msg, m *list.Model) tea.Cmd {
	return func() tea.Msg {
		return updateViewMsg(m.SelectedItem().(Snippet))
	}
}
func (d snippetDelegate) Render(w io.Writer, m list.Model, index int, item list.Item) {
	s, ok := item.(Snippet)
	if !ok {
		return
	}
	if index == m.Index() {
		fmt.Fprintln(w, "  "+DefaultStyles.Snippets.Focused.SelectedTitle.Render(s.Title))
		fmt.Fprint(w, "  "+DefaultStyles.Snippets.Focused.SelectedSubtitle.Render(humanize.Time(s.Date)))
		return
	}
	fmt.Fprintln(w, "  "+DefaultStyles.Snippets.Focused.UnselectedTitle.Render(s.Title))
	fmt.Fprint(w, "  "+DefaultStyles.Snippets.Focused.UnselectedSubtitle.Render(humanize.Time(s.Date)))
}

type Folder string

func (f Folder) FilterValue() string {
	return string(f)
}

type folderDelegate struct{}

func (d folderDelegate) Height() int {
	return 1
}

func (d folderDelegate) Spacing() int {
	return 0
}

func (d folderDelegate) Update(msg tea.Msg, m *list.Model) tea.Cmd {
	return nil
}

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
