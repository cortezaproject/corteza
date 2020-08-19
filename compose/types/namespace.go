package types

import (
	"database/sql/driver"
	"encoding/json"
	"github.com/cortezaproject/corteza-server/store"
	"time"

	"github.com/pkg/errors"

	"github.com/cortezaproject/corteza-server/pkg/permissions"
	"github.com/cortezaproject/corteza-server/pkg/rh"
)

type (
	Namespace struct {
		ID      uint64        `json:"namespaceID,string"`
		Name    string        `json:"name"`
		Slug    string        `json:"slug"`
		Enabled bool          `json:"enabled"`
		Meta    NamespaceMeta `json:"meta"`

		CreatedAt time.Time  `json:"createdAt,omitempty"`
		UpdatedAt *time.Time `json:"updatedAt,omitempty"`
		DeletedAt *time.Time `json:"deletedAt,omitempty"`
	}

	NamespaceFilter struct {
		Query string `json:"query"`
		Slug  string `json:"slug"`
		Name  string `json:"name"`

		Deleted rh.FilterState `json:"deleted"`

		// Check fn is called by store backend for each resource found function can
		// modify the resource and return false if store should not return it
		//
		// Store then loads additional resources to satisfy the paging parameters
		Check func(*Namespace) (bool, error) `json:"-"`

		// Standard helpers for paging and sorting
		store.Sorting
		store.Paging
	}

	NamespaceMeta struct {
		Subtitle    string `json:"subtitle,omitempty"    yaml:",omitempty"`
		Description string `json:"description,omitempty" yaml:",omitempty"`
	}
)

// Resource returns a system resource ID for this type
func (n Namespace) PermissionResource() permissions.Resource {
	return NamespacePermissionResource.AppendID(n.ID)
}

func (n Namespace) DynamicRoles(userID uint64) []uint64 {
	return nil
}

// FindByHandle finds namespace by it's handle/slug
func (set NamespaceSet) FindByHandle(handle string) *Namespace {
	for i := range set {
		if set[i].Slug == handle {
			return set[i]
		}
	}

	return nil
}

func (nm *NamespaceMeta) Scan(value interface{}) error {
	//lint:ignore S1034 This typecast is intentional, we need to get []byte out of a []uint8
	switch value.(type) {
	case nil:
		*nm = NamespaceMeta{}
	case []uint8:
		b := value.([]byte)
		if err := json.Unmarshal(b, nm); err != nil {
			return errors.Wrapf(err, "Can not scan '%v' into NamespaceMeta", string(b))
		}
	}

	return nil
}

func (nm NamespaceMeta) Value() (driver.Value, error) {
	return json.Marshal(nm)
}
