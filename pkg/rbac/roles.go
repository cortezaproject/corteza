package rbac

import (
	"github.com/cortezaproject/corteza-server/pkg/slice"
)

type (
	ctxRoleCheckFn func(map[string]interface{}) bool

	// role information, adapted for the needs of RBAC package
	role struct {
		// all RBAC rules refer to role ID
		id uint64

		// for debugging and logging
		handle string

		// role type that will allow us
		kind roleKind

		check ctxRoleCheckFn
	}

	roleKind int

	partRoles []map[uint64]bool
)

const (
	CommonRole roleKind = iota
	AnonymousRole
	AuthenticatedRole
	ContextRole
	BypassRole
)

// partitions roles by kind
func partitionRoles(rr ...*role) partRoles {
	out := make([]map[uint64]bool, len(roleKindsByPriority()))
	for _, r := range rr {
		if out[r.kind] == nil {
			out[r.kind] = make(map[uint64]bool)
		}

		out[r.kind][r.id] = true
	}

	return out
}

// Returns slice of role types by priority
//
// Priority is important here. We want to have
// stable RBAC check behaviour and ability
// to override allow/deny depending on how niche the role (type) is:
//  - bypass always stake precedence
//  - context (eg owners) are more niche than common
func roleKindsByPriority() []roleKind {
	return []roleKind{
		BypassRole,
		ContextRole,
		CommonRole,
		AuthenticatedRole,
		AnonymousRole,
	}
}

// compare list of session roles (ids) with preloaded roles and calculate the final list
func getContextRoles(sRoles []uint64, res Resource, preloadedRoles []*role) (out partRoles) {
	var (
		mm   = slice.ToUint64BoolMap(sRoles)
		attr = make(map[string]interface{})
	)

	if ar, ok := res.(resourceDicter); ok {
		// if resource implements Dict() fn, we can use it to
		// collect attributes, used for expr. evaluation and contextual role gathering
		attr = ar.Dict()

	}

	out = make([]map[uint64]bool, len(roleKindsByPriority()))
	for _, r := range preloadedRoles {
		if r.kind == ContextRole {
			if r.check == nil {
				// expression not defined, skip contextual role
				continue
			}

			if !r.check(attr) {
				// add role to the list ONLY of expression evaluated true
				continue
			}
		} else if !mm[r.id] {
			// skip all other types of roles that user from session is not member of
			continue
		}

		if out[r.kind] == nil {
			out[r.kind] = make(map[uint64]bool)
		}

		out[r.kind][r.id] = true
	}

	return
}
