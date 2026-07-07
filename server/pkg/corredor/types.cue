package corredor

import "github.com/crusttech/human/server/codegen/schema"

bundle: schema.#TypeBundle & {
	package:   "corredor"
	outputDir: "pkg/corredor"
	types: {
		"script": { noIdField: true }
	}
}
