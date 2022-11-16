package federation

import (
	"github.com/cortezaproject/corteza/server/codegen/schema"
)

node: schema.#resource & {
	rbac: {
		operations: {
			"manage": description:        "Manage federation node"
			"module.create": description: "Create shared module"
		}
	}
}
