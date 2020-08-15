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

		Deleted  rh.FilterState `json:"deleted"`
		Archived rh.FilterState `json:"archived"`

		Sort string `json:"sort"`

		// Standard paging fields & helpers
		rh.PageFilter

		// Resource permission check filter
		IsReadable *permissions.ResourceFilter `json:"-"`
	}

	RoleMetrics struct {
		Total         uint   `json:"total"`
		Valid         uint   `json:"valid"`
		Deleted       uint   `json:"deleted"`
		Archived      uint   `json:"archived"`
		DailyCreated  []uint `json:"dailyCreated"`
		DailyDeleted  []uint `json:"dailyDeleted"`
		DailyUpdated  []uint `json:"dailyUpdated"`
		DailyArchived []uint `json:"dailyArchived"`
	}
)

// Resource returns a resource ID for this type
func (r *Role) PermissionResource() permissions.Resource {
	return RolePermissionResource.AppendID(r.ID)
}

func (r *Role) DynamicRoles(userID uint64) []uint64 {
	return nil
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
