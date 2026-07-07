package types

import "github.com/crusttech/human/server/codegen/schema"

bundle: schema.#TypeBundle & {
	package:   "types"
	outputDir: "discovery/types"
	types: {
		"resource-activity": {}
	}
}
