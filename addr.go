package main

import (
	"log"
	"os"
	"path/filepath"
	"runtime"
)

var port string

func host() string {
	if port == "" {
		port = os.Getenv("PORT")
		if port == "" {
			port = ":5000"
		}
		// httpAddr = flag.String("addr", ":8080", "Listen address")
		// flag.Parse()
	}
	return port
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
