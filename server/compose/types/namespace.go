package types

import (
	"github.com/crusttech/human/server/pkg/filter"
	labelTypes "github.com/crusttech/human/server/pkg/label/types"
)

type (
	// Namespace is the resource struct; it is generated into namespace.gen.go.

	NamespaceFilter struct {
		NamespaceID []string `json:"namespaceID"`
		TenantID    uint64   `json:"tenantID,string,omitempty"`
		ProjectID   uint64   `json:"projectID,string,omitempty"`

		Query string `json:"query"`
		Slug  string `json:"slug"`
		Name  string `json:"name"`

		LabeledIDs []uint64                         `json:"-"`
		Labels     map[string]labelTypes.LabelValue `json:"labels,omitempty"`

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
		HideSidebar bool   `json:"hideSidebar"`

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

// Dict exposes namespace attributes for RBAC contextual role evaluation.
func (n Namespace) Dict() map[string]interface{} {
	return map[string]interface{}{
		"ID":             n.ID,
		"namespaceID":    n.ID,
		"slug":           n.Slug,
		"name":           n.Name,
		"enabled":        n.Enabled,
		"createdAt":      n.CreatedAt,
		"createdByAgent": n.CreatedByAgent,
		"updatedAt":      n.UpdatedAt,
		"deletedAt":      n.DeletedAt,
	}
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
