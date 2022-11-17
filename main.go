package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/caarlos0/env/v6"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/mattn/go-isatty"
	"github.com/sahilm/fuzzy"
	"golang.org/x/exp/maps"
	"golang.org/x/exp/slices"
)

var defaultSnippetFileFormat = `[ { "folder": "%s", "title": "%s", "tags": [], "date": "2022-11-12T15:04:05Z", "favorite": false, "file": "nap.txt", "language": "%s" } ]`

func main() {
	config := readConfig()
	snippets := readSnippets(config)
	stdin := readStdin()
	if stdin != "" {
		saveSnippet(stdin, config, snippets)
		return
	}

	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "list":
			listSnippets(snippets)
		default:
			snippet := findSnippet(os.Args[1], snippets)
			fmt.Print(snippet.Content(isatty.IsTerminal(os.Stdout.Fd())))
		}
		return
	}

	err := runInteractiveMode(config, snippets)
	if err != nil {
		fmt.Println("Alas, there's been an error", err)
	}

}

func parseName(s string) (string, string, string) {
	var (
		folder    = defaultSnippetFolder
		name      = defaultSnippetName
		language  = defaultLanguage
		remaining string
	)

	tokens := strings.Split(s, "/")
	if len(tokens) > 1 {
		folder = tokens[0]
		remaining = tokens[1]
	} else {
		remaining = tokens[0]
	}

	tokens = strings.Split(remaining, ".")
	if len(tokens) > 1 {
		name = tokens[0]
		language = tokens[1]
	} else {
		name = tokens[0]
	}

	return folder, name, language
}

func readStdin() string {
	stat, err := os.Stdin.Stat()
	if err != nil {
		return ""
	}

	if stat.Mode()&os.ModeCharDevice != 0 {
		return ""
	}

	reader := bufio.NewReader(os.Stdin)
	var b strings.Builder

	for {
		r, _, err := reader.ReadRune()
		if err != nil && err == io.EOF {
			break
		}
		_, err = b.WriteRune(r)
		if err != nil {
			return ""
		}
	}

	return b.String()
}

func readConfig() Config {
	config := Config{Home: defaultHome()}
	if err := env.Parse(&config); err != nil {
		return Config{
			Home:            defaultHome(),
			File:            defaultSnippetFileName,
			Theme:           "dracula",
			DefaultLanguage: "go",
		}
	}
	return config
}

func readSnippets(config Config) []Snippet {
	var snippets []Snippet
	file := filepath.Join(config.Home, config.File)
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
		content := fmt.Sprintf(defaultSnippetFileFormat, defaultSnippetFolder, defaultSnippetName, config.DefaultLanguage)
		_, _ = f.WriteString(content)
		dir = []byte(content)
	}
	err = json.Unmarshal(dir, &snippets)
	if err != nil {
		fmt.Printf("Unable to unmarshal %s file, %+v\n", file, err)
		return snippets
	}
	return snippets
}

func saveSnippet(content string, config Config, snippets []Snippet) {
	// Save snippet to location
	var name string = defaultSnippetName
	if len(os.Args) > 1 {
		name = strings.Join(os.Args[1:], " ")
	}

	folder, name, language := parseName(name)
	file := fmt.Sprintf("%s-%s.%s", folder, name, language)
	err := os.WriteFile(filepath.Join(config.Home, file), []byte(content), 0644)
	if err != nil {
		fmt.Println("unable to create snippet")
		return
	}

	// Add snippet metadata
	snippet := Snippet{
		Folder:   folder,
		Date:     time.Now(),
		Name:     name,
		File:     file,
		Language: language,
	}

	snippets = append([]Snippet{snippet}, snippets...)
	b, err := json.Marshal(snippets)
	if err != nil {
		fmt.Println("Could not mashal latest snippet data.", err)
		return
	}
	err = os.WriteFile(filepath.Join(config.Home, config.File), b, os.ModePerm)
	if err != nil {
		fmt.Println("Could not save snippets file.", err)
	}
}

func listSnippets(snippets []Snippet) {
	for _, snippet := range snippets {
		fmt.Println(snippet)
	}
}

func findSnippet(search string, snippets []Snippet) Snippet {
	matches := fuzzy.FindFrom(os.Args[1], Snippets{snippets})
	if len(matches) > 0 {
		return snippets[matches[0].Index]
	}
	return Snippet{}
}

func runInteractiveMode(config Config, snippets []Snippet) error {
	var folders = make(map[Folder][]list.Item)
	var items []list.Item
	for _, snippet := range snippets {
		folders[Folder(snippet.Folder)] = append(folders[Folder(snippet.Folder)], list.Item(snippet))
	}
	if len(items) <= 0 {
		items = append(items, list.Item(defaultSnippet))
	}

	var folderItems []list.Item
	foldersSlice := maps.Keys(folders)
	slices.Sort(foldersSlice)
	for _, folder := range foldersSlice {
		folderItems = append(folderItems, list.Item(Folder(folder)))
	}
	if len(folderItems) <= 0 {
		folderItems = append(folderItems, list.Item(Folder(defaultSnippetFolder)))
	}
	folderList := list.New(folderItems, folderDelegate{}, 0, 0)
	folderList.Title = "Folders"

	folderList.SetShowHelp(false)
	folderList.SetFilteringEnabled(false)
	folderList.SetShowStatusBar(false)
	folderList.DisableQuitKeybindings()
	folderList.Styles.NoItems = lipgloss.NewStyle().Margin(0, 2).Foreground(gray)
	folderList.SetStatusBarItemName("folder", "folders")

	content := viewport.New(80, 0)

	lists := map[Folder]*list.Model{}

	for folder, items := range folders {
		lists[folder] = newList(items, 20)
	}

	m := &Model{
		Lists:        lists,
		Folders:      folderList,
		Code:         content,
		ContentStyle: DefaultStyles.Content.Blurred,
		ListStyle:    DefaultStyles.Snippets.Focused,
		FoldersStyle: DefaultStyles.Folders.Blurred,
		keys:         DefaultKeyMap,
		help:         help.New(),
		config:       config,
		inputs: []textinput.Model{
			newTextInput(defaultSnippetFolder + " "),
			newTextInput(defaultSnippetName + " "),
			newTextInput(config.DefaultLanguage),
		},
		tagsInput: newTextInput("Tags"),
	}
	p := tea.NewProgram(m, tea.WithAltScreen())
	model, err := p.Run()
	if err != nil {
		return err
	}
	fm, ok := model.(*Model)
	if !ok {
		return err
	}
	var allSnippets []list.Item
	for _, list := range fm.Lists {
		allSnippets = append(allSnippets, list.Items()...)
	}
	if len(allSnippets) <= 0 {
		allSnippets = []list.Item{defaultSnippet}
	}
	b, err := json.Marshal(allSnippets)
	if err != nil {
		return err
	}
	err = os.WriteFile(filepath.Join(config.Home, config.File), b, os.ModePerm)
	if err != nil {
		return err
	}
	return nil
}

func newList(items []list.Item, height int) *list.Model {
	snippetList := list.New(items, snippetDelegate{}, 25, height)
	snippetList.SetShowHelp(false)
	snippetList.SetShowFilter(true)
	snippetList.Title = "Snippets"

	snippetList.FilterInput.Prompt = "Find: "
	snippetList.FilterInput.PromptStyle = lipgloss.NewStyle().Foreground(white).MarginLeft(1)
	snippetList.FilterInput.TextStyle = lipgloss.NewStyle().Foreground(white).Background(blue)
	snippetList.FilterInput.Cursor.Style = lipgloss.NewStyle().Foreground(white).Background(blue)
	snippetList.FilterInput.CursorStyle = lipgloss.NewStyle().Foreground(white)
	snippetList.FilterInput.Cursor.TextStyle = lipgloss.NewStyle().Foreground(white).Background(blue)
	snippetList.Styles.NoItems = lipgloss.NewStyle().Margin(0, 2).Foreground(gray)
	snippetList.SetStatusBarItemName("snippet", "snippets")
	snippetList.DisableQuitKeybindings()
	snippetList.Styles.Title = DefaultStyles.Snippets.Blurred.Title
	snippetList.Styles.TitleBar = DefaultStyles.Snippets.Blurred.TitleBar

	return &snippetList
}

func newTextInput(placeholder string) textinput.Model {
	i := textinput.New()
	i.Prompt = ""
	i.PromptStyle = lipgloss.NewStyle().Margin(0, 1)
	i.Placeholder = placeholder
	i.PlaceholderStyle = lipgloss.NewStyle().Foreground(brightBlack)
	i.CursorStyle = lipgloss.NewStyle().Foreground(blue)
	return i
}
