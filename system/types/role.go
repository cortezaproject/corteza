package types

import (
	"time"

	"github.com/crusttech/crust/internal/rules"
)

type (
	// Role - An organisation may have many roles. Roles may have many channels available. Access to channels may be shared between roles.
	Role struct {
		ID         uint64     `json:"roleID,string" db:"id"`
		Name       string     `json:"name" db:"name"`
		Handle     string     `json:"handle" db:"handle"`
		CreatedAt  time.Time  `json:"createdAt,omitempty" db:"created_at"`
		UpdatedAt  *time.Time `json:"updatedAt,omitempty" db:"updated_at"`
		ArchivedAt *time.Time `json:"archivedAt,omitempty" db:"archived_at"`
		DeletedAt  *time.Time `json:"deletedAt,omitempty" db:"deleted_at"`
	}

	RoleFilter struct {
		Query string
	}
)

// Resource returns a system resource ID for this type
func (r *Role) Resource() rules.Resource {
	resource := rules.Resource{
		Scope: "team",
	}
	if r != nil {
		resource.ID = r.ID
		resource.Name = r.Name
	}
	return resource
}
