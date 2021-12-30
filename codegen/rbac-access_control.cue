package codegen

import (
	"github.com/cortezaproject/corteza-server/codegen/schema"
	"github.com/cortezaproject/corteza-server/app"
)

rbacAccessControl:
	[...schema.#codegen] &
	[
		for cmp in app.corteza.components {
			template: "gocode/rbac/$component_access_control.go.tpl"
			output:   "\(cmp.ident)/service/access_control.gen.go"
			payload: {
				package: "service"
				imports: [
					"\"github.com/cortezaproject/corteza-server/\(cmp.ident)/types\"",
				]

				// All possible RBAC operations on component and resources
				// flattened
				operations: [
					for res in cmp.resources for op in res.rbac.operations {
						"op":          op.handle
						const:         "types.\(res.expIdent)ResourceType"
						resFunc:       "types.\(res.expIdent)RbacResource"
						goType:        "types.\(res.expIdent)"
						description:   op.description
						checkFuncName: op.checkFuncName

						if len(res.parents) > 0 {
							references: [ for p in res.parents {p}, {param: "id", refField: "ID"}]
						}
					},
					for op in cmp.rbac.operations {
						"op":          op.handle
						const:         "types.ComponentResourceType"
						resFunc:       "types.ComponentRbacResource"
						goType:        "types.Component"
						description:   op.description
						checkFuncName: op.checkFuncName
						component:     true
					},
				]

				// Operation/resource validators, grouped by resource
				validation: [
					for res in cmp.resources {
						label:    res.ident
						const:    "types.\(res.expIdent)ResourceType"
						funcName: "rbac\(res.expIdent)ResourceValidator"
						if len(res.parents) > 0 {
							references: [ for p in res.parents {p.refField}, "ID"]
						}
						operations: [ for op in res.rbac.operations {op.handle}]
					},
					{
						label:    "\(cmp.ident) component"
						const:    "types.ComponentResourceType"
						funcName: "rbacComponentResourceValidator"
						operations: [ for op in cmp.rbac.operations {op.handle}]
					},
				]
			}
		},
	]
