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
	File string `env:"NAP_FILE" envDefault:"snippets.json" yaml:"file"`

	DefaultLanguage string `env:"NAP_DEFAULT_LANGUAGE" envDefault:"go" yaml:"default_language"`

	Theme string `env:"NAP_THEME" envDefault:"dracula" yaml:"theme"`

	PrimaryColor        string `env:"NAP_PRIMARY_COLOR" envDefault:"#AFBEE1" yaml:"primary_color"`
	PrimaryColorSubdued string `env:"NAP_PRIMARY_COLOR_SUBDUED" envDefault:"#64708D" yaml:"primary_color_subdued"`
	BrightGreenColor    string `env:"NAP_BRIGHT_GREEN" envDefault:"#BCE1AF" yaml:"bright_green_color"`
	GreenColor          string `env:"NAP_GREEN" envDefault:"#527251" yaml:"green_color"`
	BrightRedColor      string `env:"NAP_BRIGHT_RED" envDefault:"#E49393" yaml:"bright_red_color"`
	RedColor            string `env:"NAP_RED" envDefault:"#A46060" yaml:"red_color"`
	ForegroundColor     string `env:"NAP_FOREGROUND" envDefault:"7" yaml:"foreground_color"`
	BackgroundColor     string `env:"NAP_BACKGROUND" envDefault:"0" yaml:"background_color"`
	GrayColor           string `env:"NAP_GRAY" envDefault:"240" yaml:"gray_color"`
	BlackColor          string `env:"NAP_BLACK" envDefault:"#373b41" yaml:"black_color"`
	WhiteColor          string `env:"NAP_WHITE" envDefault:"#FFFFFF" yaml:"white_color"`
}

// default helpers for the configuration.
// We use $XDG_DATA_HOME to avoid cluttering the user's home directory.
func defaultHome() string { return filepath.Join(xdg.DataHome, "nap") }
