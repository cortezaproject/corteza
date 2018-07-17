package main

import (
	"github.com/namsral/flag"
	_ "github.com/joho/godotenv/autoload"

	"github.com/crusttech/crust/rbac"
)

type configuration struct {
	httpAddr string
	dbDSN string
}

func flags(prefix string) configuration {
	var config configuration

	p := func(s string) string {
		return prefix + "-" + s
	}

	flag.StringVar(&config.addr, p("http-addr"), ":3000", "Listen address for HTTP server")
	flag.StringVar(&config.dsn, p("db-dsn"), "crust:crust@tcp(db1:3306)/crust?collation=utf8mb4_general_ci", "DSN for database connection")
	rbac.Flags()
	flag.Parse()
	return config
}