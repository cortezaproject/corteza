package system

import (
	"github.com/cortezaproject/corteza-server/def/schema"
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
