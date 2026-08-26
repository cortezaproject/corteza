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
