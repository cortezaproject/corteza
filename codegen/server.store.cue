package codegen

import (
	"strings"
	"list"
	"github.com/cortezaproject/corteza-server/app"
	"github.com/cortezaproject/corteza-server/codegen/schema"
)

_StoreResource: {
	res     = "res":     schema.#Resource
	typePkg = "typePkg": string

	result: {
		ident:          res.store.ident
		identPlural:    res.store.identPlural
		expIdent:       res.store.expIdent
		expIdentPlural: res.store.expIdentPlural
		goType:         "\(typePkg).\(res.expIdent)"
		goSetType:      "\(typePkg).\(res.expIdent)Set"
		goFilterType:   "\(typePkg).\(res.filter.expIdent)"

		struct: [ for f in res.struct if f.store {
			"ident":      f.ident
			"expIdent":   f.expIdent
			"storeIdent": f.storeIdent
			"name":       f.name
			"primaryKey": f.primaryKey
			"ignoreCase": f.ignoreCase
			"goType":     strings.Replace(f.goType, "types.", "\(typePkg).", 1)
		}]

		filter: {
			// query fields as defined in struct
			"query":        [ for name in res.filter.query   {res.struct[name]}],

			// filter by nil state as defined in filter
			"byNilState":   [ for name in res.filter.byNilState {res.filter.struct[name]}]

			// filter by false as defined in filter
			"byFalseState": [ for name in res.filter.byFalseState {res.filter.struct[name]}]

			// filter by value as defined in filter
			// @todo this should be pulled from the struct
			"byValue":      [ for name in res.filter.byValue {res.filter.struct[name]}]
			"byLabel": res.features.labels
		}

		auxIdent:  "aux\(expIdent)"
		auxStruct: struct

		features: {
			paging: res.features.paging
			sorting: res.features.sorting
			checkFn: res.features.checkFn
		}

		api: {
			if res.store.api != _|_ {
				_base: {
					"ident":         res.store.ident
					"expStoreIdent": res.store.expIdentPlural
					"goType":        goType
					"goFilterType":  goFilterType
					"auxIdent":      auxIdent
				}

				deleteByPK: {
					primaryKeys: [ for f in res.struct if f.primaryKey {f} ]
					_pkExpNames: strings.Join([ for f in primaryKeys { f.expIdent } ], "")
					"expFnIdent":  "Delete\(res.store.expIdent)By\(_pkExpNames)"
				}

				lookups: [
					for l in res.store.api.lookups {
						_base

						"expFnIdent": l.expIdent

						if (l.description != _|_) {
							description: "// \(l.expIdent) " + strings.Join(strings.Split(l.description, "\n"), "\n// ")
						}

						// Copy all relevant fields from the struct
						"args": [
							for name in l.fields {
								let f = res.struct[name]

								"ident":  f.ident
								"storeIdent":  f.storeIdent
								"goType": f.goType
								"ignoreCase": f.ignoreCase
							},
						]

						"nullConstraint": l.nullConstraint
						"returnType":     "\(goType)"
						"collectionFnIdent": "\(res.store.ident)Collection"
					},
				]

				// all additional store functions we need for this resource
				functions: [
					for f in res.store.api.functions {
						_base

						"expFnIdent": f.expIdent

						if (f.description != _|_) {
							description: "// \(f.expIdent) " + strings.Join(strings.Split(f.description, "\n"), "\n// ")
						}

						"args": [ for a in f.args {
							"ident":  a.ident
							"goType": strings.Replace(a.goType, "types.", "\(typePkg).", 1)
							"spread": a.spread
						}]
						"return": [ for r in f.return {strings.Replace(r, "types.", "\(typePkg).", 1)}]
					},
				]

				sortableFields: {
					_base

					"fnIdent": "sortable\(expIdent)Fields"

					fields: {
						for f in res.struct if f.sortable || f.unique || f.primaryKey {
							{
								"\(strings.ToLower(f.name))":  f.name
								"\(strings.ToLower(f.ident))": f.name
							}
						}
					}
				}

				collectCursorValues: {
					_base

					"fnIdent": "collect\(expIdent)CursorValues"

					fields: [ for f in res.struct if f.sortable || f.unique || f.primaryKey {f} ]
					primaryKeys: [ for f in res.struct if f.primaryKey {f} ]
				}

				checkConstraints: {
					_base

					"fnIdent": "check\(expIdent)Constraints"

					checks: [
						for lookup in res.store.api.lookups if lookup.constraintCheck {
							lookupFnIdent: lookup.expIdent
							fields: [ for name in lookup.fields {res.struct[name]}]
							nullConstraint: [
								for f in res.struct if list.Contains(lookup.nullConstraint, f.name) {
									"expIdent": f.expIdent
								},
							]
						},
					]
				}
			}
		}

		settings: res.store.settings
	}
}

// Codegen template payload, reused for multiple outputs
_payload: {
	package: string | *"store"

	imports: {
		// per-component type imports
		for cmp in app.corteza.components for res in cmp.resources if res.store != _|_ {
			"github.com/cortezaproject/corteza-server/\(cmp.ident)/types": "\(cmp.ident)Type"
		}

		for res in app.corteza.resources if res.store != _|_ {
			"\(res.package.import)": "\(res.package.ident)Type"
		}

		for cmp in app.corteza.components for res in cmp.resources for i in res.imports {
			"\(i.import)": ""
		}
	}

	types: {
		// for each resource in every store with store and actions defined
		for cmp in app.corteza.components for res in cmp.resources if res.store != _|_ {
			// use _Store resource as a function (https://cuetorials.com/patterns/functions/)
			// and pass res(ource) and type-package string in as "arguments"
			"\(res.store.ident)": { _StoreResource & { "res": res, "typePkg": "\(cmp.ident)Type" } }.result
		},

		for res in app.corteza.resources if res.store != _|_ {
			"\(res.store.ident)": { _StoreResource & { "res": res, "typePkg": "\(res.package.ident)Type" } }.result
		}
	}
}

[...schema.#codegen] &
[
	{
		"bulk": [
			{
				"template": "gocode/store/interfaces.go.tpl"
				"output":   "store/interfaces.gen.go"
			}, {
				"template": "gocode/store/rdbms/rdbms.go.tpl"
				"output":   "store/adapters/rdbms/rdbms.gen.go"
			}, {
				"template": "gocode/store/rdbms/aux_types.go.tpl"
				"output":   "store/adapters/rdbms/aux_types.gen.go"
			}, {
				"template": "gocode/store/rdbms/queries.go.tpl"
				"output":   "store/adapters/rdbms/queries.gen.go"
			}, {
				"template": "gocode/store/rdbms/filters.go.tpl"
				"output":   "store/adapters/rdbms/filters.gen.go"
			}, {
				"template": "gocode/store/tests/all.go.tpl"
				"output":   "store/tests/all_test.go"
			},
		]

		"payload":  { _payload }
	},
]
