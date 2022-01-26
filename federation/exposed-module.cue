package federation

import (
	"github.com/cortezaproject/corteza-server/codegen/schema"
)

exposedModule: schema.#resource & {
	parents: [
		{handle: "node"},
	]

	rbac: {
		operations: {
				"manage": description: "Manage exposed module module"
		}
	}
}
