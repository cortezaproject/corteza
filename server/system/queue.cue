package system

import (
	"github.com/cortezaproject/corteza/server/codegen/schema"
)

queue: schema.#resource & {
	rbac: {
		operations: {
			"read": description:        "Read queue"
			"update": description:      "Update queue"
			"delete": description:      "Delete queue"
			"queue.read": description:  "Read from queue"
			"queue.write": description: "Write to queue"
		}
	}
}
