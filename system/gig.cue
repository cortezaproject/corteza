package system

import (
	"github.com/cortezaproject/corteza-server/codegen/schema"
)

gig: schema.#resource & {
	rbac: {
		operations: {
			read: description:     "Read gig"
			update: description:   "Update gig"
			delete: description:   "Delete gig"
			undelete: description: "Undelete gig"
			exec: description:     "Execute gig"
		}
	}
}
