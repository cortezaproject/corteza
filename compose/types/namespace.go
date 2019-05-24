package types

import (
	"time"

	"github.com/cortezaproject/corteza-server/internal/permissions"
	"github.com/jmoiron/sqlx/types"
)

type (
	Namespace struct {
		ID      uint64         `db:"id"        json:"namespaceID,string"`
		Name    string         `db:"name"      json:"name"`
		Slug    string         `db:"slug"      json:"slug"`
		Enabled bool           `db:"enabled"   json:"enabled"`
		Meta    types.JSONText `db:"meta"      json:"meta"`

		CreatedAt time.Time  `db:"created_at"  json:"createdAt,omitempty"`
		UpdatedAt *time.Time `db:"updated_at"  json:"updatedAt,omitempty"`
		DeletedAt *time.Time `db:"deleted_at"  json:"deletedAt,omitempty"`
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
	NamespaceCRM uint64 = 88714882739863655
)

// Resource returns a system resource ID for this type
func (n Namespace) PermissionResource() permissions.Resource {
	return NamespacePermissionResource.AppendID(n.ID)
}
