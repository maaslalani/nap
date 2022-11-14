package main

import (
	"path/filepath"

	"github.com/adrg/xdg"
)

// Config holds the configuration options for the application.
//
// At the moment, it is quite limited, only supporting the home folder and the
// file name of the metadata.
type Config struct {
	Home string `env:"SNOOZE_HOME"`
	File string `env:"SNOOZE_FILE" envDefault:"snippets.json"`
}

// default helpers for the configuration.
// We use $XDG_DATA_HOME to avoid cluttering the user's home directory.
func defaultHome() string { return filepath.Join(xdg.DataHome, "snooze") }
