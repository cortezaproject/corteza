package automation

import (
	"github.com/cortezaproject/corteza/server/codegen/schema"
)

component: schema.#component & {
	handle: "automation"

	resources: {
		"workflow": workflow
		"session":  session
		"trigger":  trigger

		"ng-automation": ng_automation
		"trigger-definition": trigger_definition
	}

	rbac: operations: {
		"grant": description:                        "Manage automation permissions"
		"workflow.create": description:              "Create workflows"
		"triggers.search": description:              "List, search or filter triggers"
		"sessions.search": description:              "List, search or filter sessions"
		"workflows.search": description:             "List, search or filter workflows"
		"ng-automation.create": description:         "Create workflows"
		"ng-automations.search": description:        "List, search or filter workflows"
		"trigger-definition.create": description:    "Create trigger definitions"
		"trigger-definitions.search": description:   "List, search or filter trigger definitions"
		"resource-translations.manage": description: "List, search, create, or update resource translations"
	}
}
