package rdbms

import (
	sqlx "github.com/jmoiron/sqlx/types"
)

// A temporary alias for old type that was used handle JSON encoded data
type rawJson = sqlx.JSONText
