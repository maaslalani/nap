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
	Home string `env:"NAP_HOME"`
	File string `env:"NAP_FILE" envDefault:"snippets.json"`

	DefaultLanguage string `env:"NAP_DEFAULT_LANGUAGE" envDefault:"go"`

	Theme string `env:"NAP_THEME" envDefault:"dracula"`

	PrimaryColor        string `env:"NAP_PRIMARY_COLOR" envDefault:"#AFBEE1"`
	PrimaryColorSubdued string `env:"NAP_PRIMARY_COLOR_SUBDUED" envDefault:"#64708D"`
	BrightGreenColor    string `env:"NAP_BRIGHT_GREEN" envDefault:"#BCE1AF"`
	GreenColor          string `env:"NAP_GREEN" envDefault:"#527251"`
	BrightRedColor      string `env:"NAP_BRIGHT_RED" envDefault:"#E49393"`
	RedColor            string `env:"NAP_RED" envDefault:"#A46060"`
	ForegroundColor     string `env:"NAP_FOREGROUND" envDefault:"7"`
	BackgroundColor     string `env:"NAP_BACKGROUND" envDefault:"0"`
	GrayColor           string `env:"NAP_GRAY" envDefault:"240"`
	BlackColor          string `env:"NAP_BLACK" envDefault:"#373b41"`
	WhiteColor          string `env:"NAP_WHITE" envDefault:"#FFFFFF"`
}

// default helpers for the configuration.
// We use $XDG_DATA_HOME to avoid cluttering the user's home directory.
func defaultHome() string { return filepath.Join(xdg.DataHome, "nap") }
