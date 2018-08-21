package main

import (
	_ "github.com/joho/godotenv/autoload"
	"github.com/namsral/flag"
)

type configuration struct {
	httpAddr        string
	dbDSN           string
	monitorInterval int
}

func flags(prefix string, mountFlags ...func()) configuration {
	var config configuration

	p := func(s string) string {
		return prefix + "-" + s
	}

	flag.StringVar(&config.httpAddr, p("http-addr"), ":3000", "Listen address for HTTP server")
	flag.StringVar(&config.dbDSN, p("db-dsn"), "crust:crust@tcp(db1:3306)/crust?collation=utf8mb4_general_ci", "DSN for database connection")
	flag.IntVar(&config.monitorInterval, "monitor-interval", 300, "Monitor interval (seconds, 0 = disable)")

	for _, mount := range mountFlags {
		mount()
	}

	flag.Parse()
	return config
}
