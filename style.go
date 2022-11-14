package main

import "github.com/charmbracelet/lipgloss"

const primaryColor = lipgloss.Color("#67ecff")
const primaryColorSubdued = lipgloss.Color("#439dab")
const secondaryColor = lipgloss.Color("63")

const brightGreen = lipgloss.Color("34")
const green = lipgloss.Color("29")

const brightRed = lipgloss.Color("#f11635")
const red = lipgloss.Color("#d62c45")

const gray = lipgloss.Color("240")
const black = lipgloss.Color("0")
const brightBlack = lipgloss.Color("237")
const white = lipgloss.Color("#FFF")

type SnippetsStyle struct {
	Focused SnippetsBaseStyle
	Blurred SnippetsBaseStyle
}

type FoldersStyle struct {
	Focused FoldersBaseStyle
	Blurred FoldersBaseStyle
}

type ContentStyle struct {
	Focused ContentBaseStyle
	Blurred ContentBaseStyle
}

type SnippetsBaseStyle struct {
	Base               lipgloss.Style
	Title              lipgloss.Style
	TitleBar           lipgloss.Style
	SelectedSubtitle   lipgloss.Style
	UnselectedSubtitle lipgloss.Style
	SelectedTitle      lipgloss.Style
	UnselectedTitle    lipgloss.Style
}

type FoldersBaseStyle struct {
	Base       lipgloss.Style
	Title      lipgloss.Style
	TitleBar   lipgloss.Style
	Selected   lipgloss.Style
	Unselected lipgloss.Style
}

type ContentBaseStyle struct {
	Base       lipgloss.Style
	Title      lipgloss.Style
	Separator  lipgloss.Style
	LineNumber lipgloss.Style
}

// Styles is the struct of all styles for the application.
type Styles struct {
	Snippets SnippetsStyle
	Folders  FoldersStyle
	Content  ContentStyle
}

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
			Base:       lipgloss.NewStyle().Margin(0, 1),
			Title:      lipgloss.NewStyle().Background(primaryColorSubdued).Foreground(white).Margin(0, 0, 1, 1).Padding(0, 1),
			Separator:  lipgloss.NewStyle().Foreground(white).Margin(0, 0, 0, 1),
			LineNumber: lipgloss.NewStyle().Foreground(brightBlack),
		},
		Blurred: ContentBaseStyle{
			Base:       lipgloss.NewStyle().Margin(0, 1),
			Title:      lipgloss.NewStyle().Background(black).Foreground(gray).Margin(0, 0, 1, 1).Padding(0, 1),
			Separator:  lipgloss.NewStyle().Foreground(black).Margin(0, 0, 0, 1),
			LineNumber: lipgloss.NewStyle().Foreground(black),
		},
	},
}
