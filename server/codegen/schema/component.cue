package schema

import (
	"strings"
)

#TypeEntry: #_base & {
	handle: handle, ident: ident, expIdent: expIdent

	noIdField:         bool   | *false
	labelResourceType: string | *""
}

// #TypeBundle is a lightweight type registry for packages that are not
// component-based (pkg/*, discovery/types). Carries only the Go package
// name and a map of type entries — no resources, RBAC, or model.
#TypeBundle: {
	package:   string
	// outputDir: path relative to server root where generated files are written
	outputDir: string
	types: {
		[key=_]: {"handle": key} & #TypeEntry
	}
}

#component: #_base & {
	// copy field values from #_base
	handle: handle, ident: ident, expIdent: expIdent

	label:    strings.ToTitle(ident)
	platform: #baseHandle

	resources: {
		[key=_]: {"handle": key, "component": handle, "platform": platform} & #Resource
	}

	// standalone types: need Set/label generation but are not full resources
	types: {
		[key=_]: {"handle": key} & #TypeEntry
	} | *{}

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
