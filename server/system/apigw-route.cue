package system

import (
	"github.com/cortezaproject/corteza/server/codegen/schema"
)

apigwRoute: schema.#resource & {
	rbac: {
		operations: {
			read: description:   "Read API Gateway route"
			update: description: "Update API Gateway route"
			delete: description: "Delete API Gateway route"
		}
	}
}
