package compose

import (
	"github.com/cortezaproject/corteza-server/def/schema"
)

component: schema.#component & {
	ident: "compose"

	resources: {
		"namespace":    namespace
		"module":       module
		"module-field": moduleField
		"record":       record
		"page":         page
		"chart":        chart
	}

	rbac: operations: {
		"settings.read": description:                "Read settings"
		"settings.manage": description:              "Manage settings"
		"namespace.create": description:             "Create namespace"
		"namespaces.search": description:            "List, search or filter namespaces"
		"resource-translations.manage": description: "List, search, create, or update resource translations"
	}
}
