package options

import (
	"strings"
)

func (o *DBOpt) Defaults() {
	if o.DSN != "" && !strings.Contains(o.DSN, "://") {
		// Make sure DSN is compatible with new requirements
		o.DSN = "mysql://" + o.DSN
	}
}

func (o DBOpt) IsSQLite() bool {
	return strings.HasPrefix(o.DSN, "sqlite3")
}
