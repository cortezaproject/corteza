package main

import (
	_ "github.com/joho/godotenv/autoload"
	"github.com/namsral/flag"
)

func flags(prefix string, mountFlags ...func(...string)) {
	for _, mount := range mountFlags {
		mount(prefix)
	}
	flag.Parse()
}
