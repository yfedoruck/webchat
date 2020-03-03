package main

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type config struct {
	ClientID     string `json:"ClientID"`
	ClientSecret string `json:"ClientSecret"`
}

func (c *config) set(service string) {
	file, err := os.Open(basePath() + filepath.FromSlash("/config/"+service+".json"))
	check(err)

	err = json.NewDecoder(file).Decode(c)
	check(err)
}
