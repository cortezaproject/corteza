package types

import (
	"time"

	"github.com/jmoiron/sqlx/types"

	"github.com/crusttech/crust/internal/rules"
)

type (
	Chart struct {
		ID     uint64         `json:"chartID,string" db:"id"`
		Name   string         `json:"name" db:"name"`
		Config types.JSONText `json:"config" db:"config"`

		NamespaceID uint64 `json:"namespaceID,string" db:"rel_namespace,string"`

		CreatedAt time.Time  `db:"created_at" json:"createdAt,omitempty"`
		UpdatedAt *time.Time `db:"updated_at" json:"updatedAt,omitempty"`
		DeletedAt *time.Time `db:"deleted_at" json:"deletedAt,omitempty"`
	}

	ChartFilter struct {
		Query   string `json:"query"`
		Page    uint   `json:"page"`
		PerPage uint   `json:"perPage"`
		Sort    string `json:"sort"`
		Count   uint   `json:"count"`
	}
)

// Resource returns a system resource ID for this type
func (c Chart) PermissionResource() rules.Resource {
	return ChartPermissionResource.AppendID(c.ID)
}
