package types

import (
	"time"

	"github.com/cortezaproject/corteza-server/pkg/filter"

	"github.com/jmoiron/sqlx/types"

	"github.com/cortezaproject/corteza-server/pkg/rbac"
)

type (
	Module struct {
		ID     uint64         `json:"moduleID,string"  yaml:"-"`
		Handle string         `json:"handle"`
		Name   string         `json:"name"`
		Meta   types.JSONText `json:"meta" yaml:",omitempty"`
		Fields ModuleFieldSet `json:"fields" yaml:"-"`

		NamespaceID uint64 `json:"namespaceID,string" yaml:"-"`

		CreatedAt time.Time  `json:"createdAt,omitempty" yaml:",omitempty"`
		UpdatedAt *time.Time `json:"updatedAt,omitempty" yaml:",omitempty"`
		DeletedAt *time.Time `json:"deletedAt,omitempty" yaml:",omitempty"`
	}

	ModuleFilter struct {
		NamespaceID uint64 `json:"namespaceID,string"`
		Query       string `json:"query"`
		Handle      string `json:"handle"`
		Name        string `json:"name"`

		Deleted filter.State `json:"deleted"`

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
func (m Module) RBACResource() rbac.Resource {
	return ModuleRBACResource.AppendID(m.ID)
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
