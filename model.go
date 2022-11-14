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
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const maxPane = 3

type pane int

const (
	snippetPane pane = iota
	contentPane
	folderPane
)

type state int

const (
	navigatingState state = iota
	deletingState
	creatingState
	copyingState
	pastingState
	quittingState
	editingState
	editingTagsState
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
	// the input for snippet name and folder
	folderInput   textinput.Model
	titleInput    textinput.Model
	languageInput textinput.Model
	tagsInput     textinput.Model
	// the current active pane of focus.
	pane pane
	// the current state / action of the application.
	State state
	// stying for components
	ListStyle    SnippetsBaseStyle
	FoldersStyle FoldersBaseStyle
	ContentStyle ContentBaseStyle
}

// Init initialzes the application model.
func (m *Model) Init() tea.Cmd {
	rand.Seed(time.Now().Unix())

	m.Folders.Styles.Title = m.FoldersStyle.Title
	m.Folders.Styles.TitleBar = m.FoldersStyle.TitleBar
	m.List.Styles.Title = m.ListStyle.Title
	m.List.Styles.TitleBar = m.ListStyle.TitleBar

	m.updateKeyMap()

	return func() tea.Msg {
		return updateContentMsg(m.selectedSnippet())
	}
}

// updateContentMsg tells the application to update the content view with the
// given snippet.
type updateContentMsg Snippet

// updateContent instructs the application to fetch the latest contents of the
// snippet file.
//
// This is useful after a Paste or Edit.
func (m *Model) updateContent() tea.Cmd {
	return func() tea.Msg {
		return updateContentMsg(m.selectedSnippet())
	}
}

// changeStateMsg tells the application to enter a different state.
type changeStateMsg struct{ newState state }

// changeState returns a Cmd to enter a different state.
func changeState(newState state) tea.Cmd {
	return func() tea.Msg {
		return changeStateMsg{newState}
	}
}

// Update updates the model based on user interaction.
func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case updateContentMsg:
		return m.updateContentView(msg)
	case changeStateMsg:
		var cmd tea.Cmd

		if m.State == msg.newState {
			break
		}

		wasEditing := m.State == editingState
		wasPasting := m.State == pastingState
		wasCreating := m.State == creatingState
		m.State = msg.newState
		m.resetTitleBar()
		m.updateKeyMap()
		m.updateActivePane(msg)

		switch msg.newState {
		case navigatingState:
			if wasPasting || wasCreating {
				return m, m.updateContent()
			}

			if wasEditing {
				m.folderInput.Blur()
				m.titleInput.Blur()
				m.languageInput.Blur()
				i := m.List.Index()
				snippet := m.List.SelectedItem().(Snippet)
				m.List.RemoveItem(i)
				if m.titleInput.Value() != "" {
					snippet.Name = m.titleInput.Value()
				} else {
					snippet.Name = defaultSnippetName
				}
				if m.folderInput.Value() != "" {
					snippet.Folder = m.folderInput.Value()
				} else {
					snippet.Folder = defaultSnippetFolder
				}
				snippet.Language = m.languageInput.Value()
				m.List.InsertItem(i, snippet)
			}
		case pastingState:
			content, err := clipboard.ReadAll()
			if err != nil {
				return m, changeState(navigatingState)
			}
			f, err := os.OpenFile(m.selectedSnippetFilePath(), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
			if err != nil {
				return m, changeState(navigatingState)
			}
			f.WriteString(content)
			return m, changeState(navigatingState)
		case deletingState:
		case editingState:
			m.pane = contentPane
			snippet := m.selectedSnippet()
			m.folderInput.SetValue(snippet.Folder)
			if snippet.Name == defaultSnippetName {
				m.titleInput.SetValue("")
				// We add a space at the end because the cursor for the placeholder will be at
				// the beginning and we need to add some margin at the end.
				m.titleInput.Placeholder = defaultSnippetName + " "
			} else {
				m.titleInput.SetValue(snippet.Name)
			}
			m.languageInput.SetValue(snippet.Language)
			m.folderInput.Blur()
			m.languageInput.Blur()
			m.titleInput.CursorEnd()
			cmd = m.titleInput.Focus()
		case creatingState:
		case copyingState:
			m.pane = snippetPane
			m.updateActivePane(msg)
			m.List.Styles.TitleBar.Background(green)
			m.List.Title = "Copied " + m.selectedSnippet().Name + "!"
			m.ListStyle.SelectedTitle.Foreground(brightGreen)
			m.ListStyle.SelectedSubtitle.Foreground(green)
			return m, tea.Tick(time.Second, func(time.Time) tea.Msg { return changeStateMsg{navigatingState} })
		}

		m.updateKeyMap()
		m.updateActivePane(msg)

		return m, cmd
	case tea.WindowSizeMsg:
		m.height = msg.Height - 4
		m.List.SetHeight(m.height)
		m.Folders.SetHeight(m.height)
		m.Code.Height = m.height
		m.Code.Width = msg.Width - m.List.Width() - m.Folders.Width() - 20
		m.LineNumbers.Height = m.height
		m.LineNumbers.Width = 5
		return m, nil
	case tea.KeyMsg:
		if m.List.FilterState() == list.Filtering {
			break
		}

		if m.State == deletingState {
			switch {
			case key.Matches(msg, m.keys.Confirm):
				m.List.RemoveItem(m.List.Index())
				m.resetTitleBar()
				m.State = navigatingState
				m.updateKeyMap()
				return m, func() tea.Msg {
					return updateContentMsg(m.selectedSnippet())
				}
			case key.Matches(msg, m.keys.Quit, m.keys.Cancel):
				m.resetTitleBar()
				m.State = navigatingState
			}
			return m, nil
		} else if m.State == copyingState {
			m.resetTitleBar()
			m.State = navigatingState
			break
		} else if m.State == editingState {
			if msg.String() == "esc" || msg.String() == "enter" {
				return m, func() tea.Msg {
					return changeStateMsg{navigatingState}
				}
			}
			var cmd tea.Cmd
			var cmds []tea.Cmd
			m.folderInput, cmd = m.folderInput.Update(msg)
			cmds = append(cmds, cmd)
			m.titleInput, cmd = m.titleInput.Update(msg)
			cmds = append(cmds, cmd)
			m.languageInput, cmd = m.languageInput.Update(msg)
			cmds = append(cmds, cmd)
			return m, tea.Batch(cmds...)
		}

		switch {
		case key.Matches(msg, m.keys.NextPane):
			m.nextPane()
		case key.Matches(msg, m.keys.PreviousPane):
			m.previousPane()
		case key.Matches(msg, m.keys.Quit):
			m.State = quittingState
			return m, tea.Quit
		case key.Matches(msg, m.keys.NewSnippet):
			m.State = creatingState
			return m, func() tea.Msg {
				folder := defaultSnippetFolder
				folderItem := m.Folders.SelectedItem()
				if folderItem != nil && folderItem.FilterValue() != "" {
					folder = folderItem.FilterValue()
				}
				file := fmt.Sprintf("snooze-%d.txt", rand.Intn(1000000))
				_, _ = os.Create(filepath.Join(m.config.Home, folder, file))
				m.List.InsertItem(m.List.Index(), Snippet{Name: defaultSnippetName, Date: time.Now(), File: file, Language: m.config.DefaultLanguage, Tags: []string{}, Folder: folder})
				return changeStateMsg{navigatingState}
			}
		case key.Matches(msg, m.keys.PasteSnippet):
			return m, changeState(pastingState)
		case key.Matches(msg, m.keys.RenameSnippet):
			return m, changeState(editingState)
		case key.Matches(msg, m.keys.CopySnippet):
			return m, func() tea.Msg {
				content, err := os.ReadFile(m.selectedSnippetFilePath())
				if err != nil {
					return changeStateMsg{navigatingState}
				}
				clipboard.WriteAll(string(content))
				return changeStateMsg{copyingState}
			}
		case key.Matches(msg, m.keys.DeleteSnippet):
			m.pane = snippetPane
			m.updateActivePane(msg)
			m.State = deletingState
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
	m.pane = (m.pane + 1) % maxPane
}

// previousPane sets the previous pane to be active.
func (m *Model) previousPane() {
	m.pane--
	if m.pane < 0 {
		m.pane = maxPane - 1
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
		return updateContentMsg(m.selectedSnippet())
	})
}

func (m *Model) noContentHints() []keyHint {
	return []keyHint{
		{m.keys.EditSnippet, "edit contents"},
		{m.keys.PasteSnippet, "paste clipboard"},
		{m.keys.RenameSnippet, "rename"},
		{m.keys.SetFolder, "set folder"},
		{m.keys.SetLanguage, "set language"},
		{m.keys.TagSnippet, "set tags"},
	}
}

// updateContentView updates the content view with the correct content based on
// the active snippet or display the appropriate error message / hint message.
func (m *Model) updateContentView(msg updateContentMsg) (tea.Model, tea.Cmd) {
	if len(m.List.Items()) <= 0 {
		m.displayKeyHint([]keyHint{
			{m.keys.NewSnippet, "create a new snippet."},
		})
		return m, nil
	}

	var b bytes.Buffer
	content, err := os.ReadFile(filepath.Join(m.config.Home, msg.Folder, msg.File))
	if err != nil {
		m.displayKeyHint(m.noContentHints())
		return m, nil
	}

	if string(content) == "" {
		m.displayKeyHint(m.noContentHints())
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

type keyHint struct {
	binding key.Binding
	help    string
}

// displayKeyHint updates the content viewport with instructions on the
// relevent key binding that the user should most likely press.
func (m *Model) displayKeyHint(hints []keyHint) {
	m.LineNumbers.SetContent(strings.Repeat("  ~ \n", len(hints)))
	var s strings.Builder
	for _, hint := range hints {
		s.WriteString(
			fmt.Sprintf("%s %s\n",
				m.ContentStyle.EmptyHintKey.Render(hint.binding.Help().Key),
				m.ContentStyle.EmptyHint.Render("• "+hint.help),
			))
	}
	m.Code.SetContent(s.String())
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
	switch m.pane {
	case folderPane:
		m.ListStyle = DefaultStyles.Snippets.Blurred
		m.ContentStyle = DefaultStyles.Content.Blurred
		m.FoldersStyle = DefaultStyles.Folders.Focused
		m.Folders, cmd = m.Folders.Update(msg)
		cmds = append(cmds, cmd)
	case snippetPane:
		m.ListStyle = DefaultStyles.Snippets.Focused
		m.ContentStyle = DefaultStyles.Content.Blurred
		m.FoldersStyle = DefaultStyles.Folders.Blurred
		m.List, cmd = m.List.Update(msg)
		cmds = append(cmds, cmd)
	case contentPane:
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
	m.List.Title = "Snippets"
	if m.pane == snippetPane {
		m.List.Styles.TitleBar.Background(primaryColorSubdued)
		m.ListStyle.SelectedTitle.Foreground(primaryColor)
		m.ListStyle.SelectedSubtitle.Foreground(primaryColorSubdued)
	}
}

// updateKeyMap disables or enables the keys based on the current state of the
// snippet list.
func (m *Model) updateKeyMap() {
	hasItems := len(m.List.VisibleItems()) > 0
	isFiltering := m.List.FilterState() == list.Filtering
	isEditing := m.State == editingState
	m.keys.DeleteSnippet.SetEnabled(hasItems && !isFiltering && !isEditing)
	m.keys.CopySnippet.SetEnabled(hasItems && !isFiltering && !isEditing)
	m.keys.PasteSnippet.SetEnabled(hasItems && !isFiltering && !isEditing)
	m.keys.EditSnippet.SetEnabled(hasItems && !isFiltering && !isEditing)
	m.keys.NewSnippet.SetEnabled(!isFiltering && !isEditing)
}

// selectedSnippet returns the currently selected snippet.
func (m *Model) selectedSnippet() Snippet {
	item := m.List.SelectedItem()
	if item == nil {
		return Snippet{Name: "No Snippets", Folder: defaultSnippetFolder, Language: m.config.DefaultLanguage}
	}
	return item.(Snippet)
}

// View returns the view string for the application model.
func (m *Model) View() string {
	if m.State == quittingState {
		return ""
	}

	var (
		folder   = m.ContentStyle.Title.Render(m.selectedSnippet().Folder)
		name     = m.ContentStyle.Title.Render(m.selectedSnippet().Name)
		language = m.ContentStyle.Title.Render(m.selectedSnippet().Language)
	)

	if m.State == editingState {
		folder = m.folderInput.View()
		name = m.titleInput.View()
		language = m.languageInput.View()
	}

	return lipgloss.JoinVertical(
		lipgloss.Top,
		lipgloss.JoinHorizontal(
			lipgloss.Left,
			m.FoldersStyle.Base.Render(m.Folders.View()),
			m.ListStyle.Base.Render(m.List.View()),
			lipgloss.JoinVertical(lipgloss.Top,
				lipgloss.JoinHorizontal(lipgloss.Left,
					folder,
					m.ContentStyle.Separator.Render("/"),
					name,
					m.ContentStyle.Separator.Render("."),
					language,
				),
				lipgloss.JoinHorizontal(lipgloss.Left,
					m.ContentStyle.LineNumber.Render(m.LineNumbers.View()),
					m.ContentStyle.Base.Render(strings.ReplaceAll(m.Code.View(), "\t", strings.Repeat(" ", tabSpaces))),
				),
			),
		),
		"  "+m.help.View(m.keys),
	)
}
