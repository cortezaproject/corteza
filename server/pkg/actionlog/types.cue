package actionlog

import "github.com/crusttech/human/server/codegen/schema"

bundle: schema.#TypeBundle & {
	package:   "actionlog"
	outputDir: "pkg/actionlog"
	types: {
		"action": {}
	}
}
