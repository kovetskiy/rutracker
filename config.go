package main

import (
	"github.com/kovetskiy/ko"
)

type Config struct {
	Username string `toml:"username" required:"true"`
	Password string `toml:"password" required:"true"`
	BaseURL  string `toml:"base_url" required:"true" default:"http://rutracker.org/"`
}

func LoadConfig(path string) (*Config, error) {
	var unit Config
	return &unit, ko.Load(path, &unit)
}
