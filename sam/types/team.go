package types

import (
	"fmt"
	"time"
)

type (
	// Teams - An organisation may have many teams. Teams may have many channels available. Access to channels may be shared between teams.
	Team struct {
		ID         uint64     `json:"id" db:"id"`
		Name       string     `json:"name" db:"name"`
		Handle     string     `json:"handle" db:"handle"`
		CreatedAt  time.Time  `json:"createdAt,omitempty" db:"created_at"`
		UpdatedAt  *time.Time `json:"updatedAt,omitempty" db:"updated_at"`
		ArchivedAt *time.Time `json:"archivedAt,omitempty" db:"archived_at"`
		DeletedAt  *time.Time `json:"deletedAt,omitempty" db:"deleted_at"`
	}

	TeamFilter struct {
		Query string
	}
)

// Scope returns permissions group that for this type
func (r *Team) Scope() string {
	return "team"
}

// Resource returns a RBAC resource ID for this type
func (r *Team) Resource() string {
	return fmt.Sprintf("%s:%d", r.Scope(), r.ID)
}

// Operation returns a RBAC resource-scoped role name for an operation
func (r *Team) Operation(name string) string {
	return fmt.Sprintf("%s/%s", r.Resource(), name)
}
