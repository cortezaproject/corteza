package codegen

import (
	"github.com/cortezaproject/corteza-server/app"
	"github.com/cortezaproject/corteza-server/codegen/schema"
)


[...schema.#codegen] &
[
	for cmp in app.corteza.components {
		template: "gocode/dal/$component_model.go.tpl"
		output:   "\(cmp.ident)/model/models.gen.go"
		payload: {
			package: "model"

			imports: [
				"\"github.com/cortezaproject/corteza-server/\(cmp.ident)/types\"",
			]

			cmpIdent: cmp.ident
			// Operation/resource validators, grouped by resource
			models: [
				for res in cmp.resources if res.model.attributes != _|_ {
					var:     "\(res.expIdent)"
					resType: "types.\(res.expIdent)ResourceType"

					ident:      res.model.ident
					attributes: [
						for attr in res.model.attributes {
							attr

							dal: {
								attr.dal

								if attr.dal.default != _|_ {
									quotedDefault: attr.dal.type == "String"
								}
							}
						}
					]
				},
			]
		}
	},
]
