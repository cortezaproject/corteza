package schema

import (
	"strings"
)

#component: {
	ident:    #baseHandle
	expIdent: #expIdent | *strings.ToTitle(ident)
	label:    strings.ToTitle(ident)
	platform: #baseHandle

	resources: {
		[key=_]: {handle: key, "component": ident, "platform": platform} & #resource
	}

	// All known RBAC operations for this component
	rbac: #rbacComponent & {
		resource: type: platform + "::" + ident

		operations: {
			grant: {
				description: "Manage \(ident) permissions"
			}
		}
	}
}
