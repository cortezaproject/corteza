package types

import (
	"fmt"

	"github.com/crusttech/crust/system/types"
)

type (
	Team types.Team
)

// Scope returns permissions group that for this type
func (r *Team) Scope() string {
	return "team"
}

// Resource returns a RBAC resource ID for this type
func (r *Team) Resource() string {
	return fmt.Sprintf("%s:%d", r.Scope(), r.ID)
}

// Operation returns a RBAC resource-scoped role name for an operation
func (r *Team) Operation(name string) string {
	return fmt.Sprintf("%s/%s", r.Resource(), name)
}
