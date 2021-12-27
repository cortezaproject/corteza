package system

import (
	"github.com/cortezaproject/corteza-server/def/schema"
)

user: schema.#resource & {
	//  fields: {
	//   ID: schema.IdField
	//   handle: schema.HandleField
	//   email: { unique: true }
	//   kind: {}
	//   meta: {
	//    note: type: string
	//    sub: {
	//     sub: { "non-unique-string-named-sub": {} }
	//    }
	//   }
	//  }

	rbac: {
		operations: {
			"read": description:         "Read user"
			"update": description:       "Update user"
			"delete": description:       "Delete user"
			"suspend": description:      "Suspend user"
			"unsuspend": description:    "Unsuspend user"
			"email.unmask": description: "Unmask email"
			"name.unmask": description:  "Unmask name"
			"impersonate": description:  "Impersonate user"
		}
	}
}
