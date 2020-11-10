package options

import (
	"strings"
)

type (
	DBOpt struct {
		DSN string `env:"DB_DSN"`
	}
)

func DB(pfix string) (o *DBOpt) {
	o = &DBOpt{
		DSN: "sqlite3://file::memory:?cache=shared&mode=memory",
	}

	fill(o)

	if o.DSN != "" && !strings.Contains(o.DSN, "://") {
		// Make sure DSN is compatible with new requirements
		o.DSN = "mysql://" + o.DSN
	}

	return
}
