package options

import (
	"github.com/cortezaproject/corteza/server/codegen/schema"
)

seeder: schema.#optionsGroup & {
	handle: "seeder"
	options: {
		log_enabled: {
			type:        "bool"
			description: "Enable extra logging // fixme add some more description"
		}
	}
	title: "Seeder"
}
