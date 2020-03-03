package main

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type config struct {
	ClientID     string `json:"ClientID"`
	ClientSecret string `json:"ClientSecret"`
	RedirectURL  string `json:"RedirectURL"`
}

func (c *config) set(service string) {
	file, err := os.Open(basePath() + filepath.FromSlash("/config/"+service+".json"))
	check(err)

	decoder := json.NewDecoder(file)

	err = decoder.Decode(c)
	check(err)
}
