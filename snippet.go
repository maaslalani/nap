package main

import (
	"time"
)

const defaultSnippetFolder = "misc"

const defaultSnippetName = "Untitled Snippet"

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
