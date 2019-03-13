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

var dbs map[string]*Database

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
	if dbs == nil {
		dbs = make(map[string]*Database)
	}

	name := "default"
	if len(prefix) > 0 {
		name = prefix[0]
	}

	if db := dbs[name]; db != nil {
		return db
	}

	p := func(s string) string {
		if len(prefix) > 0 {
			return prefix[0] + "-" + s
		}
		return s
	}

	db := new(Database)
	flag.StringVar(&db.DSN, p("db-dsn"), "", "DSN for database connection (e.g. user:pass@tcp(db1:3306)/dbname?collation=utf8mb4_general_ci)")
	flag.StringVar(&db.Profiler, p("db-profiler"), "", "Profiler for DB queries (none, stdout)")
	dbs[name] = db
	return db
}
