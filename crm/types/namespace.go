package types

import (
	"github.com/crusttech/crust/internal/rules"
)

type (
	Namespace struct {
		ID string `json:"id,string" db:"id"`
	}
)

// Resource returns a system resource ID for this type
func (r *Namespace) Resource() rules.Resource {
	resource := rules.Resource{
		Service: "compose",
		// Hardcoded single namespace (CRM) for now
		Scope: "namespace:crm",
	}

	return resource
}
