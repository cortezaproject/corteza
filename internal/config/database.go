package config

import (
	"github.com/namsral/flag"
	"github.com/pkg/errors"
)

type (
	Database struct {
		DSN      string
		Profiler string
	}
)

var db *Database

func (c *Database) Validate() error {
	if c == nil {
		return nil
	}
	if c.DSN == "" {
		return errors.New("No DB DSN is set, can't connect to database")
	}
	return nil
}

func (*Database) Init(prefix ...string) *Database {
	if db != nil {
		return db
	}

	p := func(s string) string {
		return prefix[0] + "-" + s
	}

	db = new(Database)
	flag.StringVar(&db.DSN, p("db-dsn"), "crust:crust@tcp(db1:3306)/crust?collation=utf8mb4_general_ci", "DSN for database connection")
	flag.StringVar(&db.Profiler, p("db-profiler"), "", "Profiler for DB queries (none, stdout)")
	return db
}
