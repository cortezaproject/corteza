package types

/* If you want to edit this file by hand, remove codegen/[project]/types/index.php */

import (
	"github.com/jmoiron/sqlx/types"
)

type (
	// Modules - CRM module definitions
	Module struct {
		ID     uint64         `json:"id" db:"id"`
		Name   string         `json:"name" db:"name"`
		Fields types.JSONText `json:"json" db:"json"`
	}

	// Modules - CRM module definitions
	Content struct {
		ID       uint64         `json:"id" db:"id"`
		ModuleID uint64         `json:"module_id" db:"module_id"`
		Fields   types.JSONText `json:"json" db:"json"`
	}
)
