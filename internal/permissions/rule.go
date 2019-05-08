package permissions

import (
	"fmt"
)

type (
	Rule struct {
		RoleID    uint64    `json:"roleID,string" db:"rel_role"`
		Resource  Resource  `json:"resource"      db:"resource"`
		Operation Operation `json:"operation"     db:"operation"`
		Access    Access    `json:"value,string"  db:"value"`
	}
)

const (
	// Allow - Operation over a resource is allowed
	Allow Access = 1

	// Deny - Operation over a resource is denied
	Deny = 0

	// Inherit - Operation over a resource is not defined, inherit
	Inherit = -1
)

func (r Rule) String() string {
	return fmt.Sprintf("%s %d to %s on %s", r.Access, r.RoleID, r.Operation, r.Resource)
}

func (r Rule) Equals(cmp *Rule) bool {
	return r.RoleID == cmp.RoleID &&
		r.Resource == cmp.Resource &&
		r.Operation == cmp.Operation
}
