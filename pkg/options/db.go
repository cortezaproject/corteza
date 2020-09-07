package options

import (
	"strings"
	"time"
)

type (
	DBOpt struct {
		DSN string `env:"DB_DSN"`
	}
)

func DB(pfix string) (o *DBOpt) {
	const delay = 15 * time.Second
	const maxTries = 100

	o = &DBOpt{
		DSN: "mysql://corteza:corteza@tcp(db:3306)/corteza?collation=utf8mb4_general_ci",
	}

	fill(o)

	if !strings.Contains(o.DSN, "://") {
		// Make sure DSN is compatible with new requirements
		o.DSN = "mysql://" + o.DSN
	}

	return
}
