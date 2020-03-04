package main

import (
	"encoding/json"
	"os"
	"path/filepath"
)

var addr addrConf

type addrConf struct {
	Host string `json:"Host"`
}

func host() string {
	if addr.Host != "" {
		return addr.Host
	}
	file, err := os.Open(basePath() + filepath.FromSlash("/config/addr.json"))
	check(err)

	a := addrConf{}
	err = json.NewDecoder(file).Decode(&a)
	check(err)

	addr.Host = a.Host
	return addr.Host
}
