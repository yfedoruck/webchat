package main

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"
	"runtime"
)

func port() string {
	p := os.Getenv("PORT")
	if p == "" {
		p = "5000"
	}
	return p
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

type AddrConf struct {
	Host     string `json:"Host"`
	Socket   string `json:"socket"`
}

var addr *AddrConf

func Conf() *AddrConf {
	if addr != nil {
		return addr
	}
	file, err := os.Open(basePath() + filepath.FromSlash("/config/"+Domain()+"/addr.json"))
	check(err)

	addr = &AddrConf{}
	err = json.NewDecoder(file).Decode(addr)
	check(err)

	return addr
}

func Domain() string {
	domain := os.Getenv("USERDOMAIN")
	if domain == "home" {
		return "local"
	}
	return "heroku"
}