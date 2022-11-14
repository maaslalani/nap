package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/caarlos0/env/v6"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"golang.org/x/exp/maps"
)

var defaultSnippetFileFormat = `[ { "folder": "%s", "title": "%s", "tags": [], "date": "2022-11-12T15:04:05Z", "favorite": false, "file": "snooze.txt", "language": "%s" } ]`

func main() {
	config := Config{Home: defaultHome()}
	if err := env.Parse(&config); err != nil {
		fmt.Println("Unable to unmarshal config", err)
	}

	var snippets []Snippet
	file := filepath.Join(config.Home, config.File)
	dir, err := os.ReadFile(file)
	if err != nil {
		// File does not exist, create one.
		err := os.MkdirAll(filepath.Join(config.Home, defaultSnippetFolder), os.ModePerm)
		if err != nil {
			fmt.Printf("Unable to create directory %s, %+v", config.Home, err)
		}
		f, err := os.Create(file)
		if err != nil {
			fmt.Printf("Unable to create file %s, %+v", file, err)
		}
		content := fmt.Sprintf(defaultSnippetFileFormat, defaultSnippetFolder, defaultSnippetName, config.DefaultLanguage)
		_, _ = f.WriteString(content)
		dir = []byte(content)
	}
	err = json.Unmarshal(dir, &snippets)
	if err != nil {
		fmt.Printf("Unable to unmarshal %s file, %+v\n", file, err)
		return
	}

	var folders = make(map[string]int)
	var items []list.Item
	for _, snippet := range snippets {
		folders[snippet.Folder]++
		items = append(items, list.Item(snippet))
	}
	snippetList := list.New(items, snippetDelegate{}, 0, 0)

	var folderItems []list.Item
	for _, folder := range maps.Keys(folders) {
		folderItems = append(folderItems, list.Item(Folder(folder)))
	}
	folderList := list.New(folderItems, folderDelegate{}, 0, 0)
	folderList.Title = "Folders"

	folderList.SetShowHelp(false)
	folderList.SetFilteringEnabled(false)
	folderList.SetShowStatusBar(false)
	folderList.DisableQuitKeybindings()
	folderList.Styles.NoItems = lipgloss.NewStyle().Margin(0, 2).Foreground(gray)
	folderList.SetStatusBarItemName("folder", "folders")

	snippetList.SetShowHelp(false)
	snippetList.SetShowFilter(true)
	snippetList.Title = "Snippets"

	snippetList.FilterInput.Prompt = "Find: "
	snippetList.FilterInput.CursorStyle = lipgloss.NewStyle().Foreground(primaryColor)
	snippetList.FilterInput.PromptStyle = lipgloss.NewStyle().Foreground(white).MarginLeft(1)
	snippetList.FilterInput.TextStyle = lipgloss.NewStyle().Foreground(white).Background(primaryColorSubdued)
	snippetList.Styles.NoItems = lipgloss.NewStyle().Margin(0, 2).Foreground(gray)
	snippetList.SetStatusBarItemName("snippet", "snippets")
	snippetList.DisableQuitKeybindings()

	content := viewport.New(80, 0)

	m := &Model{
		List:         snippetList,
		Folders:      folderList,
		Code:         content,
		ContentStyle: DefaultStyles.Content.Blurred,
		ListStyle:    DefaultStyles.Snippets.Focused,
		FoldersStyle: DefaultStyles.Folders.Blurred,
		keys:         DefaultKeyMap,
		help:         help.New(),
		config:       config,
		inputs:       []textinput.Model{newTextInput(), newTextInput(), newTextInput()},
		tagsInput:    newTextInput(),
	}
	p := tea.NewProgram(m, tea.WithAltScreen())
	model, err := p.Run()
	if err != nil {
		fmt.Println("Alas, there was an error.", err)
		return
	}

	fm, ok := model.(*Model)
	if !ok {
		fmt.Println("Alas, there was an error.", err)
		return
	}

	b, err := json.Marshal(fm.List.Items())
	if err != nil {
		fmt.Println("Could not mashal latest snippet data.", err)
		return
	}
	err = os.WriteFile(filepath.Join(config.Home, config.File), b, os.ModePerm)
	if err != nil {
		fmt.Println("Could not save snippets file.", err)
		return
	}
}

func newTextInput() textinput.Model {
	i := textinput.New()
	i.Prompt = ""
	i.PromptStyle = lipgloss.NewStyle().Margin(0, 1)
	i.TextStyle = lipgloss.NewStyle().MarginBottom(1)
	i.CursorStyle = lipgloss.NewStyle().Foreground(primaryColor)
	return i
}
