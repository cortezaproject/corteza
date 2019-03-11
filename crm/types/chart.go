package types

import (
	"time"

	"github.com/jmoiron/sqlx/types"

	"github.com/crusttech/crust/internal/rules"
)

type (
	// Modules - CRM module definitions
	Chart struct {
		ID     uint64         `json:"chartID,string" db:"id"`
		Name   string         `json:"name" db:"name"`
		Config types.JSONText `json:"config" db:"config"`

		CreatedAt time.Time  `db:"created_at" json:"createdAt,omitempty"`
		UpdatedAt *time.Time `db:"updated_at" json:"updatedAt,omitempty"`
		DeletedAt *time.Time `db:"deleted_at" json:"deletedAt,omitempty"`
	}
)

// Resource returns a system resource ID for this type
func (r *Chart) Resource() rules.Resource {
	resource := rules.Resource{
		Service: "compose",
		Scope:   "page",
		ID:      r.ID,
	}

	return resource
}
