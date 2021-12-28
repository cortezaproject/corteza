package codegen

import (
	"github.com/cortezaproject/corteza-server/app"
	"github.com/cortezaproject/corteza-server/codegen/schema"
)

rbacTypes:
	[...schema.#codegen] &
	[
		for cmp in app.corteza.components {
			template: "gocode/rbac/types.go.tpl"
			output:   "\(cmp.ident)/types/rbac.gen.go"
			payload: {
				package: "types"

				cmpIdent: cmp.ident
				// Operation/resource validators, grouped by resource
				types: [
					for res in cmp.resources {
						const:   "\(res.expIdent)ResourceType"
						type:    res.fqrn
						resFunc: "\(res.expIdent)RbacResource"
						tplFunc: "\(res.expIdent)RbacResourceTpl"
						attFunc: "\(res.expIdent)RbacAttributes"
						goType:  res.expIdent

						if len(res.parents) > 0 {
							references: [ for p in res.parents {p}, {param: "id", refField: "ID"}]
						}
					},
					{
						const:     "ComponentResourceType"
						type:      cmp.fqrn
						resFunc:   "ComponentRbacResource"
						tplFunc:   "ComponentRbacResourceTpl"
						attFunc:   "ComponentRbacAttributes"
						goType:    "Component"
						component: true
					},
				]
			}
		},
	]
