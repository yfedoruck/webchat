package main

import (
	"flag"
	"log"
)

func check(err error) {
	if err != nil {
		log.Panic(err)
	}
}

func flags() *string {
	addr := flag.String("addr", ":8080", "The addr of the application")

	flag.Parse()
	return addr
}
