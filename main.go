package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/caarlos0/env/v6"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"golang.org/x/exp/maps"
)

const defaultSnippetFileContent = `[ { "folder": "", "title": "Untitled Snippet", "tags": [], "date": "2022-11-12T15:04:05Z", "favorite": false, "file": "snooze.txt", "language": "go" } ]`

func main() {
	config := Config{Home: defaultHome(), File: defaultFile()}
	if err := env.Parse(&config); err != nil {
		fmt.Println("Unable to unmarshal config", err)
	}

	var snippets []Snippet
	file := config.Home + "/" + config.File
	dir, err := os.ReadFile(file)
	if err != nil {
		// File does not exist, create one.
		err := os.MkdirAll(config.Home, os.ModePerm)
		if err != nil {
			fmt.Printf("Unable to create directory %s, %+v", config.Home, err)
		}
		f, err := os.Create(file)
		if err != nil {
			fmt.Printf("Unable to create file %s, %+v", file, err)
		}
		_, _ = f.WriteString(defaultSnippetFileContent)
		dir = []byte(defaultSnippetFileContent)
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

	snippetList.SetShowHelp(false)
	snippetList.SetShowFilter(true)
	snippetList.Title = "Snippets"

	snippetList.FilterInput.Prompt = "Find: "
	snippetList.FilterInput.CursorStyle = lipgloss.NewStyle().Foreground(primaryColor)
	snippetList.FilterInput.PromptStyle = lipgloss.NewStyle().Foreground(white).MarginLeft(1)
	snippetList.FilterInput.TextStyle = lipgloss.NewStyle().Foreground(white).Background(primaryColorSubdued)
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
	}

	b, err := json.Marshal(fm.List.Items())

	fmt.Println(string(b))
}
