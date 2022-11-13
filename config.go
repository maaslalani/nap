package main

type Config struct {
	Home string `env:"SNOOZE_HOME" envDefault:".snooze" envExpand:"true"`
	File string `env:"SNOOZE_FILE" envDefault:"snippets.json"`
}
