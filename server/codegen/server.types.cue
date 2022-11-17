package codegen

import (
	"github.com/cortezaproject/corteza/server/app"
	"github.com/cortezaproject/corteza/server/codegen/schema"
)

[...schema.#codegen] &
[
	for cmp in app.corteza.components {
		template: "gocode/types/$component_resources.go.tpl"
		output:   "\(cmp.ident)/types/resources.gen.go"
		payload: {
			package: "types"

			cmpIdent: cmp.ident
			// Operation/resource validators, grouped by resource
			types: [
				for res in cmp.resources {
					const:   "\(res.expIdent)ResourceType"
					type:    res.fqrt
				},
				{
					const:     "ComponentResourceType"
					type:      cmp.fqrt
				},
			]
		}
	},
]
