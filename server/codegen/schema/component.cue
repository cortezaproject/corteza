package schema

import (
	"strings"
)

#component: #_base & {
	// copy field values from #_base
	handle: handle, ident: ident, expIdent: expIdent

	label:    strings.ToTitle(ident)
	platform: #baseHandle

	resources: {
		[key=_]: {"handle": key, "component": handle, "platform": platform} & #Resource
	}

	// component-scoped merge of every resource's type registry -> one types pkg.
	// Same-named entries across resources unify (matching shapes) or conflict (error).
	_typeDefs: {
		for _, r in resources if r.types != _|_ {
			for tn, t in r.types.defs {(tn): t}
		}
	}

	fqrt: platform + "::" + handle

	envoy: {
		omit: bool | *false
	}

	// All known RBAC operations for this component
	rbac: #rbacComponent & {
		operations: {
			grant: {
				description: "Manage \(handle) permissions"
			}
		}
	}
}
