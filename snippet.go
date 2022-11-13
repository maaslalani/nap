package main

import (
	"time"
)

// Snippet represents a snippet of code in a language.
// It is nested within a folder and can be tagged with metadata.
type Snippet struct {
	Tags     []string  `json:"tags"`
	Folder   string    `json:"folder"`
	Date     time.Time `json:"date"`
	Favorite bool      `json:"favorite"`
	Title    string    `json:"title"`
	File     string    `json:"file"`
	Language string    `json:"language"`
}

const untitledSnippet = `package main

import "fmt"

// Untitled Snippet

func main() {
	fmt.Println("Press 'e' to edit this snippet.")
	fmt.Println("Press 'r' to rename this snippet.")
}
`
