package compose

import (
	"github.com/cortezaproject/corteza-server/codegen/schema"
)

record: schema.#resource & {
	parents: [
		{handle: "namespace"},
		{handle: "module"},
	]

	rbac: {
		operations: {
			"read": {}
			"update": {}
			"delete": {}
		}
	}
}
