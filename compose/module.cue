package compose

import (
	"github.com/cortezaproject/corteza-server/codegen/schema"
)

module: schema.#resource & {
	handle: "module"
	parents: [
		{handle: "namespace"},
	]

	rbac: {
		operations: {
			"read": {}
			"update": {}
			"delete": {}
			"record.create": description:  "Create record"
			"records.search": description: "List, search or filter records"
		}
	}

	locale: {
		keys: {
			"name": {}
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
