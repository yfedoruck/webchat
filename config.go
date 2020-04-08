package main

import (
	"encoding/json"
	"github.com/yfedoruck/webchat/pkg/env"
	"github.com/yfedoruck/webchat/pkg/fail"
	"os"
	"path/filepath"
)

type config struct {
	ClientID     string `json:"ClientID"`
	ClientSecret string `json:"ClientSecret"`
}

func (c *config) set(service string) {
	file, err := os.Open(env.AppPath() + filepath.FromSlash("/config/"+service+".json"))
	fail.Check(err)

	err = json.NewDecoder(file).Decode(c)
	fail.Check(err)
}
