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
				for res in cmp.resources if (res.model.attributes != _|_)  {
					var:     "\(res.expIdent)"
					resType: "types.\(res.expIdent)ResourceType"

					ident: res.model.ident

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

			 		if res.model.indexes != _|_ {
						indexes: [
							for index in res.model.indexes {
								if index.primary {
									ident: "PRIMARY",
								}

								if !index.primary {
									ident: index.ident,
									unique: index.unique,
								}

								type: index.type,

								predicate?: index.predicate,

								fields: [
									for field in index.fields {
									  // craft a handy string that will yield a descriptive error
									  // when referencing an unexisting attribute
									  "model (\(res.model.ident)) index (\(index.ident)) field attribute reference (\(field.attribute)) validation":
									   	res.model.attributes[field.attribute].ident

									  "attribute":    res.model.attributes[field.attribute].expIdent
										"modifiers"?:   field.modifiers
										"sort"?:        field.sort
										"nulls"?:       field.nulls
									},
								]
							}
						]
					}
				},
			]
		}
	},
]
