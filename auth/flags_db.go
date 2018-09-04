package auth

import (
	"github.com/namsral/flag"
	"github.com/pkg/errors"
)

type (
	dbFlags struct {
		dsn      string
		profiler string
	}
)

func (c *dbFlags) validate() error {
	if c.dsn == "" {
		return errors.New("No DB DSN is set, can't connect to database")
	}
	return nil
}

func (c *dbFlags) flags(prefix ...string) {
	p := func(s string) string {
		return prefix[0] + "-" + s
	}

	flag.StringVar(&c.dsn, p("db-dsn"), "crust:crust@tcp(db1:3306)/crust?collation=utf8mb4_general_ci", "DSN for database connection")
	flag.StringVar(&c.profiler, p("db-profiler"), "", "Profiler for DB queries (none, stdout)")
}
