package main

type Config struct {
	Home string `env:"SNOOZE_HOME" envDefault:"${XDG_DATA_HOME}/.snooze" envExpand:"true"`
	File string `env:"SNOOZE_FILE" envDefault:"snippets.json"`
}
