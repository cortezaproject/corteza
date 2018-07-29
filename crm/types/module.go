package types

import (
	"github.com/jmoiron/sqlx/types"
)

type (
	// Modules - CRM module definitions
	Module struct {
		ID   uint64 `db:"id"`
		Name string `db:"name"`
	}

	// Modules - CRM module definitions
	ModuleContentRow struct {
		ID       uint64         `db:"id"`
		ModuleID uint64         `db:"module_id"`
		Fields   types.JSONText `db:"fields"`
	}
)
