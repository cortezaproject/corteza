package options

import (
	"github.com/cortezaproject/corteza-server/codegen/schema"
)

DB: schema.#optionsGroup & {
	handle: "DB"
	options: {
		DSN: {
			defaultValue: "sqlite3://file::memory:?cache=shared&mode=memory"
			description:  "Database connection string."
		}
	}
	title: "Connection to data store backend"
}
