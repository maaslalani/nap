package main

import (
	"path/filepath"

	"github.com/adrg/xdg"
)

type Config struct {
	Home string `env:"SNOOZE_HOME"`
	File string `env:"SNOOZE_FILE"`
}

func defaultHome() string {
	return filepath.Join(xdg.DataHome, "snooze")
}

func defaultFile() string {
	return "snippets.json"
}
