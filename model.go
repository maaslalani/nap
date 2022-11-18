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
	"golang.org/x/exp/maps"
	"golang.org/x/exp/slices"
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

type input int

const (
	folderInput input = iota
	nameInput
	languageInput
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
	// the List of snippets to display to the user.
	Lists map[Folder]*list.Model
	// the list of Folders to display to the user.
	Folders list.Model
	// the viewport of the Code snippet.
	Code        viewport.Model
	LineNumbers viewport.Model
	// the input for snippet folder, name, language
	activeInput input
	inputs      []textinput.Model
	tagsInput   textinput.Model
	// the current active pane of focus.
	pane pane
	// the current state / action of the application.
	state state
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

type updateFoldersMsg struct {
	items               []list.Item
	selectedFolderIndex int
}

// updateFolders returns a Cmd to  tell the application that there are possible
// folder changes to update.
func (m *Model) updateFolders() tea.Cmd {
	return func() tea.Msg {
		msg := m.updateFoldersView()
		return msg
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
	case updateFoldersMsg:
		setItemsCmd := m.Folders.SetItems(msg.items)
		m.Folders.Select(msg.selectedFolderIndex)
		var cmd tea.Cmd
		m.Folders, cmd = m.Folders.Update(msg)
		return m, tea.Batch(setItemsCmd, cmd)
	case updateContentMsg:
		return m.updateContentView(msg)
	case changeStateMsg:
		m.List().SetDelegate(snippetDelegate{m.ListStyle, msg.newState})

		var cmd tea.Cmd

		if m.state == msg.newState {
			break
		}

		wasEditing := m.state == editingState
		wasPasting := m.state == pastingState
		wasCreating := m.state == creatingState
		m.state = msg.newState
		m.updateKeyMap()
		m.updateActivePane(msg)

		switch msg.newState {
		case navigatingState:
			if wasPasting || wasCreating {
				return m, m.updateContent()
			}

			if wasEditing {
				m.blurInputs()
				i := m.List().Index()
				snippet := m.selectedSnippet()
				if m.inputs[nameInput].Value() != "" {
					snippet.Name = m.inputs[nameInput].Value()
				} else {
					snippet.Name = defaultSnippetName
				}
				if m.inputs[folderInput].Value() != "" {
					snippet.Folder = m.inputs[folderInput].Value()
				} else {
					snippet.Folder = defaultSnippetFolder
				}
				if m.inputs[languageInput].Value() != "" {
					snippet.Language = m.inputs[languageInput].Value()
				} else {
					snippet.Language = m.config.DefaultLanguage
				}
				newFile := fmt.Sprintf("%s-%s.%s", snippet.Folder, snippet.Name, snippet.Language)
				_ = os.Rename(m.selectedSnippetFilePath(), filepath.Join(m.config.Home, newFile))
				snippet.File = newFile
				setCmd := m.List().SetItem(i, snippet)
				m.pane = snippetPane
				cmd = tea.Batch(setCmd, m.updateFolders(), m.updateContent())
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
			m.state = deletingState
		case editingState:
			m.pane = contentPane
			snippet := m.selectedSnippet()
			m.inputs[folderInput].SetValue(snippet.Folder)
			if snippet.Name == defaultSnippetName {
				m.inputs[nameInput].SetValue("")
			} else {
				m.inputs[nameInput].SetValue(snippet.Name)
			}
			m.inputs[languageInput].SetValue(snippet.Language)
			cmd = m.focusInput(m.activeInput)
		case creatingState:
		case copyingState:
			m.pane = snippetPane
			m.state = copyingState
			m.updateActivePane(msg)
			cmd = tea.Tick(time.Second, func(t time.Time) tea.Msg {
				return changeStateMsg{navigatingState}
			})
		}

		m.updateKeyMap()
		m.updateActivePane(msg)
		return m, cmd
	case tea.WindowSizeMsg:
		m.height = msg.Height - 4
		for _, li := range m.Lists {
			li.SetHeight(m.height)
		}
		m.Folders.SetHeight(m.height)
		m.Code.Height = m.height
		m.LineNumbers.Height = m.height
		m.Code.Width = msg.Width - m.List().Width() - m.Folders.Width() - 20
		m.LineNumbers.Width = 5
		return m, nil
	case tea.KeyMsg:
		if m.List().FilterState() == list.Filtering {
			break
		}

		if m.state == deletingState {
			switch {
			case key.Matches(msg, m.keys.Confirm):
				m.List().RemoveItem(m.List().Index())
				m.state = navigatingState
				m.updateKeyMap()
				return m, tea.Batch(changeState(navigatingState), func() tea.Msg {
					return updateContentMsg(m.selectedSnippet())
				})
			case key.Matches(msg, m.keys.Quit, m.keys.Cancel):
				return m, changeState(navigatingState)
			}
			return m, nil
		} else if m.state == copyingState {
			return m, changeState(navigatingState)
		} else if m.state == editingState {
			if msg.String() == "esc" || msg.String() == "enter" {
				return m, changeState(navigatingState)
			}
			var cmd tea.Cmd
			var cmds []tea.Cmd
			for i := range m.inputs {
				m.inputs[i], cmd = m.inputs[i].Update(msg)
				cmds = append(cmds, cmd)
			}
			return m, tea.Batch(cmds...)
		}

		switch {
		case key.Matches(msg, m.keys.NextPane):
			m.nextPane()
		case key.Matches(msg, m.keys.PreviousPane):
			m.previousPane()
		case key.Matches(msg, m.keys.Quit):
			m.state = quittingState
			return m, tea.Quit
		case key.Matches(msg, m.keys.NewSnippet):
			m.state = creatingState
			return m, m.createNewSnippetFile()
		case key.Matches(msg, m.keys.PasteSnippet):
			return m, changeState(pastingState)
		case key.Matches(msg, m.keys.RenameSnippet):
			m.activeInput = nameInput
			return m, changeState(editingState)
		case key.Matches(msg, m.keys.ChangeFolder):
			m.pane = snippetPane
			cmd := m.updateActivePane(msg)
			return m, cmd
		case key.Matches(msg, m.keys.ToggleHelp):
			m.help.ShowAll = !m.help.ShowAll

			var newHeight int
			if m.help.ShowAll {
				newHeight = m.height - 4
			} else {
				newHeight = m.height
			}
			m.List().SetHeight(newHeight)
			m.Folders.SetHeight(newHeight)
			m.Code.Height = newHeight
			m.LineNumbers.Height = newHeight
		case key.Matches(msg, m.keys.SetFolder):
			m.activeInput = folderInput
			return m, changeState(editingState)
		case key.Matches(msg, m.keys.SetLanguage):
			m.activeInput = languageInput
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
			m.List().Title = "Delete? (y/N)"
			return m, changeState(deletingState)
		case key.Matches(msg, m.keys.EditSnippet):
			return m, m.editSnippet()
		}
	}

	m.updateKeyMap()
	cmd := m.updateActivePane(msg)
	return m, cmd
}

// blurInputs blurs all the inputs.
func (m *Model) blurInputs() {
	for i := range m.inputs {
		m.inputs[i].Blur()
	}
}

// focusInput focuses the speficied input and blurs the rest.
func (m *Model) focusInput(i input) tea.Cmd {
	m.blurInputs()
	m.inputs[i].CursorEnd()
	return m.inputs[i].Focus()
}

// selectedSnippetFilePath returns the file path of the snippet that is
// currently selected.
func (m *Model) selectedSnippetFilePath() string {
	return filepath.Join(m.config.Home, m.selectedSnippet().File)
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
	}
}

// updateFolderView updates the folders list to display the current folders.
func (m *Model) updateFoldersView() tea.Msg {
	var selectedFolder Folder
	selectedFolderIndex := m.Folders.Index()
	for folder, li := range m.Lists {
		for i, item := range li.Items() {
			snippet, ok := item.(Snippet)
			if !ok {
				continue
			}
			f := Folder(snippet.Folder)
			_, ok = m.Lists[f]
			if !ok {
				m.Lists[f] = newList([]list.Item{}, m.height, m.ListStyle)
				selectedFolder = f
			}
			if f != folder {
				li.RemoveItem(i)
				m.Lists[f].InsertItem(0, item)
				selectedFolder = f
			}
		}
	}
	var folderItems []list.Item

	foldersSlice := maps.Keys(m.Lists)
	slices.Sort(foldersSlice)
	for i, folder := range foldersSlice {
		folderItems = append(folderItems, Folder(folder))
		if folder == selectedFolder {
			selectedFolderIndex = i
		}
	}

	return updateFoldersMsg{
		items:               folderItems,
		selectedFolderIndex: selectedFolderIndex,
	}
}

// updateContentView updates the content view with the correct content based on
// the active snippet or display the appropriate error message / hint message.
func (m *Model) updateContentView(msg updateContentMsg) (tea.Model, tea.Cmd) {
	if len(m.List().Items()) <= 0 {
		m.displayKeyHint([]keyHint{
			{m.keys.NewSnippet, "create a new snippet."},
		})
		return m, nil
	}

	var b bytes.Buffer
	content, err := os.ReadFile(filepath.Join(m.config.Home, msg.File))
	if err != nil {
		m.displayKeyHint(m.noContentHints())
		return m, nil
	}

	if string(content) == "" {
		m.displayKeyHint(m.noContentHints())
		return m, nil
	}

	// b.WriteString(string(content))
	err = quick.Highlight(&b, string(content), msg.Language, "terminal16m", m.config.Theme)
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
		m.ListStyle = DefaultStyles(m.config).Snippets.Blurred
		m.ContentStyle = DefaultStyles(m.config).Content.Blurred
		m.FoldersStyle = DefaultStyles(m.config).Folders.Focused
		m.Folders, cmd = m.Folders.Update(msg)
		m.updateKeyMap()
		cmds = append(cmds, cmd, m.updateContent())
	case snippetPane:
		m.ListStyle = DefaultStyles(m.config).Snippets.Focused
		m.ContentStyle = DefaultStyles(m.config).Content.Blurred
		m.FoldersStyle = DefaultStyles(m.config).Folders.Blurred
		*m.List(), cmd = (*m.List()).Update(msg)
		cmds = append(cmds, cmd)
	case contentPane:
		m.ListStyle = DefaultStyles(m.config).Snippets.Blurred
		m.ContentStyle = DefaultStyles(m.config).Content.Focused
		m.FoldersStyle = DefaultStyles(m.config).Folders.Blurred
		m.Code, cmd = m.Code.Update(msg)
		cmds = append(cmds, cmd)
		m.LineNumbers, cmd = m.LineNumbers.Update(msg)
		cmds = append(cmds, cmd)
	}
	m.List().SetDelegate(snippetDelegate{m.ListStyle, m.state})
	m.Folders.Styles.TitleBar = m.FoldersStyle.TitleBar
	m.Folders.Styles.Title = m.FoldersStyle.Title

	return tea.Batch(cmds...)
}

// updateKeyMap disables or enables the keys based on the current state of the
// snippet list.
func (m *Model) updateKeyMap() {
	hasItems := len(m.List().VisibleItems()) > 0
	isFiltering := m.List().FilterState() == list.Filtering
	isEditing := m.state == editingState
	m.keys.DeleteSnippet.SetEnabled(hasItems && !isFiltering && !isEditing)
	m.keys.CopySnippet.SetEnabled(hasItems && !isFiltering && !isEditing)
	m.keys.PasteSnippet.SetEnabled(hasItems && !isFiltering && !isEditing)
	m.keys.EditSnippet.SetEnabled(hasItems && !isFiltering && !isEditing)
	m.keys.NewSnippet.SetEnabled(!isFiltering && !isEditing)
	m.keys.ChangeFolder.SetEnabled(m.pane == folderPane)
}

// selectedSnippet returns the currently selected snippet.
func (m *Model) selectedSnippet() Snippet {
	item := m.List().SelectedItem()
	if item == nil {
		return defaultSnippet
	}
	return item.(Snippet)
}

// selected folder returns the currently selected folder.
func (m *Model) selectedFolder() Folder {
	item := m.Folders.SelectedItem()
	if item == nil {
		return "misc"
	}
	return item.(Folder)
}

// List returns the active list.
func (m *Model) List() *list.Model {
	return m.Lists[m.selectedFolder()]
}

// createNewSnippet creates a new snippet file and adds it to the the list.
func (m *Model) createNewSnippetFile() tea.Cmd {
	return func() tea.Msg {
		folder := defaultSnippetFolder
		folderItem := m.Folders.SelectedItem()
		if folderItem != nil && folderItem.FilterValue() != "" {
			folder = folderItem.FilterValue()
		}

		file := fmt.Sprintf("%s-snippet-%d.%s", folder, rand.Intn(1000000), m.config.DefaultLanguage)
		_, _ = os.Create(filepath.Join(m.config.Home, file))

		newSnippet := Snippet{
			Name:     defaultSnippetName,
			Date:     time.Now(),
			File:     file,
			Language: m.config.DefaultLanguage,
			Tags:     []string{},
			Folder:   folder,
		}

		m.List().InsertItem(m.List().Index(), newSnippet)
		return changeStateMsg{navigatingState}
	}
}

// View returns the view string for the application model.
func (m *Model) View() string {
	if m.state == quittingState {
		return ""
	}

	var (
		folder   = m.ContentStyle.Title.Render(m.selectedSnippet().Folder)
		name     = m.ContentStyle.Title.Render(m.selectedSnippet().Name)
		language = m.ContentStyle.Title.Render(m.selectedSnippet().Language)
		titleBar = m.ListStyle.TitleBar.Render("Snippets")
	)

	if m.state == editingState {
		folder = m.inputs[folderInput].View()
		name = m.inputs[nameInput].View()
		language = m.inputs[languageInput].View()
	} else if m.state == copyingState {
		titleBar = m.ListStyle.CopiedTitleBar.Render("Copied Snippet!")
	} else if m.state == deletingState {
		titleBar = m.ListStyle.DeletedTitleBar.Render("Delete Snippet? (y/N)")
	} else if m.List().SettingFilter() {
		titleBar = m.ListStyle.TitleBar.Render(m.List().FilterInput.View())
	}

	return lipgloss.JoinVertical(
		lipgloss.Top,
		lipgloss.JoinHorizontal(
			lipgloss.Left,
			m.FoldersStyle.Base.Render(m.Folders.View()),
			m.ListStyle.Base.Render(titleBar+m.List().View()),
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
		marginStyle.Render(m.help.View(m.keys)),
	)
}
