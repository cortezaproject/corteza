package compose

import (
	"github.com/cortezaproject/corteza-server/codegen/schema"
)

namespace: schema.#resource & {
	rbac: {
		operations: {
			"read": {}
			"update": {}
			"delete": {}
			"manage": description:         "Access to namespace admin panel"
			"module.create": description:  "Create module on namespace"
			"modules.search": description: "List, search or filter module on namespace"
			"chart.create": description:   "Create chart on namespace"
			"charts.search": description:  "List, search or filter chart on namespace"
			"page.create": description:    "Create page on namespace"
			"pages.search": description:   "List, search or filter pages on namespace"
		}
	}

	locale: {
		keys: {
			name: {}
			metaSubtitle: {
				path: ["meta", "subtitle"]
			}
			metaDescription: {
				path: ["meta", "description"]
			}
		}
	}
}
