package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/alecthomas/chroma/quick"
	"github.com/atotto/clipboard"
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

const (
	NavigatingState int = iota
	DeletingState
	CreatingState
	CopyingState
	QuittingState
)

// Model represents the state of the application.
// It contains all the snippets organized in folders.
type Model struct {
	// the config map.
	config Config
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
	// the current state / action of the application.
	State int
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
	m.updateKeyMap()
	return func() tea.Msg {
		return updateViewMsg(m.selectedSnippet())
	}
}

// updateViewMsg tells the application to update the content view with the
// given snippet.
type updateViewMsg Snippet

// Update updates the model based on user interaction.
func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case updateViewMsg:
		return m.updateContentView(msg)
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
		if m.State == DeletingState {
			switch {
			case key.Matches(msg, m.keys.Confirm):
				m.List.RemoveItem(m.List.Index())
				m.resetTitleBar()
				m.State = NavigatingState
				m.updateKeyMap()
				return m, func() tea.Msg {
					return updateViewMsg(m.selectedSnippet())
				}
			case key.Matches(msg, m.keys.Quit, m.keys.Cancel):
				m.resetTitleBar()
				m.State = NavigatingState
			}
			return m, nil
		} else if m.State == CopyingState {
			m.resetTitleBar()
			m.State = NavigatingState
			break
		}
		switch {
		case key.Matches(msg, m.keys.NextPane):
			m.nextPane()
		case key.Matches(msg, m.keys.PreviousPane):
			m.previousPane()
		case key.Matches(msg, m.keys.Quit):
			m.State = QuittingState
			return m, tea.Quit
		case key.Matches(msg, m.keys.NewSnippet):
			m.State = CreatingState
			folder := defaultFolder
			folderItem := m.Folders.SelectedItem()
			if folderItem != nil && folderItem.FilterValue() != "" {
				folder = folderItem.FilterValue()
			}
			rand.Seed(time.Now().Unix())
			file := fmt.Sprintf("snooze-%d.go", rand.Intn(1000000))
			_, _ = os.Create(filepath.Join(m.config.Home, folder, file))
			m.List.InsertItem(m.List.Index(), Snippet{Title: "Untitled Snippet", Date: time.Now(), File: file, Language: "Go", Tags: []string{}, Folder: folder})
		case key.Matches(msg, m.keys.CopySnippet):
			m.State = CopyingState
			content, err := os.ReadFile(m.selectedSnippetFilePath())
			if err != nil {
				return m, nil
			}
			clipboard.WriteAll(string(content))
			m.List.Styles.TitleBar.Background(green)
			m.List.Title = "Copied " + m.selectedSnippet().Title + "!"
			m.ListStyle.SelectedTitle.Foreground(brightGreen)
			m.ListStyle.SelectedSubtitle.Foreground(green)
			return m, tea.Tick(time.Second/2, func(t time.Time) tea.Msg {
				if m.State == CopyingState {
					m.resetTitleBar()
				}
				return tea.KeyMsg{}
			})
		case key.Matches(msg, m.keys.DeleteSnippet):
			m.Pane = SnippetPane
			m.updateActivePane(msg)
			m.State = DeletingState
			m.List.Styles.TitleBar.Background(red)
			m.List.Title = "Delete snippet? (y/N)"
			m.ListStyle.SelectedTitle.Foreground(brightRed)
			m.ListStyle.SelectedSubtitle.Foreground(red)
		case key.Matches(msg, m.keys.EditSnippet):
			return m, m.editSnippet()
		}
	}

	m.updateKeyMap()
	cmd := m.updateActivePane(msg)
	return m, cmd
}

func (m *Model) selectedSnippetFilePath() string {
	return filepath.Join(m.config.Home, m.selectedSnippet().Folder, m.selectedSnippet().File)
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

// editSnippet opens the editor with the selected snippet file path.
func (m *Model) editSnippet() tea.Cmd {
	editor := os.Getenv("EDITOR")
	if editor == "" {
		editor = "vim"
	}
	cmd := exec.Command(editor, m.selectedSnippetFilePath())
	return tea.ExecProcess(cmd, func(err error) tea.Msg {
		return updateViewMsg(m.selectedSnippet())
	})
}

// updateContentView updates the content view with the correct content based on
// the active snippet or display the appropriate error message / hint message.
func (m *Model) updateContentView(msg updateViewMsg) (tea.Model, tea.Cmd) {
	if len(m.List.Items()) <= 0 {
		m.displayKeyHint("No Snippets.", "Press", "n", "to create a new snippet.")
		return m, nil
	}

	var b bytes.Buffer
	content, err := os.ReadFile(filepath.Join(m.config.Home, msg.Folder, msg.File))
	if err != nil {
		m.displayKeyHint("No Content.", "Press", "e", "to edit snippet.")
		return m, nil
	}

	if string(content) == "" {
		m.displayKeyHint("No Content.", "Press", "e", "to edit snippet.")
		return m, nil
	}

	err = quick.Highlight(&b, string(content), msg.Language, "terminal16m", "dracula")
	if err != nil {
		m.displayError("Unable to highlight file.")
		return m, nil
	}

	s := b.String()
	m.writeLineNumbers(lipgloss.Height(s))
	m.Code.SetContent(s)
	return m, nil
}

// displayKeyHint updates the content viewport with instructions on the
// relevent key binding that the user should most likely press.
func (m *Model) displayKeyHint(title, prefix, key, suffix string) {
	m.LineNumbers.SetContent(" ~ \n ~ ")
	m.Code.SetContent(fmt.Sprintf("%s\n%s %s %s",
		m.ContentStyle.EmptyHint.Render(title),
		m.ContentStyle.EmptyHint.Render(prefix),
		m.ContentStyle.EmptyHintKey.Render(key),
		m.ContentStyle.EmptyHint.Render(suffix),
	))
}

// displayError updates the content viewport with the error message provided.
func (m *Model) displayError(error string) {
	m.LineNumbers.SetContent(" ~ ")
	m.Code.SetContent(fmt.Sprintf("%s",
		m.ContentStyle.EmptyHint.Render(error),
	))
}

// writeLineNumbers writes the number of line numbers to the line number
// viewport.
func (m *Model) writeLineNumbers(n int) {
	var lineNumbers strings.Builder
	for i := 1; i < n; i++ {
		lineNumbers.WriteString(fmt.Sprintf("%3d \n", i))
	}
	m.LineNumbers.SetContent(lineNumbers.String() + "  ~ \n")
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

// resetTitleBar resets the title bar to original navigating state.
func (m *Model) resetTitleBar() {
	m.List.Styles.TitleBar.Background(primaryColorSubdued)
	m.ListStyle.SelectedTitle.Foreground(primaryColor)
	m.ListStyle.SelectedSubtitle.Foreground(primaryColorSubdued)
	m.List.Title = "Snippets"
}

// updateKeyMap disables or enables the keys based on the current state of the
// snippet list.
func (m *Model) updateKeyMap() {
	hasItems := len(m.List.VisibleItems()) > 0
	isFiltering := m.List.FilterState() == list.Filtering
	m.keys.DeleteSnippet.SetEnabled(hasItems && !isFiltering)
	m.keys.CopySnippet.SetEnabled(hasItems && !isFiltering)
	m.keys.EditSnippet.SetEnabled(hasItems && !isFiltering)
	m.keys.NewSnippet.SetEnabled(!isFiltering)
}

// selectedSnippet returns the currently selected snippet.
func (m *Model) selectedSnippet() Snippet {
	item := m.List.SelectedItem()
	if item == nil {
		return Snippet{Title: "No Snippets", Folder: defaultFolder}
	}
	return item.(Snippet)
}

// View returns the view string for the application model.
func (m *Model) View() string {
	if m.State == QuittingState {
		return ""
	}
	return lipgloss.JoinVertical(
		lipgloss.Top,
		lipgloss.JoinHorizontal(
			lipgloss.Left,
			m.FoldersStyle.Base.Render(m.Folders.View()),
			m.ListStyle.Base.Render(m.List.View()),
			lipgloss.JoinVertical(lipgloss.Top,
				lipgloss.JoinHorizontal(lipgloss.Left,
					m.ContentStyle.Title.Render(m.selectedSnippet().Folder),
					m.ContentStyle.Separator.Render("/"),
					m.ContentStyle.Title.Render(m.selectedSnippet().Title)),
				lipgloss.JoinHorizontal(lipgloss.Left,
					m.ContentStyle.LineNumber.Render(m.LineNumbers.View()),
					m.ContentStyle.Base.Render(strings.ReplaceAll(m.Code.View(), "\t", strings.Repeat(" ", tabSpaces))),
				),
			),
		),
		"  "+m.help.View(m.keys),
	)
}
