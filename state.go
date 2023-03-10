package main

import (
	"encoding/json"
	"errors"
	"io/fs"
	"os"

	"github.com/adrg/xdg"
)

// State is application state between runs
type State struct {
	CurrentFolder int
}

// Save saves the state of the application
func (s State) Save() error {
	fi, err := os.Create(defaultState())
	if err != nil {
		return err
	}
	defer fi.Close()
	return json.NewEncoder(fi).Encode(s)

}

// defaultState returns the default state path
func defaultState() string {
	if c := os.Getenv("NAP_STATE"); c != "" {
		return c
	}
	statePath, err := xdg.StateFile("nap/state.json")
	if err != nil {
		return "state.json"
	}
	return statePath
}

// readState returns the application state
func readState() State {
	var s State
	fi, err := os.Open(defaultState())
	if err != nil && !errors.Is(err, fs.ErrNotExist) {
		return s
	}
	defer fi.Close()

	if err := json.NewDecoder(fi).Decode(&s); err != nil {
		return s
	}

	return s
}
