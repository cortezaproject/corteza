package types

import (
	"time"

	"github.com/cortezaproject/corteza-server/pkg/permissions"
	"github.com/cortezaproject/corteza-server/pkg/rh"
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
		RoleID   []uint64 `json:"roleID"`
		MemberID uint64   `json:"memberID"`

		Query string `json:"query"`

		Handle string `json:"handle"`
		Name   string `json:"name"`

		IncDeleted  bool `json:"incDeleted"`
		IncArchived bool `json:"incArchived"`

		Sort string `json:"sort"`

		// Standard paging fields & helpers
		rh.PageFilter

		// Resource permission check filter
		IsReadable *permissions.ResourceFilter `json:"-"`
	}
)

// Resource returns a resource ID for this type
func (r *Role) PermissionResource() permissions.Resource {
	return RolePermissionResource.AppendID(r.ID)
}

// FindByHandle finds role by it's handle
func (set RoleSet) FindByHandle(handle string) *Role {
	for i := range set {
		if set[i].Handle == handle {
			return set[i]
		}
	}

	return nil
}
