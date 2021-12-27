package compose

import (
	"github.com/cortezaproject/corteza-server/def/schema"
)

chart: schema.#resource & {
	rbac: {
		resource: references: [ "namespaceID", "ID"]

		operations: {
			"read": {}
		  "update": {}
		  "delete": {}
		}
	}

	// locale:
	//   resource:
	//     references: [ namespace, ID ]
	//   keys:
	//     - name
}
