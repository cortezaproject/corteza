package system

import (
	"github.com/cortezaproject/corteza-server/def/schema"
)

authClient: schema.#resource & {
	rbac: {
		operations: {
			read: description:      "Read authorization client"
			update: description:    "Update authorization client"
			delete: description:    "Delete authorization client"
			authorize: description: "Authorize authorization client"
		}
	}
}
