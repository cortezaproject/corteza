package types

import (
	"fmt"
	"time"
)

type (
	// Organisations - Organisations represent a top-level grouping entity. There may be many organisations defined in a single deployment.
	Organisation struct {
		ID         uint64     `json:"id" db:"id"`
		FQN        string     `json:"fqn" db:"fqn"`
		Name       string     `json:"name" db:"name"`
		CreatedAt  time.Time  `json:"createdAt,omitempty" db:"created_at"`
		UpdatedAt  *time.Time `json:"updatedAt,omitempty" db:"updated_at"`
		ArchivedAt *time.Time `json:"archivedAt,omitempty" db:"archived_at"`
		DeletedAt  *time.Time `json:"deletedAt,omitempty" db:"deleted_at"`
	}

	OrganisationFilter struct {
		Query string
	}
)

// Scope returns permissions group that for this type
func (r *Organisation) Scope() string {
	return "organisation"
}

// Resource returns a RBAC resource ID for this type
func (r *Organisation) Resource() string {
	return fmt.Sprintf("%s:%d", r.Scope(), r.ID)
}

// Operation returns a RBAC resource-scoped role name for an operation
func (r *Organisation) Operation(name string) string {
	return fmt.Sprintf("%s/%s", r.Resource(), name)
}
