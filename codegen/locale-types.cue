package codegen

import (
	"github.com/cortezaproject/corteza-server/app"
	"github.com/cortezaproject/corteza-server/codegen/schema"
	"strings"
)

localeTypes:
	[...schema.#codegen] &
	[
		// wrapped with additional for loop to trim out templates with empty types list
		for tpl in [
			for cmp in app.corteza.components {
				template: "gocode/locale/$component_types.go.tpl"
				output:   "\(cmp.ident)/types/locale.gen.go"
				payload: {
					package: "types"

					resources: [
						for res in cmp.resources if res.locale != _|_ {
							expIdent: res.expIdent
							const:    res.locale.resource.const
							type:     res.locale.resource.type

							references: [ for p in res.parents {p}, {param: "id", refField: "ID"}]

							extended: res.locale.extended

							keys: [ for key in res.locale.keys if key.handlerFunc == _|_ {
								struct: key.struct
								field:  strings.ToTitle(key.name)

								path: strings.Join([ for p in key.expandedPath {
									if p.var {"{{\(p.part)}}"}
									if !p.var {p.part}
								}], ".")

								if !key.customHandler {
									fieldPath: strings.Join([ for p in key.expandedPath {
										strings.ToTitle(p.part)
									}], ".")
								}

								"extended": extended
								if key.decodeFunc != _|_ {decodeFunc: key.decodeFunc}
								if key.encodeFunc != _|_ {encodeFunc: key.encodeFunc}

							}]
						},
					]
				}
			},
			// skip empty type lists
		] if len(tpl.payload.resources) > 0 {tpl}]
