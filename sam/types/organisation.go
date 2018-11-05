package types

import (
	"fmt"

	"github.com/crusttech/crust/system/types"
)

type (
	// Organisations - Organisations represent a top-level grouping entity. There may be many organisations defined in a single deployment.
	Organisation types.Organisation
)

// Scope returns permissions group that for this type
func (r *Organisation) Scope() string {
	return "organisation"
}

// Resource returns a RBAC resource ID for this type
func (r *Organisation) Resource() string {
	return fmt.Sprintf("%s:%d", r.Scope(), r.ID)
}

// Operation returns a RBAC resource-scoped role name for an operation
func (r *Organisation) Operation(name string) string {
	return fmt.Sprintf("%s/%s", r.Resource(), name)
}
