package federation

import (
	"github.com/cortezaproject/corteza/server/codegen/schema"
)

component: schema.#component & {
	handle: "federation"

	resources: {
		"node":           node
		"node-sync":      nodeSync
		"exposed-module": exposedModule
		"shared-module":  sharedModule
		"module-mapping": moduleMapping
	}

	rbac: operations: {
		"grant": description:           "Manage federation permissions"
		"pair": description:            "Pair federation nodes"
		"settings.read": description:   "Read settings"
		"settings.manage": description: "Manage settings"
		"node.create": description:     "Create new federation node"
		"nodes.search": description:    "List, search or filter federation nodes"
	}
}
