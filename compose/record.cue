package compose

import (
	"github.com/cortezaproject/corteza-server/def/schema"
)

record: schema.#resource & {
	rbac: {
		resource: references: [ "namespaceID", "moduleID", "ID"]

		operations: {
			"read": {}
			"update": {}
			"delete": {}
		}
	}
}
