package types

import (
	"time"

	"github.com/cortezaproject/corteza-server/internal/permissions"
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

// Resource returns a resource ID for this type
func (r *Role) PermissionResource() permissions.Resource {
	return RolePermissionResource.AppendID(r.ID)
}
