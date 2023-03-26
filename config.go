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
	Home string `env:"NAP_HOME" yaml:"home"`
	File string `env:"NAP_FILE" yaml:"file"`

	DefaultLanguage string `env:"NAP_DEFAULT_LANGUAGE" yaml:"default_language"`

	Theme string `env:"NAP_THEME" yaml:"theme"`

	PrimaryColor        string `env:"NAP_PRIMARY_COLOR" yaml:"primary_color"`
	PrimaryColorSubdued string `env:"NAP_PRIMARY_COLOR_SUBDUED" yaml:"primary_color_subdued"`
	BrightGreenColor    string `env:"NAP_BRIGHT_GREEN" yaml:"bright_green"`
	GreenColor          string `env:"NAP_GREEN" yaml:"green"`
	BrightRedColor      string `env:"NAP_BRIGHT_RED" yaml:"bright_red"`
	RedColor            string `env:"NAP_RED" yaml:"red"`
	ForegroundColor     string `env:"NAP_FOREGROUND" yaml:"foreground"`
	BackgroundColor     string `env:"NAP_BACKGROUND" yaml:"background"`
	GrayColor           string `env:"NAP_GRAY" yaml:"gray"`
	BlackColor          string `env:"NAP_BLACK" yaml:"black"`
	WhiteColor          string `env:"NAP_WHITE" yaml:"white"`
}

func newConfig() Config {
	return Config{
		Home:                defaultHome(),
		File:                "snippets.json",
		DefaultLanguage:     defaultLanguage,
		Theme:               "dracula",
		PrimaryColor:        "#AFBEE1",
		PrimaryColorSubdued: "#64708D",
		BrightGreenColor:    "#BCE1AF",
		GreenColor:          "#527251",
		BrightRedColor:      "#E49393",
		RedColor:            "#A46060",
		ForegroundColor:     "15",
		BackgroundColor:     "235",
		GrayColor:           "241",
		BlackColor:          "#373b41",
		WhiteColor:          "#FFFFFF",
	}
}

// default helpers for the configuration.
// We use $XDG_DATA_HOME to avoid cluttering the user's home directory.
func defaultHome() string { return filepath.Join(xdg.DataHome, "nap") }
