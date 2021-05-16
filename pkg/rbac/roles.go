package rbac

type (
	// role information, adapted for the needs of RBAC package
	role struct {
		// all RBAC rules refer to role ID
		id uint64

		// for debugging and logging
		handle string

		// role type that will allow us
		kind roleKind
	}

	roleKind int

	roles []*role

	partRoles []map[uint64]bool
)

const (
	CommonRole = iota
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

func roleKindsByPriority() []int {
	return []int{
		BypassRole,
		ContextRole,
		CommonRole,
		AuthenticatedRole,
		AnonymousRole,
	}
}
