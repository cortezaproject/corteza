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
		[key=_]: {"handle": key, "component": handle, "platform": platform} & #resource
	}

	fqrn: platform + "::" + handle

	// All known RBAC operations for this component
	rbac: #rbacComponent & {
		operations: {
			grant: {
				description: "Manage \(handle) permissions"
			}
		}
	}
}
