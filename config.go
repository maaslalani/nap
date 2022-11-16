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
	Home            string `env:"NAP_HOME"`
	File            string `env:"NAP_FILE" envDefault:"snippets.json"`
	Theme           string `env:"NAP_THEME" envDefault:"dracula"`
	DefaultLanguage string `env:"NAP_DEFAULT_LANGUAGE" envDefault:"go"`
}

// default helpers for the configuration.
// We use $XDG_DATA_HOME to avoid cluttering the user's home directory.
func defaultHome() string { return filepath.Join(xdg.DataHome, "nap") }
