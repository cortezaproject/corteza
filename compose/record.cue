package compose

import (
	"github.com/cortezaproject/corteza-server/def/schema"
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
