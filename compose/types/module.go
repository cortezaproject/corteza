package types

import (
	"github.com/cortezaproject/corteza-server/pkg/filter"
	"time"

	"github.com/jmoiron/sqlx/types"

	"github.com/cortezaproject/corteza-server/pkg/permissions"
	"github.com/cortezaproject/corteza-server/pkg/rh"
)

type (
	Module struct {
		ID     uint64         `json:"moduleID,string"`
		Handle string         `json:"handle"`
		Name   string         `json:"name"`
		Meta   types.JSONText `json:"meta"`
		Fields ModuleFieldSet `json:"fields" db:"-"`

		NamespaceID uint64 `json:"namespaceID,string"`

		CreatedAt time.Time  `json:"createdAt,omitempty"`
		UpdatedAt *time.Time `json:"updatedAt,omitempty"`
		DeletedAt *time.Time `json:"deletedAt,omitempty"`
	}

	ModuleFilter struct {
		NamespaceID uint64 `json:"namespaceID,string"`
		Query       string `json:"query"`
		Handle      string `json:"handle"`
		Name        string `json:"name"`

		Deleted rh.FilterState `json:"deleted"`

		// Check fn is called by store backend for each resource found function can
		// modify the resource and return false if store should not return it
		//
		// Store then loads additional resources to satisfy the paging parameters
		Check func(*Module) (bool, error) `json:"-"`

		// Standard helpers for paging and sorting
		filter.Sorting
		filter.Paging
	}
)

// Resource returns a system resource ID for this type
func (m Module) PermissionResource() permissions.Resource {
	return ModulePermissionResource.AppendID(m.ID)
}

func (m Module) DynamicRoles(userID uint64) []uint64 {
	return nil
}

func (m Module) Clone() *Module {
	c := &m
	c.Fields = m.Fields.Clone()
	return c
}

// FindByHandle finds module by it's handle
func (set ModuleSet) FindByHandle(handle string) *Module {
	for i := range set {
		if set[i].Handle == handle {
			return set[i]
		}
	}

	return nil
}
