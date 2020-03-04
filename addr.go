package main

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"
	"runtime"
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

var baseDir string

func basePath() string {
	if baseDir != "" {
		return baseDir
	}
	_, b, _, ok := runtime.Caller(0)
	if !ok {
		log.Panic("Caller error")
	}

	baseDir = filepath.Dir(b)
	return baseDir
}
