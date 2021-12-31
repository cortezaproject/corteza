package system

import (
	"github.com/cortezaproject/corteza-server/codegen/schema"
)

role: schema.#resource & {
	rbac: {
		operations: {
			read: description:             "Read role"
			update: description:           "Update role"
			delete: description:           "Delete role"
			"members.manage": description: "Manage members"
		}}}
