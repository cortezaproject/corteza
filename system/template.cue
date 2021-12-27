package system

import (
	"github.com/cortezaproject/corteza-server/def/schema"
)

template: schema.#resource & {
	rbac: {
		operations: {
			read: description:   "Read template"
			update: description: "Update template"
			delete: description: "Delete template"
			render: description: "Render template"
		}
	}
}
