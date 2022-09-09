package codegen

import (
	"github.com/cortezaproject/corteza-server/app"
	"github.com/cortezaproject/corteza-server/codegen/schema"
)


_dalModelFn: {
	res="res": schema.#Resource

	"var":     string | *"\(res.expIdent)"
	"resType": string | *"types.\(res.expIdent)ResourceType"

	"ident": res.model.ident

	"attributes": [
		for attr in res.model.attributes if attr.dal != _|_ {
			attr

			"dal": {
				attr.dal

				// Only these field types support "has-default" flag
				if (attr.dal.type & ( "ID" | "Ref" | "Number" | "Boolean" | "Enum" | "JSON" )) != _|_ {
					"hasDefault": attr.dal.default != _|_
				}

				if attr.dal.default != _|_ {
					"quotedDefault": attr.dal.type == "String"
				}
			}
		}
	]

	"indexes": {[string]: close({
		"ident":      string
		"type":       string
		"unique":     bool | *false
		"predicate"?: string
		"fields": [{
			"attribute": string
			"modifers"?: [string, ...]
			"sort"?:		 string
			"nulls"?:    string
		}, ...]
	})} | *null

	if res.model.indexes != _|_ {
		"indexes": {
			for index in res.model.indexes {
				"\(index.ident)": {

				  if index.primary == true {
				  	"ident": "PRIMARY"
				  }

				  if index.primary == false {
				  	"ident": index.ident
						"unique": bool | index.unique
				  }

				  "type": index.type

					if index.predicate != _|_ {
				  	"predicate": index.predicate
				  }

				  "fields": [
				  	for field in index.fields {
				  		"attribute": res.model.attributes[field.attribute].expIdent

				  		if field.modifiers != _|_ {
				  			"modifiers": field.modifiers
				  		}
				  		if field.sort != _|_ {
				  			"sort": field.sort
				  		}
				  		if field.nulls != _|_ {
				  			"nulls": field.nulls
				  		}
				  	},
				  ]
				}
			}
		}
	}
},

[...schema.#codegen] &
[
	{
		template: "gocode/dal/$component_model.go.tpl"
		output:   "system/model/corteza.gen.go"
		payload: {
			package: "model"

			imports: [
				for res in app.resources if (res.model.attributes != _|_)  {
					"\(res.package.ident)type \"\(res.package.import)\"",

				}
			]

			// Operation/resource validators, grouped by resource
			models: {
				for res in app.resources if (res.model.attributes != _|_) {
					"\(res.ident)": {
					  _out: _dalModelFn & { "res": res }

					  "var":        _out.var
					  "resType":    "\(res.package.ident)type.\(res.expIdent)ResourceType"
					  "ident":      _out.ident
					  "attributes": _out.attributes
					  "indexes":    _out.indexes
					},
				},
			}
		}
	},
	for cmp in app.corteza.components {
		template: "gocode/dal/$component_model.go.tpl"
		output:   "\(cmp.ident)/model/models.gen.go"
		payload: {
			package: "model"

			imports: [
				"\"github.com/cortezaproject/corteza-server/\(cmp.ident)/types\"",
			]

			// Operation/resource validators, grouped by resource
			models: {
				for res in cmp.resources if (res.model.attributes != _|_) {
					"\(res.ident)": {

						_out: _dalModelFn & { "res": res }

					  "var":        _out.var
					  "resType":    _out.resType
					  "ident":      _out.ident
			  		"attributes": _out.attributes
				  	"indexes":    _out.indexes
				  }
				},
			}
		}
	},

	for cmp in app.corteza.components {
		template: "gocode/dal/$component_init.go.tpl"
		output:   "\(cmp.ident)/model/init.gen.go"
		payload: {
			package: "model"
		}
	},
]
