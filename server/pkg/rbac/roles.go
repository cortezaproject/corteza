package rbac

import (
	"github.com/cortezaproject/corteza/server/pkg/logger"
	"github.com/cortezaproject/corteza/server/pkg/slice"
	"go.uber.org/zap"
)

const (
	roleKinds = 5
)

type (
	ctxRoleCheckFn func(map[string]interface{}) bool

	// role information, adapted for the needs of RBAC package
	Role struct {
		// all RBAC rules refer to role ID
		id uint64

		// for debugging and logging
		handle string

		// role type that will allow us
		kind roleKind

		check ctxRoleCheckFn

		// compatible resource types
		crtypes map[string]bool
	}

	roleKind int

	partRoles [roleKinds]map[uint64]bool
)

const (
	CommonRole roleKind = iota
	AnonymousRole
	AuthenticatedRole
	ContextRole
	BypassRole
)

func (k roleKind) String() string {
	switch k {
	case BypassRole:
		return "bypass"
	case ContextRole:
		return "context"
	case CommonRole:
		return "common"
	case AuthenticatedRole:
		return "authenticated"
	case AnonymousRole:
		return "anonymous"
	default:
		panic("unknown role kind")
	}
}

func (k roleKind) Make(id uint64, handle string) *Role {
	return &Role{
		kind:   k,
		id:     id,
		handle: handle,
	}
}

func MakeContextRole(id uint64, handle string, fn ctxRoleCheckFn, tt ...string) *Role {
	return &Role{
		kind:    ContextRole,
		id:      id,
		handle:  handle,
		check:   fn,
		crtypes: slice.ToStringBoolMap(tt),
	}
}

// partitions roles by kind
func partitionRoles(rr ...*Role) partRoles {
	out := [roleKinds]map[uint64]bool{}
	for _, r := range rr {
		if out[r.kind] == nil {
			out[r.kind] = make(map[uint64]bool)
		}

		out[r.kind][r.id] = true
	}

	return out
}

func (p partRoles) LogFields() (ff []zap.Field) {
	for _, k := range []roleKind{BypassRole, ContextRole, CommonRole, AuthenticatedRole, AnonymousRole} {
		ii := make([]uint64, 0, len(p[k]))
		for r := range p[k] {
			ii = append(ii, r)
		}

		ff = append(ff, logger.Uint64s(k.String(), ii))
	}

	return
}

// counts roles per type
func statRoles(rr ...*Role) (stats map[roleKind]int) {
	stats = make(map[roleKind]int)
	for _, r := range rr {
		stats[r.kind]++
	}

	return
}

// compare list of session roles (ids) with preloaded roles and calculate the final list
func getContextRoles(s Session, res Resource, preloadedRoles ...*Role) (out partRoles) {
	var (
		mm    = slice.ToUint64BoolMap(s.Roles())
		scope = make(map[string]interface{})
	)

	out = [roleKinds]map[uint64]bool{}

	{
		// if this is an anonymous user (has one role of kind anonymous)
		// ensure role-kind integrity and skip complex checks
		for _, r := range preloadedRoles {
			if !mm[r.id] {
				// skip roles that are not in the security session
				// skip all other types of roles that user from session is not member of
				continue
			}

			if r.kind != AnonymousRole {
				continue
			}

			if out[AnonymousRole] == nil {
				out[AnonymousRole] = make(map[uint64]bool)
			}

			out[AnonymousRole][r.id] = true
		}

		if len(out[AnonymousRole]) > 0 {
			return
		}
	}

	if d, ok := res.(resourceDicter); ok {
		// if resource implements Dict() fn, we can use it to
		// collect attributes, used for expression evaluation and contextual role gathering
		scope["resource"] = d.Dict()
	}

	scope["userID"] = s.Identity()

	for _, r := range preloadedRoles {
		if r.kind == ContextRole {
			if hasWildcards(res.RbacResource()) {
				// if resource has wildcards, we can't use it for contextual role evaluation
				//
				// this exception causes RBAC trace requests that can have wildcard
				// resources to ignore this role
				//
				// without skipping contextual roles like this
				// check function on role is highly likely to fail to evaluate properly
				continue
			}

			if len(r.crtypes) == 0 || !r.crtypes[ResourceType(res.RbacResource())] {
				// resource type not compatible with this contextual role
				continue
			}

			if r.check == nil {
				// expression not defined, skip contextual role
				continue
			}

			if !r.check(scope) {
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
