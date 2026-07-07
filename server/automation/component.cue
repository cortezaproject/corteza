package automation

import (
	"github.com/crusttech/human/server/codegen/schema"
)

component: schema.#component & {
	handle: "automation"

	resources: {
		"workflow": workflow
		"session":  session
		"trigger":  trigger

		"ng-automation": ng_automation
	}

	types: {
		"trigger-constraint":    { noIdField: true }
		"workflow-path":         { noIdField: true }
		"workflow-issue":        { noIdField: true }
		"workflow-step":         {}
		"state":                 {}
		"ng-automation-issue":   { noIdField: true }
		"ng-automation-trigger": {}
		"ng-automation-step":    {}
		"ng-automation-path":    { noIdField: true }
	}

	rbac: operations: {
		"grant": description:                        "Manage automation permissions"
		"workflow.create": description:              "Create workflows"
		"triggers.search": description:              "List, search or filter triggers"
		"sessions.search": description:              "List, search or filter sessions"
		"workflows.search": description:             "List, search or filter workflows"
		"ng-automation.create": description:         "Create workflows"
		"ng-automations.search": description:        "List, search or filter workflows"
		"resource-translations.manage": description: "List, search, create, or update resource translations"
	}
}
