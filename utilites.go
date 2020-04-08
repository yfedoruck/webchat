package main

import (
	"flag"
)

func flags() *string {
	addr := flag.String("addr", ":8080", "The addr of the application")

	flag.Parse()
	return addr
}
