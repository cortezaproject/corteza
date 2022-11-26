package options

import (
	"github.com/cortezaproject/corteza/server/codegen/schema"
)

DB: schema.#optionsGroup & {
	handle: "DB"
	options: {
		DSN: {
			defaultValue: "sqlite3://file::memory:?cache=shared&mode=memory"
			description:  "Database connection string."
		}
		allow_destructive_schema_changes: {
			type:          "bool"
			defaultGoExpr: "false"
			description:   "Allow for irreversible changes to the database schema such as dropping columns and tables."
		}
	}
	title: "Connection to data store backend"
}
