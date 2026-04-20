package codegen

import (
	"github.com/crusttech/human/server/app"
	"github.com/crusttech/human/server/codegen/schema"
)

[...schema.#codegen] &
[
	for cmp in app.human.components {
		template: "gocode/rbac/$component_types_rbac.go.tpl"
		output:   "\(cmp.ident)/types/rbac.gen.go"
		payload: {
			package: "types"

			cmpIdent: cmp.ident
			// Operation/resource validators, grouped by resource
			types: [
				for res in cmp.resources if res.rbac != _|_ {
					const:   "\(res.expIdent)ResourceType"
					type:    res.fqrt
					resFunc: "\(res.expIdent)RbacResource"
					tplFunc: "\(res.expIdent)RbacResourceTpl"
					attFunc: "\(res.expIdent)RbacAttributes"
					goType:  res.expIdent

					references: [ for p in res.parents {p}, {param: "id", refField: "ID"}]
				},
				{
					const:     "ComponentResourceType"
					type:      cmp.fqrt
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
