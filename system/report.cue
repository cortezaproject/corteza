package system

import (
	"github.com/cortezaproject/corteza-server/codegen/schema"
)

report: schema.#resource & {
	rbac: {
		operations: {
			read: description:   "Read report"
			update: description: "Update report"
			delete: description: "Delete report"
			run: description:    "Run report"
		}
	}
	// locale:
	//   extended: true
	//   keys:
	//     - { path: name,    field: "Meta.Name" }
	//     - { path: description, field: "Meta.Description" }
	//     - { name: block title, path: "block.{{blockID}}.title", custom: true }
	//     - { name: block description, path: "block.{{blockID}}.description", custom: true }
}
