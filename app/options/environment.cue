package options

import (
	"github.com/cortezaproject/corteza-server/codegen/schema"
)

environment: schema.#optionsGroup & {
	handle: "environment"
	options: {
		environment: {
			default: "\"production\""
			env:     "ENVIRONMENT"
		}
	}
	title: "Environment"
}
