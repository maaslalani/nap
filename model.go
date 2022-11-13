package main

import (
	"bytes"
	"fmt"
	"os"
	"strings"

	"github.com/alecthomas/chroma/quick"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
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
	// the key map.
	keys KeyMap
	// the help model.
	help help.Model
	// the height of the terminal.
	height int
	// the working directory.
	Workdir string
	// code Snippets.
	Snippets []Snippet
	// the List of snippets to display to the user.
	List list.Model
	// the list of Folders to display to the user.
	Folders list.Model
	// the viewport of the Code snippet.
	Code        viewport.Model
	LineNumbers viewport.Model
	// the current active pane of focus.
	Pane int
	// the current snippet being displayed.
	ActiveSnippet Snippet

	// stying for components
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
		m.ActiveSnippet = Snippet(msg)
		content, err := os.ReadFile(configDir + "/" + msg.Folder + "/" + msg.File)
		if err != nil {
			m.LineNumbers.SetContent(" ~ ")
			m.Code.SetContent(b.String() + "Error: unable to read file.")
			return m, nil
		}
		err = quick.Highlight(&b, string(content), msg.Language, "terminal16m", "dracula")
		var lineNumbers strings.Builder
		height := lipgloss.Height(b.String())
		for i := 0; i < height; i++ {
			lineNumbers.WriteString(fmt.Sprintf("%3d │ \n", i+1))
		}
		m.LineNumbers.SetContent(lineNumbers.String())
		m.Code.SetContent(b.String())
		return m, nil
	case tea.WindowSizeMsg:
		m.height = msg.Height
		m.List.SetHeight(msg.Height - 4)
		m.Folders.SetHeight(msg.Height - 4)
		m.Code.Height = msg.Height - 4
		m.Code.Width = msg.Width - m.List.Width() - m.Folders.Width() - 20
		m.LineNumbers.Height = msg.Height - 4
		m.LineNumbers.Width = 5
		return m, nil
	case tea.KeyMsg:
		if m.List.FilterState() == list.Filtering {
			break
		}
		switch {
		case key.Matches(msg, m.keys.NextPane):
			m.nextPane()
		case key.Matches(msg, m.keys.PreviousPane):
			m.previousPane()
		case key.Matches(msg, m.keys.Quit):
			return m, tea.Quit
		}
	}

	cmd := m.updateActivePane(msg)
	return m, cmd
}

// nextPane sets the next pane to be active.
func (m *Model) nextPane() {
	m.Pane = (m.Pane + 1) % MaxPane
}

// previousPane sets the previous pane to be active.
func (m *Model) previousPane() {
	m.Pane--
	if m.Pane < 0 {
		m.Pane = MaxPane - 1
	}
}

const tabSpaces = 4

// updateActivePane updates the currently active pane.
func (m *Model) updateActivePane(msg tea.Msg) tea.Cmd {
	var cmds []tea.Cmd
	var cmd tea.Cmd
	switch m.Pane {
	case FolderPane:
		m.ListStyle = DefaultStyles.Snippets.Blurred
		m.ContentStyle = DefaultStyles.Content.Blurred
		m.FoldersStyle = DefaultStyles.Folders.Focused
		m.Folders, cmd = m.Folders.Update(msg)
		cmds = append(cmds, cmd)
	case SnippetPane:
		m.ListStyle = DefaultStyles.Snippets.Focused
		m.ContentStyle = DefaultStyles.Content.Blurred
		m.FoldersStyle = DefaultStyles.Folders.Blurred
		m.List, cmd = m.List.Update(msg)
		cmds = append(cmds, cmd)
	case ContentPane:
		m.ListStyle = DefaultStyles.Snippets.Blurred
		m.ContentStyle = DefaultStyles.Content.Focused
		m.FoldersStyle = DefaultStyles.Folders.Blurred
		m.Code, cmd = m.Code.Update(msg)
		cmds = append(cmds, cmd)
		m.LineNumbers, cmd = m.LineNumbers.Update(msg)
		cmds = append(cmds, cmd)
	}
	m.List.Styles.TitleBar = m.ListStyle.TitleBar
	m.List.Styles.Title = m.ListStyle.Title
	m.Folders.Styles.TitleBar = m.FoldersStyle.TitleBar
	m.Folders.Styles.Title = m.FoldersStyle.Title

	return tea.Batch(cmds...)
}

// View returns the view string for the application model.
func (m *Model) View() string {
	return lipgloss.JoinVertical(
		lipgloss.Top,
		lipgloss.JoinHorizontal(
			lipgloss.Left,
			m.FoldersStyle.Base.Render(m.Folders.View()),
			m.ListStyle.Base.Render(m.List.View()),
			lipgloss.JoinVertical(lipgloss.Top,
				m.ContentStyle.Title.Render(m.ActiveSnippet.Title),
				lipgloss.JoinHorizontal(lipgloss.Left,
					m.ContentStyle.LineNumber.Render(m.LineNumbers.View()),
					m.ContentStyle.Base.Render(strings.ReplaceAll(m.Code.View(), "\t", strings.Repeat(" ", tabSpaces))),
				),
			),
		),
		"  "+m.help.View(m.keys),
	)
}
