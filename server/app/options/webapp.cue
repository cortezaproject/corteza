package options

import (
	"github.com/cortezaproject/corteza/server/codegen/schema"
)

webapp: schema.#optionsGroup & {
	handle: "webapp"

	options: {
		scss_dir_path: {
			description:  "Path to custom SCSS source files directory"
		}
	}
}
