package compose

import (
	"github.com/cortezaproject/corteza-server/def/schema"
)

page: schema.#resource & {
	rbac: {
		resource: references: [ "namespaceID", "ID"]

		operations: {
			"read": {}
			"update": {}
			"delete": {}
		}
	}

	//locale:
	//  resource:
	//    references: [ namespace, ID ]
	//
	//  extended: true
	//  keys:
	//    - title
	//    - description
	//    - { name: blockTitle,                 path: "pageBlock.{{blockID}}.title",                         custom: true }
	//    - { name: blockDescription,           path: "pageBlock.{{blockID}}.description",                   custom: true }
	//    - { name: blockAutomationButtonlabel, path: "pageBlock.{{blockID}}.button.{{buttonID}}.label", custom: true }
}
