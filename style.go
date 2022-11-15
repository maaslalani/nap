package main

import "github.com/charmbracelet/lipgloss"

// primary color and primary color subdued are the colors used for the
// selection purposes.
const primaryColor = lipgloss.Color("#B294BB")
const primaryColorSubdued = lipgloss.Color("#85678F")

// colors for the application.
const black = lipgloss.Color("#282a2e")
const blue = lipgloss.Color("#5f819d")
const brightBlack = lipgloss.Color("#373b41")
const brightBlue = lipgloss.Color("#81a2be")
const brightCyan = lipgloss.Color("#8abeb7")
const brightGreen = lipgloss.Color("#b5bd68")
const brightMagenta = lipgloss.Color("#b294bb")
const brightRed = lipgloss.Color("#CC6666")
const brightWhite = lipgloss.Color("#c5c8c6")
const brightYellow = lipgloss.Color("#f0c674")
const cyan = lipgloss.Color("#5e8d87")
const gray = lipgloss.Color("240")
const green = lipgloss.Color("#8c9440")
const magenta = lipgloss.Color("#85678f")
const red = lipgloss.Color("#954D4D")
const white = lipgloss.Color("#fff")
const yellow = lipgloss.Color("#de935f")

// SnippetsStyle is the style struct to handle the focusing and blurring of the
// snippets pane in the application.
type SnippetsStyle struct {
	Focused SnippetsBaseStyle
	Blurred SnippetsBaseStyle
}

// FoldersStyle is the style struct to handle the focusing and blurring of the
// folders pane in the application.
type FoldersStyle struct {
	Focused FoldersBaseStyle
	Blurred FoldersBaseStyle
}

// ContentStyle is the style struct to handle the focusing and blurring of the
// content pane in the application.
type ContentStyle struct {
	Focused ContentBaseStyle
	Blurred ContentBaseStyle
}

// SnippetsBaseStyle holds the neccessary styling for the snippets pane of
// the application.
type SnippetsBaseStyle struct {
	Base               lipgloss.Style
	Title              lipgloss.Style
	TitleBar           lipgloss.Style
	SelectedSubtitle   lipgloss.Style
	UnselectedSubtitle lipgloss.Style
	SelectedTitle      lipgloss.Style
	UnselectedTitle    lipgloss.Style
}

// FoldersBaseStyle holds the neccessary styling for the folders pane of
// the application.
type FoldersBaseStyle struct {
	Base       lipgloss.Style
	Title      lipgloss.Style
	TitleBar   lipgloss.Style
	Selected   lipgloss.Style
	Unselected lipgloss.Style
}

// ContentBaseStyle holds the neccessary styling for the content pane of the
// application.
type ContentBaseStyle struct {
	Base         lipgloss.Style
	Title        lipgloss.Style
	Separator    lipgloss.Style
	LineNumber   lipgloss.Style
	EmptyHint    lipgloss.Style
	EmptyHintKey lipgloss.Style
}

// Styles is the struct of all styles for the application.
type Styles struct {
	Snippets SnippetsStyle
	Folders  FoldersStyle
	Content  ContentStyle
}

var marginStyle = lipgloss.NewStyle().Margin(0, 2)

// DefaultStyles is the default implementation of the styles struct for all
// styling in the application.
var DefaultStyles = Styles{
	Snippets: SnippetsStyle{
		Focused: SnippetsBaseStyle{
			Base:               lipgloss.NewStyle().Width(35),
			Title:              lipgloss.NewStyle().Padding(0, 1).Foreground(white),
			TitleBar:           lipgloss.NewStyle().Background(primaryColorSubdued).Width(35-2).Margin(0, 1, 1, 1),
			SelectedSubtitle:   lipgloss.NewStyle().Foreground(primaryColorSubdued),
			UnselectedSubtitle: lipgloss.NewStyle().Foreground(lipgloss.Color("237")),
			SelectedTitle:      lipgloss.NewStyle().Foreground(primaryColor),
			UnselectedTitle:    lipgloss.NewStyle().Foreground(lipgloss.Color("240")),
		},
		Blurred: SnippetsBaseStyle{
			Base:               lipgloss.NewStyle().Width(35),
			Title:              lipgloss.NewStyle().Padding(0, 1).Foreground(gray),
			TitleBar:           lipgloss.NewStyle().Background(black).Width(35-2).Margin(0, 1, 1, 1),
			SelectedSubtitle:   lipgloss.NewStyle().Foreground(black),
			UnselectedSubtitle: lipgloss.NewStyle().Foreground(black),
			SelectedTitle:      lipgloss.NewStyle().Foreground(black),
			UnselectedTitle:    lipgloss.NewStyle().Foreground(black),
		},
	},
	Folders: FoldersStyle{
		Focused: FoldersBaseStyle{
			Base:       lipgloss.NewStyle().Width(20),
			Title:      lipgloss.NewStyle().Padding(0, 1).Foreground(white),
			TitleBar:   lipgloss.NewStyle().Background(primaryColorSubdued).Width(20-2).Margin(0, 1, 1, 1),
			Selected:   lipgloss.NewStyle().Foreground(primaryColor),
			Unselected: lipgloss.NewStyle().Foreground(gray),
		},
		Blurred: FoldersBaseStyle{
			Base:       lipgloss.NewStyle().Width(20),
			Title:      lipgloss.NewStyle().Padding(0, 1).Foreground(gray),
			TitleBar:   lipgloss.NewStyle().Background(black).Width(20-2).Margin(0, 1, 1, 1),
			Selected:   lipgloss.NewStyle().Foreground(black),
			Unselected: lipgloss.NewStyle().Foreground(black),
		},
	},
	Content: ContentStyle{
		Focused: ContentBaseStyle{
			Base:         lipgloss.NewStyle().Margin(0, 1),
			Title:        lipgloss.NewStyle().Background(primaryColorSubdued).Foreground(white).Margin(0, 0, 1, 1).Padding(0, 1),
			Separator:    lipgloss.NewStyle().Foreground(white).Margin(0, 0, 1, 1),
			LineNumber:   lipgloss.NewStyle().Foreground(brightBlack),
			EmptyHint:    lipgloss.NewStyle().Foreground(gray),
			EmptyHintKey: lipgloss.NewStyle().Foreground(primaryColor),
		},
		Blurred: ContentBaseStyle{
			Base:         lipgloss.NewStyle().Margin(0, 1),
			Title:        lipgloss.NewStyle().Background(black).Foreground(gray).Margin(0, 0, 1, 1).Padding(0, 1),
			Separator:    lipgloss.NewStyle().Foreground(gray).Margin(0, 0, 1, 1),
			LineNumber:   lipgloss.NewStyle().Foreground(black),
			EmptyHint:    lipgloss.NewStyle().Foreground(gray),
			EmptyHintKey: lipgloss.NewStyle().Foreground(primaryColor),
		},
	},
}
