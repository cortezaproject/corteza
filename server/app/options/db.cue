package options

import (
	"github.com/crusttech/human/server/codegen/schema"
)

DB: schema.#optionsGroup & {
	handle: "DB"
	options: {
		DSN: {
			defaultValue: "sqlite3://file::memory:?cache=shared&mode=memory"
			description: """
				Database connection string.

				Connection pool and retry settings can be added to the query string of the DSN.
				They are read by the server and removed before the DSN is passed to the database driver,
				so they can be combined with the driver's own parameters (for example `postgres://user:pass@host/human?sslmode=disable&*connMaxOpen=50`).
				The same parameters work in the DSN of a DAL connection.

				* `*connMaxOpen` (default `256`): maximum number of open connections to the database. Keep it below the server's connection limit (PostgreSQL defaults to 100).
				* `*connMaxIdle` (default `32`): maximum number of idle connections kept in the pool.
				* `*connMaxLifetime` (default `10m`): maximum time a connection may be reused.
				* `*connTryTimeout` (default `30s`): timeout for a single connection attempt on startup.
				* `*connTryBackoffDelay` (default `10s`): delay between failed connection attempts.
				* `*connMaxTries` (default `99`): number of connection attempts before giving up.
				* `*connTryPatience` (default `0`): time window in which failed connection attempts are not logged as errors.
				"""
		}
		
	}
	title: "Connection to data store backend"
}
