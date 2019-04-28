package types

import (
	"time"

	"github.com/jmoiron/sqlx/types"

	"github.com/crusttech/crust/internal/rules"
)

type (
	Namespace struct {
		ID      uint64         `json:"namespaceID,string" db:"id"`
		Name    string         `json:"name"               db:"name"`
		Slug    string         `json:"slug"               db:"slug"`
		Enabled bool           `json:"enabled"            db:"enabled"`
		Meta    types.JSONText `json:"meta"               db:"meta"`

		CreatedAt time.Time  `json:"createdAt,omitempty" db:"created_at"`
		UpdatedAt *time.Time `json:"updatedAt,omitempty" db:"updated_at"`
		DeletedAt *time.Time `json:"deletedAt,omitempty" db:"deleted_at"`
	}

	NamespaceFilter struct {
		Query   string `json:"query"`
		Page    uint   `json:"page"`
		PerPage uint   `json:"perPage"`
		Sort    string `json:"sort"`
		Count   uint   `json:"count"`
	}
)

const (
	NamespaceCRM uint64 = 10000000
)

// Resource returns a system resource ID for this type
func (n Namespace) PermissionResource() rules.Resource {
	return NamespacePermissionResource.AppendID(n.ID)
}
