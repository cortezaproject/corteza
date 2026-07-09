package codegen

import (
	"strings"
	"github.com/crusttech/human/server/app"
	"github.com/crusttech/human/server/codegen/schema"
)

// Produces a single JSON file cataloguing every resource in every component,
// including model fields and nested struct/enum/slice type definitions.
// Pattern mirrors _TypeResource in server.types.cue.

[...schema.#codegen] &
[
	{
		template: "json/resource_schema.json.tpl"
		output:   "pkg/codegen/resource_schema.gen.json"
		syntax:   "json"
		payload: {
			components: [
				for cmp in app.human.components {
					ident:  cmp.ident
					label:  cmp.label
					resources: [
						for res in cmp.resources {
							handle:       res.handle
							expIdent:     res.expIdent
							resourceType: res.fqrt

							fields: [
								for attr in res.model.attributes
								if attr.sortableJSON == _|_ {
									expIdent: attr.expIdent
									goType:   strings.Replace(attr.goType, "types.", "", -1)

									_dalType: [ if attr.dal != _|_ {attr.dal.type}, "" ][0]
									_isID:   _dalType == "ID" && attr.name == "id"
									_isRef:  !_isID && attr.goType == "uint64" && (_dalType == "ID" || _dalType == "Ref")
									_isTime: attr.goType == "time.Time" || attr.goType == "*time.Time"

									jsonTag: [
										if attr.json != _|_ if (attr.json & string) != _|_ {"json:\"\(attr.json)\""},
										if attr.json != _|_ if (attr.json & bool) != _|_ if (attr.json & bool) {"json:\"\(attr.name)\""},
										if attr.json != _|_ if (attr.json & bool) == _|_ if (attr.json & string) == _|_ {
											let _j = {field: string | *attr.ident, omitEmpty: bool | *false, "string": bool | *false} & attr.json
											let _oe = [ if _j.omitEmpty {",omitempty"}, "" ][0]
											let _st = [ if _j.string {",string"}, "" ][0]
											"json:\"\(_j.field)\(_oe)\(_st)\""
										},
										if attr.json == _|_ if _isID {"json:\"\(res.ident)ID,string\""},
										if attr.json == _|_ if _isRef {"json:\"\(attr.ident),string,omitempty\""},
										if attr.json == _|_ if _isTime {"json:\"\(attr.ident),omitempty\""},
										"json:\"\(attr.ident)\"",
									][0]
								}
							]

							if res.types != _|_ {
								structDefs: [
									for d in res.types.defs
									if (d & {fields: _}) != _|_ {
										name: d.name
										doc:  [ if d.doc != _|_ {d.doc}, "" ][0]
										fields: [ for f in d.fields {
											let _base = [
												if f.type == _|_ {""},
												if (f.type & string) != _|_ {f.type},
												if (f.type & string) == _|_ {f.type.name},
											][0]
											let _ptr   = [ if f.ptr {"*"}, "" ][0]
											let _slice = [ if f.slice {"[]"}, "" ][0]
											name:    f.name
											goType:  [ if f.goType != _|_ {f.goType}, "\(_ptr)\(_slice)\(_base)" ][0]
											jsonTag: [ if f.json != _|_ {"json:\"\(f.json)\""}, "json:\"\(strings.ToCamel(f.name))\"" ][0]
											doc:     [ if f.doc != _|_ {f.doc}, "" ][0]
										}]
									}
								]
								enumDefs: [
									for d in res.types.defs
									if (d & {values: _}) != _|_ {
										name:   d.name
										goKind: d.kind
										doc:    [ if d.doc != _|_ {d.doc}, "" ][0]
										values: [ for v in d.values {
											constIdent: [ if strings.HasPrefix(v.ident, d.name) {v.ident}, d.name + v.ident ][0]
											value:      v.value
											doc:        [ if v.doc != _|_ {v.doc}, "" ][0]
										}]
									}
								]
								sliceDefs: [
									for d in res.types.defs
									if (d & {elem: _}) != _|_ {
										d.name
									}
								]
							}
						}
					]
				}
			]
		}
	},
]
