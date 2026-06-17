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

	fqrt: platform + "::" + handle

	envoy: {
		omit: bool | *false
	}

	// REST endpoint definitions (ported from the legacy per-component rest.yaml).
	// Consumed by server.rest.cue to generate rest/handlers, rest/request and
	// the opt-in rest/<entrypoint>.gen.go controllers.
	rest?: #rest

	// All known RBAC operations for this component
	rbac: #rbacComponent & {
		operations: {
			grant: {
				description: "Manage \(handle) permissions"
			}
		}
	}
}
