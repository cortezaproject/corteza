package compose

import (
	"github.com/cortezaproject/corteza-server/codegen/schema"
)

component: schema.#component & {
	handle: "compose"

	resources: {
		"chart":        chart
		"module":       module
		"module-field": moduleField
		"namespace":    namespace
		"page":         page
		"record":       record
	}

	rbac: operations: {
		"settings.read": description:                "Read settings"
		"settings.manage": description:              "Manage settings"
		"namespace.create": description:             "Create namespace"
		"namespaces.search": description:            "List, search or filter namespaces"
		"resource-translations.manage": description: "List, search, create, or update resource translations"
	}
}
