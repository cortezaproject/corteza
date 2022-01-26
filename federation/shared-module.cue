package federation

import (
	"github.com/cortezaproject/corteza-server/codegen/schema"
)

sharedModule: schema.#resource & {
	parents: [
		{handle: "node"},
	]

	rbac: {
		operations: {
				"map": description: "Map shared module"
		}
	}
}
