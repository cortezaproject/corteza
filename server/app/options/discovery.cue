package options

import (
	"github.com/cortezaproject/corteza/server/codegen/schema"
)

discovery: schema.#optionsGroup & {
	handle: "discovery"
	options: {
		enabled: {
			type:          "bool"
			defaultGoExpr: "false"
			description:   "Enable discovery endpoints"
		}
		debug: {
			type:          "bool"
			defaultGoExpr: "false"
			description:   "Enable discovery related activity info"
		}
		corteza_domain: {
			type:        "string"
			description: "Indicates host of corteza compose webapp"
		}
		base_url: {
			type:        "string"
			description: "Indicates host of corteza discovery server"
		}
		jwt_secret: {
			type:        "string"
			description: "JWT secret shared with the discovery searcher for signing search tokens"
		}

	}
	title: "Discovery"
}
