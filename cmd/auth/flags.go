package main

import (
	_ "github.com/joho/godotenv/autoload"
	"github.com/namsral/flag"
)

type configuration struct {
	monitorInterval int
}

func flags(prefix string, mountFlags ...func(...string)) configuration {
	var config configuration

	flag.IntVar(&config.monitorInterval, "monitor-interval", 300, "Monitor interval (seconds, 0 = disable)")

	for _, mount := range mountFlags {
		mount(prefix)
	}

	flag.Parse()
	return config
}
