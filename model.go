package main

import (
	"bytes"
	"os"

	"github.com/alecthomas/chroma/quick"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const MaxPane = 3

const (
	SnippetPane int = iota
	ContentPane
	FolderPane
)

// Model represents the state of the application.
// It contains all the snippets organized in folders.
type Model struct {
	// the working directory.
	Workdir string
	// code Snippets.
	Snippets []Snippet
	// the List of snippets to display to the user.
	List list.Model
	// the list of Folders to display to the user.
	Folders list.Model
	// the viewport of the Code snippet.
	Code viewport.Model
	// the current active pane of focus.
	Active int

	ListStyle    SnippetsBaseStyle
	FoldersStyle FoldersBaseStyle
	ContentStyle ContentBaseStyle
}

// Init initialzes the application model.
func (m *Model) Init() tea.Cmd {
	m.List.Styles.TitleBar = m.ListStyle.TitleBar
	m.List.Styles.Title = m.ListStyle.Title
	m.Folders.Styles.TitleBar = m.FoldersStyle.TitleBar
	m.Folders.Styles.Title = m.FoldersStyle.Title
	return func() tea.Msg {
		if len(m.Snippets) > 0 {
			return updateViewMsg(m.Snippets[0])
		}
		return nil
	}
}

// updateViewMsg tells the application to update the content view with the
// given snippet.
type updateViewMsg Snippet

// Update updates the model based on user interaction.
func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case updateViewMsg:
		var b bytes.Buffer
		b.WriteString(m.ContentStyle.Title.Render(msg.Title))
		b.WriteRune('\n')
		content, err := os.ReadFile(".leaf/" + msg.Folder + "/" + msg.File)
		if err != nil {
			m.Code.SetContent(b.String() + "Error: unable to read file.")
			return m, nil
		}
		err = quick.Highlight(&b, string(content), msg.Language, "terminal16m", "dracula")
		m.Code.SetContent(b.String())
		return m, nil
	case tea.WindowSizeMsg:
		m.List.SetHeight(msg.Height - 1)
		m.Folders.SetHeight(msg.Height - 1)
		m.Code.Height = msg.Height - 1
		m.Code.Width = msg.Width
		return m, nil
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyTab:
			m.nextPane()
		case tea.KeyShiftTab:
			m.previousPane()
		}
	}

	cmd := m.updateActivePane(msg)
	return m, cmd
}

// nextPane sets the next pane to be active.
func (m *Model) nextPane() {
	m.Active = (m.Active + 1) % MaxPane
}

// previousPane sets the previous pane to be active.
func (m *Model) previousPane() {
	m.Active--
	if m.Active < 0 {
		m.Active = MaxPane - 1
	}
}

// updateActivePane updates the currently active pane.
func (m *Model) updateActivePane(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd
	switch m.Active {
	case FolderPane:
		m.ListStyle = DefaultStyles.Snippets.Blurred
		m.ContentStyle = DefaultStyles.Content.Blurred
		m.FoldersStyle = DefaultStyles.Folders.Focused
		m.Folders, cmd = m.Folders.Update(msg)
	case SnippetPane:
		m.ListStyle = DefaultStyles.Snippets.Focused
		m.ContentStyle = DefaultStyles.Content.Blurred
		m.FoldersStyle = DefaultStyles.Folders.Blurred
		m.List, cmd = m.List.Update(msg)
	case ContentPane:
		m.ListStyle = DefaultStyles.Snippets.Blurred
		m.ContentStyle = DefaultStyles.Content.Focused
		m.FoldersStyle = DefaultStyles.Folders.Blurred
		m.Code, cmd = m.Code.Update(msg)
	}
	m.List.Styles.TitleBar = m.ListStyle.TitleBar
	m.List.Styles.Title = m.ListStyle.Title
	m.Folders.Styles.TitleBar = m.FoldersStyle.TitleBar
	m.Folders.Styles.Title = m.FoldersStyle.Title
	return cmd
}

// View returns the view string for the application model.
func (m *Model) View() string {
	return lipgloss.JoinHorizontal(
		lipgloss.Left,
		m.FoldersStyle.Base.Render(m.Folders.View()),
		m.ListStyle.Base.Render(m.List.View()),
		m.ContentStyle.Base.Render(m.Code.View()),
	)
}
