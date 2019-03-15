package types

import (
	"github.com/crusttech/crust/internal/rules"
)

type (
	Namespace struct {
		ID uint64 `json:"id,string" db:"id"`
	}
)

const (
	NamespaceCRM uint64 = 10000000
)

// Resource returns a system resource ID for this type
func (n Namespace) PermissionResource() rules.Resource {
	return NamespacePermissionResource.AppendID(n.ID)
}
