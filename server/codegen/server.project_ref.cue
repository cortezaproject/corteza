package codegen

import (
	"github.com/crusttech/human/server/app"
	"github.com/crusttech/human/server/codegen/schema"
)

[...schema.#codegen] &
[
	for cmp in app.human.components {
		template: "gocode/project_ref/$component_types_project_ref.go.tpl"
		output:   "\(cmp.ident)/types/project_ref.gen.go"
		payload: {
			package: "types"

			cmpIdent: cmp.ident
			// Resources that denormalise their owning project (schema.ProjectRefField)
			types: [
				for res in cmp.resources if res.model.attributes.project_id != _|_ {
					goType: res.expIdent
					field:  res.model.attributes.project_id.expIdent
				},
			]
		}
	},
]
