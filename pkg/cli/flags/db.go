package flags

import (
	"github.com/spf13/cobra"
)

type (
	DBOpt struct {
		DSN      string
		Profiler string
	}
)

func DB(cmd *cobra.Command, pfix string) (o *DBOpt) {
	o = &DBOpt{}

	BindString(cmd, &o.DSN,
		pFlag(pfix, "db-dsn"), "corteza:corteza@tcp(db:3306)/corteza?collation=utf8mb4_general_ci",
		"DSN for database connection")

	BindString(cmd, &o.Profiler,
		pFlag(pfix, "db-profiler"), "none",
		"Profiler for DB queries (none, stdout, logger)")

	return
}
