package main

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/alecthomas/chroma/v2/quick"
)

// default values for empty state.
const (
	defaultSnippetFolder   = "misc"
	defaultLanguage        = "go"
	defaultSnippetName     = "Untitled Snippet"
	defaultSnippetFileName = "nap.txt"
)

// defaultSnippetFileFormat is the file to use for an empty snippets.json file.
var defaultSnippetFileFormat = `[ { "folder": "%s", "title": "%s", "tags": [], "date": "2022-11-12T15:04:05Z", "favorite": false, "file": "nap.txt", "language": "%s" } ]`

// defaultSnippet is a snippet with all of the default values, used for when
// there are no snippets available.
var defaultSnippet = Snippet{
	Name:     defaultSnippetName,
	Folder:   defaultSnippetFolder,
	Language: defaultLanguage,
	File:     defaultSnippetFileName,
	Date:     time.Now(),
}

// Snippet represents a snippet of code in a language.
// It is nested within a folder and can be tagged with metadata.
type Snippet struct {
	Tags     []string  `json:"tags"`
	Folder   string    `json:"folder"`
	Date     time.Time `json:"date"`
	Favorite bool      `json:"favorite"`
	Name     string    `json:"title"`
	File     string    `json:"file"`
	Language string    `json:"language"`
}

// String returns the folder/name.ext of the snippet.
func (s Snippet) String() string {
	return fmt.Sprintf("%s/%s.%s", s.Folder, s.Name, s.Language)
}

// LegacyPath returns the legacy path <folder>-<file>
func (s Snippet) LegacyPath() string {
	return s.File
}

// Path returns the path <folder>/<file>
func (s Snippet) Path() string {
	return filepath.Join(s.Folder, s.File)
}

// Content returns the snippet contents.
func (s Snippet) Content(highlight bool) string {
	config := readConfig()
	file := filepath.Join(config.Home, s.Path())
	content, err := os.ReadFile(file)
	if err != nil {
		return ""
	}

	if !highlight {
		return string(content)
	}

	var b bytes.Buffer
	err = quick.Highlight(&b, string(content), s.Language, "terminal16m", config.Theme)
	if err != nil {
		return string(content)
	}
	return b.String()
}

// Snippets is a wrapper for a snippets array to implement the fuzzy.Source
// interface.
type Snippets struct {
	snippets []Snippet
}

// String returns the string of the snippet at the specified position i
func (s Snippets) String(i int) string {
	return s.snippets[i].String()
}

// Len returns the length of the snippets array.
func (s Snippets) Len() int {
	return len(s.snippets)
}
