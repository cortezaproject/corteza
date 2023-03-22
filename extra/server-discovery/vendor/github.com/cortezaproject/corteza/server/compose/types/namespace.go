package types

import (
	"database/sql/driver"
	"encoding/json"
	"github.com/cortezaproject/corteza/server/pkg/filter"
	"github.com/cortezaproject/corteza/server/pkg/sql"
	"time"
)

type (
	Namespace struct {
		ID      uint64        `json:"namespaceID,string"`
		Slug    string        `json:"slug"`
		Enabled bool          `json:"enabled"`
		Meta    NamespaceMeta `json:"meta"`

		Labels map[string]string `json:"labels,omitempty"`

		CreatedAt time.Time  `json:"createdAt,omitempty"`
		UpdatedAt *time.Time `json:"updatedAt,omitempty"`
		DeletedAt *time.Time `json:"deletedAt,omitempty"`

		// Warning: value of this field is now handled via resource-translation facility
		//          struct field is kept for the convenience for now since it allows us
		//          easy encoding/decoding of the outgoing/incoming values
		Name string `json:"name"`
	}

	NamespaceFilter struct {
		NamespaceID []uint64 `json:"namespaceID"`

		Query string `json:"query"`
		Slug  string `json:"slug"`
		Name  string `json:"name"`

		LabeledIDs []uint64          `json:"-"`
		Labels     map[string]string `json:"labels,omitempty"`

		Deleted filter.State `json:"deleted"`

		// Check fn is called by store backend for each resource found function can
		// modify the resource and return false if store should not return it
		//
		// Store then loads additional resources to satisfy the paging parameters
		Check func(*Namespace) (bool, error) `json:"-"`

		// Standard helpers for paging and sorting
		filter.Sorting
		filter.Paging
	}

	NamespaceMeta struct {
		// Temporary icon & logo URLs
		// @todo rework this when we rework attachment management
		Icon        string `json:"icon,omitempty"`
		IconID      uint64 `json:"iconID,string"`
		Logo        string `json:"logo,omitempty"`
		LogoID      uint64 `json:"logoID,string"`
		LogoEnabled bool   `json:"logoEnabled,omitempty"`

		// Warning: value of this field is now handled via resource-translation facility
		//          struct field is kept for the convenience for now since it allows us
		//          easy encoding/decoding of the outgoing/incoming values
		Subtitle string `json:"subtitle,omitempty"`

		// Warning: value of this field is now handled via resource-translation facility
		//          struct field is kept for the convenience for now since it allows us
		//          easy encoding/decoding of the outgoing/incoming values
		Description string `json:"description,omitempty"`
	}
)

func (n Namespace) Clone() *Namespace {
	c := &n
	return c
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

func (nm *NamespaceMeta) Scan(src any) error          { return sql.ParseJSON(src, nm) }
func (nm NamespaceMeta) Value() (driver.Value, error) { return json.Marshal(nm) }
