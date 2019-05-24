package types

import (
	"time"

	"github.com/cortezaproject/corteza-server/internal/permissions"
	"github.com/jmoiron/sqlx/types"
)

type (
	Module struct {
		ID     uint64         `json:"moduleID,string" db:"id"`
		Name   string         `json:"name" db:"name"`
		Meta   types.JSONText `json:"meta" db:"json"`
		Fields ModuleFieldSet `json:"fields" db:"-"`

		NamespaceID uint64 `json:"namespaceID,string" db:"rel_namespace"`

		CreatedAt time.Time  `db:"created_at" json:"createdAt,omitempty"`
		UpdatedAt *time.Time `db:"updated_at" json:"updatedAt,omitempty"`
		DeletedAt *time.Time `db:"deleted_at" json:"deletedAt,omitempty"`
	}

	ModuleFilter struct {
		NamespaceID uint64 `json:"namespaceID,string"`
		Query       string `json:"query"`
		Page        uint   `json:"page"`
		PerPage     uint   `json:"perPage"`
		// Sort    string `json:"sort"`
		Count uint `json:"count"`
	}
)

// Resource returns a system resource ID for this type
func (m Module) PermissionResource() permissions.Resource {
	return ModulePermissionResource.AppendID(m.ID)
}
