package types

import "github.com/crusttech/human/server/codegen/schema"

bundle: schema.#TypeBundle & {
	package:   "types"
	outputDir: "pkg/label/types"
	types: {
		"label": { noIdField: true }
	}
}
