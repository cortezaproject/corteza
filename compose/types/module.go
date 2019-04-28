package types

import (
	"time"

	"github.com/jmoiron/sqlx/types"

	"github.com/crusttech/crust/internal/rules"
)

type (
	Module struct {
		ID     uint64         `json:"moduleID,string" db:"id"`
		Name   string         `json:"name" db:"name"`
		Meta   types.JSONText `json:"meta" db:"json"`
		Fields ModuleFieldSet `json:"fields" db:"-"`
		Page   *Page          `json:"page,omitempty"`

		NamespaceID uint64 `json:"namespaceID,string" db:"rel_namespace"`

		CreatedAt time.Time  `db:"created_at" json:"createdAt,omitempty"`
		UpdatedAt *time.Time `db:"updated_at" json:"updatedAt,omitempty"`
		DeletedAt *time.Time `db:"deleted_at" json:"deletedAt,omitempty"`
	}
)

// Resource returns a system resource ID for this type
func (m Module) PermissionResource() rules.Resource {
	return ModulePermissionResource.AppendID(m.ID)
}
