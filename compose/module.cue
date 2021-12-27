package compose

import (
	"github.com/cortezaproject/corteza-server/def/schema"
)

module: schema.#resource & {
	rbac: {
		resource: references: [ "namespaceID", "ID"]

		operations: {
			"read": {}
			"update": {}
			"delete": {}
			"record.create": description:  "Create record"
			"records.search": description: "List, search or filter records"
		}
	}

	//locale:
	//  resource:
	//    references: [ namespace, ID ]
	//
	//  extended: true
	//  keys:
	//    - name
}
