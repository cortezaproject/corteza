package app

// Registers all supported store backends
//
// SQLite is intentionally ignored here
import (
	_ "github.com/cortezaproject/corteza-server/store/mysql"
	_ "github.com/cortezaproject/corteza-server/store/pgsql"
)
