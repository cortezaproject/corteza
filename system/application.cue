package system

import (
	"github.com/cortezaproject/corteza-server/codegen/schema"
)

application: schema.#resource & {
	rbac: {
		operations: {
			read:
				description: "Read application"
			update:
				description: "Update application"
			delete:
				description: "Delete application"
		}
	}
}
