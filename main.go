package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"golang.org/x/exp/maps"
)

const configDir = ".snooze"
const configFile = "snippets.json"

func main() {
	config, err := os.ReadFile(configDir + "/" + configFile)
	var snippets []Snippet
	err = json.Unmarshal(config, &snippets)
	if err != nil {
		fmt.Println("Unable to unmarshal snippets.json file", err)
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
	snippetList.SetShowHelp(false)
	snippetList.SetShowFilter(true)
	snippetList.Title = "Snippets"

	snippetList.FilterInput.Prompt = "Find: "
	snippetList.FilterInput.CursorStyle = lipgloss.NewStyle().Foreground(primaryColor)
	snippetList.FilterInput.PromptStyle = lipgloss.NewStyle().Foreground(white).MarginLeft(1)
	snippetList.FilterInput.TextStyle = lipgloss.NewStyle().Foreground(white).Background(primaryColorSubdued)

	content := viewport.New(80, 0)

	m := &Model{
		Snippets:     snippets,
		List:         snippetList,
		Folders:      folderList,
		Code:         content,
		ContentStyle: DefaultStyles.Content.Blurred,
		ListStyle:    DefaultStyles.Snippets.Focused,
		FoldersStyle: DefaultStyles.Folders.Blurred,
		keys:         DefaultKeyMap,
	}
	p := tea.NewProgram(m, tea.WithAltScreen())
	_, err = p.Run()
	if err != nil {
		fmt.Println("Alas, there was an error.", err)
		return
	}
}
