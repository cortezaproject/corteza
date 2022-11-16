package app

// Registers all supported store backends
import (
	_ "github.com/cortezaproject/corteza/server/store/mysql"
	_ "github.com/cortezaproject/corteza/server/store/postgres"
	_ "github.com/cortezaproject/corteza/server/store/sqlite3"
)
