package automation

import (
	"github.com/cortezaproject/corteza/server/codegen/schema"
)

component: schema.#component & {
	handle: "automation"

	resources: {
		"workflow": workflow
	}

	rbac: operations: {
		"grant": description:                        "Manage automation permissions"
		"workflow.create": description:              "Create workflows"
		"triggers.search": description:              "List, search or filter triggers"
		"sessions.search": description:              "List, search or filter sessions"
		"workflows.search": description:             "List, search or filter workflows"
		"resource-translations.manage": description: "List, search, create, or update resource translations"
	}
}
